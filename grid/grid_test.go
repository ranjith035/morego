package grid

import (
	"bytes"
	"testing"
)

func TestAuthAndRBAC(t *testing.T) {
	auth := NewAuthManager()

	// Register roles
	auth.RegisterKey("admin_secret_key", "tenant_a", RoleAdmin)
	auth.RegisterKey("viewer_secret_key", "tenant_b", RoleViewer)

	// 1. Authenticate correct key
	session, err := auth.Authenticate("admin_secret_key")
	if err != nil {
		t.Fatalf("Failed to authenticate valid key: %v", err)
	}
	if session.TenantID != "tenant_a" || session.Role != RoleAdmin {
		t.Errorf("Mismatched authenticated session: %+v", session)
	}

	// Authenticate invalid key
	_, err = auth.Authenticate("invalid_key")
	if err == nil {
		t.Fatal("Expected invalid API key to fail authentication, but it passed")
	}

	// 2. Validate RBAC Permissions
	adminSession := UserSession{Role: RoleAdmin}
	viewerSession := UserSession{Role: RoleViewer}

	// Admin executes Developer-level operation
	err = auth.CheckPermission(adminSession, RoleDeveloper)
	if err != nil {
		t.Errorf("Admin failed to pass Developer level permissions check: %v", err)
	}

	// Viewer executes Developer-level operation
	err = auth.CheckPermission(viewerSession, RoleDeveloper)
	if err == nil {
		t.Fatal("Expected Viewer to fail Developer level permissions check, but it passed")
	}
}

func TestDeviceFarmAllocations(t *testing.T) {
	farm := NewDeviceFarm()

	farm.AddDevice(DeviceDetails{
		ID:       "android_1",
		Platform: "android",
		Model:    "Pixel 7 Pro",
		Status:   StatusIdle,
	})
	farm.AddDevice(DeviceDetails{
		ID:       "ios_1",
		Platform: "ios",
		Model:    "iPhone 15 Pro",
		Status:   StatusIdle,
	})

	if count := farm.DeviceCount(); count != 2 {
		t.Errorf("Expected farm size 2, got %d", count)
	}

	// 1. Lease Android device
	leasedDev, err := farm.AcquireDevice("tenant_a", "android")
	if err != nil {
		t.Fatalf("Failed to acquire android device: %v", err)
	}
	if leasedDev.ID != "android_1" || leasedDev.Status != StatusBusy || leasedDev.ReservedByTenant != "tenant_a" {
		t.Errorf("Unexpected leased device attributes: %+v", leasedDev)
	}

	// 2. Try to lease same device type (Android) again - should fail
	_, err = farm.AcquireDevice("tenant_a", "android")
	if err == nil {
		t.Fatal("Expected second lease request for Android to fail, but it passed")
	}

	// 3. Return leased device
	err = farm.ReleaseDevice("android_1")
	if err != nil {
		t.Fatalf("Failed to release device: %v", err)
	}

	// 4. Try leasing Android again - should succeed now
	leasedDev, err = farm.AcquireDevice("tenant_a", "android")
	if err != nil {
		t.Fatalf("Failed to acquire device after releasing it: %v", err)
	}
	if leasedDev.ID != "android_1" {
		t.Errorf("Expected 'android_1', got %s", leasedDev.ID)
	}
}

func TestFrameBroadcaster(t *testing.T) {
	broadcaster := NewFrameBroadcaster()

	// Connect subscribers
	client1 := broadcaster.Subscribe()
	client2 := broadcaster.Subscribe()

	testFrame := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46} // Fake JPEG header

	// Broadcast
	broadcaster.Broadcast(testFrame)

	// Verify subscriber 1 received
	select {
	case frame := <-client1:
		if !bytes.Equal(frame, testFrame) {
			t.Errorf("Subscriber 1 received incorrect frame content")
		}
	default:
		t.Error("Subscriber 1 did not receive broadcasted frame")
	}

	// Verify subscriber 2 received
	select {
	case frame := <-client2:
		if !bytes.Equal(frame, testFrame) {
			t.Errorf("Subscriber 2 received incorrect frame content")
		}
	default:
		t.Error("Subscriber 2 did not receive broadcasted frame")
	}

	// Clean up
	broadcaster.Unsubscribe(client1)
	broadcaster.Unsubscribe(client2)
}

func TestTunnelManager(t *testing.T) {
	tm := NewTunnelManager()

	// 1. Establish tunnel
	err := tm.EstablishTunnel("tunnel_123", "tenant_abc")
	if err != nil {
		t.Fatalf("Failed to establish tunnel: %v", err)
	}
	if !tm.IsTunnelActive("tunnel_123") {
		t.Error("Tunnel is active but IsTunnelActive returned false")
	}

	// 2. Establish tunnel with same ID again - should fail
	err = tm.EstablishTunnel("tunnel_123", "tenant_abc")
	if err == nil {
		t.Fatal("Expected establishing active tunnel to fail, but it passed")
	}

	// 3. Close tunnel
	err = tm.CloseTunnel("tunnel_123")
	if err != nil {
		t.Fatalf("Failed to close tunnel: %v", err)
	}
	if tm.IsTunnelActive("tunnel_123") {
		t.Error("Tunnel closed but IsTunnelActive returned true")
	}
}
