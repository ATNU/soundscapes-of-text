package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func init() {
	rootCmd.AddCommand(webserverCmd)
}

var webserverCmd = &cobra.Command{
	Use:   "webserver",
	Short: "Launches a RESTful Web API for backend behaviour",
	Long: `Webserver creates a RESTful http api for backend bevahiour.
	The following requests can be handled:
	- /languages Get all supported languages
	- /voice/{voice} Get all voices of specified language
	- /demo/{id} Get a demo sample of the specified voice`,
	Run: func(cmd *cobra.Command, args []string) {
		WebServer()
	},
}

// WebServer initiates a HTTP webserver for providing
// a RESTful API
//
// The following paths are handled:
// - GET -> languages
// - GET -> tags
// - GET -> voices/{languageCoce}
// - GET -> demo/{voiceID}
// - POST -> generate
//
// Any unrouted request will return 404 error
func WebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/languages", handleLanguages).Methods("GET")
	r.HandleFunc("/tags", handleTags).Methods("GET")
	r.HandleFunc("/voices/{voice}", handleVoices).Methods("GET")
	r.HandleFunc("/demo/{id}", handleDemo).Methods("GET")
	r.HandleFunc("/generate", handleGenerate).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// handleLanguages returns all supported AWS Polly languages
// Response served as application/JSON with format (eg):
//	{
//		"LanguageName": "English American",
//		"LanguageCode": "en-US"
//	},
//
// Any errors contacting AWS will be returned to client
func handleLanguages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("Request received " + vars["voice"])

}

// handleTags returns all defined tags to be used for
// generating SSML tts encodings
// Response served as application/JSON with format (eg):
//	{
//		"Name": "Happy",
//		"tags": {
//					"tone": "x",
//					"break": "2s"
//				}
//	},
//
// Any errors generated will be returned to client
func handleTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Tag Request received")

	// Pre defined set of tags stored in memory
	// Retreieve file and build up struct slice
	// JSON write slice

}

// handleVoices returns all supported AWS Polly voices
// Response served as application/JSON with format (eg):
//	{
//		"Gender": "Male",
//		"ID": "Brian",
//		"LanguageCode": "en-GB",
//		"LanguageName": "English",
//		"Name": "Brian",
//	},
//
// Any errors generated will be returned to client
func handleVoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Voice request received: " + vars["voice"])
	v := GetVoices(vars["voice"])

	en := json.NewEncoder(w)
	en.Encode(v.Voices)
}

// handleDemo returns a demonstration of an AWS Polly voice
// speaking "Hi my name is: {voice}"
// Response served as audio/mpeg with format (eg):
//
// Each voice demo is cached to avoid regeneration overheads
//
// Any errors encountered will return a 500 error to client
func handleDemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("INFO: Voice request received " + vars["id"])

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path.Join(viper.GetString("assets.demoPath"), vars["id"]) + ".mp3"); os.IsNotExist(err) {
		log.Println("INFO: No demo cache available - generating one")
		f, err := Generate(("Hi my name is " + mux.Vars(r)["id"]), mux.Vars(r)["id"],
			path.Join(viper.GetString("assets.demoPath"), mux.Vars(r)["id"]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Println("Error found", err)
			return
		}
		fi, err = f.Stat()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Println("Error found", err)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")

	http.ServeFile(w, r, path.Join(viper.GetString("assets.demoPath"), fi.Name()))
}

// handleGenerate returns a text-to-speach encoding of the provided
// request body
// Response served as audio/mpeg with format (eg):
//
// Any errors encountered will return a 500 error to client
func handleGenerate(w http.ResponseWriter, r *http.Request) {
	f, err := Generate(r.FormValue("ssml"), r.FormValue("voice"),
		path.Join(viper.GetString("assets.ttsPath"), fmt.Sprint(time.Now().Unix())))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println(err)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println(err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	w.Header().Set("Access-Control-Allow-Origin", viper.GetString("webserver.clientAddr"))

	http.ServeFile(w, r, path.Join(viper.GetString("assets.ttsPath"), fi.Name()))
}
