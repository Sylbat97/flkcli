package set

import (
	"flkcli/config"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/people"
	"gopkg.in/masci/flickr.v3/photosets"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "list",
	Short: "List photo sets",
	Long: `List photo sets of a user
			If no user is specified, the currents user sets are listed`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get token from config
		token, err := config.GetTokenConfig()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		// Get api from config
		api, err := config.GetAPIConfig()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		var userid string

		client := flickr.NewFlickrClient(api.Key, api.Secret)
		client.OAuthToken = token.OAuthToken
		client.OAuthTokenSecret = token.OAuthTokenSecret

		username, _ := cmd.Flags().GetString("username")

		// If no args set user id to empty else set it to first arg
		if len(args) == 0 {
			userid = ""
		} else if username == "" {
			userid = args[0]
		} else {
			fmt.Print("You can't specify a username and a user id at the same time. Please specify only one.")
			return
		}

		if username != "" {
			user_response, err := people.FindByUsername(client, username)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			userid = user_response.User.Id
		}

		response, err := photosets.GetList(client, true, userid, 0)

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		fmt.Printf("Total sets: %d\n", response.Photosets.Total)
		for _, item := range response.Photosets.Items {
			fmt.Printf("- %s\n", item.Title)
		}

	},
}

func init() {
	SetCmd.AddCommand(addCmd)
	SetCmd.PersistentFlags().String("username", "", "The username of the user to list the sets")

}
