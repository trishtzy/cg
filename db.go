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

var DB *gorm.DB

func db_conn() (*gorm.DB, error) {
	// Connect to DB and migrate model
	var err error
	dsn := "host=db user=countryadmin password=password dbname=countrydb port=5432 sslmode=disable TimeZone=Asia/Singapore"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return DB, nil
}

func db_close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	sqlDB.Close()
}

func db_setup() error {
	DB.AutoMigrate(&model.Country{})

	// Seed data
	resp, err := http.Get("https://gist.githubusercontent.com/rusty-key/659db3f4566df459bd59c8a53dc9f71f/raw/4127f9550ef063121c564025f6d27dceeb279623/counties.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var countries []model.Country
	err = json.Unmarshal(body, &countries)
	if err != nil {
		return err
	}

	result := DB.Create(&countries)
	if result.Error != nil {
		return result.Error
	}
	log.Println("Successfully setup database")
	return nil
}
