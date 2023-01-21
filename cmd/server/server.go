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
type LoggingContent struct {
	Url    string
	Amount int
}

func main() {
	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		logRequest(req)

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
		logRequest(req)

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

	statsHandler := func(w http.ResponseWriter, req *http.Request) {
		logRequest(req)

		db, err := sql.Open("mysql",
			"go:password@tcp(172.30.124.56:3306)/logging")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if err = db.Ping(); err != nil {
			panic(err)
		}

		statsBytes, err := os.ReadFile("./public/stats.html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		tmpl, err := template.New("test").Parse(string(statsBytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		rows, err := db.Query("SELECT URL, COUNT(URL) FROM logging GROUP BY URL;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var (
			args   []LoggingContent
			url    string
			amount int
		)
		for rows.Next() {
			err = rows.Scan(&url, &amount)
			if err != nil {
				panic(err)
			}
			args = append(args, LoggingContent{url, amount})
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, args); err != nil {
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
	http.HandleFunc("/stats", statsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logRequest(req *http.Request) {
	db, err := sql.Open("mysql",
		"go:password@tcp(172.30.124.56:3306)/logging")
	if err != nil {
		log.Fatal("Error initialising DB connection", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO logging(url, useragent, time) VALUES(?, ?, NOW())")
	if err != nil {
		log.Fatal("Error preparing SQL insert statement", err)
	}

	_, err = stmt.Exec(req.URL.String(), req.UserAgent())
	if err != nil {
		log.Fatal(err)
	}
}
