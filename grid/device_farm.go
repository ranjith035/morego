package grid

import (
	"errors"
	"sync"
)

type DeviceStatus string

const (
	StatusIdle    DeviceStatus = "idle"
	StatusBusy    DeviceStatus = "busy"
	StatusOffline DeviceStatus = "offline"
)

// DeviceDetails tracks attributes of physical or virtual machines in the cloud farm.
type DeviceDetails struct {
	ID               string
	Platform         string // android or ios
	Model            string
	Status           DeviceStatus
	ReservedByTenant string
}

// DeviceFarm coordinates device leases and reservation cycles.
type DeviceFarm struct {
	mu      sync.Mutex
	devices map[string]*DeviceDetails
}

// NewDeviceFarm instantiates a DeviceFarm registry.
func NewDeviceFarm() *DeviceFarm {
	return &DeviceFarm{
		devices: make(map[string]*DeviceDetails),
	}
}

// AddDevice adds a device descriptor into the registry.
func (df *DeviceFarm) AddDevice(d DeviceDetails) {
	df.mu.Lock()
	defer df.mu.Unlock()
	df.devices[d.ID] = &d
}

// AcquireDevice leases a matching idle device.
func (df *DeviceFarm) AcquireDevice(tenantID string, platform string) (*DeviceDetails, error) {
	df.mu.Lock()
	defer df.mu.Unlock()

	for _, d := range df.devices {
		if d.Platform == platform && d.Status == StatusIdle {
			d.Status = StatusBusy
			d.ReservedByTenant = tenantID
			return d, nil
		}
	}

	return nil, errors.New("no available idle devices matching request")
}

// ReleaseDevice returns a leased device back to the pool.
func (df *DeviceFarm) ReleaseDevice(deviceID string) error {
	df.mu.Lock()
	defer df.mu.Unlock()

	d, exists := df.devices[deviceID]
	if !exists {
		return errors.New("device not found in farm registry")
	}

	d.Status = StatusIdle
	d.ReservedByTenant = ""
	return nil
}

// DeviceCount returns the registry size.
func (df *DeviceFarm) DeviceCount() int {
	df.mu.Lock()
	defer df.mu.Unlock()
	return len(df.devices)
}
