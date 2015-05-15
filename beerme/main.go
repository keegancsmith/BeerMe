package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

type Suip struct {
	Drink string `json:"drink"`
	Team string `json:"team"`
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
	flagHost := flag.String("host", "localhost:8080", "host:port")
	flag.Parse()

	http.Handle("/",
		http.FileServer(http.Dir(staticDir())))

	http.HandleFunc("/suip",
		func (w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				var suip Suip
				body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
				if err != nil {
					log.Println("Failed to read request")
				} else if err := r.Body.Close(); err != nil {
					log.Println("Failed to close request")
				} else if err := json.Unmarshal(body, &suip); err != nil {
					log.Println(string(body))
					log.Println("Failed to decode drink")
				} else {
					log.Println("Another happy customer", suip)
				}
			}
			fmt.Fprintf(w, "\"Enjoy your beverage\"")
		})

	log.Println("Listening on", *flagHost)
	err := http.ListenAndServe(*flagHost, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err);
	}
}
