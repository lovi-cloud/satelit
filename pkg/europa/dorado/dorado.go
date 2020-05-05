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
		return nil, fmt.Errorf("failed to create dorado backend Client: %w", err)
	}

	hmds, err := client.GetHyperMetroDomains(context.Background(), dorado.NewSearchQueryName(doradoConfig.HyperMetroDomainName))
	if err != nil {
		return nil, fmt.Errorf("failed to get hyper metro domain (name: %s): %w", doradoConfig.HyperMetroDomainName, err)
	}
	if len(hmds) != 1 {
		return nil, errors.New(fmt.Sprintf("founf multiple HyperMetro Domain in same name (name: %s)", doradoConfig.HyperMetroDomainName))
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
		return nil, fmt.Errorf("failed to create volume (name: %s): %w", name.String(), err)
	}

	volume, err := d.toVolume(hmp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert europa.volume (ID: %s): %w", hmp.ID, err)
	}

	return volume, nil
}

// ListVolume return list of volume by Dorado
func (d *Dorado) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume list: %w", err)
	}

	var vs []europa.Volume
	for _, hmp := range hmps {
		v, err := d.toVolume(&hmp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert europa.volume (ID: %s): %w", hmp.ID, err)
		}
		vs = append(vs, *v)
	}

	return vs, nil
}

// GetVolume get volume from Dorado
func (d *Dorado) GetVolume(ctx context.Context, id string) (*europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, dorado.NewSearchQueryId(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes (ID: %s): %w", id, err)
	}
	fmt.Printf("%+v\n", hmps)
	if len(hmps) != 1 {
		return nil, fmt.Errorf("found multiple volumes in same name (ID: %s)", id)
	}

	volume := hmps[0]
	v, err := d.toVolume(&volume)
	if err != nil {
		return nil, fmt.Errorf("failed to get volume (ID: %s): %w", volume.ID, err)
	}
	return v, nil
}

// DeleteVolume delete volume by Dorado
func (d *Dorado) DeleteVolume(ctx context.Context, id string) error {
	err := d.client.DeleteVolume(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete volume (ID: %s): %w", id, err)
	}

	return nil
}

// AttachVolume attach to hostname by Dorado
func (d *Dorado) AttachVolume(ctx context.Context, id string, hostname string) error {
	iqn, err := d.datastore.GetIQN(ctx, hostname)
	if err != nil {
		return fmt.Errorf("failed to get iqn (hostname: %s): %w", hostname, err)
	}

	err = d.client.AttachVolume(ctx, id, hostname, iqn)
	if err != nil {
		return fmt.Errorf("failed to attach volume (ID: %s): %w", id, err)
	}

	return nil
}

// DetachVolume detach volume by Dorado
func (d *Dorado) DetachVolume(ctx context.Context, id string) error {
	err := d.client.DetachVolume(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to detach volume (ID: %s): %w", id, err)
	}

	return nil
}

func (d *Dorado) toVolume(hmp *dorado.HyperMetroPair) (*europa.Volume, error) {
	v := &europa.Volume{}
	v.ID = hmp.ID

	c, err := strconv.Atoi(hmp.CAPACITYBYTE)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CAPACITYBYTE (ID: %s): %w", hmp.ID, err)
	}

	v.CapacityGB = uint32(c / dorado.CapacityUnit)
	return v, nil
}
