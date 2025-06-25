package kodepos

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Region struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Url  string `json:"url"`
}

func (a *kodePosImpl) Region(provinceKey string) ([]*Region, error) {
	var err error
	result := []*Region{}

	if provinceKey == "" {
		return result, fmt.Errorf("province key is required")
	}

	baseUrl := a.config.GetBaseUrl()
	endpoint := "kodepos"
	uri, err := url.JoinPath(baseUrl, endpoint, provinceKey)
	if err != nil {
		return result, err
	}

	req, err := a.api.NewRequest(http.MethodGet, uri, nil, nil, nil)
	if err != nil {
		return result, err
	}

	err = a.api.SendRequest(req, func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		// log.Println(doc.Html())

		doc.Find("div.col-12.col-md-4 a").Each(func(i int, s *goquery.Selection) {
			content := s.Text()

			name := strings.Trim(strings.ReplaceAll(content, "\n", ""), " ")
			key := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			uri := s.AttrOr("href", "")

			region := Region{
				Name: name,
				Key:  key,
				Url:  uri,
			}

			result = append(result, &region)
		})

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
