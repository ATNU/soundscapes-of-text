package cmd

import (
	"encoding/gob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

type tag struct {
	Name  string
	Color string
}

func init() {
	addCmd.AddCommand(addTagCmd)
}

var addTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a new  tags",
	Long:  `Create a new tag`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Insufficient parameters provided")
		}
		AddTag(args[0], args[1])
	},
}

// AddTag creates a new tag to be used for ssml encoding
// Configuration:
// - assets.tagsPath
func AddTag(name, colour string) {
	f, err := os.Create(path.Join(viper.GetString("assets.tagsPath"), name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(Tag{Name: name, Colour: colour})
	if err != nil {
		log.Fatal(err)
	}
}
