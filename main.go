package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/trishtzy/cg/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		r  map[string]interface{}
		wg sync.WaitGroup
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

	var countries []model.Country
	result := db.Find(&countries)
	if result.Error != nil {
		log.Println(result.Error.Error())
	}

	for i, country := range countries {
		wg.Add(1)
		go func(i int, country model.Country) {
			defer wg.Done()

			b, err := json.Marshal(country)
			if err != nil {
				log.Fatalln(err)
			}
			req := esapi.IndexRequest{
				Index:        "testa",
				DocumentType: "search_as_you_type",
				DocumentID:   strconv.FormatUint(uint64(country.ID), 10),
				Body:         strings.NewReader(string(b)),
				Refresh:      "true",
			}

			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response while indexing: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), country.ID)
			} else {
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing indexRequest response body: %s", err)
				}
				log.Printf("%v", r)
			}
		}(i, country)
	}
	wg.Wait()
}
