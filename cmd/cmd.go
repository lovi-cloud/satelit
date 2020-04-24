package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/whywaita/satelit/api"

	"github.com/pkg/errors"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/europa/dorado"
)

var conf = flag.String("conf", "./configs/satelit.yaml", "set gateway config")

type Satelit struct {
	Europa europa.Europa
}

func init() {
	flag.Parse()
	err := config.Load(conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	logger.New(config.GetValue().LogLevel)
}

func NewSatelit() (*Satelit, error) {
	doradoBackend, err := dorado.New(config.GetValue().Dorado)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Dorado Backend")
	}

	err = teleskop.New(config.GetValue().Teleskop.Endpoints)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create teleskop agent")
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

	resp, err := teleskop.GetClient(config.GetValue().Teleskop.Endpoints[0]).GetISCSIQualifiedName(context.Background(), &pb.GetISCSIQualifiedNameRequest{})
	if err != nil {
		logger.Logger.Error(err.Error())
		return 1
	}
	fmt.Printf("%+v\n", resp.Iqn)

	return 0
}
