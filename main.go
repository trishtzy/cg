package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		r map[string]interface{}
	)

	// Connect to DB and migrate model
	dsn := "host=localhost user=countryadmin password=password dbname=countrydb port=5432 sslmode=disable TimeZone=Asia/Singapore"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("SETUP") == "true" {
		err = db_setup(db)
		if err != nil {
			log.Println(err)
		}
	}

	// cfg := elasticsearch.Config{
	// 	Addresses: []string{
	// 		"http://localhost:9200",
	// 	},
	// }
	es, _ := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

}
