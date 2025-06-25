package kodepos

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Province struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Url  string `json:"url"`
}

func (a *kodePosImpl) Province() ([]*Province, error) {
	var err error
	result := []*Province{}

	baseUrl := a.config.GetBaseUrl()
	endpoint := "/kodepos"
	uri, err := url.JoinPath(baseUrl, endpoint)
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

		doc.Find(".row a").Each(func(i int, s *goquery.Selection) {
			content := s.Text()

			name := strings.Trim(strings.ReplaceAll(content, "\n", ""), " ")
			key := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			uri := s.AttrOr("href", "")

			province := Province{
				Name: name,
				Key:  key,
				Url:  uri,
			}

			result = append(result, &province)
		})

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, err
}
