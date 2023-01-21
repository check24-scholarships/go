package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	type PublicData struct {
		Time string
	}

	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		indexBytes, err := os.ReadFile("./public/index.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		publicData := PublicData{time.Now().Format("01-02-2006 15:04:05")}

		tmpl, err := template.New("test").Parse(string(indexBytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, publicData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		indexTemplated := tpl.String()

		_, err = io.WriteString(w, indexTemplated)
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
