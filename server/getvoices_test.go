package main_test

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	sot "github.com/mattnolf/sot"
	"testing"
)

// DescribeVoices mock AWS call throws error if code invalid
func (m *mockPollyClient) DescribeVoices(p *polly.DescribeVoicesInput) (*polly.DescribeVoicesOutput, error) {
	supportedCodes := []string{
		"da-DK",
		"nl-NL",
		"en-AU",
		"en-GB",
		"en-IN",
		"en-US",
		"en-GB-WLS",
		"fr-FR",
		"fr-CA",
		"de-DE",
		"is-IS",
		"it-IT",
		"ja-JP",
		"ko-KR",
		"nb-NO",
		"pl-PL",
		"pt-BR",
		"pt-PT",
		"ro-RO",
		"ru-RU",
		"es-ES",
		"es-US",
		"sv-SE",
		"tr-TR",
		"cy-GB",
	}

	for _, val := range supportedCodes {
		if *p.LanguageCode == val {
			return nil, nil
		}
	}
	return nil, errors.New("Invalid languagecode")
}

func TestGetVoices(t *testing.T) {
	sot.SetupConfig()
	mockSvc := &mockPollyClient{}
	tt := []struct {
		code       string
		polly      pollyiface.PollyAPI
		shouldPass bool
	}{
		{"en-GB", mockSvc, true},
		{"uh-HH", mockSvc, false},
	}
	for _, tc := range tt {
		t.Run(tc.code, func(t *testing.T) {
			_, err := sot.GetVoices(tc.code, tc.polly)
			if err != nil && tc.shouldPass {
				t.Errorf("%s failed and should have passed: ", tc.code)
			}
			if err == nil && !tc.shouldPass {
				t.Errorf("%s passed and should have failed: ", tc.code)
			}
		})
	}
}
