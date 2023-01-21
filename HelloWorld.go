package main

import (
	"bytes"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

	/*db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {

	}*/

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

	search := func(query string) {
		fmt.Println(query)
	}

	searchHTMLHandler := func(w http.ResponseWriter, req *http.Request) {
		indexBytes, err := os.ReadFile("./public/search.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		q := req.URL.Query().Get("q")
		if q != "" {
			search(q)
		}

		_, err = io.WriteString(w, string(indexBytes))
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHTMLHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
