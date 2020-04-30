package cmd

import (
	"flag"
	"log"
	"os"

	"github.com/whywaita/satelit/pkg/datastore/memory"

	"github.com/whywaita/satelit/pkg/api"

	"github.com/pkg/errors"
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

func NewSatelit() (*api.SatelitServer, error) {
	//c := config.GetValue().MySQLConfig
	//ds, err := mysql.New(&c)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to create mysql connection")
	//}
	ds := memory.New() // For development

	doradoBackend, err := dorado.New(config.GetValue().Dorado, ds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Dorado Backend")
	}

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create teleskop agent")
	}

	return &api.SatelitServer{
		Europa: doradoBackend,
	}, nil
}
