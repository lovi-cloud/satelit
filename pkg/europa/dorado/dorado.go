package dorado

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	agentpb "github.com/whywaita/satelit/api"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/go-dorado-sdk/dorado"
	"github.com/whywaita/go-os-brick/osbrick"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/qcow2"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/europa"
	"go.uber.org/zap"
)

// A Dorado is backend of europa by Dorado
type Dorado struct {
	client    *dorado.Client
	datastore datastore.Datastore

	hyperMetroDomainID string
	storagePoolName    string
	portGroupName      string
	local              Device
	remote             Device
}

// A Device is IDs of devices
type Device struct {
	storagePoolID int
	portGroupID   int
}

// New create Dorado backend
func New(doradoConfig config.Dorado, datastore datastore.Datastore) (*Dorado, error) {
	ctx := context.Background()
	client, err := dorado.NewClient(
		doradoConfig.LocalIps,
		doradoConfig.RemoteIps,
		doradoConfig.Username,
		doradoConfig.Password,
		doradoConfig.PortGroupName,
		zap.NewStdLog(logger.Logger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dorado backend Client: %w", err)
	}

	hmds, err := client.GetHyperMetroDomains(context.Background(), dorado.NewSearchQueryName(doradoConfig.HyperMetroDomainName))
	if err != nil {
		return nil, fmt.Errorf("failed to get hyper metro domain (name: %s): %w", doradoConfig.HyperMetroDomainName, err)
	}
	if len(hmds) != 1 {
		return nil, fmt.Errorf("not only one HyperMetro Domain in same name (name: %s)", doradoConfig.HyperMetroDomainName)
	}

	localStoragePoolIDs, err := client.LocalDevice.GetStoragePools(ctx, dorado.NewSearchQueryName(doradoConfig.StoragePoolName))
	if err != nil {
		return nil, fmt.Errorf("failed to get local storage pool id: %w", err)
	}
	if len(localStoragePoolIDs) != 1 {
		return nil, fmt.Errorf("not only one StoragePool ID in same name (name: %s)", doradoConfig.StoragePoolName)
	}

	localPortgroupIDs, err := client.LocalDevice.GetPortGroups(ctx, dorado.NewSearchQueryName(doradoConfig.PortGroupName))
	if err != nil {
		return nil, fmt.Errorf("failed to get local port group id: %w", err)
	}
	if len(localPortgroupIDs) != 1 {
		return nil, fmt.Errorf("not only one port group id: in same name (name: %s)", doradoConfig.PortGroupName)
	}

	l := Device{
		storagePoolID: localStoragePoolIDs[0].ID,
		portGroupID:   localPortgroupIDs[0].ID,
	}

	remoteStoragePoolIDs, err := client.RemoteDevice.GetStoragePools(ctx, dorado.NewSearchQueryName(doradoConfig.StoragePoolName))
	if err != nil {
		return nil, fmt.Errorf("failed to get remote storage pool id: %w", err)
	}
	if len(remoteStoragePoolIDs) != 1 {
		return nil, fmt.Errorf("found multiple StoragePool ID in same name (name: %s)", doradoConfig.StoragePoolName)
	}

	remotePortgroupIDs, err := client.RemoteDevice.GetPortGroups(ctx, dorado.NewSearchQueryName(doradoConfig.PortGroupName))
	if err != nil {
		return nil, fmt.Errorf("failed to get local port group id: %w", err)
	}
	if len(remotePortgroupIDs) != 1 {
		return nil, fmt.Errorf("not only one port group id: in same name (name: %s)", doradoConfig.PortGroupName)
	}

	r := Device{
		storagePoolID: remoteStoragePoolIDs[0].ID,
		portGroupID:   remotePortgroupIDs[0].ID,
	}

	return &Dorado{
		client:             client,
		datastore:          datastore,
		hyperMetroDomainID: hmds[0].ID,

		local:  l,
		remote: r,
	}, nil
}

// ListVolume return list of volume by Dorado
func (d *Dorado) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	var ids []string
	for _, hmp := range hmps {
		ids = append(ids, hmp.ID)
	}

	vs, err := d.datastore.ListVolume(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list from datastore: %w", err)
	}
	return vs, nil
}

// GetVolume get volume from Dorado
func (d *Dorado) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	hmp, err := d.client.GetHyperMetroPair(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes (ID: %s): %w", id, err)
	}
	logger.Logger.Debug(fmt.Sprintf("successfully retrieves HyperMetroPair: %+v", hmp))

	vd, err := d.datastore.GetVolume(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume (ID: %s): %w", id, err)
	}
	logger.Logger.Debug(fmt.Sprintf("successfully retrieves volume from datastore: %+v", vd))

	v, err := d.toVolume(hmp, vd.Attached, vd.HostName)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume (ID: %s): %w", hmp.ID, err)
	}
	return v, nil
}

