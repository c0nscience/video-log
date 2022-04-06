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
	log.Println(string(b))
	if len(b) > 0 {
		_ = json.Unmarshal(b, Settings)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		bytes, err := json.Marshal(Settings)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(bytes)
	}
}
