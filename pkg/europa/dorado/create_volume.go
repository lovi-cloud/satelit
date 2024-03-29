package dorado

import (
	"context"
	"fmt"

	"github.com/lovi-cloud/go-os-brick/osbrick"
	"github.com/lovi-cloud/satelit/pkg/europa"
	uuid "github.com/satori/go.uuid"
)

// CreateVolume create raw volume
func (d *Dorado) CreateVolume(ctx context.Context, name uuid.UUID, capacityGB int) (*europa.Volume, error) {
	hmp, err := d.client.CreateVolumeRaw(ctx, name, capacityGB, d.storagePoolName, d.hyperMetroDomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", name.String(), err)
	}

	volume, err := d.toVolume(hmp, false, "")
	if err != nil {
		return nil, fmt.Errorf("failed to convert europa.volume (ID: %s): %w", hmp.ID, err)
	}

	err = d.datastore.PutVolume(ctx, *volume)
	if err != nil {
		return nil, fmt.Errorf("failed to put volume to datastore (ID: %s): %w", hmp.ID, err)
	}

	return volume, nil
}

// CreateVolumeFromImage create volume that copied image
func (d *Dorado) CreateVolumeFromImage(ctx context.Context, name uuid.UUID, capacityGB int, imageID uuid.UUID) (*europa.Volume, error) {
	image, err := d.GetImage(imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	hmp, err := d.client.CreateVolumeFromSource(ctx, name, capacityGB, d.storagePoolName, d.hyperMetroDomainID, image.CacheVolumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume from source: %w", err)
	}

	volume, err := d.toVolume(hmp, false, "")
	if err != nil {
		return nil, fmt.Errorf("failed to convert europa.Volume: %w", err)
	}
	volume.BaseImageID = imageID

	err = d.datastore.PutVolume(ctx, *volume)
	if err != nil {
		return nil, fmt.Errorf("failed to put volume to datastore (ID: %s): %w", hmp.ID, err)
	}

	return volume, nil
}

// attachVolumeLocalhost attach volume satelit host.
func (d *Dorado) attachVolumeSatelit(ctx context.Context, hyperMetroPairID string) (int, string, error) {
	hmp, err := d.client.GetHyperMetroPair(ctx, hyperMetroPairID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get hyper metro pair: %w", err)
	}

	targetPortalIPs, err := d.client.GetPortalIPAddresses(ctx, d.local.portGroupID, d.remote.portGroupID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get portal ip addresses: %w", err)
	}

	hostLUNID, err := d.GetHostLUNIDLocalhost(ctx, hmp)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get host lun ID in localhost: %w", err)
	}

	// mount to satelit server
	deviceName, err := osbrick.ConnectMultipathVolume(ctx, targetPortalIPs, hostLUNID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to connect device: %w", err)
	}

	return hostLUNID, deviceName, nil
}
