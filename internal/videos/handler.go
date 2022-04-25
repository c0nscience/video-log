package videos

import (
	"encoding/json"
	"fmt"
	"github.com/c0nscience/video-log/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type video struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

const (
	FetchUrl = "/videos"
)

func Get(w http.ResponseWriter, r *http.Request) {
	dir := config.Settings.Dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	videoDesc := map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
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
		d, _ := time.Parse("2006-01-02-15-04-05", name)
		fd := d.Format(time.ANSIC)
		videos = append(videos, video{
			Name:        name,
			Description: desc,
			Date:        fd,
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

func Post(w http.ResponseWriter, r *http.Request) {
	dir := config.Settings.Dir
	b, _ := ioutil.ReadAll(r.Body)
	v := video{}
	_ = json.Unmarshal(b, &v)

	err := ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", dir, v.Name), []byte(v.Description), 0644)
	if err != nil {
		log.Printf("could not update the file: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	dir := config.Settings.Dir
	b, _ := ioutil.ReadAll(r.Body)
	v := video{}
	_ = json.Unmarshal(b, &v)

	_ = os.Remove(fmt.Sprintf("%s/%s.txt", dir, v.Name))
	_ = os.Remove(fmt.Sprintf("%s/%s.mkv", dir, v.Name))

	w.WriteHeader(http.StatusOK)
}
