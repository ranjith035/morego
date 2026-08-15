package grid

import (
	"errors"
	"sync"
)

// UserRole defines tenancy permission limits.
type UserRole int

const (
	RoleViewer UserRole = iota
	RoleDeveloper
	RoleAdmin
)

// UserSession stores authenticated client contexts.
type UserSession struct {
	APIKey   string
	TenantID string
	Role     UserRole
}

// AuthManager regulates keys credentials validation and RBAC.
type AuthManager struct {
	mu   sync.RWMutex
	keys map[string]UserSession
}

// NewAuthManager constructs a thread-safe AuthManager.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		keys: make(map[string]UserSession),
	}
}

// RegisterKey associates an API key with client details.
func (am *AuthManager) RegisterKey(apiKey string, tenantID string, role UserRole) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.keys[apiKey] = UserSession{
		APIKey:   apiKey,
		TenantID: tenantID,
		Role:     role,
	}
}

// Authenticate verifies key membership.
func (am *AuthManager) Authenticate(apiKey string) (UserSession, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	session, exists := am.keys[apiKey]
	if !exists {
		return UserSession{}, errors.New("authentication failed: invalid api key")
	}
	return session, nil
}

// CheckPermission asserts user's access boundaries.
func (am *AuthManager) CheckPermission(session UserSession, requiredRole UserRole) error {
	if session.Role < requiredRole {
		return errors.New("authorization failed: insufficient permissions to perform action")
	}
	return nil
}
