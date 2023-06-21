package urlerror

type URLShortnerError struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}
