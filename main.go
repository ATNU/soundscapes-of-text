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

// Generate mp3 file for provided ssml file
// AWS configuration files:
// - ~/.aws/
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

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"),
		TextType: aws.String("ssml"), Text: aws.String(s), VoiceId: aws.String(viper.GetString("voice"))}

	output, err := p.SynthesizeSpeech(input)
	if err != nil {
		log.Fatal(err)
	}

	names := strings.Split(os.Args[1], ".")
	mp3File := names[0] + ".mp3"
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
