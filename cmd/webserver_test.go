package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattnolf/polly/cmd"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func InitConfig() {
	viper.SetConfigName(".cobra")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func TestHandleLanguages(t *testing.T) {

}

func TestHandleTags(t *testing.T) {

}

// TestHandleVoice asserts that get requests to AWS return successfully
func TestHandleVoices(t *testing.T) {
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

			rec := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/voices/{voice}", cmd.HandleVoices)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK && tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, rec.Code, http.StatusOK)
			}

			if rec.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, rec.Code, http.StatusInternalServerError)
			}
		})
	}
}

// TestHandleDemo asserts that a demo a synchronous encoding task is sent to AWS
// and the result returned
func TestHandleDemo(t *testing.T) {

	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"Brian", true},
		{"MattyLad", false},
	}

	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			path := fmt.Sprintf("/demo/%s", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/demo/{id}", cmd.HandleDemo)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK && tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, rec.Code, http.StatusOK)
			}

			if rec.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					tc.routeVariable, rec.Code, http.StatusInternalServerError)
			}
		})
	}
}

// TestHandleGenerateS3 asserts that asynchronous encoding tasks are successfully sent to AWS
// and that resource URI returned when ready
func TestHandleGenerateS3(t *testing.T) {
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

func TestHandleGenerateFile(t *testing.T) {

}
