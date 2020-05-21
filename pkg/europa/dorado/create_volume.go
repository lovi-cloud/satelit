package dorado

import (
	"context"
	"fmt"
	"os"

	"github.com/whywaita/go-dorado-sdk/dorado"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/go-os-brick/osbrick"
	"github.com/whywaita/satelit/pkg/europa"
)

// CreateVolumeRaw create raw volume
func (d *Dorado) CreateVolumeRaw(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	hmp, err := d.client.CreateVolume(ctx, name, capacity, d.storagePoolName, d.hyperMetroDomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", name.String(), err)
	}

	volume, err := d.toVolume(hmp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert europa.volume (ID: %s): %w", hmp.ID, err)
	}

	return volume, nil
}

// CreateVolumeImage create volume that copied image
func (d *Dorado) CreateVolumeImage(ctx context.Context, name uuid.UUID, capacity int, imageID string) (*europa.Volume, error) {
	// create raw volume
	hmp, err := d.client.CreateVolume(ctx, name, capacity, d.storagePoolName, d.hyperMetroDomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", name.String(), err)
	}

	// raw device attach to localhost
	_, err = d.attachVolumeLocalhost(ctx, hmp)
	if err != nil {
		return nil, fmt.Errorf("failed to attach volume: %w", err)
	}

	// image device attach to localhost

	// exec qemu-img convert

	// unmount raw device and image device

	return nil, nil
}

// attachVolumeLocalhost attach volume to localhost (= running satelit host).
func (d *Dorado) attachVolumeLocalhost(ctx context.Context, hmp *dorado.HyperMetroPair) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	err = d.AttachVolume(ctx, hmp.ID, hostname)
	if err != nil {
		return "", fmt.Errorf("failed to attach volume: %w", err)
	}

	targetPortalIPs, err := d.client.GetPortalIPAddresses(ctx, d.local.portGroupID, d.remote.portGroupID)
	if err != nil {
		return "", fmt.Errorf("failed to get portal ip addresses: %w", err)
	}

	hostLUNID, err := d.GetHostLUNIDLocalhost(ctx, hmp)
	if err != nil {
		return "", fmt.Errorf("failed to get host lun ID in localhost: %w", err)
	}

	// mount to satelit server
	deviceName, err := osbrick.ConnectMultipathVolume(ctx, targetPortalIPs, hostLUNID)
	if err != nil {
		return "", fmt.Errorf("failed to connect device: %w", err)
	}

	return deviceName, nil
}

// getHostgroupLocalhost get hostgroup and host in local device and remote device.
// if not exist in device, then create host / hostgroup object.
func (d *Dorado) getHostgroupLocalhost(ctx context.Context) (*dorado.HostGroup, *dorado.Host, *dorado.HostGroup, *dorado.Host, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get hostname: %w", err)
	}

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
