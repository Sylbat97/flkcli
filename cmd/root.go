package cmd

import (
	"flkcli/flkutils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
)

var Client flickr.FlickrClient

var InitializeClient = func() error {
	newClient, err := flkutils.GetFlickrClient()
	if err != nil {
		return fmt.Errorf("failed to get flickr client: %w", err)
	}
	Client = *newClient
	return nil
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "flkcli",
	Short: "Flickr cli",
	Long:  `This CLI tool is used to interact with the Flickr API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
