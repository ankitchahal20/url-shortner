package models

type URLInfo struct {
	OriginalURL string `json:"originalurl"`
	ShortURL    string `json:"shorturl,omitempty"`
}
