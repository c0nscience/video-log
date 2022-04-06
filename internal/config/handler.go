package config

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type settings struct {
	Dir string `json:"dir"`
}

var Settings = &settings{}

const (
	FetchUrl   = "/config"
	configFile = ".vl.config"
)

func Fetch() func(http.ResponseWriter, *http.Request) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		b, _ := json.Marshal(Settings)
		ioutil.WriteFile(configFile, b, fs.ModePerm)
	}

	b, _ := ioutil.ReadFile(configFile)
	log.Println("file contents", string(b))
	if len(b) > 0 {
		_ = json.Unmarshal(b, Settings)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			bytes, err := json.Marshal(Settings)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
			w.Write(bytes)
		}

		if r.Method == http.MethodPost {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			err = ioutil.WriteFile(configFile, b, fs.ModePerm)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			err = json.Unmarshal(b, Settings)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}
