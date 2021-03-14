package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/trishtzy/cg/model"
)

func es_setup() error {
	var (
		r  map[string]interface{}
		wg sync.WaitGroup
	)
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		// log.Fatalf("Error creating the client: %s", err)
		return err
	}
	res, err := es.Info()
	if err != nil {
		// log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		// log.Fatalf("Error: %s", res.String())
		return errors.New(res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		// log.Fatalf("Error parsing the response body: %s", err)
		return err
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	var countries []model.Country
	result := DB.Find(&countries)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return err
	}

	for i, country := range countries {
		wg.Add(1)
		go indexCountry(i, country, &wg, es)
	}
	wg.Wait()
	return nil
}

func indexCountry(i int, country model.Country, wg *sync.WaitGroup, es *elasticsearch.Client) {
	defer wg.Done()

	b, err := json.Marshal(country)
	if err != nil {
		log.Fatalln(err)
	}
	req := esapi.IndexRequest{
		Index:        "country",
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

	var r map[string]interface{}
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), country.ID)
	} else {
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing indexRequest response body: %s", err)
		}
		log.Printf("%v", r)
	}
}
