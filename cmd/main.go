package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ranjith035/morego/ai"
	"github.com/ranjith035/morego/core"
	"github.com/ranjith035/morego/drivers"
	pb "github.com/ranjith035/morego/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Concrete implementation of core.Container
type containerImpl struct {
	mu         sync.RWMutex
	components map[string]interface{}
}

func (c *containerImpl) Register(name string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.components[name] = value
}

func (c *containerImpl) Resolve(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.components[name]
	return val, ok
}

// Delegator that implements core.ActionEngine by forwarding calls to active driver session
type actionEngineDelegator struct {
	server *server
}

func (a *actionEngineDelegator) Name() string {
	return "ActionEngine"
}

func (a *actionEngineDelegator) Init(ctx context.Context) error {
	return nil
}

func (a *actionEngineDelegator) Shutdown(ctx context.Context) error {
	return nil
}

func (a *actionEngineDelegator) GetSource(ctx context.Context, sessionID string, format string) (string, error) {
	a.server.mu.Lock()
	driver, exists := a.server.sessions[sessionID]
	a.server.mu.Unlock()
	if !exists {
		return "", fmt.Errorf("session %s not found", sessionID)
	}
	return driver.GetSource(ctx, format)
}

func (a *actionEngineDelegator) Click(ctx context.Context, sessionID string, locator *core.Locator) error {
	return nil
}

func (a *actionEngineDelegator) ClickAt(ctx context.Context, sessionID string, pt core.Point) error {
	return nil
}

func (a *actionEngineDelegator) Fill(ctx context.Context, sessionID string, locator *core.Locator, text string) error {
	return nil
}

func (a *actionEngineDelegator) Swipe(ctx context.Context, sessionID string, start, end core.Point, duration time.Duration) error {
	return nil
}

func (a *actionEngineDelegator) Screenshot(ctx context.Context, sessionID string, locator *core.Locator) ([]byte, error) {
	return nil, nil
}

func (a *actionEngineDelegator) ExecuteScript(ctx context.Context, sessionID string, script string, args []string) (string, error) {
	return "", nil
}

type server struct {
	pb.UnimplementedSessionServiceServer
	pb.UnimplementedDriverServiceServer
	pb.UnimplementedAIServiceServer

	mu          sync.Mutex
	sessions    map[string]drivers.Driver
	idGen       int
	container   core.Container
	waitEngine  core.WaitEngine
	locatorEng  core.LocatorEngine
}

func newServer() *server {
	s := &server{
		sessions: make(map[string]drivers.Driver),
	}

	container := &containerImpl{components: make(map[string]interface{})}
	s.container = container

	// Wire ActionEngine
	ae := &actionEngineDelegator{server: s}
	container.Register("ActionEngine", ae)

	// Wire LocatorEngine
	le := core.NewLocatorEngine(container)
	container.Register("LocatorEngine", le)
	s.locatorEng = le

	// Wire WaitEngine
	we := core.NewWaitEngine(container)
	container.Register("WaitEngine", we)
	s.waitEngine = we

	return s
}

func toCoreLocator(pbLoc *pb.Locator) *core.Locator {
	if pbLoc == nil {
		return nil
	}

	strategy := pbLoc.Strategy.String()
	strategy = strings.TrimPrefix(strategy, "LOCATOR_STRATEGY_")

	coreLoc := &core.Locator{
		Strategy: core.LocatorStrategy(strategy),
		Selector: pbLoc.Selector,
		Index:    int(pbLoc.Index),
	}

	if pbLoc.Parent != nil {
		coreLoc.Parent = toCoreLocator(pbLoc.Parent)
	}

	for _, c := range pbLoc.Constraints {
		dir := c.Direction.String()
		dir = strings.TrimPrefix(dir, "RELATIVE_DIRECTION_")
		coreLoc.Constraints = append(coreLoc.Constraints, core.RelativeConstraint{
			Direction: core.RelativeDirection(dir),
			Target:    toCoreLocator(c.Target),
			Distance:  int(c.DistancePx),
		})
	}

	return coreLoc
}

