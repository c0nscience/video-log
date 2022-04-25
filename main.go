package main

import (
	"embed"
	"fmt"
	"github.com/c0nscience/video-log/internal/config"
	"github.com/c0nscience/video-log/internal/videos"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed public
var embeddedFiles embed.FS

const port = 8050

func main() {

	fsys, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc(videos.FetchUrl, videos.Get).Methods("GET")
	r.HandleFunc(videos.FetchUrl, videos.Post).Methods("POST")
	r.HandleFunc(videos.FetchUrl, videos.Delete).Methods("DELETE")
	r.HandleFunc(config.FetchUrl, config.Get).Methods("GET")
	r.HandleFunc(config.FetchUrl, config.Post).Methods("POST")
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
