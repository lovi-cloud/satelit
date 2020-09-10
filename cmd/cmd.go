package cmd

import (
	"flag"
	"fmt"

	isucon_sshkey "github.com/whywaita/isucon-sshkey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	doradoBackend, err := dorado.New(config.GetValue().Dorado, ds)
	if err != nil {
		return nil, fmt.Errorf("failed to create Dorado Backend: %w", err)
	}

	ipamBackend := ipam.New(ds)

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleskop agent: %w", err)
	}

	libvirtBackend := libvirt.New(ds)

	schedulerBackend := scheduler.New(ds)

	return &api.SatelitServer{
		Europa:    doradoBackend,
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

	p := config.GetValue().Portal
	client, err := isucon_sshkey.NewClient(p.Endpoint, p.HMACSecretKey, logger.Logger)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create ISUCON portal client: %+v", err)
	}

	return &api.SatelitDatastore{
		Datastore: ds,
		Client:    client,
	}, nil
}
