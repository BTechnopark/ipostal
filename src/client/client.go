package client

import (
	"net/http"
	"time"
)

var ClientApi *http.Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
	Timeout: 30 * time.Second,
}
