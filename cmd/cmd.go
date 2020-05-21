package cmd

import (
	"flag"
	"fmt"

	"github.com/whywaita/go-os-brick/osbrick"
	"go.uber.org/zap"

	"github.com/whywaita/satelit/pkg/api"
	"github.com/whywaita/satelit/pkg/datastore/mysql"

	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/europa/dorado"
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

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to create teleskop agent: %w", err)
	}

	return &api.SatelitServer{
		Europa:    doradoBackend,
		Datastore: ds,
	}, nil
}
