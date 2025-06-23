package api

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BTechnopark/ipostal/src/model"
	"github.com/BTechnopark/ipostal/src/parser"
	"github.com/PuerkitoBio/goquery"
)

func (a *iPostalApiImpl) FindPostalCode(kodepos string) (model.ListPostalCode, error) {
	var err error
	result := model.ListPostalCode{}

	err = a.cache.Get(kodepos, &result)
	if err == nil {
		return result, nil
	}

	base := a.config.GetBaseUrl()
	endpoint := "/CariKodepos"
	uri, err := url.JoinPath(base, endpoint)
	if err != nil {
		return result, err
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
		return result, err
	}

	err = a.SendRequest(req, func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}
		p := parser.NewIPostalParser(doc)
		data, err := p.PostalCode()
		if err != nil {
			return err
		}

		result = data

		return nil
	})
	if err != nil {
		return result, err
	}

	err = a.cache.Set(kodepos, result, time.Minute)
	if err != nil {
		return result, err
	}

	return result, err
}