// DeleteVolume delete volume by Dorado
func (d *Dorado) DeleteVolume(ctx context.Context, id string) error {
	err := d.client.DeleteVolume(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete volume (ID: %s): %w", id, err)
	}

	err = d.datastore.DeleteVolume(id)
	if err != nil {
		return fmt.Errorf("failed to delete volume from datastore (ID: %s): %w", id, err)
	}

	return nil
}

// AttachVolumeTeleskop attach volume to hostname (running teleskop) by Dorado
// return (host lun id, attached device name, error)
func (d *Dorado) AttachVolumeTeleskop(ctx context.Context, id string, hostname string) (int, string, error) {
	hmp, err := d.client.GetHyperMetroPair(ctx, id)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get HyperMetroPair: %w", err)
	}

	iqn, err := d.datastore.GetIQN(ctx, hostname)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get iqn (hostname: %s): %w", hostname, err)
	}

	// create dorado mappingview object
	err = d.client.AttachVolume(ctx, id, hostname, iqn)
	if err != nil {
		return 0, "", fmt.Errorf("failed to attach volume (ID: %s): %w", id, err)
	}

	targetPortalIPs, err := d.client.GetPortalIPAddresses(ctx, d.local.portGroupID, d.remote.portGroupID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get portal ip addresses: %w", err)
	}

	hostLUNID, err := d.GetHostLUNID(ctx, hmp, hostname)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get host lun id (ID: %s): %w", id, err)
	}

	req := &agentpb.ConnectBlockDeviceRequest{
		PortalAddresses: targetPortalIPs,
		HostLunId:       uint32(hostLUNID),
	}

	teleskopClient, err := teleskop.GetClient(hostname)
	if err != nil {
		return 0, "", fmt.Errorf("failed to retrieve teleskop client: %w", err)
	}
	resp, err := teleskopClient.ConnectBlockDevice(ctx, req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to connect block device to teleskop: %w", err)
	}

	volume, err := d.datastore.GetVolume(id)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get volume (ID: %s): %w", id, err)
	}
	volume.Attached = true
	volume.HostName = hostname
	err = d.datastore.PutVolume(*volume)
	if err != nil {
		return 0, "", fmt.Errorf("failed to update volume record (ID: %s): %w", id, err)
	}

	return hostLUNID, resp.DeviceName, nil
}

// AttachVolumeSatelit attach volume to satelit by Dorado
// return (host lun id, attached device name, error)
func (d *Dorado) AttachVolumeSatelit(ctx context.Context, hyperMetroPairID string, hostname string) (int, string, error) {
	iqn, err := d.datastore.GetIQN(ctx, hostname)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get iqn (hostname: %s): %w", hostname, err)
	}

	// create dorado mappingview object
	err = d.client.AttachVolume(ctx, hyperMetroPairID, hostname, iqn)
	if err != nil {
		return 0, "", fmt.Errorf("failed to attach volume (ID: %s): %w", hyperMetroPairID, err)
	}

	//  running satelit server (not teleskop)
	hostLUNID, deviceName, err := d.attachVolumeSatelit(ctx, hyperMetroPairID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to attach volume to localhost (ID: %s): %w", hyperMetroPairID, err)
	}

	volume, err := d.datastore.GetVolume(hyperMetroPairID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get volume (ID: %s): %w", hyperMetroPairID, err)
	}
	volume.Attached = true
	volume.HostName = hostname
	err = d.datastore.PutVolume(*volume)
	if err != nil {
		return 0, "", fmt.Errorf("failed to update volume record (ID: %s): %w", hyperMetroPairID, err)
	}

	return hostLUNID, deviceName, nil
}

// DetachVolume detach volume by Dorado
func (d *Dorado) DetachVolume(ctx context.Context, hyperMetroPairID string) error {
	err := d.client.DetachVolume(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to detach volume (ID: %s): %w", hyperMetroPairID, err)
	}

	// TODO: send to detach operation

	volume, err := d.datastore.GetVolume(hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get volume (ID: %s): %w", hyperMetroPairID, err)
	}

	volume.Attached = false
	volume.HostName = ""

	err = d.datastore.PutVolume(*volume)
	if err != nil {
		return fmt.Errorf("failed to update volume record (ID: %s): %w", hyperMetroPairID, err)
	}

	return nil
}

// DetachVolumeSatelit detach volume from satelit server
func (d *Dorado) DetachVolumeSatelit(ctx context.Context, hyperMetroPairID string, hostLUNID int) error {
	err := d.client.DetachVolume(ctx, hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to delete dorado attach mapping (ID: %s): %w", hyperMetroPairID, err)
	}

	// detach volume
	targetPortalIPs, err := d.client.GetPortalIPAddresses(ctx, d.local.portGroupID, d.remote.portGroupID)
	if err != nil {
		return fmt.Errorf("failed to get portal ip addresses: %w", err)
	}

	err = osbrick.DisconnectVolume(ctx, targetPortalIPs, hostLUNID)
	if err != nil {
		return fmt.Errorf("failed to detach volume: %w", err)
	}

	volume, err := d.datastore.GetVolume(hyperMetroPairID)
	if err != nil {
		return fmt.Errorf("failed to get volume (ID: %s): %w", hyperMetroPairID, err)
	}

	volume.Attached = false
	volume.HostName = ""

	err = d.datastore.PutVolume(*volume)
	if err != nil {
		return fmt.Errorf("failed to update volume record (ID: %s): %w", hyperMetroPairID, err)
	}

	return nil
}

