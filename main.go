package main

import (
	"fmt"
	"log"
	"os"

	"github.com/whywaita/satelit/cmd"
)

var (
	revision string
)

func init() {
	fmt.Println(fmt.Sprintf("satelit revision: %s", revision))
}

func main() {
	app, err := cmd.NewSatelit()
	if err != nil {
		log.Fatal(err)
	}

	// Run
	os.Exit(app.Run())
}