// CreateSession starts a new session with capabilities on a target device.
func (s *server) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Default to "adb" driver for real Android automation
	driverType := req.DeviceId
	if driverType == "" || driverType == "pixel_6_pro" {
		driverType = "adb"
	}

	driver, err := drivers.DefaultRegistry.CreateDriver(driverType)
	if err != nil {
		// Fallback to adb if requested device is generic
		driver, err = drivers.DefaultRegistry.CreateDriver("adb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create driver: %v", err)
		}
	}

	caps := req.Capabilities
	if caps == nil {
		caps = make(map[string]string)
	}
	caps["device_id"] = req.DeviceId
	if req.AppId != "" {
		caps["app_id"] = req.AppId
		caps["bundle_id"] = req.AppId
	}

	err = driver.Connect(ctx, caps)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to connect driver: %v", err)
	}

	if req.AppId != "" {
		err = driver.LaunchApp(ctx, req.AppId, nil, nil)
		if err != nil {
			_ = driver.Disconnect(ctx)
			return nil, status.Errorf(codes.FailedPrecondition, "failed to launch app %q: %v", req.AppId, err)
		}
	}

	s.idGen++
	sessionID := fmt.Sprintf("session_%d_%d", s.idGen, time.Now().Unix())
	s.sessions[sessionID] = driver

	fmt.Printf("[Server] Session %s created for device %q using driver %T\n", sessionID, req.DeviceId, driver)

	return &pb.CreateSessionResponse{
		SessionId:    sessionID,
		Capabilities: caps,
		Status:       "ACTIVE",
	}, nil
}

// CloseSession terminates an active session.
func (s *server) CloseSession(ctx context.Context, req *pb.CloseSessionRequest) (*pb.CloseSessionResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.SessionId]
	if exists {
		delete(s.sessions, req.SessionId)
	}
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "session %s not found", req.SessionId)
	}

	_ = driver.Disconnect(ctx)
	fmt.Printf("[Server] Session %s closed\n", req.SessionId)

	return &pb.CloseSessionResponse{}, nil
}

// FindElement searches the active UI hierarchy using the WaitEngine and maps it back to a driver-native reference.
func (s *server) FindElement(ctx context.Context, req *pb.FindElementRequest) (*pb.FindElementResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.SessionId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "session %s not found", req.SessionId)
	}

	// 1. Convert protobuf locator to core locator
	coreLoc := toCoreLocator(req.Locator)

	timeout := time.Duration(req.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	// 2. Perform WaitEngine run to wait for element visiblity/stability
	elem, err := s.waitEngine.WaitForState(ctx, req.SessionId, coreLoc, core.WaitStateVisible, core.WaitOptions{
		Timeout: timeout,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "element not found or not stable: %v", err)
	}

	// 3. Map the resolved core.Element back to a driver-native element ID
	var elemID string
	var findErr error

	// If Android driver (adb), locate using the BOUNDS strategy
	if strings.Contains(fmt.Sprintf("%T", driver), "ADBDriver") {
		boundsStr := fmt.Sprintf("[%d,%d][%d,%d]", elem.Bounds.X, elem.Bounds.Y, elem.Bounds.X+elem.Bounds.Width, elem.Bounds.Y+elem.Bounds.Height)
		elemID, findErr = driver.FindElement(ctx, "BOUNDS", boundsStr)
	} else {
		// Fallback to absolute XPath strategy
		elemID, findErr = driver.FindElement(ctx, "XPATH", elem.ID)
	}

	if findErr != nil {
		return nil, status.Errorf(codes.Internal, "failed to map resolved element to driver reference: %v", findErr)
	}

	return &pb.FindElementResponse{
		Element: &pb.ElementDescriptor{
			ElementId: elemID,
		},
		Found: true,
	}, nil
}

