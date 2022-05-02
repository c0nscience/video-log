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

func Load() {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		b, _ := json.Marshal(Settings)
		_ = ioutil.WriteFile(configFile, b, fs.ModePerm)
	}

	b, _ := ioutil.ReadFile(configFile)
	log.Println("config loaded", string(b))
	if len(b) > 0 {
		_ = json.Unmarshal(b, Settings)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(Settings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	_, _ = w.Write(bytes)
}

func Post(w http.ResponseWriter, r *http.Request) {
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
