package grid

import (
	"errors"
	"sync"
)

// TunnelSession maps network adapters between local environments and cloud networks.
type TunnelSession struct {
	ID       string
	TenantID string
	Active   bool
}

// TunnelManager coordinates tunnels.
type TunnelManager struct {
	mu      sync.Mutex
	tunnels map[string]*TunnelSession
}

// NewTunnelManager constructs a TunnelManager.
func NewTunnelManager() *TunnelManager {
	return &TunnelManager{
		tunnels: make(map[string]*TunnelSession),
	}
}

// EstablishTunnel registers a network proxy tunnel.
func (tm *TunnelManager) EstablishTunnel(tunnelID string, tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if t, exists := tm.tunnels[tunnelID]; exists && t.Active {
		return errors.New("tunnel session is already active")
	}

	tm.tunnels[tunnelID] = &TunnelSession{
		ID:       tunnelID,
		TenantID: tenantID,
		Active:   true,
	}
	return nil
}

// CloseTunnel marks the network tunnel session inactive.
func (tm *TunnelManager) CloseTunnel(tunnelID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, exists := tm.tunnels[tunnelID]
	if !exists {
		return errors.New("tunnel session not found")
	}

	t.Active = false
	return nil
}

// IsTunnelActive checks the status of a tunnel.
func (tm *TunnelManager) IsTunnelActive(tunnelID string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t, exists := tm.tunnels[tunnelID]
	return exists && t.Active
}
