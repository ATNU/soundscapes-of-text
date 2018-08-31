package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable}))
var p *polly.Polly

// HandleLanguages returns all supported AWS Polly languages
// Response served as application/JSON with format (eg):
func HandleLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Language Request received")

	en := json.NewEncoder(w)
	en.Encode(GetLanguages())
}

// HandleVoices returns all supported AWS Polly voices
func HandleVoices(w http.ResponseWriter, r *http.Request) {
	p = polly.New(sess)

	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	log.Println("INFO: Voice request received: " + vars["voice"])
	v, err := GetVoices(vars["voice"], p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400 - Something bad happened!"))
		log.Println(err)
		return
	}

	en := json.NewEncoder(w)
	en.Encode(v.Voices)
}

// HandleDemo returns a demonstration of an AWS Polly voice
// speaking "Hi my name is: {voice}"
//
// Each voice demo is cached to avoid regeneration overheads
func HandleDemo(w http.ResponseWriter, r *http.Request) {
	p = polly.New(sess)
	vars := mux.Vars(r)
	log.Println("INFO: Voice request received " + vars["id"])

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(path.Join(viper.GetString("assets.demoPath"), vars["id"]) + ".mp3"); os.IsNotExist(err) {
		log.Println("INFO: No demo cache available - generating one")
		f, err := Generate(("Hi my name is " + mux.Vars(r)["id"]), mux.Vars(r)["id"],
			path.Join(viper.GetString("assets.demoPath"), mux.Vars(r)["id"]), p)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fi, err = f.Stat()
		if err != nil {
			log.Println("This error2")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(w, r, path.Join(viper.GetString("assets.demoPath"), fi.Name()))
}

// HandleGenerateS3 returns the S3 URL of the text-to-speech encoding task
// generated for the request body
// The generated URL is only returned once the resource is generated and s3 URI returns 200
func HandleGenerateS3(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Generate request received")

	p = polly.New(sess)
	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	f, err := GenerateToS3(string(bod[:]), r.URL.Query().Get("voice"), p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var timeout int
	for resp, err := http.Head(f); err == nil; {
		if resp.StatusCode != 200 || timeout < viper.GetInt("s3.maxRetryCount") {
			time.Sleep(2 * time.Second)
			resp, err = http.Head(f)
			timeout = +1
		}
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", viper.GetString("webserver.clientAddr"))
		w.Write([]byte(f))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}

// HandleGenerateFile returns a text-to-speech encoding of the provided
// request body
// deprecated
func HandleGenerateFile(w http.ResponseWriter, r *http.Request) {
	p = polly.New(sess)
	f, err := Generate(r.FormValue("ssml"), r.FormValue("voice"),
		path.Join(viper.GetString("assets.ttsPath"), fmt.Sprint(time.Now().Unix())), p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400 - Something bad happened!"))
		log.Println(err)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400 - Something bad happened!"))
		log.Println(err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	w.Header().Set("Access-Control-Allow-Origin", viper.GetString("webserver.clientAddr"))

	http.ServeFile(w, r, path.Join(viper.GetString("assets.ttsPath"), fi.Name()))
}
