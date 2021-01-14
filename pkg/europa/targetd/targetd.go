package targetd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/lovi-cloud/go-os-brick/osbrick"

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
	portalIP    string
}

// New is create instance of targetd.Client
func New(apiURL, username, password, poolName, backendName, portalIP string, ds datastore.Datastore) (*Targetd, error) {
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
		portalIP:    portalIP,
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
func (t *Targetd) AttachVolumeSatelit(ctx context.Context, volName string, hostname string) (int, string, error) {
	iqn, err := osbrick.GetIQN(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get iqn (hostname: %s): %w", hostname, err)
	}

	newHostLUNID, err := t.GetHostLUNID(ctx, volName, iqn)
	if err != nil {
		return 0, "", fmt.Errorf("failed to retrieve host lun ID: %w", err)
	}
	if err := t.client.CreateExport(ctx, t.pool.Name, volName, newHostLUNID, iqn); err != nil {
		return 0, "", fmt.Errorf("failed to create export: %w", err)
	}
	defer func() {
		poolName := t.pool.Name
		if err != nil {
			if err := t.client.DestroyExport(ctx, poolName, volName, iqn); err != nil {
				logger.Logger.Warn(fmt.Sprintf("failed to DestroyExport: %+v", err))
			}
		}
	}()

	vol, err := t.GetVolume(ctx, volName)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get volume: %w", err)
	}

	deviceName, err := osbrick.ConnectSinglePathVolume(ctx, t.portalIP, newHostLUNID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to attach iSCSI session: %w", err)
	}

	vol.Attached = true
	vol.HostName = hostname
	vol.HostLUNID = newHostLUNID
	if err := t.datastore.PutVolume(ctx, *vol); err != nil {
		return 0, "", fmt.Errorf("failed to update volume record (name: %s): %w", volName, err)
	}

	return newHostLUNID, deviceName, nil
}

// DetachVolume detach volume
func (t *Targetd) DetachVolume(ctx context.Context, volName string) error {
	return nil
}

// DetachVolumeSatelit detach volume from satelit server
func (t *Targetd) DetachVolumeSatelit(ctx context.Context, volName string, hostLUNID int) error {
	volume, err := t.datastore.GetVolume(ctx, volName)
	if err != nil {
		return fmt.Errorf("failed to retrieve ")
	}
	if volume.Attached == false {
		return fmt.Errorf("not attached in satelit: %w", err)
	}

	iqn, err := osbrick.GetIQN(ctx)
	if err != nil {
		return fmt.Errorf("failed to get iqn: %w", err)
	}

	if err := osbrick.DisconnectSinglePathVolume(ctx, t.portalIP, hostLUNID); err != nil {
		return fmt.Errorf("failed to disconnect iSCSI session: %w", err)
	}
	if err := t.client.DestroyExport(ctx, t.pool.Name, volName, iqn); err != nil {
		return fmt.Errorf("failed to destroy export from targetd: %w", err)
	}

	volume.Attached = false
	volume.HostName = ""
	volume.HostLUNID = 0

	if err := t.datastore.PutVolume(ctx, *volume); err != nil {
		return fmt.Errorf("failed to update volume record (name: %s): %w", volName, err)
	}

	return nil
}

// GetImage return image by id
func (t *Targetd) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	image, err := t.datastore.GetImage(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to getr image from datastore: %w", err)
	}

	return image, nil
}

// ListImage retrieves all images
func (t *Targetd) ListImage() ([]europa.BaseImage, error) {
	images, err := t.datastore.ListImage()
	if err != nil {
		return nil, fmt.Errorf("failed to get images from datastore: %w", err)
	}

	return images, nil
}

// UploadImage upload to qcow2 image file
func (t *Targetd) UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*europa.BaseImage, error) {
	// image write to tmpfile
	tmpfile, err := ioutil.TempFile("", "satelit")
	if err != nil {
		return nil, fmt.Errorf("failde to create temporary file: %w", err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err = tmpfile.Write(image); err != nil {
		return nil, fmt.Errorf("failed to write image to temporary file: %w", err)
	}
	if err = tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	// mount new volume
	u := uuid.NewV4()
	v, err := t.CreateVolume(ctx, u, imageSizeGB)
	if err != nil {
		return nil, fmt.Errorf("faeiled to create volume (name: %s): %w", u, err)
	}
	defer func() {
		volumeID := v.ID
		if err != nil {
			if err := t.DeleteVolume(ctx, volumeID); err != nil {
				logger.Logger.Warn(fmt.Sprintf("failed to DeleteVolume: %+v", err))
			}
		}
	}()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	hostLUNID, deviceName, err := t.AttachVolumeSatelit(ctx, v.ID, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to attach volume: %w", err)
	}

	// exec qemu-img convert
	if err := osbrick.QEMUToRaw(ctx, tmpfile.Name(), deviceName); err != nil {
		return nil, fmt.Errorf("failed to convert raw format: %w", err)
	}

	bi := &europa.BaseImage{
		UUID:          u,
		Name:          name,
		Description:   description,
		CacheVolumeID: v.ID,
	}

	if err := t.DetachVolumeSatelit(ctx, v.ID, hostLUNID); err != nil {
		return nil, fmt.Errorf("failed to detach volume: %w", err)
	}

	return bi, nil
}

// DeleteImage delete image
func (t *Targetd) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return nil
}

// GetHostLUNID return host LUN id. (not created)
func (t *Targetd) GetHostLUNID(ctx context.Context, volName, iqn string) (int, error) {
	latestHostLUNID, err := t.getExistLatestHostLUNID(ctx, iqn)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve exist host LUN ID: %w", err)
	}
	newHostLUNID := latestHostLUNID + 1

	return newHostLUNID, nil
}

func (t *Targetd) getExistLatestHostLUNID(ctx context.Context, iqn string) (int, error) {
	exports, err := t.getExports(ctx, iqn)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve exports: %w", err)
	}

	var hostLUNIDs []int
	for _, e := range exports {
		hostLUNIDs = append(hostLUNIDs, e.LUN)
	}
	if len(hostLUNIDs) == 0 {
		return 0, nil
	}

	sort.Slice(hostLUNIDs, func(i, j int) bool {
		return hostLUNIDs[i] < hostLUNIDs[j]
	})

	return hostLUNIDs[len(hostLUNIDs)-1], nil
}

func (t *Targetd) getExports(ctx context.Context, iqn string) ([]targetd.Export, error) {
	exports, err := t.client.ListExport(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list exports: %w", err)
	}

	var result []targetd.Export
	for _, e := range exports {
		if e.InitiatorWwn == iqn {
			result = append(result, e)
		}
	}

	return result, nil
}
