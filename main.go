package main

import (
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=countryadmin password=password dbname=countrydb port=5432 sslmode=disable TimeZone=Asia/Singapore"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	resp, err := http.Get("https://gist.githubusercontent.com/rusty-key/659db3f4566df459bd59c8a53dc9f71f/raw/4127f9550ef063121c564025f6d27dceeb279623/counties.json")
	countries := make([]Country, 100)
	if err != nil {
		println(err)
	}
}
