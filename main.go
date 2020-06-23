package main

import (
	"fmt"
	"log"
	"os"

	"github.com/whywaita/satelit/internal/logger"
	"golang.org/x/sync/errgroup"

	"github.com/whywaita/satelit/cmd"
)

var (
	revision string
)

func init() {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("satelit revision: %s", revision))
}

func main() {
	app, err := cmd.NewSatelit()
	if err != nil {
		log.Fatal(err)
	}

	datastore, err := cmd.NewSatelitDatastore()
	if err != nil {
		log.Fatal(err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		return app.Run()
	})
	eg.Go(func() error {
		return datastore.Run()
	})
	if err := eg.Wait(); err != nil {
		logger.Logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
