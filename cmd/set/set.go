package set

import (
	"flkcli/cmd"
	"flkcli/flkutils"
	"fmt"

	"github.com/spf13/cobra"
)

var asUsername bool

// GetUserFromArgs returns the user id from the arguments and resolves it if necessary
var getUserIDFromArgs = func(args []string) (userid string, err error) {
	if len(args) == 0 {
		return "", nil
	}

	userid = args[0]

	if asUsername {
		userid, err = flkutils.ResolveId(&cmd.Client, userid)
	}

	if err != nil {
		return "", fmt.Errorf("failed to resolve user id: %w", err)
	}

	return userid, nil
}

var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Manage sets",
	Long:  `Manage sets`,
}

func init() {
	cmd.RootCmd.AddCommand(SetCmd)
}
