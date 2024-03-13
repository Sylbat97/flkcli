package set

import (
	"flkcli/cmd"

	"github.com/spf13/cobra"
)

// barCmd represents the bar command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Manage sets",
	Long:  `Manage sets`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("set called")
	// },
}

func init() {
	cmd.RootCmd.AddCommand(SetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// barCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// barCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
