package cmd

import (
	"flag"
	"fmt"

	"github.com/lovi-cloud/go-os-brick/osbrick"
	"github.com/lovi-cloud/satelit/internal/client/teleskop"
	"github.com/lovi-cloud/satelit/internal/logger"
	"github.com/lovi-cloud/satelit/pkg/api"
	"github.com/lovi-cloud/satelit/pkg/config"
	"github.com/lovi-cloud/satelit/pkg/datastore/mysql"
	"github.com/lovi-cloud/satelit/pkg/europa"
	"github.com/lovi-cloud/satelit/pkg/europa/targetd"
	"github.com/lovi-cloud/satelit/pkg/ganymede/libvirt"
	"github.com/lovi-cloud/satelit/pkg/ipam/ipam"
	"github.com/lovi-cloud/satelit/pkg/scheduler/scheduler"
	"go.uber.org/zap"
)

var conf = flag.String("conf", "./configs/satelit.yaml", "set satelit config")

func init() {
	flag.Parse()
	err := config.Load(conf)
	if err != nil {
		panic(err)
	}
	logger.New(config.GetValue().LogLevel)
	stdlogger, err := zap.NewStdLogAt(logger.Logger, zap.DebugLevel)
	if err != nil {
		panic(err)
	}
	osbrick.SetLogger(stdlogger)
}

// NewSatelit create SatelitServer instance.
func NewSatelit() (*api.SatelitServer, error) {
	c := config.GetValue().MySQLConfig
	ds, err := mysql.New(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to create mysql connection: %w", err)
	}

	targetds := map[string]europa.Europa{}
	for _, c := range config.GetValue().Targetd {
		targetdBackend, err := targetd.New(
			c.APIEndpoint,
			c.Username,
			c.Password,
			c.PoolName,
			c.BackendName,
			ds,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create targetd backend: %w", err)
		}
		targetds[c.BackendName] = targetdBackend
	}

	ipamBackend := ipam.New(ds)

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleskop agent: %w", err)
	}

	libvirtBackend := libvirt.New(ds)

	schedulerBackend := scheduler.New(ds)

	return &api.SatelitServer{
		Europa:    targetds,
		IPAM:      ipamBackend,
		Datastore: ds,
		Ganymede:  libvirtBackend,
		Scheduler: schedulerBackend,
	}, nil
}

// NewSatelitDatastore create SatelitDatastoreServer instance.
func NewSatelitDatastore() (*api.SatelitDatastore, error) {
	c := config.GetValue().MySQLConfig
	ds, err := mysql.New(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to create mysql connection: %w", err)
	}

	return &api.SatelitDatastore{
		Datastore: ds,
	}, nil
}
