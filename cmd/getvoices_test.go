package cmd_test

import (
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
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"en-GB", true},
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

			if rec.Code == http.StatusOK && !tc.shouldPass {
				if rec.Code == http.StatusOK && !tc.shouldPass {
					t.Errorf("handler on routeVariable %s: got %v want %v",
						tc.routeVariable, rec.Code, http.StatusInternalServerError)
				}
			}

			if rec.Code == http.StatusOK && tc.shouldPass {
				expected := `[{"Gender":"Female","Id":"Emma","LanguageCode":"en-GB","LanguageName":"British English","Name":"Emma"},{"Gender":"Male","Id":"Brian","LanguageCode":"en-GB","LanguageName":"British English","Name":"Brian"},{"Gender":"Female","Id":"Amy","LanguageCode":"en-GB","LanguageName":"British English","Name":"Amy"}]`

				if expected == rec.Body.String() {
					t.Errorf("handler returned unexpected body: got %s want %s",
						rec.Body.String(), expected)
				}
			}
		})
	}
}
