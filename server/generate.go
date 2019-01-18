package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	"github.com/spf13/viper"
)

// Generate creates a text-to-speech encoding of the provided body of text
// using AWS Polly
//
// Parameters:
// - string body of text
// - string id of AWS Polly voice
// - string path of output file
// - pollyiface.PollyAPI Polly instance with valid session
//
// Configuration:
// - outputType
//
// Returns a pointer to the generated file and any errors generated
func Generate(body, id, path string, svc pollyiface.PollyAPI) (*os.File, error) {
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String(viper.GetString("outputType")),
		TextType: aws.String("text"),
		Text:     aws.String(body),
		VoiceId:  aws.String(id)}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	outFile, err := os.Create(path + ".mp3")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return outFile, nil
}

// GenerateFromFile creates a text-to-speech encoding of the provided text file
// The generated file is placed in the root folder
//
// Parameters:
// - pollyiface.PollyAPI Polly instance with valid session
//
// Configuration:
// - input
// - voice
// - assets.ttsPath
//
// Returns any errors generated
func GenerateFromFile(svc pollyiface.PollyAPI) error {
	contents, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		return err
	}
	s := string(contents[:])
	aws.String(viper.GetString("voice"))
	_, err = Generate(s, viper.GetString("voice"),
		path.Join(viper.GetString("assets.ttsPath"), fmt.Sprint(time.Now().Unix())), svc)

	if err != nil {
		return err
	}

	return nil
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
// - pollyiface.PollyAPI Polly instance with valid session
//
// Returns the ID of the asynchronous task
func GenerateToS3(body, id string, svc pollyiface.PollyAPI) (string, error) {
	task := new(polly.StartSpeechSynthesisTaskInput)

	task.SetOutputFormat(viper.GetString("s3.outputFormat"))
	task.SetOutputS3BucketName(viper.GetString("s3.bucketName"))
	task.SetText(body)
	task.SetTextType("ssml")
	task.SetSnsTopicArn(viper.GetString("sns.pollyTopicName"))
	task.SetVoiceId(id)

	o, err := svc.StartSpeechSynthesisTask(task)
	if err != nil {
		return "", err
	}
	return *o.SynthesisTask.OutputUri, nil
}
