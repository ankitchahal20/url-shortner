package main

import (
	"log"

	"github.com/ankit/project/url-shortner/docs"
	"github.com/ankit/project/url-shortner/url-shortner/db"
	"github.com/ankit/project/url-shortner/url-shortner/server"
	"github.com/ankit/project/url-shortner/url-shortner/service"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

// @title URL Shortner
// @version 1.0
// @description This is a URL Shortner service. For a given long URL, it gives you a short URL. You can visit the GitHub repository at https://github.com/ankitchahal20/url-shortner

// @contact.name Ankit Chahal
// @contact.url none
// @contact.email https://github.com/ankitchahal20/url-shortner

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1/urlshortner
// @query.collection.format multi
// @schemes        http
func main() {
	logger, _ := zap.NewDevelopment()
	logger.Info("Main started")

	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Host = "0.0.0.0:8080"
	postgres, err := db.New()
	if err != nil {
		log.Fatal("Unable to connect to DB : ", err)
	}
	service.NewURLShortner(postgres)
	server.Start()
}
