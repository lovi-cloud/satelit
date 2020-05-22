package memory

import (
	"context"
	"errors"
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

// GetImage return image by id from on memory
func (m *Memory) GetImage(imageID string) (*europa.BaseImage, error) {
	m.mutex.Lock()
	i, ok := m.images[imageID]
	m.mutex.Unlock()

	if ok == false {
		return nil, errors.New("not found")
	}

	return &i, nil
}

// GetImages return all images
func (m *Memory) GetImages() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	m.mutex.Lock()
	for _, v := range m.images {
		images = append(images, v)
	}
	m.mutex.Unlock()

	return images, nil
}

// PutImage write image
func (m *Memory) PutImage(image europa.BaseImage) error {
	m.mutex.Lock()
	m.images[string(image.ID)] = image
	m.mutex.Unlock()

	return nil
}

// DeleteImage delete image
func (m *Memory) DeleteImage(imageID string) error {
	m.mutex.Lock()
	delete(m.images, imageID)
	m.mutex.Unlock()

	return nil
}
