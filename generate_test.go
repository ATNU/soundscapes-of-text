package main_test

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	sot "github.com/mattnolf/sot"
	"regexp"
	"testing"
)

type mockPollyClient struct {
	pollyiface.PollyAPI
}

func (m *mockPollyClient) StartSpeechSynthesisTask(p *polly.StartSpeechSynthesisTaskInput) (*polly.StartSpeechSynthesisTaskOutput, error) {
	if *p.VoiceId != "Brian" {
		return nil, errors.New("Bad voice")
	}
	match, err := regexp.MatchString("[A-Za-z0-9]+", *p.Text)
	if match == false || err != nil {
		return nil, errors.New("Invalid text")
	}
	out := new(polly.StartSpeechSynthesisTaskOutput)
	out.SetSynthesisTask(new(polly.SynthesisTask))
	out.SynthesisTask.SetOutputUri("bla")
	return out, nil
}

func TestGenerateToS3(t *testing.T) {
	sot.SetupConfig()
	mockSvc := &mockPollyClient{}
	tt := []struct {
		body       string
		id         string
		polly      pollyiface.PollyAPI
		shouldPass bool
	}{
		{"Hello World", "Brian", mockSvc, true},
		{"Hello World", "NotName", mockSvc, false},
		{"<>", "Brian", mockSvc, false},
	}
	for _, tc := range tt {
		t.Run(tc.body, func(t *testing.T) {
			_, err := sot.GenerateToS3(tc.body, tc.id, tc.polly)
			if err != nil && tc.shouldPass {
				t.Errorf("%s failed and should have passed: ", tc.body)
			}
			if err == nil && !tc.shouldPass {
				t.Errorf("%s passed and should have failed: ", tc.body)
			}
		})
	}
}
