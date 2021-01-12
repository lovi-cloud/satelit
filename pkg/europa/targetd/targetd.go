package targetd

import (
	"context"
	"fmt"

	"github.com/lovi-cloud/go-targetd/targetd"
	"github.com/lovi-cloud/satelit/internal/logger"
	"github.com/lovi-cloud/satelit/pkg/datastore"
	"github.com/lovi-cloud/satelit/pkg/europa"
	"go.uber.org/zap"

	uuid "github.com/satori/go.uuid"
)

// Targetd is implement open-iscsi/targetd
type Targetd struct {
	client    *targetd.Client
	datastore datastore.Datastore

	pool        targetd.Pool
	backendName string
}

// New is create instance of targetd.Client
func New(apiURL, username, password, poolName, backendName string, ds datastore.Datastore) (*Targetd, error) {
	client, err := targetd.New(
		apiURL,
		username,
		password,
		nil,
		zap.NewStdLog(logger.Logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client of targetd: %w", err)
	}

	pool, err := client.GetPool(context.Background(), poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve list of pool: %w", err)
	}

	return &Targetd{
		client:      client,
		datastore:   ds,
		pool:        *pool,
		backendName: backendName,
	}, nil
}

// CreateVolume create raw volume
func (t *Targetd) CreateVolume(ctx context.Context, name uuid.UUID, capacityGB int) (*europa.Volume, error) {
	size := capacityGB * 1024 * 1024 * 1024

	if err := t.client.CreateVolume(ctx, t.pool.Name, name.String(), size); err != nil {
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", name.String(), err)
	}
	v, err := t.client.GetVolume(ctx, t.pool.Name, name.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get volume (name: %s): %w", name.String(), err)
	}

	vol, err := t.toVolume(*v, false, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get europa.Volume: %w", err)
	}

	if err := t.datastore.PutVolume(ctx, *vol); err != nil {
		return nil, fmt.Errorf("failed to put volume to datastore (ID: %s): %w", vol.ID, err)
	}

	return vol, nil
}

// CreateVolumeFromImage create volume that copied image
func (t *Targetd) CreateVolumeFromImage(ctx context.Context, name uuid.UUID, capacityGB int, imageID uuid.UUID) (*europa.Volume, error) {
	return nil, nil
}

// DeleteVolume delete volume
func (t *Targetd) DeleteVolume(ctx context.Context, id string) error {
	// TODO: check attach status
	if err := t.client.DestroyVolume(ctx, t.pool.Name, id); err != nil {
		return fmt.Errorf("failed to delete volume (ID: %s): %w", id, err)
	}

	if err := t.datastore.DeleteVolume(ctx, id); err != nil {
		return fmt.Errorf("failed to delete volume from datastore (ID: %s): %w", id, err)
	}

	return nil
}

// ListVolume return list of volume
func (t *Targetd) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	vols, err := t.client.GetVolumeList(ctx, t.pool.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	var ids []string
	for _, vol := range vols {
		ids = append(ids, vol.Name)
	}

	vs, err := t.datastore.ListVolume(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list from datastore: %w", err)
	}
	return vs, nil
}

// GetVolume get volume
func (t *Targetd) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	vol, err := t.client.GetVolume(ctx, t.pool.Name, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes (ID: %s): %w", id, err)
	}
	logger.Logger.Debug(fmt.Sprintf("successfully retrieves volume: %+v", vol))

	vd, err := t.datastore.GetVolume(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume (ID: %s): %w", id, err)
	}
	logger.Logger.Debug(fmt.Sprintf("successfully retrieves volume from datastore: %+v", vd))

	return t.toVolume(*vol, vd.Attached, vd.HostName)
}

func (t *Targetd) toVolume(vol targetd.Volume, isAttached bool, hostname string) (*europa.Volume, error) {
	v := &europa.Volume{}
	v.ID = vol.Name

	gb := vol.Size / (1024 * 1024 * 1024)
	v.CapacityGB = uint32(gb)

	v.Attached = isAttached
	v.HostName = hostname
	v.BackendName = t.backendName

	return v, nil
}

// AttachVolumeTeleskop attach volume to hostname (running teleskop)
// return (host lun id, attached device name, error)
func (t *Targetd) AttachVolumeTeleskop(ctx context.Context, id string, hostname string) (int, string, error) {
	return 0, "", nil
}

// AttachVolumeSatelit attach volume to satelit
// return (host lun id, attached device name, error)
func (t *Targetd) AttachVolumeSatelit(ctx context.Context, id string, hostname string) (int, string, error) {
	return 0, "", nil
}

// DetachVolume detach volume
func (t *Targetd) DetachVolume(ctx context.Context, id string) error {
	return nil
}

// GetImage return image by id
func (t *Targetd) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	return nil, nil
}

// ListImage retrieves all images
func (t *Targetd) ListImage() ([]europa.BaseImage, error) {
	return nil, nil
}

// UploadImage upload to qcow2 image file
func (t *Targetd) UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*europa.BaseImage, error) {
	return nil, nil
}

// DeleteImage delete image
func (t *Targetd) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return nil
}
