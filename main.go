package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"

	"github.com/spf13/viper"

	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Generate text-to-speech encoding of input file
// Pass file to encode as command line argument
// - example: ./polly mytext.xml
// The following file extensions are currently compatible:
// - .txt for plaintext
// - .xml for speech synthesis markup language
// Returns a synthesis of text in the following possible file types
// - .mp3
// - .ogg
// - .pcm
func main() {
	cfgInit()
	if len(os.Args) != 2 {
		log.Fatal("No input file provided")
	}
	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	s := string(contents[:])
	names := strings.Split(os.Args[1], ".")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"),
		Text: aws.String(s), VoiceId: aws.String(viper.GetString("voice"))}

	switch names[1] {
	case "txt":
		input.SetTextType("text")
	case "xml":
		input.SetTextType("ssml")
	}

	output, err := p.SynthesizeSpeech(input)
	if err != nil {
		log.Fatal(err)
	}

	mp3File := names[0] + viper.GetString("output")
	outFile, err := os.Create(mp3File)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		log.Fatal(err)
	}
}

func cfgInit() {
	viper.SetConfigName("cfg")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
