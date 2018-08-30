package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	getCmd.AddCommand(getVoicesCmd)
}

var getVoicesCmd = &cobra.Command{
	Use:   "voices",
	Short: "get [languageCode]",
	Long:  `Gets all available voices for specified language code`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a Language Code")
		}
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable}))

		p := polly.New(sess)
		LogVoices(args[0], p)
	},
}

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

// LogVoices returns available AWS voices for specified language code
// Parameters:
// - string language code
func LogVoices(lan string, svc pollyiface.PollyAPI) {
	log.Println(GetVoices(lan, svc))
}
