package memory

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/europa"
)

type Memory struct {
	Volumes map[string]europa.Volume
	Mu      sync.RWMutex
}

func New() (*Memory, error) {
	var vs map[string]europa.Volume
	m := Memory{
		Volumes: vs,
		Mu:      sync.RWMutex{},
	}

	return &m, nil
}

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
func (m *Memory) DeleteVolume(ctx context.Context, name uuid.UUID) error {
	m.Mu.Lock()
	delete(m.Volumes, "name")
	m.Mu.Unlock()

	return nil
}
func (m *Memory) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	var vs []europa.Volume

	m.Mu.RLock()
	for _, v := range m.Volumes {
		vs = append(vs, v)
	}

	return vs, nil
}

func (m *Memory) AttachVolume(ctx context.Context, name uuid.UUID, hostname string) (*europa.Volume, error) {
	m.Mu.RLock()
	v, ok := m.Volumes[name.String()]
	if ok == false {
		return nil, errors.New("not found")
	}

	if v.Attached == true {
		return nil, errors.New("already attached")
	}
	m.Mu.RUnlock()

	m.Mu.Lock()
	newV := v
	newV.Attached = true
	newV.HostName = hostname
	m.Volumes[name.String()] = v
	m.Mu.Unlock()

	return &v, nil
}
