package cmd

import (
	"encoding/gob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Tag defines the required information for an SSML tag
type Tag struct {
	Name   string
	Colour string
	Tags   []string
}

func init() {
	getCmd.AddCommand(getTagsCmd)
}

var getTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Gets all available tags",
	Long:  `Get all available tags `,
	Run: func(cmd *cobra.Command, args []string) {
		GetTags()
	},
}

// GetTags retreives all defined tags from persistant memory store
func GetTags() {
	dir, err := ioutil.ReadDir(viper.GetString("assets.tagsPath"))
	if err != nil {
		log.Println(err)
		return
	}

	ts := make([]Tag, 0)

	for _, val := range dir {
		inFile, err := os.Open(path.Join(viper.GetString("assets.tagsPath"), val.Name()))
		if err != nil {
			log.Fatal(err)
		}
		enc := gob.NewDecoder(inFile)
		var t Tag
		enc.Decode(&t)
		ts = append(ts, t)
	}
	log.Println(ts)
}
