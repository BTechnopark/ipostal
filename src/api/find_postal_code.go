package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (a *iPostalApiImpl) FindPostalCode(kodepos string) (*goquery.Document, error) {
	var err error

	base := a.config.GetBaseUrl()
	endpoint := "/CariKodepos"
	uri, err := url.JoinPath(base, endpoint)
	if err != nil {
		return nil, err
	}

	if kodepos == "" {
		kodepos = "0"
	}

	payload := url.Values{
		"kodepos": {kodepos},
	}
	rawPayload := payload.Encode()

	headers := map[string]string{
		"Origin":  base,
		"Referer": uri,
	}

	req, err := a.Request(http.MethodPost, uri, nil, strings.NewReader(rawPayload), headers)
	if err != nil {
		return nil, err
	}

	var doc *goquery.Document
	err = a.SendRequest(req, func(resp *http.Response) error {
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		return nil
	})

	return doc, err
}