// Click triggers touch click events.
func (s *server) Click(ctx context.Context, req *pb.ClickRequest) (*pb.ClickResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.DriverId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "driver session %s not found", req.DriverId)
	}

	var err error
	if req.GetElementId() != "" {
		err = driver.Click(ctx, req.GetElementId())
	} else if req.GetCoordinates() != nil {
		err = driver.ClickAt(ctx, int(req.GetCoordinates().X), int(req.GetCoordinates().Y))
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "click target is unspecified")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "click failed: %v", err)
	}

	return &pb.ClickResponse{Success: true}, nil
}

// Fill inputs keyboard characters.
func (s *server) Fill(ctx context.Context, req *pb.FillRequest) (*pb.FillResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.DriverId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "driver session %s not found", req.DriverId)
	}

	err := driver.Fill(ctx, req.ElementId, req.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fill failed: %v", err)
	}

	return &pb.FillResponse{Success: true}, nil
}

// Swipe executes drag/swipe gestures.
func (s *server) Swipe(ctx context.Context, req *pb.SwipeRequest) (*pb.SwipeResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.DriverId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "driver session %s not found", req.DriverId)
	}

	err := driver.Swipe(ctx, int(req.Start.X), int(req.Start.Y), int(req.End.X), int(req.End.Y), time.Duration(req.DurationMs)*time.Millisecond)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "swipe failed: %v", err)
	}

	return &pb.SwipeResponse{Success: true}, nil
}

// Screenshot captures the physical device screen.
func (s *server) Screenshot(ctx context.Context, req *pb.ScreenshotRequest) (*pb.ScreenshotResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.DriverId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "driver session %s not found", req.DriverId)
	}

	data, err := driver.Screenshot(ctx, req.ElementId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "screenshot failed: %v", err)
	}

	return &pb.ScreenshotResponse{
		ImageData: data,
		Success:   true,
	}, nil
}

// GetSource returns layout trees.
func (s *server) GetSource(ctx context.Context, req *pb.GetSourceRequest) (*pb.GetSourceResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.DriverId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "driver session %s not found", req.DriverId)
	}

	src, err := driver.GetSource(ctx, req.Format)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get source failed: %v", err)
	}

	return &pb.GetSourceResponse{
		SourceData: src,
		Success:    true,
	}, nil
}

// SelfHealLocator evaluates element trees to propose healed queries.
func (s *server) SelfHealLocator(ctx context.Context, req *pb.SelfHealRequest) (*pb.SelfHealResponse, error) {
	hist := ai.PastNodeHistory{
		Class:       "",
		Text:        req.OriginalLocator.Selector,
		ResourceID:  req.OriginalLocator.Selector,
		ContentDesc: "",
	}

	res, err := ai.HealLocator(req.OriginalLocator.Selector, hist, req.ViewHierarchy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "healing computation error: %v", err)
	}

	if res == nil {
		return &pb.SelfHealResponse{
			Reasoning: "No matching visual elements found for healing.",
		}, nil
	}

	var pbStrat pb.LocatorStrategy
	switch res.Strategy {
	case "RESOURCE_ID":
		pbStrat = pb.LocatorStrategy_LOCATOR_STRATEGY_RESOURCE_ID
	case "TEXT":
		pbStrat = pb.LocatorStrategy_LOCATOR_STRATEGY_TEXT
	case "ACCESSIBILITY_ID":
		pbStrat = pb.LocatorStrategy_LOCATOR_STRATEGY_ACCESSIBILITY_ID
	default:
		pbStrat = pb.LocatorStrategy_LOCATOR_STRATEGY_UNSPECIFIED
	}

	return &pb.SelfHealResponse{
		HealedLocator: &pb.Locator{
			Strategy: pbStrat,
			Selector: res.HealedSelector,
		},
		Confidence: float32(res.Confidence),
		Reasoning:  "AI healed locator successfully via selector property similarity matching.",
	}, nil
}

