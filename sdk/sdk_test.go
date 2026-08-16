package sdk

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/ranjith035/morego/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockSessionServer struct {
	pb.UnimplementedSessionServiceServer
	lastSessionID string
}

func (s *mockSessionServer) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	s.lastSessionID = "mock_session_123"
	return &pb.CreateSessionResponse{
		SessionId: s.lastSessionID,
	}, nil
}

func (s *mockSessionServer) CloseSession(ctx context.Context, req *pb.CloseSessionRequest) (*pb.CloseSessionResponse, error) {
	return &pb.CloseSessionResponse{}, nil
}

type mockDriverServer struct {
	pb.UnimplementedDriverServiceServer
	lastElementID string
	clickedID     string
	filledID      string
	filledVal     string
	swiped        bool
}

func (s *mockDriverServer) FindElement(ctx context.Context, req *pb.FindElementRequest) (*pb.FindElementResponse, error) {
	s.lastElementID = "elem_abc_999"
	return &pb.FindElementResponse{
		Element: &pb.ElementDescriptor{
			ElementId: s.lastElementID,
		},
		Found: true,
	}, nil
}

func (s *mockDriverServer) Click(ctx context.Context, req *pb.ClickRequest) (*pb.ClickResponse, error) {
	s.clickedID = req.GetElementId()
	return &pb.ClickResponse{Success: true}, nil
}

func (s *mockDriverServer) Fill(ctx context.Context, req *pb.FillRequest) (*pb.FillResponse, error) {
	s.filledID = req.ElementId
	s.filledVal = req.Value
	return &pb.FillResponse{Success: true}, nil
}

func (s *mockDriverServer) Swipe(ctx context.Context, req *pb.SwipeRequest) (*pb.SwipeResponse, error) {
	s.swiped = true
	return &pb.SwipeResponse{Success: true}, nil
}

func TestSDKFluentAPI(t *testing.T) {
	// 1. Setup local gRPC listener
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	sessionSrv := &mockSessionServer{}
	driverSrv := &mockDriverServer{}

	pb.RegisterSessionServiceServer(server, sessionSrv)
	pb.RegisterDriverServiceServer(server, driverSrv)

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("Mock server exited: %v", err)
		}
	}()
	defer server.Stop()

	// 2. Establish SDK Connection using DialOption to connect to bufconn listener
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	// Wrap connections inside our SDK Device structure
	device := &Device{
		conn:          conn,
		sessionClient: pb.NewSessionServiceClient(conn),
		driverClient:  pb.NewDriverServiceClient(conn),
		deviceID:      "pixel_6_pro",
	}

	// 3. Execute Fluent API Chain
	session, err := device.NewSession(ctx, "com.example.app", nil)
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}

	if session.ID() != "mock_session_123" {
		t.Errorf("Expected session ID 'mock_session_123', got %q", session.ID())
	}

	// Click element
	err = session.Locator("ACCESSIBILITY_ID", "submit_btn").Click(ctx)
	if err != nil {
		t.Errorf("Click action failed: %v", err)
	}
	if driverSrv.clickedID != "elem_abc_999" {
		t.Errorf("Mock server did not receive correct click element ID: %q", driverSrv.clickedID)
	}

	// Fill element
	err = session.Locator("TEXT", "Enter text").Fill(ctx, "John Doe")
	if err != nil {
		t.Errorf("Fill action failed: %v", err)
	}
	if driverSrv.filledID != "elem_abc_999" || driverSrv.filledVal != "John Doe" {
		t.Errorf("Mock server did not receive correct fill inputs: ID=%q, Val=%q", driverSrv.filledID, driverSrv.filledVal)
	}

	// Swipe
	err = session.Swipe(ctx, 10, 20, 10, 80, 300*time.Millisecond)
	if err != nil {
		t.Errorf("Swipe action failed: %v", err)
	}
	if !driverSrv.swiped {
		t.Error("Mock server did not receive Swipe call")
	}

	// Close Session
	err = session.Close(ctx)
	if err != nil {
		t.Errorf("Close Session failed: %v", err)
	}
}

func TestSDKNewLocators(t *testing.T) {
	session := &Session{}
	
	loc := session.GetByText("Save").Above(session.GetByRole("Button")).Nth(1)
	
	protoLoc := loc.toProto()
	
	if protoLoc.Strategy != pb.LocatorStrategy_LOCATOR_STRATEGY_TEXT {
		t.Errorf("Expected strategy TEXT, got %v", protoLoc.Strategy)
	}
	if protoLoc.Selector != "Save" {
		t.Errorf("Expected selector 'Save', got %v", protoLoc.Selector)
	}
	if protoLoc.Index != 1 {
		t.Errorf("Expected index 1, got %d", protoLoc.Index)
	}
	
	if len(protoLoc.Constraints) != 1 {
		t.Fatalf("Expected 1 constraint, got %d", len(protoLoc.Constraints))
	}
	
	constraint := protoLoc.Constraints[0]
	if constraint.Direction != pb.RelativeDirection_RELATIVE_DIRECTION_ABOVE {
		t.Errorf("Expected direction ABOVE, got %v", constraint.Direction)
	}
	if constraint.Target.Strategy != pb.LocatorStrategy_LOCATOR_STRATEGY_ROLE || constraint.Target.Selector != "Button" {
		t.Errorf("Expected target Strategy ROLE and Selector 'Button', got Strategy %v Selector %q", constraint.Target.Strategy, constraint.Target.Selector)
	}
}

