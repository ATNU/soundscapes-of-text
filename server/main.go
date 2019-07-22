package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

const cfgName string = ".cfg"

func main() {
	SetupConfig()
	WebServer()
}

// SetupConfig loads configuration file and watched for modifications
func SetupConfig() {
	viper.SetConfigName(cfgName)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Failed to find configuration file")
	}
	log.Println("Successfully found configuration file - test deploy to azure cloud" + cfgName)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuration file modified\nRestarting connection...")
	})
}
