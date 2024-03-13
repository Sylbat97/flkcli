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
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Setup called")
		// Get apikey and apisecret
		apikey, _ := cmd.Flags().GetString("apikey")
		apisecret, _ := cmd.Flags().GetString("apisecret")
		// Call CreateConfigFile
		config.SetApiConfig(apisecret, apikey)
	},
}

func init() {
	SetupCmd.PersistentFlags().String("apikey", "", "Your Flickr API key")
	SetupCmd.MarkPersistentFlagRequired("apikey")
	SetupCmd.PersistentFlags().String("apisecret", "", "Your Flickr API secret")
	SetupCmd.MarkPersistentFlagRequired("apisecret")
	cmd.RootCmd.AddCommand(SetupCmd)
}
