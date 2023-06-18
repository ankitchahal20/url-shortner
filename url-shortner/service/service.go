package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/ankit/project/url-shortner/url-shortner/constants"
	"github.com/ankit/project/url-shortner/url-shortner/db"
	"github.com/ankit/project/url-shortner/url-shortner/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var counter int64 = 100000000000

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
	return func(context *gin.Context) {
		var urlInfo models.URLInfo
		if err := context.ShouldBindBodyWith(&urlInfo, binding.JSON); err == nil {
			fmt.Println("URL Info ", urlInfo.OriginalURL)
			// Validate request body
			if urlInfo.OriginalURL == "" {
				err := errors.New("invalid request received")
				context.JSON(http.StatusBadRequest, gin.H{"Received url is empty": err.Error()})
				return
			}
			err := urlShortnerClient.createShortURL(urlInfo)
			if err != nil {
				context.Writer.WriteHeader(http.StatusInternalServerError)
			} else {
				context.Writer.WriteHeader(http.StatusOK)
			}
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}

	}
}

func (service *URLShortnerService) createShortURL(urlInfo models.URLInfo) error {

	fmt.Println("Request Reached till service layer")

	// logic to genrate a short url for a given large url
	randomNum := service.rangeIn(counter, 999999999999)

	/*
		MD5 => Message Digest
		input=> any string
		output=128bit string => 16B => 32charters

		MD1, MD2, MD3,.....MD100 => first 7 characters can be same.

		SHA-1
		input=> any string
		output=160bit string => 20B => 40charters
		SHA1, SHA2, SHA2,.....SHA3 => first 7 characters can be same.

		base62
		a-z 26
		A-Z 26
		0-9 10
		---------
			62 characters.

		62|10009nkdaanfksdu73y8y399393
		  | quetiont                | rem
		   quetiont2| rem2
		   quetiont3| rem3
	*/

	shortUrl := fmt.Sprintf("%s/%s", constants.Domain, service.base62Encode(randomNum))
	if shortUrl != "" {
		urlInfo.ShortURL = shortUrl
	} else {
		return errors.New("genrate short url is empty")
	}
	fmt.Println("url info : ", urlInfo)
	counter += 1
	err := service.SqlDb.CreateShortURL(urlInfo)
	return err
}

func (service URLShortnerService) base62Encode(randomNum int64) string {
	hashStr := ""

	for randomNum > 0 {
		hashStr = fmt.Sprintf("%s%c", hashStr, constants.Base62[randomNum%62])
		randomNum = randomNum / 62
	}

	return hashStr
}

func (service URLShortnerService) rangeIn(low, hi int64) int64 {
	return low + rand.Int63n(hi-low)
}

// This function register the particular adapter based on providerID
func GetOriginalURL() func(ctx *gin.Context) {
	return func(context *gin.Context) {
		var urlInfo models.URLInfo
		if err := context.ShouldBindBodyWith(&urlInfo, binding.JSON); err == nil {
			fmt.Println("URL Info ", urlInfo.ShortURL)
			// Validate request body
			if urlInfo.ShortURL == "" {
				err := errors.New("invalid request received")
				context.JSON(http.StatusBadRequest, gin.H{"Received url is empty": err.Error()})
				return
			}
			originalURL, err := urlShortnerClient.getShortURL(urlInfo)
			if err != nil {
				context.Writer.WriteHeader(http.StatusInternalServerError)
			} else {
				http.Redirect(context.Writer, context.Request, originalURL, http.StatusMovedPermanently)
				context.Writer.WriteHeader(http.StatusOK)
			}
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (url *URLShortnerService) getShortURL(urlInfo models.URLInfo) (string, error) {
	originalURL, err := url.SqlDb.GetOriginalURL(urlInfo)
	if err != nil {
		return "", err
	}
	fmt.Println("OriginalURL : ", originalURL)
	return originalURL, nil
}
