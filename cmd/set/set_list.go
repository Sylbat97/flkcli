package set

import (
	"fmt"

	"github.com/spf13/cobra"

	"flkcli/cmd"
	"flkcli/flkutils"
)

// The list command is used to list the photo sets of a user
var listCmd = &cobra.Command{
	Use:   "list [userid]",
	Short: "List photo sets",
	Long: `List photo sets of a user
		If no user is specified, the currents user sets are listed`,
	Args: cobra.MaximumNArgs(1),
	Run: func(command *cobra.Command, args []string) {
		// Initialize the client
		err := cmd.InitializeClient()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		// Get the user id from the arguments
		userid, err := getUserIDFromArgs(args)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		// Get the list of sets
		total, photoSetItems, err := flkutils.ListSets(&cmd.Client, userid)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		fmt.Printf("Total sets: %d\n", total)
		for _, item := range photoSetItems {
			fmt.Printf("- %s\n", item.Title)
		}

	},
}

func init() {
	listCmd.PersistentFlags().BoolVar(&asUsername, "as-username", false, "Treat the user id as a username")
	SetCmd.AddCommand(listCmd)
}
