package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ankit/project/url-shortner/url-shortner/models"
)

type postgres struct{ db *sql.DB }

type URLService interface {
	CreateShortURL(models.URLInfo) error
	GetOriginalURL(string) (string, error)
}

func (p postgres) CreateShortURL(urlInfo models.URLInfo) error {

	query := `insert into url(originalurl, shorturl) values($1,$2)`
	fmt.Println("Query : ", query)
	_, err := p.db.Exec(query, urlInfo.OriginalURL, urlInfo.ShortURL)
	if err != nil {
		log.Println("unable to insert url info in table : ", err)
		return err
	}
	return err

}

func (p postgres) GetOriginalURL(shortURL string) (string, error) {
	query := `select originalurl from url where shorturl=$1`
	rows, err := p.db.Query(query, shortURL)
	if err != nil {
		log.Println("unable to perform select opertion on the url table : ", err)
		return "", err
	}
	fmt.Println("ROWS : ", rows)
	var url string
	for rows.Next() {
		err := rows.Scan(&url)
		fmt.Println("URL : ", url)
		if err != nil {
			log.Printf("error scanning row: %v ", err)
			return "", err
		}
		if url != "" {
			return url, nil
		}
	}
	return url, errors.New("no row found in DB for the given short url")
}
