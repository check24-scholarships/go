package main

import (
	"bytes"
	"database/sql"
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

	db, err := sql.Open("mysql",
		"go:password@tcp(172.30.124.56:3306)/products")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}

	search := func(query string, order string) QueryResult {
		rows, err := db.Query("SELECT * FROM products WHERE name LIKE '%" + query + "%' ORDER BY price " + order)
		if err != nil {
			return QueryResult{[]DBContent{}}
		}
		defer rows.Close()

		var (
			args  []DBContent
			name  string
			price int
			image string
		)
		for rows.Next() {
			err = rows.Scan(&name, &price, &image)
			if err != nil {
				panic(err)
			}
			args = append(args, DBContent{name, price, image})
		}

		return QueryResult{args}
	}

	searchHandler := func(w http.ResponseWriter, req *http.Request) {
		searchBytes, err := os.ReadFile("./public/search.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		q := req.URL.Query().Get("q")
		o := req.URL.Query().Get("o")
		publicData := QueryResult{[]DBContent{}}

		if q != "" && o != "" {
			publicData = search(q, o)
		}

		tmpl, err := template.New("test").Parse(string(searchBytes))
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
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
