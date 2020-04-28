package main

import (
	"log"
	"os"

	"github.com/whywaita/satelit/cmd"
)

var (
	revision string
)

func main() {
	app, err := cmd.NewSatelit()
	if err != nil {
		log.Fatal(err)
	}

	// Run
	os.Exit(app.Run())
}
