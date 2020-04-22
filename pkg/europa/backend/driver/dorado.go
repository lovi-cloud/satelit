package driver

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/go-dorado-sdk/dorado"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/europa"
	"go.uber.org/zap"
)

type DoradoBackend struct {
	client *dorado.Client

	storagePoolName    string
	hyperMetroDomainId string
}

func NewDoradoBackend(doradoConfig config.Dorado) (*DoradoBackend, error) {
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

	return &DoradoBackend{
		client:             client,
		storagePoolName:    doradoConfig.StoragePoolName,
		hyperMetroDomainId: hmds[0].ID,
	}, nil
}

func (d *DoradoBackend) CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	hmp, err := d.client.CreateVolume(ctx, name, capacity, d.storagePoolName, d.hyperMetroDomainId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create volume")
	}

	return d.toVolume(hmp), nil
}

func (d *DoradoBackend) ListVolume(ctx context.Context) ([]europa.Volume, error) {
	hmps, err := d.client.GetHyperMetroPairs(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get volume list")
	}

	var vs []europa.Volume
	for _, hmp := range hmps {
		vs = append(vs, *d.toVolume(&hmp))
	}

	return vs, nil
}

func (d *DoradoBackend) DeleteVolume(ctx context.Context, name uuid.UUID) error {
	hmps, err := d.client.GetHyperMetroPairs(ctx, dorado.NewSearchQueryName(name.String()))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get volume. name: %s", name.String()))
	}
	if len(hmps) != 1 {
		return errors.New(fmt.Sprintf("found multiple volume in same name. name: %s", name.String()))
	}
	err = d.client.DeleteVolume(ctx, hmps[0].ID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete volume. id: %s", hmps[0].ID))
	}

	return nil
}

func (d *DoradoBackend) toVolume(hmp *dorado.HyperMetroPair) *europa.Volume {
	return &europa.Volume{
		ID: hmp.ID,
	}
}
