package driver

import (
	"context"

	"github.com/whywaita/satelit/internal/logger"

	"go.uber.org/zap"

	"github.com/whywaita/satelit/internal/config"

	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"

	"github.com/whywaita/go-dorado-sdk/dorado"
	"github.com/whywaita/satelit/pkg/europa"
)

type DoradoBackend struct {
	client *dorado.Client

	storagePoolId      string
	hyperMetroDomainId string
}

func NewDoradoBackend(doradoConfig config.Dorado) (*DoradoBackend, error) {
	db := &DoradoBackend{}

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

	db.client = client

	return db, nil
}

func (d *DoradoBackend) CreateVolume(ctx context.Context, name uuid.UUID, capacity int) (*europa.Volume, error) {
	hmp, err := d.client.CreateVolume(ctx, name, capacity, d.storagePoolId, d.hyperMetroDomainId)
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

func (d *DoradoBackend) toVolume(hmp *dorado.HyperMetroPair) *europa.Volume {
	return &europa.Volume{
		ID: hmp.ID,
	}
}
