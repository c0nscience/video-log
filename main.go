package main

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/c0nscience/video-log/internal/config"
	"github.com/c0nscience/video-log/internal/videos"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed public
var embeddedFiles embed.FS

const port = 8050

type Timestamp struct {
	Path string `json:"path"`
}

func main() {

	fsys, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic(err)
	}

	//C:\Video\Rendered\One and Two\Poor Mans UV Mapping - Result.csv
	r := mux.NewRouter()
	r.HandleFunc(videos.FetchUrl, videos.Get).Methods("GET")
	r.HandleFunc(videos.FetchUrl, videos.Post).Methods("POST")
	r.HandleFunc(videos.FetchUrl, videos.Delete).Methods("DELETE")
	r.HandleFunc(config.FetchUrl, config.Get).Methods("GET")
	r.HandleFunc(config.FetchUrl, config.Post).Methods("POST")
	r.HandleFunc("/tools/timestamps", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("could not read request body", err)
			return
		}

		t := Timestamp{}
		err = json.Unmarshal(b, &t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("could not unmarshal request body", err)
			return
		}

		file, err := os.Open(t.Path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("could not open file: '%s'. Error: %v", t.Path, err)
			return
		}
		defer file.Close()

		stat, _ := file.Stat()
		log.Printf("size: %d", stat.Size())

		re := csv.NewReader(file)
		re.LazyQuotes = true
		records, err := re.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("could not read file '%s' into csv. Error: %v", t.Path, err)
			return
		}

		result := []string{}
		for _, record := range records[1:] {
			tmsp := record[6][3 : len(record[6])-3]
			note := record[19]

			result = append(result, fmt.Sprintf("%s %s", tmsp, note))
		}

		res := struct {
			Timestamp string `json:"timestamp"`
		}{
			Timestamp: strings.Join(result, "\n"),
		}

		b, err = json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("could not marshal response", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.FS(fsys)))

	config.Load()

	err = browser.OpenURL(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		log.Fatal("could not open the browser", err)
	}

	log.Println("Started")
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
