package main

import (
	"log"
	"os"

	"github.com/whywaita/satelit/cmd"
)

func main() {
	app, err := cmd.NewSatelit()
	if err != nil {
		log.Fatal(err)
	}

	// Run
	os.Exit(app.Run())
}