// GetImage return image by id
func (d *Dorado) GetImage(imageID uuid.UUID) (*europa.BaseImage, error) {
	image, err := d.datastore.GetImage(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to getr image from datastore: %w", err)
	}

	return image, nil
}

// ListImage retrieves all images
func (d *Dorado) ListImage() ([]europa.BaseImage, error) {
	images, err := d.datastore.ListImage()
	if err != nil {
		return nil, fmt.Errorf("failed to get images from datastore: %w", err)
	}

	return images, nil
}

// UploadImage upload to qcow2 image file
func (d *Dorado) UploadImage(ctx context.Context, image []byte, name, description string, imageSizeGB int) (*europa.BaseImage, error) {
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
	hmp, err := d.client.CreateVolumeRaw(ctx, u, imageSizeGB, d.storagePoolName, d.hyperMetroDomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", u.String(), err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	hostLUNID, deviceName, err := d.AttachVolumeSatelit(ctx, hmp.ID, hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to attach volume: %w", err)
	}
	defer func() {
		d.DetachVolumeSatelit(ctx, hmp.ID, hostLUNID)
	}()

	// exec qemu-img convert
	err = qcow2.ToRaw(ctx, tmpfile.Name(), deviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to convert raw format: %w", err)
	}

	bi := &europa.BaseImage{
		UUID:          u,
		Name:          name,
		Description:   description,
		CacheVolumeID: hmp.ID,
	}

	return bi, nil
}

// DeleteImage delete image by Dorado
func (d *Dorado) DeleteImage(ctx context.Context, id uuid.UUID) error {
	image, err := d.datastore.GetImage(id)
	if err != nil {
		return fmt.Errorf("failed to get image from datastore: %w", err)
	}

	// TODO: check attached status

	err = d.DeleteVolume(ctx, image.CacheVolumeID)
	if err != nil {
		return fmt.Errorf("failed to delete image cache volume: %w", err)
	}

	err = d.datastore.DeleteImage(id)
	if err != nil {
		return fmt.Errorf("failed to delete datastore: %w", err)
	}

	return nil
}

func (d *Dorado) toVolume(hmp *dorado.HyperMetroPair, isAttached bool, hostname string) (*europa.Volume, error) {
	v := &europa.Volume{}
	v.ID = hmp.ID

	c, err := strconv.Atoi(hmp.CAPACITYBYTE)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CAPACITYBYTE (ID: %s): %w", hmp.ID, err)
	}

	v.Attached = isAttached
	v.HostName = hostname

	v.CapacityGB = uint32(c / dorado.CapacityUnit)
	return v, nil
}

// getHostgroupLocalhost get hostgroup and host in local device and remote device.
// if not exist in device, then create host / hostgroup object.
func (d *Dorado) getHostgroupLocalhost(ctx context.Context, hostname string) (*dorado.HostGroup, *dorado.Host, *dorado.HostGroup, *dorado.Host, error) {
	lHostgroup, lHost, err := d.client.LocalDevice.GetHostGroupForce(ctx, hostname)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get local host group: %w", err)
	}
	rHostgroup, rHost, err := d.client.RemoteDevice.GetHostGroupForce(ctx, hostname)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get remote host group: %w", err)
	}

	return lHostgroup, lHost, rHostgroup, rHost, nil
}

// GetHostLUNID return host LUN id.
func (d *Dorado) GetHostLUNID(ctx context.Context, hmp *dorado.HyperMetroPair, hostname string) (int, error) {
	_, lHost, _, _, err := d.getHostgroupLocalhost(ctx, hostname)
	if err != nil {
		return -1, fmt.Errorf("failed to create host group: %w", err)
	}

	// NOTE(whywaita): this code except same host lun ID between local device and remote device.
	// if diff in local and remote, need to fix go-os-brick
	hostLUNID, err := d.client.LocalDevice.GetHostLUNID(ctx, hmp.LOCALOBJID, lHost.ID)
	if err != nil {
		return -1, fmt.Errorf("failed to get host lun id: %w", err)
	}

	return hostLUNID, nil
}

// GetHostLUNIDLocalhost return host LUN id in localhost.
func (d *Dorado) GetHostLUNIDLocalhost(ctx context.Context, hmp *dorado.HyperMetroPair) (int, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, fmt.Errorf("failed to get hostname: %w", err)
	}

	return d.GetHostLUNID(ctx, hmp, hostname)
}
