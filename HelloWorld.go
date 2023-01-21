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

type DBContent struct {
	Name  string
	Price int
	Image string
}
type QueryResult struct {
	Results []DBContent
}

type TimeStruct struct {
	Time string
}

func main() {
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

		publicData := TimeStruct{time.Now().Format("01-02-2006 15:04:05")}

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

	search := func(query string) QueryResult {
		fmt.Println(query)
		return QueryResult{[]DBContent{
			{
				"hehe",
				69,
				"https://picsum.photos/200",
			},
		}}
	}

	searchHTMLHandler := func(w http.ResponseWriter, req *http.Request) {
		searchBytes, err := os.ReadFile("./public/search.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			panic(err)
		}

		q := req.URL.Query().Get("q")
		publicData := QueryResult{[]DBContent{}}

		if q != "" {
			publicData = search(q)
		}

		tmpl, err := template.New("test").Parse(string(searchBytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, publicData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}

		indexTemplated := tpl.String()

		_, err = io.WriteString(w, indexTemplated)
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHTMLHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
