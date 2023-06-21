package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/ankit/project/url-shortner/url-shortner/constants"
	"github.com/ankit/project/url-shortner/url-shortner/db"
	"github.com/ankit/project/url-shortner/url-shortner/models"
	"github.com/ankit/project/url-shortner/url-shortner/urlerror"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var counter int64 = 100000000000

var urlShortnerClient *URLShortnerService

type URLShortnerService struct {
	repo db.URLService
}

func NewURLShortner(conn db.URLService) {
	urlShortnerClient = &URLShortnerService{
		repo: conn,
	}
}

// @Description Fetches a short URL for a given long URL
// @Summary Fetches a short URL for a given long URL
// @Accept       json
// @Param RequestFields body models.URLInfo true "Request Fields"
// @Failure					500				{object}	error.URLShortnerError	"Internal Server Error"
// @Failure					403				{object}	error.URLShortnerError	"Forbidden"
// @Failure					401				{object}	error.URLShortnerError	"Unauthorized"
// @Failure					400				{object}	error.URLShortnerError	"Bad Request"
// @Success					200				{object}	error.URLShortnerError	"OK"
// @Router /v1/urlshortner/create [POST]
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
			shortURL, err := urlShortnerClient.createShortURL(urlInfo)
			if err != nil {
				context.Writer.WriteHeader(http.StatusInternalServerError)
			} else {
				context.JSON(http.StatusCreated, shortURL)
				//context.Writer.WriteHeader(http.StatusOK)
			}
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}

	}
}

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
	  	|______________________________
		| quotient1	| rem1
		| quotient2	| rem2
		| quotient3	| rem3
*/
func (service *URLShortnerService) createShortURL(urlInfo models.URLInfo) (string, *urlerror.URLShortnerError) {
	fmt.Println("Request Reached till service layer")

	// logic to genrate a short url for a given large url
	randomNum := service.rangeIn(counter, 999999999999)

	shortUrl := fmt.Sprintf("%s/%s", constants.Domain, service.base62Encode(randomNum))
	if shortUrl != "" {
		urlInfo.ShortURL = shortUrl
	} else {
		return "", &urlerror.URLShortnerError{
			Status:  "Service Unavailable",
			Code:    http.StatusInternalServerError,
			Message: "unable to generate short url for the given url",
		}
	}
	fmt.Println("url info : ", urlInfo)
	counter += 1
	err := service.repo.CreateShortURL(urlInfo)
	if err != nil {
		if errors.Is(err, db.ErrUnableToInsertARow) {
			return "", &urlerror.URLShortnerError{
				Status:  "Service Unavailable",
				Code:    http.StatusInternalServerError,
				Message: "unable to generate short url for the given url",
			}
		}

		//return "", &urlerror.URLShortnerError{}

	}
	return shortUrl, nil
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

// @Description Get orginal URL for a short URL
// @Summary Get orginal URL for a short URL
// @Accept       json
// @Param RequestFields body models.URLInfo true "Request Fields"
// @Failure					500				{object}	error.URLShortnerError	"Internal Server Error"
// @Failure					403				{object}	error.URLShortnerError	"Forbidden"
// @Failure					401				{object}	error.URLShortnerError	"Unauthorized"
// @Failure					400				{object}	error.URLShortnerError	"Bad Request"
// @Success					200				{object}	error.URLShortnerError	"OK"
// @Router /v1/urlshortner [GET]
// This function register the particular adapter based on providerID
func GetOriginalURL() func(ctx *gin.Context) {
	return func(context *gin.Context) {

		shortURL := context.Param("url")
		fmt.Println("URL Info ", shortURL)
		// Validate request body
		if shortURL == "" {
			err := errors.New("invalid request received")
			context.JSON(http.StatusBadRequest, gin.H{"Received url is empty": err.Error()})
			return
		}
		originalURL, err := urlShortnerClient.getShortURL(shortURL)
		fmt.Println("originalURL : ", originalURL)
		if err != nil {
			context.Writer.WriteHeader(http.StatusInternalServerError)
		} else {
			context.Redirect(http.StatusMovedPermanently, originalURL)
			context.Writer.WriteHeader(http.StatusOK)
		}

	}
}

func (url *URLShortnerService) getShortURL(shortURL string) (string, *urlerror.URLShortnerError) {
	shortUrl := fmt.Sprintf("%s/%s", constants.Domain, shortURL)
	fmt.Println("shortUrl : ", shortUrl)
	originalURL, err := url.repo.GetOriginalURL(shortUrl)
	if err != nil {
		if errors.Is(err, db.ErrNoRowFound) {

		} else if errors.Is(err, db.ErrScanningRows) {

		} else if errors.Is(err, db.ErrUnableToSelectRows) {

		}
	}

	return originalURL, nil
}
