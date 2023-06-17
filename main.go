package main

import (
	"log"

	"github.com/ankit/project/url-shortner/url-shortner/db"
	"github.com/ankit/project/url-shortner/url-shortner/server"
	"github.com/ankit/project/url-shortner/url-shortner/service"
)

func main() {

	postgres, err := db.New()
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}
	service.NewURLShortner(postgres)
	server.Start()
}
