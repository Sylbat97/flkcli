package login

import (
	"flkcli/cmd"
	"flkcli/config"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Interactive login",
	Long:  `Interactive login`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		apiConfig, err := config.GetAPIConfig()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		client := flickr.NewFlickrClient(apiConfig.Key, apiConfig.Secret)
		requestTok, err := flickr.GetRequestToken(client)
		if err != nil {
			fmt.Printf("Cannot get request token: %s\nPlease check your api key and secret. You can set them using flkcli setup", err)
			return
		}

		url, _ := flickr.GetAuthorizeUrl(client, requestTok)
		// Print url
		fmt.Printf("Please visit this URL to authorize the application: %s\n", url)
		// Ask user to input the code
		fmt.Print("Please enter the OAuth confirmation code: ")
		var confirmationCode string
		fmt.Scanln(&confirmationCode)

		accessTok, err := flickr.GetAccessToken(client, requestTok, confirmationCode)

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		//Print that the login was successful
		fmt.Println("Login successful")
		config.SetTokenConfig(accessTok.OAuthToken, accessTok.OAuthTokenSecret)
	},
}

func init() {
	cmd.RootCmd.AddCommand(LoginCmd)
}
