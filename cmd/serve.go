package cmd

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Watches the text file and generates new tts on notify of file change",
	Long: `Serve watches for any file changes and generates a new tts file at the point of change.
WARNING, may be costly as generating new tts encodings often `,
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

// Serve starts a listener on the input text, and generates a new TTS encoding everytime the input
// text is modified.
func Serve() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add("hello.txt"); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event)
				Generate()
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()
	<-done
}
