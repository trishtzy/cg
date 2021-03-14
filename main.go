package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, err := db_conn()
	if err != nil {
		log.Println(err)
	}

	if os.Getenv("SETUP") == "true" {
		err = db_setup()
		if err != nil {
			log.Println(err)
		}
		err = es_setup()
		if err != nil {
			log.Fatalln(err)
		}
	}
	s := &server{}
	http.Handle("/", s)
	http.HandleFunc("/countries", countries)
	log.Println("Listening at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
