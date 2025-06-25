package pos_indonesia

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BTechnopark/ipostal/src/model"
	"github.com/PuerkitoBio/goquery"
)

func (a *iPostalApiImpl) SearchPostalCode(kodepos string) (model.ListPostalCode, error) {
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

	req, err := a.api.NewRequest(http.MethodPost, uri, nil, strings.NewReader(rawPayload), headers)
	if err != nil {
		return result, err
	}

	err = a.api.SendRequest(req, func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
			if err != nil {
				return
			}

			postalCode := model.PostalCode{}
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				content := s.Text()

				switch i {
				case 1:
					postalCode.PostalCode = content
				case 2:
					postalCode.Village = content
				case 3:
					postalCode.District = content
				case 4:
					postalCode.Region = content
				case 5:
					postalCode.Province = content
				}
			})

			result = append(result, &postalCode)
		})

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
