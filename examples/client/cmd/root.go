package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "satelit client",
}

// flag variables
var (
	satelitAddress string
)

// Execute execute command in root
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&satelitAddress, "address", "", "", "address of satelit, <host>:9262")
	if err := rootCmd.MarkPersistentFlagRequired("address"); err != nil {
		panic(err)
	}
}
