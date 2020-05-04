package dorado

import (
	"context"
	"fmt"

	"strconv"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/go-dorado-sdk/dorado"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/europa"
	"go.uber.org/zap"
)

// A Dorado is backend of europa by Dorado
type Dorado struct {
	client    *dorado.Client
	datastore datastore.Datastore

	storagePoolName    string
	hyperMetroDomainID string
}

// New create Dorado backend
func New(doradoConfig config.Dorado, datastore datastore.Datastore) (*Dorado, error) {
	client, err := dorado.NewClient(
		doradoConfig.LocalIps[0],
		doradoConfig.RemoteIps[0],
		doradoConfig.Username,
		doradoConfig.Password,
		doradoConfig.PortGroupName,
		zap.NewStdLog(logger.Logger),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dorado backend Client")
	}

	hmds, err := client.GetHyperMetroDomains(context.Background(), dorado.NewSearchQueryName(doradoConfig.HyperMetroDomainName))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get hypermetrodomain. name: %s", doradoConfig.HyperMetroDomainName))
	}
	if len(hmds) != 1 {
		return nil, errors.New(fmt.Sprintf("founf multiple HyperMetro Domain in same name. name: %s", doradoConfig.HyperMetroDomainName))
	}

	return &Dorado{
		client:             client,
		datastore:          datastore,
		storagePoolName:    doradoConfig.StoragePoolName,
		hyperMetroDomainID: hmds[0].ID,
	}, nil
}

// CreateVolume create volume by Dorado
func (d *Dorado) CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	hmp, err := d.client.CreateVolume(ctx, name, capacity, d.storagePoolName, d.hyperMetroDomainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create volume")
	}

	volume, err := d.toVolume(hmp)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get europa.Volume. ID: %s", hmp.ID))
	}

	return volume, nil
}

// ListVolume return list of volume by Dorado
func (d *Dorado) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get volume list")
	}

	var vs []europa.Volume
	for _, hmp := range hmps {
		v, err := d.toVolume(&hmp)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get europa.Volume. ID: %s", hmp.ID))
		}
		vs = append(vs, *v)
	}

	return vs, nil
}

// GetVolume get volume from Dorado
func (d *Dorado) GetVolume(ctx context.Context, name uuid.UUID) (*europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, dorado.NewSearchQueryName(name.String()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get volumes")
	}
	if len(hmps) != 1 {
		return nil, errors.New("found multiple volumes in same name")
	}

	volume := hmps[0]
	v, err := d.toVolume(&volume)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get europa.Volume. ID: %s", volume.ID))
	}
	return v, nil
}

// DeleteVolume delete volume by Dorado
func (d *Dorado) DeleteVolume(ctx context.Context, name uuid.UUID) error {
	volume, err := d.GetVolume(ctx, name)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get Volume info. name: %s", name.String()))
	}

	err = d.client.DeleteVolume(ctx, volume.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete volume. id: %s", volume.ID))
	}

	return nil
}

// AttachVolume attach to hostname by Dorado
func (d *Dorado) AttachVolume(ctx context.Context, name uuid.UUID, hostname string) (*europa.Volume, error) {
	volume, err := d.GetVolume(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get Volume info. name: %s", name.String()))
	}

	iqn, err := d.datastore.GetIQN(ctx, hostname)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get iqn. hostname: %s", hostname))
	}

	err = d.client.AttachVolume(ctx, volume.ID, hostname, iqn)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to delete volume. id: %s", volume.ID))
	}

	v, err := d.GetVolume(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get europa.Volume. ID: %s", v.ID))
	}
	v.Attached = true
	v.HostName = hostname

	return v, nil
}

// DetachVolume detach volume by Dorado
func (d *Dorado) DetachVolume(ctx context.Context, name uuid.UUID) error {
	volume, err := d.GetVolume(ctx, name)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get Volume info. name: %s", name.String()))
	}

	err = d.client.DetachVolume(ctx, volume.ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to detach volume. name: %s", name.String()))
	}

	return nil
}

func (d *Dorado) toVolume(hmp *dorado.HyperMetroPair) (*europa.Volume, error) {
	v := &europa.Volume{}
	v.ID = hmp.ID

	c, err := strconv.Atoi(hmp.CAPACITYBYTE)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse CAPACITYBYTE. Volume.ID: %s", hmp.ID))
	}

	v.CapacityGB = c / dorado.CapacityUnit
	return v, nil
}
