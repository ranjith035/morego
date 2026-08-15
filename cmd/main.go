package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ranjith035/morego/ai"
	"github.com/ranjith035/morego/drivers"
	pb "github.com/ranjith035/morego/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedSessionServiceServer
	pb.UnimplementedDriverServiceServer
	pb.UnimplementedAIServiceServer

	mu       sync.Mutex
	sessions map[string]drivers.Driver
	idGen    int
}

func newServer() *server {
	return &server{
		sessions: make(map[string]drivers.Driver),
	}
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

	err = driver.Connect(ctx, caps)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to connect driver: %v", err)
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

// FindElement searches the active UI hierarchy and returns a unique element reference.
func (s *server) FindElement(ctx context.Context, req *pb.FindElementRequest) (*pb.FindElementResponse, error) {
	s.mu.Lock()
	driver, exists := s.sessions[req.SessionId]
	s.mu.Unlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "session %s not found", req.SessionId)
	}

	strategy := req.Locator.Strategy.String()
	strategy = strings.TrimPrefix(strategy, "LOCATOR_STRATEGY_")

	elemID, err := driver.FindElement(ctx, strategy, req.Locator.Selector)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "element not found: %v", err)
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
	return &pb.SuggestLocatorsResponse{}, nil
}

// AnalyzeFailure explains stack trace error sources.
func (s *server) AnalyzeFailure(ctx context.Context, req *pb.AnalyzeFailureRequest) (*pb.AnalyzeFailureResponse, error) {
	return &pb.AnalyzeFailureResponse{
		Analysis:       "Simulated Failure analysis matching logs.",
		RecommendedFix: "Ensure target element exists before timeout.",
	}, nil
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

	fmt.Println("Mobile Automation Server listening on gRPC port 50051...")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