// SuggestLocators provides locator healing hints.
func (s *server) SuggestLocators(ctx context.Context, req *pb.SuggestLocatorsRequest) (*pb.SuggestLocatorsResponse, error) {
	suggestions, err := ai.SuggestLocators(req.ViewHierarchy)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid view hierarchy: %v", err)
	}

	resp := &pb.SuggestLocatorsResponse{
		Suggestions: make([]*pb.SuggestLocatorsResponse_Suggestion, 0, len(suggestions)),
	}

	for _, suggestion := range suggestions {
		resp.Suggestions = append(resp.Suggestions, &pb.SuggestLocatorsResponse_Suggestion{
			Locator: &pb.Locator{
				Strategy: mapAIStrategy(suggestion.Strategy),
				Selector: suggestion.Selector,
			},
			Reason:         suggestion.Reason,
			StabilityScore: suggestion.StabilityScore,
		})
	}

	return resp, nil
}

// AnalyzeFailure explains stack trace error sources.
func (s *server) AnalyzeFailure(ctx context.Context, req *pb.AnalyzeFailureRequest) (*pb.AnalyzeFailureResponse, error) {
	return &pb.AnalyzeFailureResponse{
		Analysis:       "Simulated Failure analysis matching logs.",
		RecommendedFix: "Ensure target element exists before timeout.",
	}, nil
}

func mapAIStrategy(strategy string) pb.LocatorStrategy {
	switch strategy {
	case "ACCESSIBILITY_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_ACCESSIBILITY_ID
	case "TEST_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_TEST_ID
	case "ROLE":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_ROLE
	case "TEXT":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_TEXT
	case "PLACEHOLDER":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_PLACEHOLDER
	case "LABEL":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_LABEL
	case "RESOURCE_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_RESOURCE_ID
	case "XPATH":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_XPATH
	default:
		return pb.LocatorStrategy_LOCATOR_STRATEGY_UNSPECIFIED
	}
}

func startHTTPServer(srv *server) {
	mux := http.NewServeMux()

	// Serve the static files from the ./inspector directory
	fs := http.FileServer(http.Dir("./inspector"))
	mux.Handle("/", fs)

	// API Endpoint to get the live XML source and Base64 screenshot
	mux.HandleFunc("/api/session/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}

		srv.mu.Lock()
		var sessionID string
		for id := range srv.sessions {
			sessionID = id
			break
		}
		srv.mu.Unlock()

		if sessionID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "No active session found"}`))
			return
		}

		srv.mu.Lock()
		driver := srv.sessions[sessionID]
		srv.mu.Unlock()

		// 1. Fetch XML layout source
		xmlData, err := driver.GetSource(r.Context(), "xml")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to get layout source: %v", err)})
			return
		}

		// 2. Fetch screenshot
		imgBytes, err := driver.Screenshot(r.Context(), "")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to get screenshot: %v", err)})
			return
		}

		imgBase64 := base64.StdEncoding.EncodeToString(imgBytes)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sessionId":  sessionID,
			"xml":        xmlData,
			"screenshot": imgBase64,
		})
	})

	// API Endpoint to trigger click at specific coordinates
	mux.HandleFunc("/api/session/action/click", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}

		var req struct {
			SessionID string `json:"sessionId"`
			X         int    `json:"x"`
			Y         int    `json:"y"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		srv.mu.Lock()
		driver, exists := srv.sessions[req.SessionID]
		srv.mu.Unlock()
		if !exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Session not found"}`))
			return
		}

		err := driver.ClickAt(r.Context(), req.X, req.Y)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	})

	fmt.Println("Mobile Automation Web Inspector dashboard available at http://localhost:8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		fmt.Printf("HTTP Server error: %v\n", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	srv := newServer()

	pb.RegisterSessionServiceServer(s, srv)
	pb.RegisterDriverServiceServer(s, srv)
	pb.RegisterAIServiceServer(s, srv)

	// Start the Live Web Inspector HTTP server
	go startHTTPServer(srv)

	fmt.Println("Mobile Automation Server listening on gRPC port 50051...")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
