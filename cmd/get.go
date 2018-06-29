package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets all available voices",
	Long:  `The current version of Polly you are running`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
