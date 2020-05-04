package memory

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/europa"
)

// A Memory is backend of europa by on memory for testing.
type Memory struct {
	Volumes map[string]europa.Volume
	Mu      sync.RWMutex
}

// New create on memory backend
func New() (*Memory, error) {
	vs := map[string]europa.Volume{}
	m := Memory{
		Volumes: vs,
		Mu:      sync.RWMutex{},
	}

	return &m, nil
}

// CreateVolume write volume information to on memory
func (m *Memory) CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	v := europa.Volume{
		ID:         name.String(),
		HostName:   "",
		CapacityGB: capacity,
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	return &v, nil
}

// DeleteVolume delete volume on memory
func (m *Memory) DeleteVolume(ctx context.Context, name uuid.UUID) error {
	m.Mu.Lock()
	delete(m.Volumes, "name")
	m.Mu.Unlock()

	return nil
}

// ListVolume return list of volume on memory
func (m *Memory) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	var vs []europa.Volume

	m.Mu.RLock()
	for _, v := range m.Volumes {
		vs = append(vs, v)
	}
	m.Mu.RUnlock()

	return vs, nil
}

// GetVolume get volume on memory
func (m *Memory) GetVolume(ctx context.Context, name uuid.UUID) (*europa.Volume, error) {
	m.Mu.RLock()
	v, ok := m.Volumes[name.String()]
	m.Mu.RUnlock()
	if ok == false {
		return nil, errors.New("not found")
	}

	return &v, nil
}

// AttachVolume write attach information on memory
func (m *Memory) AttachVolume(ctx context.Context, name uuid.UUID, hostname string) (*europa.Volume, error) {
	v, err := m.GetVolume(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get volume")
	}

	if v.Attached == true {
		return nil, errors.New("already attached")
	}

	m.Mu.Lock()
	newV := v
	newV.Attached = true
	newV.HostName = hostname
	m.Volumes[name.String()] = *newV
	m.Mu.Unlock()

	return newV, nil
}

// DetachVolume delete attach information on memory
func (m *Memory) DetachVolume(ctx context.Context, name uuid.UUID) error {
	v, err := m.GetVolume(ctx, name)
	if err != nil {
		return errors.Wrap(err, "failed to get volume")
	}

	newV := v
	newV.Attached = false
	newV.HostName = ""

	m.Mu.Lock()
	m.Volumes[name.String()] = *newV
	m.Mu.Unlock()
	return nil
}
