package service

import (
	"fmt"

	"github.com/ankit/project/url-shortner/url-shortner/db"
	"github.com/gin-gonic/gin"
)

var urlShortnerClient *URLShortnerService

type URLShortnerService struct {
	SqlDb db.URLService
}

func NewURLShortner(conn db.URLService) {
	urlShortnerClient = &URLShortnerService{
		SqlDb: conn,
	}
}

// This function register the particular adapter based on providerID
func CreateShorterURL() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		urlShortnerClient.createShortURL()
	}
}

func (url *URLShortnerService) createShortURL() {
	url.SqlDb.CreateShortURL()
}

// This function register the particular adapter based on providerID
func GetOriginalURL() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("Hello from Register Handler")
		urlShortnerClient.getShortURL()
	}
}

func (url *URLShortnerService) getShortURL() {

}
