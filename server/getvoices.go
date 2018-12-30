package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
)

// GetVoices retrieved available AWS voices for specified language code
// Parameters:
// - string language code
//
// Returns a DescribeVoicesOutput containing all voices, and any errors generated
func GetVoices(lan string, svc pollyiface.PollyAPI) (*polly.DescribeVoicesOutput, error) {
	input := &polly.DescribeVoicesInput{
		LanguageCode: aws.String(lan),
	}

	result, err := svc.DescribeVoices(input)
	if err != nil {
		return nil, err
	}
	return result, err
}
