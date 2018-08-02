package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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
// - GET -> /languages
// - GET -> /tags
// - GET -> /voices/{languageCoce}
// - GET -> /demo/{voiceID}
// - POST -> /generate?voice={voiceID}
//
// Any unrouted request will return 404 error
func WebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/languages", HandleLanguages).Methods("GET")
	r.HandleFunc("/tags", HandleTags).Methods("GET")
	r.HandleFunc("/voices/{voice}", HandleVoices).Methods("GET")
	r.HandleFunc("/demo/{id}", HandleDemo).Methods("GET")
	r.HandleFunc("/generate", HandleGenerateS3).Methods("POST")

	srv := &http.Server{
		Addr:         viper.GetString("webserver.addr"),
		WriteTimeout: time.Second * viper.GetDuration("webserver.timeout.write"),
		ReadTimeout:  time.Second * viper.GetDuration("webserver.timeout.read"),
		IdleTimeout:  time.Second * viper.GetDuration("webserver.timeout.idle"),
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("webserver.timeout.cancel"))
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("INFO: Gracefully shutting down")
	os.Exit(0)
}

// HandleLanguages returns all supported AWS Polly languages
// Response served as application/JSON with format (eg):
//	{
//		"LanguageName": "English American",
//		"LanguageCode": "en-US"
//	},
//
// Any errors contacting AWS will be returned to client
func HandleLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Language Request received")

	en := json.NewEncoder(w)
	en.Encode(GetLanguages())
}

// HandleTags returns all defined tags to be used for
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
func HandleTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Tag Request received")

	// Pre defined set of tags stored in memory
	// Retreieve file and build up struct slice
	// JSON write slice

}

// HandleVoices returns all supported AWS Polly voices
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
func HandleVoices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Voice request received: " + vars["voice"])
	v, err := GetVoices(vars["voice"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println(err)
		return
	}

	en := json.NewEncoder(w)
	en.Encode(v.Voices)
}

// HandleDemo returns a demonstration of an AWS Polly voice
// speaking "Hi my name is: {voice}"
// Response served as audio/mpeg with format (eg):
//
// Each voice demo is cached to avoid regeneration overheads
//
// Any errors encountered will return a 500 error to client
func HandleDemo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("INFO: Voice request received " + vars["id"])

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path.Join(viper.GetString("assets.demoPath"), vars["id"]) + ".mp3"); os.IsNotExist(err) {
		log.Println("INFO: No demo cache available - generating one")
		log.Println(("Hi my name is " + mux.Vars(r)["id"]))
		f, err := Generate(("Hi my name is " + mux.Vars(r)["id"]), mux.Vars(r)["id"],
			path.Join(viper.GetString("assets.demoPath"), mux.Vars(r)["id"]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Println(err)
			return
		}
		fi, err = f.Stat()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Println(err)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")

	http.ServeFile(w, r, path.Join(viper.GetString("assets.demoPath"), fi.Name()))
}

// HandleGenerateS3 returns the S3 URL of the text-to-speech encoding task
// generated for the request body
// The generated URL is only returned once the resource is generated and s3 URI returns 200
//
// Any errors encountered will return a 500 error to client
func HandleGenerateS3(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Generate request received")
	log.Println("Generate: ")

	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println(err)
		return
	}

	f, err := GenerateToS3(string(bod[:]), r.URL.Query().Get("voice"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.Println(err)
		return
	}

	var timeout int
	for resp, err := http.Head(f); resp.StatusCode != 200 || timeout >= viper.GetInt("s3.maxRetryCount"); {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
			log.Println(err)
			return
		}
		time.Sleep(2 * time.Second)
		resp, err = http.Head(f)
		timeout = +1
		log.Println("Code", resp.StatusCode)
	}

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", viper.GetString("webserver.clientAddr"))
	w.Write([]byte(f))
}

// HandleGenerateFile returns a text-to-speech encoding of the provided
// request body
// Response served as audio/mpeg with format (eg):
//
// Any errors encountered will return a 500 error to client
func HandleGenerateFile(w http.ResponseWriter, r *http.Request) {
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
