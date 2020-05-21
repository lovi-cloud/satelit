package memory

import (
	"context"
	"sync"

	"github.com/whywaita/satelit/pkg/europa"
)

// A Memory is on memory datastore for testing.
type Memory struct {
	mutex *sync.Mutex

	images map[string]europa.BaseImage
}

// New create Memory
func New() *Memory {
	return &Memory{}
}

// GetIQN return IQN from on memory
func (m *Memory) GetIQN(ctx context.Context, hostname string) (string, error) {
	return "dummy iqn", nil
}

// PutImage write image
func (m *Memory) PutImage(image europa.BaseImage) error {
	m.mutex.Lock()
	m.images[image.ID] = image
	m.mutex.Unlock()

	return nil
}
