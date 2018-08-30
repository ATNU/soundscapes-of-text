package cmd_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/polly/pollyiface"
	"github.com/gorilla/mux"
	"github.com/mattnolf/polly/cmd"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"
)

type mockPollyClient struct {
	pollyiface.PollyAPI
}

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

// SynthesizeSpeech mock AWS call throws error with invalid Speech Input
func (m *mockPollyClient) SynthesizeSpeech(*polly.SynthesizeSpeechInput) (*polly.SynthesizeSpeechOutput, error) {
	// Check voiceid exists
	// Validate body
	// - ssml
	// - text
	return nil, nil
}

// StartSpeechSynthesisTask mock AWS call throws error with invalid Task Input
func (m *mockPollyClient) StartSpeechSynthesisTask(*polly.StartSpeechSynthesisTaskInput) (*polly.StartSpeechSynthesisTaskOutput, error) {
	return nil, nil
}

func InitConfig() {
	viper.SetConfigName(".cobra")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

// TestAddTag asserts that adding a new tag successfully adds a persistant data blob
func TestAddTag(t *testing.T) {
	InitConfig()

	// Remove all files in the assets > tag folder

	tags, err := cmd.GetTags()
	if err != nil {
		t.Error(err)
	}
	pre := len(tags)

	cmd.AddTag("rayman", "red")

	tags, err = cmd.GetTags()
	if err != nil {
		t.Error(err)
	}
	if pre == len(tags) {
		t.Fatal("No new tag was generated")
	}
}

// TestGetTag asserts the successful retreival of persistant data blobs
func TestGetTag(t *testing.T) {
	InitConfig()
	cmd.AddTag("rayman", "red")
	tags, err := cmd.GetTags()
	if err != nil {
		t.Error(err)
	}

	if len(tags) != 2 {
		t.Error("Did not retreive tag")
	}

	tag := tags[1]
	if tag.Name != "rayman" || tag.Colour != "red" {
		t.Error("Did not retrieve tag details correctly")
	}
}

// TestGetVoicesByCode asserts that retreiving AWS Polly voices is handled successfully
func TestGetVoices(t *testing.T) {
	InitConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"en-GB", true},
		{"bla-bla", false},
	}

	mockSvc := &mockPollyClient{}
	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			_, err := cmd.GetVoices(tc.routeVariable, mockSvc)
			if tc.shouldPass && err != nil {
				t.Errorf("handler on GetVoices %s: got %v",
					tc.routeVariable, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("handler on GetVoices %s: got %v",
					tc.routeVariable, err)
			}

		})
	}
}

// TestGenerate asserts that generation of text-to-speech
// using AWS Polly API is successful and that any errors are handled
func TestGenerate(t *testing.T) {
	InitConfig()
	tt := []struct {
		id         string
		body       string
		shouldPass bool
	}{
		{"Brian", "Hello World", true},
		{"Nobody", "Hello World", false},
		{"Brian", "Hello <invalid>World</invalid>", false},
	}
	mockSvc := &mockPollyClient{}
	for _, tc := range tt {
		t.Run(tc.id, func(t *testing.T) {
			_, err := cmd.Generate(tc.body, tc.id,
				path.Join(viper.GetString("assets.ttsPath"), fmt.Sprint(time.Now().Unix())), mockSvc)
			if tc.shouldPass && err != nil {
				t.Errorf("handler on GetVoices %s: got %v",
					tc.id, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("handler on GetVoices %s: got %v",
					tc.id, err)
			}

		})
	}
}

// TestGenerateFromFile asserts that generation of
// text-to-speech encoding using AWS Polly uses
// configured file
func TestGenerateFromFile(t *testing.T) {
	InitConfig()
	mockSvc := &mockPollyClient{}
	err := cmd.GenerateFromFile(mockSvc)

	if err != nil {
		t.Error(err)
	}
}

// TestGenerateToS3 asserts that generation of text-to-speech
// usign AWS Polly API successfully stores result in S3
func TestGenerateToS3(t *testing.T) {
	InitConfig()
	tt := []struct {
		id         string
		body       string
		shouldPass bool
	}{
		{"Brian", "Hello World", true},
		{"Nobody", "Hello World", false},
		{"Brian", "Hello <invalid>World</invalid>", false},
	}

	mockSvc := &mockPollyClient{}

	for _, tc := range tt {
		t.Run(tc.id, func(t *testing.T) {
			_, err := cmd.GenerateToS3(tc.body, tc.id, mockSvc)

			if err != nil && tc.shouldPass {
				t.Errorf("Should have passed on Generating %s: got %v",
					tc.id, err)
			}
			if err == nil && !tc.shouldPass {
				t.Errorf("Should have not have passed on Generating %s: but did",
					tc.id)
			}
		})
	}
}

// TestRetreiveFromS3 asserts that retrieval of
// text-to-speech objects from AWS S3 is successful
// and any errors are handled
func TestRetrieveFromS3(t *testing.T) {
	// @deprecated
}

// TestHandleDemo asserts that a demo a synchronous encoding task is sent to AWS
// and the result returned
func TestWebserverHandleDemo(t *testing.T) {
	InitConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"Brian", false},
		{"Amy", false},
		{"MattyLad", false},
	}

	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			path := fmt.Sprintf("/demo/%v", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/demo/{id}", cmd.HandleDemo)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK && tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, w.Code, http.StatusOK)
			}

			if w.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, w.Code, http.StatusInternalServerError)
			}
		})
	}
}

func TestWebserverHandleGenerateFile(t *testing.T) {

}

// TestHandleGenerateS3 asserts that asynchronous encoding tasks are successfully sent to AWS
// and that resource URI returned when ready
func TestWebserverHandleGenerateS3(t *testing.T) {
	InitConfig()
	r := bytes.NewReader([]byte("Hello World!"))
	req, err := http.NewRequest("POST", "/generate?voice=Brian", r)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(cmd.HandleGenerateS3)

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	client := &http.Client{}
	resp, err := client.Get(rec.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Could not retrieve tts resource: got %v want %v",
			resp.Status, http.StatusOK)
	}
}

func TestWebserverHandleLanguages(t *testing.T) {

}

func TestWebserverHandleTags(t *testing.T) {

}

// TestWebserverHandleVoice asserts that get requests to AWS return successfully
func TestWebserverHandleVoices(t *testing.T) {
	InitConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"da-DK", true},
		{"nl-NL", true},
		{"en-AU", true},
		{"en-GB", true},
		{"en-IN", true},
		{"en-US", true},
		{"en-GB-WLS", true},
		{"fr-FR", true},
		{"fr-CA", true},
		{"de-DE", true},
		{"is-IS", true},
		{"it-IT", true},
		{"ja-JP", true},
		{"ko-KR", true},
		{"nb-NO", true},
		{"pl-PL", true},
		{"pt-BR", true},
		{"pt-PT", true},
		{"ro-RO", true},
		{"ru-RU", true},
		{"es-ES", true},
		{"es-US", true},
		{"sv-SE", true},
		{"tr-TR", true},
		{"cy-GB", true},
		{"bla-bla", false},
	}

	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			path := fmt.Sprintf("/voices/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/voices/{voice}", cmd.HandleVoices)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK && tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, w.Code, http.StatusOK)
			}

			if w.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, w.Code, http.StatusInternalServerError)
			}
		})
	}
}

func TestMyThing(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable}))
	p := polly.New(sess)
	cmd.GenerateToS3("Bla bla", "Brian", p)
}
