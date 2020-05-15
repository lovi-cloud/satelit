package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

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
		log.Println(err)
		os.Exit(1)
	}
	logger.New(config.GetValue().LogLevel)
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
		Europa: doradoBackend,
	}, nil
}
