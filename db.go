package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/trishtzy/cg/model"
	"gorm.io/gorm"
)

func db_setup(db *gorm.DB) error {
	db.AutoMigrate(&model.Country{})

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

	result := db.Create(&countries)
	if result.Error != nil {
		return result.Error
	}
	log.Println("Successfully setup database")
	return nil
}
