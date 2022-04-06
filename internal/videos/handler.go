package videos

import (
	"encoding/json"
	"fmt"
	"github.com/c0nscience/video-log/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type video struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

const (
	FetchUrl = "/videos"
)

func Fetch() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := config.Settings.Dir
		log.Printf("directory \"%s\"", dir)
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
	}
}
