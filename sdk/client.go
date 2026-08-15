package sdk

import (
	"context"
	"fmt"

	pb "github.com/ranjith035/morego/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Device wraps the network connection to the Core Engine.
type Device struct {
	conn          *grpc.ClientConn
	sessionClient pb.SessionServiceClient
	driverClient  pb.DriverServiceClient
	aiClient      pb.AIServiceClient
	deviceID      string
}

// Connect dials the Go Core gRPC server and returns a Device wrapper.
func Connect(ctx context.Context, address string, deviceID string) (*Device, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Core Engine at %s: %w", address, err)
	}

	return &Device{
		conn:          conn,
		sessionClient: pb.NewSessionServiceClient(conn),
		driverClient:  pb.NewDriverServiceClient(conn),
		aiClient:      pb.NewAIServiceClient(conn),
		deviceID:      deviceID,
	}, nil
}

// Close releases the gRPC connection channel.
func (d *Device) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

// NewSession creates an execution session on the target app.
func (d *Device) NewSession(ctx context.Context, appID string, capabilities map[string]string) (*Session, error) {
	resp, err := d.sessionClient.CreateSession(ctx, &pb.CreateSessionRequest{
		DeviceId:     d.deviceID,
		AppId:        appID,
		Capabilities: capabilities,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		device:    d,
		sessionID: resp.SessionId,
		appID:     appID,
	}, nil
}

// Suggestion models locator recommendations.
type Suggestion struct {
	Strategy       string
	Selector       string
	StabilityScore float32
}

// SuggestLocators queries the AI engine for recommended stable locators.
func (d *Device) SuggestLocators(ctx context.Context, xmlHierarchy string) ([]Suggestion, error) {
	resp, err := d.aiClient.SuggestLocators(ctx, &pb.SuggestLocatorsRequest{
		SessionId:     d.deviceID,
		ViewHierarchy: xmlHierarchy,
	})
	if err != nil {
		return nil, fmt.Errorf("ai suggest locators failed: %w", err)
	}

	var suggestions []Suggestion
	for _, s := range resp.Suggestions {
		suggestions = append(suggestions, Suggestion{
			Strategy:       s.Locator.Strategy.String(),
			Selector:       s.Locator.Selector,
			StabilityScore: s.StabilityScore,
		})
	}

	return suggestions, nil
}
