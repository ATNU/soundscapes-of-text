package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattnolf/polly/cmd"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetVoicesByCode asserts that retreiving AWS Polly voices is
// handled successfully
func TestGetVoices(t *testing.T) {
	InitConfig()
	c := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"en-GB", true},
		{"bla-bla", false},
	}

	for _, val := range c {
		path := fmt.Sprintf("/voices/%s", val.routeVariable)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		r := mux.NewRouter()
		r.HandleFunc("/voices/{voice}", cmd.HandleVoices)
		r.ServeHTTP(rec, req)

		if rec.Code == http.StatusOK && !val.shouldPass {
			if rec.Code == http.StatusOK && !val.shouldPass {
				t.Errorf("handler on routeVariable %s: got %v want %v",
					val.routeVariable, rec.Code, http.StatusInternalServerError)
			}
		}

		if rec.Code == http.StatusOK && val.shouldPass {
			expected := []byte(`[{"Gender":"Female","Id":"Emma","LanguageCode":"en-GB","LanguageName":"British English","Name":"Emma"},{"Gender":"Male","Id":"Brian","LanguageCode":"en-GB","LanguageName":"British English","Name":"Brian"},{"Gender":"Female","Id":"Amy","LanguageCode":"en-GB","LanguageName":"British English","Name":"Amy"}]`)

			if bytes.Compare(expected, rec.Body.Bytes()) != 0 {
				t.Errorf("handler returned unexpected body: got %s want %s",
					rec.Body.String(), expected)
			}
		}
	}
}
