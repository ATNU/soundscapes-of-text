package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	getCmd.AddCommand(getLexiconsCmd)
}

var getLexiconsCmd = &cobra.Command{
	Use:   "lexicons",
	Short: "Gets all available lexicons",
	Long:  `The current version of Polly you are running`,
	Run: func(cmd *cobra.Command, args []string) {
		GetLexicons()
	},
}

// GetLexicons returns available AWS Lexicons
func GetLexicons() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))

	p := polly.New(sess)

	l := &polly.ListLexiconsInput{}

	result, err := p.ListLexicons(l)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
