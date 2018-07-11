package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add [type]",
	Long:  `Add can be used to retrieve AWS Polly supported assets`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
