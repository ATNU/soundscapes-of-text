package main_test

import (
	"fmt"
	"github.com/gorilla/mux"
	sot "github.com/mattnolf/sot"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleVoices(t *testing.T) {
	sot.SetupConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"en-GB", false}, // Set to true
		{"jgg", false},
	}

	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			path := fmt.Sprintf("/voices/%v", tc.routeVariable)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/voices/{voice}", sot.HandleVoices)
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

func TestHandleDemo(t *testing.T) {
	sot.SetupConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"Brian", false}, // Set to true
		{"NotBrian", false},
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
			r.HandleFunc("/demo/{id}", sot.HandleDemo)
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

func TestHandleGenerateS3(t *testing.T) {
	sot.SetupConfig()
	tt := []struct {
		routeVariable string
		shouldPass    bool
	}{
		{"Brian", false}, // Set to true
	}

	for _, tc := range tt {
		t.Run(tc.routeVariable, func(t *testing.T) {
			path := fmt.Sprintf("/generate/%v", tc.routeVariable)
			req, err := http.NewRequest("POST", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/generate/{id}", sot.HandleDemo)
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
