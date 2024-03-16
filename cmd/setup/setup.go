package setup

import (
	"flkcli/cmd"
	"flkcli/config"

	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cli",
	Long:  `Setup cli`,
	Run: func(command *cobra.Command, args []string) {
		command.Println("Setup called")
		// Get apikey and apisecret
		apikey, _ := command.Flags().GetString("apikey")
		apisecret, _ := command.Flags().GetString("apisecret")
		// Call CreateConfigFile
		config.SetApiConfig(apisecret, apikey)
	},
}

func init() {
	// Add flags to the setup command
	// The apikey and apisecret flags are required
	// The apikey and apisecret flags are used to set the api key and secret in the config file
	SetupCmd.PersistentFlags().String("apikey", "", "Your Flickr API key")
	SetupCmd.MarkPersistentFlagRequired("apikey")
	SetupCmd.PersistentFlags().String("apisecret", "", "Your Flickr API secret")
	SetupCmd.MarkPersistentFlagRequired("apisecret")
	cmd.RootCmd.AddCommand(SetupCmd)
}
