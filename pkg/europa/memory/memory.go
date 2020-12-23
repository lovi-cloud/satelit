package memory

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lovi-cloud/satelit/pkg/datastore"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/lovi-cloud/satelit/pkg/europa"
)

// A Memory is backend of europa by in-memory for testing.
type Memory struct {
	// im-memory storage engine (ex: Dorado)
	Volumes        map[string]europa.Volume // ID: Volume
	Mu             sync.RWMutex
	BackendendName string

	ds datastore.Datastore
}

// New create on memory backend
func New(ds datastore.Datastore) *Memory {
	vs := map[string]europa.Volume{}

	return &Memory{
		Volumes:        vs,
		ds:             ds,
		BackendendName: "memory",
	}
}

// CreateVolume write volume information to in-memory
func (m *Memory) CreateVolume(ctx context.Context, name uuid.UUID, capacityGB int) (*europa.Volume, error) {
	v := europa.Volume{
		ID:          name.String(),
		CapacityGB:  uint32(capacityGB),
		BackendName: m.BackendendName,
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	if err := m.ds.PutVolume(ctx, v); err != nil {
		return nil, err
	}

	return &v, nil
}

// CreateVolumeFromImage write volume info to in-memory
func (m *Memory) CreateVolumeFromImage(ctx context.Context, name uuid.UUID, capacityGB int, imageID uuid.UUID) (*europa.Volume, error) {
	v := europa.Volume{
		ID:          name.String(),
		CapacityGB:  uint32(capacityGB),
		BaseImageID: imageID,
		BackendName: m.BackendendName,
	}

	m.Mu.Lock()
	m.Volumes[v.ID] = v
	m.Mu.Unlock()

	if err := m.ds.PutVolume(ctx, v); err != nil {
		return nil, err
	}

	return &v, nil
}

// DeleteVolume delete volume in-memory
func (m *Memory) DeleteVolume(ctx context.Context, id string) error {
	m.Mu.Lock()
	delete(m.Volumes, id)
	m.Mu.Unlock()

	if err := m.ds.DeleteVolume(ctx, id); err != nil {
		return err
	}

	return nil
}

// ListVolume return list of volume in-memory
func (m *Memory) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	var ids []string

	m.Mu.RLock()
	for _, v := range m.Volumes {
		ids = append(ids, v.ID)
	}
	m.Mu.RUnlock()

	volumes, err := m.ds.ListVolume(ctx, ids)
	if err != nil {
		return nil, err
	}

	return volumes, nil
}

// GetVolume get volume in-memory
func (m *Memory) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	return m.ds.GetVolume(ctx, id)
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

	if err := m.ds.PutVolume(ctx, *newV); err != nil {
		return 0, "", err
	}

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

	if err := m.ds.PutVolume(ctx, *newV); err != nil {
		return err
	}

	return nil
}

// DetachVolumeSatelit detach volume from satelit server
func (m *Memory) DetachVolumeSatelit(ctx context.Context, hyperMetroPairID string, hostLUNID int) error {
	v, err := m.GetVolume(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get volume: %w", err)
	}
	err = m.DetachVolume(ctx, v.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetImage retrieves image
func (m *Memory) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	return m.ds.GetImage(imageID)
}

// ListImage return image from in-memory
func (m *Memory) ListImage() ([]europa.BaseImage, error) {
	return m.ds.ListImage()
}

// UploadImage upload image to in-memory
func (m *Memory) UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*europa.BaseImage, error) {
	u := uuid.NewV4()
	v, err := m.CreateVolume(ctx, u, imageSizeGB)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	hostLUNID, _, err := m.AttachVolumeSatelit(ctx, v.ID, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to attach volume: %w", err)
	}
	defer func() {
		m.DetachVolumeSatelit(ctx, v.ID, hostLUNID)
	}()

	bi := &europa.BaseImage{
		UUID:          u,
		Name:          name,
		Description:   description,
		CacheVolumeID: v.ID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := m.ds.PutImage(*bi); err != nil {
		return nil, err
	}

	return bi, nil
}

// DeleteImage delete from in-memory
func (m *Memory) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return m.ds.DeleteImage(id)
}
