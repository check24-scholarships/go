package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
)

type Product struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Image string `json:"image"`
}

func main() {
	products, err := parseJSON()
	if err != nil {
		log.Fatal(err)
	}

	maxNameLength := 0
	maxImageLength := 0
	for i := 0; i < len(products); i++ {
		if maxNameLength < len(products[i].Name) {
			maxNameLength = len(products[i].Name)
		}

		if maxImageLength < len(products[i].Image) {
			maxImageLength = len(products[i].Image)
		}
	}

	fmt.Println("Maximum Name Length: ", maxNameLength)
	fmt.Println("Maximum Image Length: ", maxImageLength)

	db, err := sql.Open("mysql",
		"go:password@tcp(172.30.124.56:3306)/products")
	if err != nil {
		log.Fatal("Error initialising DB connection", err)
	}

	stmt, err := db.Prepare("INSERT INTO products(name, price, image) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing SQL insert statement", err)
	}

	for i := 0; i < len(products); i++ {
		_, err := stmt.Exec(products[i].Name, products[i].Price, products[i].Image)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close()
}

func parseJSON() ([]Product, error) {
	jsonFile, err := os.Open("dump.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var products []Product
	if err := json.Unmarshal(byteValue, &products); err != nil {
		return nil, err
	}

	return products, nil
}
