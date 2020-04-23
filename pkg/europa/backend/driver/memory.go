package driver

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/europa"
)

type MemoryBackend struct {
	Values   []europa.Volume
	Attached map[string]bool
	Mu       sync.RWMutex
}

func NewMemoryBackend() (*MemoryBackend, error) {
	var vs []europa.Volume
	m := MemoryBackend{
		Values: vs,
		Mu:     sync.RWMutex{},
	}

	return &m, nil
}

func (m *MemoryBackend) CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	v := europa.Volume{
		ID:       name.String(),
		Attached: false,
		HostName: "",
		Capacity: capacity,
	}

	m.Mu.Lock()
	m.Values = append(m.Values, v)
	m.Mu.Unlock()

	return &v, nil
}
func (m *MemoryBackend) DeleteVolume(ctx context.Context, name uuid.UUID) error {
	i, _ := m.getVolumeByName(ctx, name)
	if i == -1 {
		return errors.New("not found")
	}

	s := append(m.Values[:i], m.Values[i+1:]...)

	m.Mu.Lock()
	m.Values = s
	m.Mu.Unlock()

	return nil
}
func (m *MemoryBackend) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	return m.Values, nil
}

func (m *MemoryBackend) AttachVolume(ctx context.Context, name uuid.UUID, hostname string) (*europa.Volume, error) {
	i, v := m.getVolumeByName(ctx, name)
	if i == -1 {
		return nil, errors.New("not found")
	}
	if v.Attached == true {
		return nil, errors.New("already attached")
	}

	newV := v
	newV.Attached = true
	newV.HostName = hostname

	m.Mu.Lock()
	m.Values[i] = *newV
	m.Mu.Unlock()

	return newV, nil
}

func (m *MemoryBackend) IsAttached(ctx context.Context, name uuid.UUID) (bool, error) {
	i, v := m.getVolumeByName(ctx, name)
	if i == -1 {
		return false, errors.New("not found")
	}
	return v.Attached, nil
}

func (m *MemoryBackend) getVolumeByName(ctx context.Context, name uuid.UUID) (index int, volume *europa.Volume) {
	for i, v := range m.Values {
		if strings.EqualFold(v.ID, name.String()) {
			return i, &v
		}
	}

	return -1, nil
}
