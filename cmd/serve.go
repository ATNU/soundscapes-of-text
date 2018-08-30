package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
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
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable}))
		p := polly.New(sess)
		err := Serve(p)
		if err != nil {
			log.Println(err)
		}
	},
}

// Serve starts a listener on the input text, and generates a new
// TTS encoding at the point the file is modified
//
// Parameters:
// - pollyiface.PollyAPI Polly instance with valid session
//
// Returns any errors encountered
func Serve(svc pollyiface.PollyAPI) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(viper.GetString("input")); err != nil {
		return err
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event)
				GenerateFromFile(svc)
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()
	<-done
	return nil
}
