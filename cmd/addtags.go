package cmd

import (
	"encoding/gob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

func init() {
	addCmd.AddCommand(addTagCmd)
}

var addTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Gets all available tags",
	Long:  `Get all available tags `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Insufficient parameters provided")
		}
		AddTag(args[0], args[1])
	},
}

// AddTag retreives all defined tags from persistant memory store
// located in cfg as
func AddTag(name, colour string) {
	f, err := os.Create(path.Join(viper.GetString("assets.tagsPath"), name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	encoder.Encode(Tag{Name: name, Colour: colour})
}
