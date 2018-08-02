package cmd_test

import (
	"github.com/mattnolf/polly/cmd"
	"log"
	"testing"
)

// TestAddTag asserts that adding a new tag successfully adds
// a persistant data blob
func TestAddTag(t *testing.T) {
	tags, err := cmd.GetTags()
	if err != nil {
		log.Println(err)
	}
	pre := len(tags)
	log.Println(pre)

	cmd.AddTag("jonyn", "red")

	tags, err = cmd.GetTags()
	if err != nil {
		t.SkipNow()
	}
	log.Println(len(tags))
	if pre == len(tags) {
		t.Fatal("No new tag was generated")
	}
}

// TestNewTagIntegrity asserts that newly generated tags maintain their integrity
// when marshalling/unmarshalling from gob
func TestAddTagIntegrity(t *testing.T) {
	// Create tag
	// Retrieve tag and inspect its values are identical
}
