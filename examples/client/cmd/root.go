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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&satelitAddress, "address", "", "", "address of satelit")
	if err := rootCmd.MarkPersistentFlagRequired("address"); err != nil {
		panic(err)
	}
}
