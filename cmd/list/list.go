package list

import (
	"flkcli/cmd"
	"flkcli/flkutils"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
)

var asUsername bool

// GetUserFromArgs returns the user id from the arguments and resolves it if necessary
var getUserIDFromArgs = func(args []string, client *flickr.FlickrClient) (userid string, err error) {
	if len(args) == 0 {
		return "", nil
	}

	userid = args[0]

	if asUsername {
		userid, err = flkutils.ResolveId(client, userid)
	}

	if err != nil {
		return "", fmt.Errorf("failed to resolve user id: %w", err)
	}

	return userid, nil
}

var SetCmd = &cobra.Command{
	Use:   "list",
	Short: "List flickr resources",
	Long: `List flickr resources
	You must specify the resource type to list
	Available resource types are: sets`,
}

func init() {
	cmd.RootCmd.AddCommand(SetCmd)
}
