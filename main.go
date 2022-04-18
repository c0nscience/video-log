package main

import (
	"embed"
	"fmt"
	"github.com/c0nscience/video-log/internal/config"
	"github.com/c0nscience/video-log/internal/videos"
	"github.com/pkg/browser"
	"io/fs"
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

	//todo use gorillamux for easier mapping to path and method
	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc(videos.FetchUrl, videos.Fetch())
	http.HandleFunc(config.FetchUrl, config.Fetch())

	log.Println("Started")

	err = browser.OpenURL(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		log.Fatal("could not open the browser", err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
