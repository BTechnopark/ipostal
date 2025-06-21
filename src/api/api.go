package api

import (
	"io"
	"net/http"
	"time"

	"github.com/BTechnopark/ipostal/src/session"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/schema"
)

func NewIPostalApi(config PostalConfig, session session.Session) IPostalApi {
	encoder := schema.NewEncoder()

	return &iPostalApiImpl{
		config:  config,
		encoder: encoder,
		session: session,
	}
}

type IPostalApi interface {
	FindPostalCode(kodepos string) (*goquery.Document, error)
}

type PostalConfig interface {
	GetBaseUrl() string
}

var clientApi *http.Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
	Timeout: 30 * time.Second,
}

type iPostalApiImpl struct {
	config  PostalConfig
	encoder *schema.Encoder
	session session.Session
}

var defaultHeader = map[string]string{
	"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	// "Accept-Encoding": "gzip, deflate, br, zstd",
	"Accept-Language": "en-US,en;q=0.9,id;q=0.8",
	"Cache-Control":   "no-cache",
	"Content-Type":    "application/x-www-form-urlencoded",
	// "Origin":                    "https://kodepos.posindonesia.co.id",
	"Pragma":   "no-cache",
	"Priority": "u=0, i",
	// "Referer":                   "https://kodepos.posindonesia.co.id/CariKodepos",
	"Sec-Ch-Ua":                 "\"Google Chrome\";v=\"137\", \"Chromium\";v=\"137\", \"Not/A)Brand\";v=\"24\"",
	"Sec-Ch-Ua-Mobile":          "?0",
	"Sec-Ch-Ua-Platform":        "Windows",
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "same-origin",
	"Sec-Fetch-User":            "?1",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
}

func (a *iPostalApiImpl) Request(method string, uri string, query any, payload io.Reader, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return nil, err
	}

	if query != nil {
		q := req.URL.Query()
		a.encoder.Encode(query, q)
		req.URL.RawQuery = q.Encode()
	}

	for key, value := range defaultHeader {
		req.Header.Set(key, value)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	a.session.AddCookiesToRequest(req)

	return req, err
}

func (a *iPostalApiImpl) SendRequest(req *http.Request, handler func(resp *http.Response) error) error {
	resp, err := clientApi.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if handler != nil {
		handler(resp)
	}

	err = a.session.Update(resp.Cookies())
	return err
}
