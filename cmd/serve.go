package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Watches the input text file and generates new tts on notify of file change",
	Long: `Serve watches for any file changes and generates a new tts encoding at the point of file change.
WARNING, may be costly as generating new tts encodings often `,
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

// Serve starts a listener on the input text, and generates a new
// TTS encoding at the point the file is modified
func Serve() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(viper.GetString("input")); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event)
				GenerateFromFile()
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()
	<-done
}
