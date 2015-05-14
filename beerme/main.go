package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

type Comment struct {
	Author string `json:"author"`
	Text string `json:"text"`
}

func staticDir() string {
	_, file, _, _ := runtime.Caller(0)
	parent, err := filepath.Abs(
		filepath.Join(filepath.Dir(file), ".."))
	if err != nil {
		log.Panic(err)
	}
	return filepath.Join(parent, "static")
}

func main() {
	comments := append(make([]Comment, 0),
		Comment{"Pete Hunt", "This is one comment"},
		Comment{"Jordan Walke", "This is *another* comment"})
	
	flagHost := flag.String("host", "localhost:8080", "host:port")
	flag.Parse()

	http.Handle("/",
		http.FileServer(http.Dir(staticDir())))

	http.HandleFunc("/comments.json",
		func (w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				var comment Comment
				body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
				if err != nil {
					log.Println("Failed to read comment")
				} else if err := r.Body.Close(); err != nil {
					log.Println("Failed to close comment stream")
				} else if err := json.Unmarshal(body, &comment); err != nil {
					log.Println(string(body))
					log.Println("Failed to decode comment")
				} else {
					log.Println("Adding comment", comment)
					comments = append(comments, comment)
				}
			}

			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			if err := enc.Encode(comments); err != nil {
				log.Fatal("Can't encode comments", err)
			}
		})

	log.Println("Listening on", *flagHost)
	err := http.ListenAndServe(*flagHost, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err);
	}
}
