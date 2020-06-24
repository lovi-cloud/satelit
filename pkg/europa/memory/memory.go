package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/satelit/pkg/europa"
)

// A Memory is backend of europa by in-memory for testing.
type Memory struct {
	Volumes map[string]europa.Volume    // ID: Volume
	Images  map[string]europa.BaseImage // ID: BaseImage
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

// CreateVolumeRaw write volume information to in-memory
func (m *Memory) CreateVolumeRaw(ctx context.Context, name uuid.UUID, capacityGB int) (*europa.Volume, error) {
	id := name.String()

	v := europa.Volume{
		ID:         id,
		CapacityGB: uint32(capacityGB),
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	return &v, nil
}

// CreateVolumeImage write volume info to in-memory
func (m *Memory) CreateVolumeImage(ctx context.Context, name uuid.UUID, capacityGB int, imageID string) (*europa.Volume, error) {
	id := name.String()

	v := europa.Volume{
		ID:          id,
		CapacityGB:  uint32(capacityGB),
		BaseImageID: imageID,
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	return &v, nil
}

// DeleteVolume delete volume in-memory
func (m *Memory) DeleteVolume(ctx context.Context, id string) error {
	m.Mu.Lock()
	delete(m.Volumes, "name")
	m.Mu.Unlock()

	return nil
}

// ListVolume return list of volume in-memory
func (m *Memory) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	var vs []europa.Volume

	m.Mu.RLock()
	for _, v := range m.Volumes {
		vs = append(vs, v)
	}
	m.Mu.RUnlock()

	return vs, nil
}

// GetVolume get volume in-memory
func (m *Memory) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	m.Mu.RLock()
	v, ok := m.Volumes[id]
	m.Mu.RUnlock()
	if ok == false {
		return nil, errors.New("not found")
	}

	return &v, nil
}

// AttachVolumeTeleskop write attach information in-memory
func (m *Memory) AttachVolumeTeleskop(ctx context.Context, id string, hostname string) (int, string, error) {
	return m.AttachVolume(ctx, id, hostname)
}

// AttachVolumeSatelit write attach information in-memory
func (m *Memory) AttachVolumeSatelit(ctx context.Context, id string, hostname string) (int, string, error) {
	return m.AttachVolume(ctx, id, hostname)
}

// AttachVolume write attach information in-memory
func (m *Memory) AttachVolume(ctx context.Context, id string, hostname string) (int, string, error) {
	v, err := m.GetVolume(ctx, id)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get volume: %w", err)
	}

	if v.Attached == true {
		return 0, "", errors.New("already attached")
	}

	m.Mu.Lock()
	newV := v
	newV.Attached = true
	newV.HostName = hostname
	m.Volumes[id] = *newV
	m.Mu.Unlock()

	return 1, "/dev/xxa", nil
}

// DetachVolume delete attach information in-memory
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

// GetImages return image from in-memory
func (m *Memory) GetImages() ([]europa.BaseImage, error) {
	var images []europa.BaseImage

	m.Mu.RLock()
	for _, v := range m.Images {
		images = append(images, v)
	}
	m.Mu.RUnlock()

	return images, nil
}

// UploadImage upload image to in-memory
func (m *Memory) UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*europa.BaseImage, error) {
	// TODO: implement
	return nil, nil
}

// DeleteImage delete from in-memory
func (m *Memory) DeleteImage(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}
