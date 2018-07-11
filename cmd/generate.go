package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a TTS encoding of the provided text",
	Long: `Generate TTS encoding using AWS Polly. Text taken from
input cfg value`,
	Run: func(cmd *cobra.Command, args []string) {
		GenerateFromFile()
	},
}

// Generate creates a text-to-speech encoding of the provided body of text
// using AWS Polly
//
// Parameters:
// - string body of text
// - string id of AWS Polly voice
// - string path of output file
//
// Returns a pointer to the generated file and any errors generated
func Generate(body, id, path string) (*os.File, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String(viper.GetString("outputtype")),
		TextType: aws.String("text"),
		Text:     aws.String(body),
		VoiceId:  aws.String(id)}

	output, err := p.SynthesizeSpeech(input)
	if err != nil {
		return nil, err
	}

	outFile, err := os.Create(path + ".mp3")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		return nil, err
	}
	return outFile, nil
}

// GenerateFromFile creates a text-to-speech encoding of the provided text file
// The generated file is placed in the root folder
//
// Any errors generated will be logged followed by graceful shutdown
func GenerateFromFile() {
	contents, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(contents[:])
	aws.String(viper.GetString("voice"))
	Generate(s, viper.GetString("voice"), "output")
}
