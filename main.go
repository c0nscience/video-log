package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/pkg/browser"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

//go:embed public
var embeddedFiles embed.FS

const port = 8050

type video struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

		videoDesc := map[string]string{}
		for _, file := range files {
			name := file.Name()
			isVideo := strings.HasSuffix(name, "mkv")
			if isVideo {
				name = strings.TrimSuffix(name, ".mkv")
			} else {
				name = strings.TrimSuffix(name, ".txt")
			}

			v, ok := videoDesc[name]
			if !ok {
				v = ""
			}

			if !isVideo {
				b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
				if err != nil {
					log.Fatal(err)
				}

				v = string(b)
			}
			videoDesc[name] = v
		}

		videos := []video{}

		for name, desc := range videoDesc {
			videos = append(videos, video{
				Name:        name,
				Description: desc,
			})
		}

		sort.Slice(videos, func(i, j int) bool {
			return videos[i].Name > videos[j].Name
		})

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

	log.Println("Started")

	err = browser.OpenURL(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
