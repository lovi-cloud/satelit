package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/europa"
)

// A Memory is backend of europa by on memory for testing.
type Memory struct {
	Volumes map[string]europa.Volume // ID: Volume
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
	id := name.String()

	v := europa.Volume{
		ID:         id,
		HostName:   "",
		CapacityGB: uint32(capacity),
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	return &v, nil
}

// DeleteVolume delete volume on memory
func (m *Memory) DeleteVolume(ctx context.Context, id string) error {
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
func (m *Memory) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	m.Mu.RLock()
	v, ok := m.Volumes[id]
	m.Mu.RUnlock()
	if ok == false {
		return nil, errors.New("not found")
	}

	return &v, nil
}

// AttachVolume write attach information on memory
func (m *Memory) AttachVolume(ctx context.Context, id string, hostname string) error {
	v, err := m.GetVolume(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	if v.Attached == true {
		return errors.New("already attached")
	}

	m.Mu.Lock()
	newV := v
	newV.Attached = true
	newV.HostName = hostname
	m.Volumes[id] = *newV
	m.Mu.Unlock()

	return nil
}

// DetachVolume delete attach information on memory
func (m *Memory) DetachVolume(ctx context.Context, id string) error {
	v, err := m.GetVolume(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}

	newV := v
	newV.Attached = false
	newV.HostName = ""

	m.Mu.Lock()
	m.Volumes[id] = *newV
	m.Mu.Unlock()
	return nil
}
