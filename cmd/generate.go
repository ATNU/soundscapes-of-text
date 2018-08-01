package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
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
// Configuration:
// - outputType
//
// Returns a pointer to the generated file and any errors generated
func Generate(body, id, path string) (*os.File, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String(viper.GetString("outputType")),
		TextType: aws.String("ssml"),
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
// Configuration:
// - input
// - voice
// - assets.ttsPath
//
// Any errors generated will be logged followed by graceful shutdown
func GenerateFromFile() {
	contents, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(contents[:])
	aws.String(viper.GetString("voice"))
	_, err = Generate(s, viper.GetString("voice"),
		path.Join(viper.GetString("assets.ttsPath"), fmt.Sprint(time.Now().Unix())))

	if err != nil {
		log.Fatal(err)
	}
}

// GenerateToS3 creates a synthesis task of the provided body of text
// using AWS Polly and stores within an AWS S3 Bucket.
// Logging of status are sent to AWS SNS.
//
// Parameters:
// - string body of text
// - string id of AWS Polly voice
//
// Configuration:
// - s3.outputFormat
// - s3.bucketName
// - sns.pollyTopicName
//
// Returns the ID of the asynchronous task
func GenerateToS3(body, id string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)
	task := new(polly.StartSpeechSynthesisTaskInput)

	task.SetOutputFormat(viper.GetString("s3.outputFormat"))
	task.SetOutputS3BucketName(viper.GetString("s3.bucketName"))
	task.SetText(body)
	task.SetTextType("text")
	task.SetSnsTopicArn(viper.GetString("sns.pollyTopicName"))
	task.SetVoiceId(id)

	o, err := p.StartSpeechSynthesisTask(task)
	if err != nil {
		return "", err
	}

	return *o.SynthesisTask.OutputUri, err
}

// retrieveFromS3 downloads an object from S3 into a local file.
// Parameters:
// - string bucket name
// - string key name
//
// Returns a pointer to the generated file
func retrieveFromS3(bucket, key string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	dl := s3manager.NewDownloader(sess)

	f, err := os.Create("test")
	if err != nil {
		log.Fatal(err)
	}

	n, err := dl.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n)
}
