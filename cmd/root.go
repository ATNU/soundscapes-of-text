//Package cmd contains all cli commands
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "pol",
	Short: "pol uses AWS Polly to generate an audio tts encoding",
	Long:  `pol uses AWS Polly to generate an audio tts encoding of .txt or s`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

// Initialise viper
func initConfig() {
	viper.SetConfigName(".cobra")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
}

// Execute a command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
