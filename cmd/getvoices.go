package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	getCmd.AddCommand(getVoicesCmd)
}

var getVoicesCmd = &cobra.Command{
	Use:   "voices",
	Short: "Gets all available voices",
	Long:  `The current version of Polly you are running`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a Language Code")
		}
		GetVoices(args[0])
	},
}

// GetVoices returns available AWS voices for specified language code
// Parameters:
// - string language code
func GetVoices(lan string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))

	p := polly.New(sess)

	input := &polly.DescribeVoicesInput{
		LanguageCode: aws.String(lan),
	}

	result, err := p.DescribeVoices(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
