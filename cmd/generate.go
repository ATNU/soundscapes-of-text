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
	"strings"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Print the version number of Polly",
	Long:  `The current version of Polly you are running`,
	Run: func(cmd *cobra.Command, args []string) {
		Generate()
	},
}

// Generate text-to-speech encoding of input file
// Returns a synthesis of text in the following possible file types
// - .mp3
// - .ogg
// - .pcm
func Generate() {

	contents, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(contents[:])

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String(viper.GetString("outputtype")),
		TextType: aws.String("ssml"),
		Text:     aws.String(s),
		VoiceId:  aws.String(viper.GetString("voice"))}

	output, err := p.SynthesizeSpeech(input)
	if err != nil {
		log.Fatal(err)
	}

	outname := strings.Split(viper.GetString("input"), ".")[0] + "." + viper.GetString("outputtype")
	log.Println(outname)
	outFile, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		log.Fatal(err)
	}
}
