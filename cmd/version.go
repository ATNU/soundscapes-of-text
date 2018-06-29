package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(verCmd)
}

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Polly",
	Long:  `The current version of Polly you are running`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Polly v0.1 -- HEAD")
	},
}
