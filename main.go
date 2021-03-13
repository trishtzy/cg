package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/trishtzy/cg/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Connect to DB and migrate model
	dsn := "host=localhost user=countryadmin password=password dbname=countrydb port=5432 sslmode=disable TimeZone=Asia/Singapore"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&model.Country{})

	// Seed data
	resp, err := http.Get("https://gist.githubusercontent.com/rusty-key/659db3f4566df459bd59c8a53dc9f71f/raw/4127f9550ef063121c564025f6d27dceeb279623/counties.json")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var countries []model.Country
	err = json.Unmarshal(body, &countries)
	if err != nil {
		log.Println(err)
	}

	result := db.Create(&countries)
	log.Println(result)
}
