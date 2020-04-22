package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/whywaita/satelit/internal/logger"

	"github.com/pkg/errors"

	"github.com/whywaita/satelit/pkg/europa"

	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/pkg/europa/backend/driver"
)

var conf = flag.String("conf", "./configs/satelit.yaml", "set gateway config")

type Satelit struct {
	Europa europa.Europa
}

func init() {
	flag.Parse()
	err := config.Load(conf)
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	logger.New(config.GetValue().LogLevel)
}

func NewSatelit() (*Satelit, error) {
	doradoBackend, err := driver.NewDoradoBackend(config.GetValue().Dorado)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Dorado Backend")
	}

	return &Satelit{
		Europa: doradoBackend,
	}, nil
}

func (s *Satelit) Run() int {
	// TODO: implement Serve
	vs, err := s.Europa.ListVolume(context.Background())
	if err != nil {
		logger.Logger.Error(err.Error())
		return 1
	}

	fmt.Printf("%+v\n", vs)
	return 0
}
