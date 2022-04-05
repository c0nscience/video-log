package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
)

//go:embed public
var embeddedFiles embed.FS

const port = 8050

func main() {

	fsys, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic(err)
	}

	dir := "D:\\Raw"

	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		videos := []string{}
		for _, file := range files {
			videos = append(videos, file.Name())
		}

		bytes, err := json.Marshal(videos)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Printf("Server running on http://localhost:%d/", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
