package cmd

import (
	"flag"
	"fmt"

	"github.com/whywaita/satelit/pkg/europa"

	"github.com/whywaita/satelit/pkg/scheduler/scheduler"

	"go.uber.org/zap"

	"github.com/whywaita/go-os-brick/osbrick"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/api"
	"github.com/whywaita/satelit/pkg/datastore/mysql"
	"github.com/whywaita/satelit/pkg/europa/dorado"
	"github.com/whywaita/satelit/pkg/ganymede/libvirt"
	"github.com/whywaita/satelit/pkg/ipam/ipam"
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

	dorados := map[string]europa.Europa{}
	for _, c := range config.GetValue().Dorado {
		doradoBackend, err := dorado.New(c, ds)
		if err != nil {
			return nil, fmt.Errorf("failed to create Dorado Backend: %w", err)
		}

		dorados[c.BackendName] = doradoBackend
	}

	ipamBackend := ipam.New(ds)

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleskop agent: %w", err)
	}

	libvirtBackend := libvirt.New(ds)

	schedulerBackend := scheduler.New(ds)

	return &api.SatelitServer{
		Europa:    dorados,
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
