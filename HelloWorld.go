package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	type PublicData struct {
		Material string
		Count    uint
	}
	indexBuffer, err := os.ReadFile("./public/index.html")
	indexString := string(indexBuffer)
	publicData := PublicData{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, publicData)
	if err != nil {
		panic(err)
	}

	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, indexString)
	}

	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
