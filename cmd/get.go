package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get [type]",
	Long:  `Get can be used to retrieve AWS Polly supported assets`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
