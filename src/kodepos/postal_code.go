package kodepos

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type KodePosData struct {
	Code       string `json:"code"`
	PostalCode string `json:"postal_code"`
	Province   string `json:"province"`
	Region     string `json:"region"`
	District   string `json:"district"`
	Village    string `json:"village"`
}

type PostalCodeQuery struct {
	Page int `json:"page" schema:"page"`
}

func (a *kodePosImpl) postalCode(provinceKey, regionKey string, query *PostalCodeQuery) ([]*KodePosData, error) {
	var err error
	result := []*KodePosData{}

	base := a.config.GetBaseUrl()
	endpoint := "kodepos"
	uri, err := url.JoinPath(base, endpoint, provinceKey, regionKey)
	if err != nil {
		return result, err
	}

	req, err := a.api.NewRequest(http.MethodGet, uri, query, nil, nil)
	if err != nil {
		return result, err
	}

	err = a.api.SendRequest(req, func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		doc.Find("tbody tr").Each(func(rowIndex int, s *goquery.Selection) {
			kodePos := KodePosData{}

			s.Find("td a").Each(func(columnIndex int, s *goquery.Selection) {
				content := s.Text()

				switch columnIndex {
				case 0:
					kodePos.Code = content
				case 1:
					kodePos.PostalCode = content
				case 2:
					kodePos.Village = content
				case 3:
					kodePos.District = content
				case 4:
					kodePos.Region = content
				case 5:
					kodePos.Province = content
				}
			})

			result = append(result, &kodePos)
		})

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, err
}

func (a *kodePosImpl) PostalCode(provinceKey, regionKey string) ([]*KodePosData, error) {
	var err error
	result := []*KodePosData{}

	if provinceKey == "" {
		return result, fmt.Errorf("province key is required")
	}

	base := a.config.GetBaseUrl()
	endpoint := "kodepos"
	uri, err := url.JoinPath(base, endpoint, provinceKey, regionKey)
	if err != nil {
		return result, err
	}

	req, err := a.api.NewRequest(http.MethodGet, uri, nil, nil, nil)
	if err != nil {
		return result, err
	}

	totalPage := 1
	err = a.api.SendRequest(req, func(resp *http.Response) error {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		doc.Find(".page-link").Each(func(i int, s *goquery.Selection) {
			if err != nil {
				return
			}

			rel := s.AttrOr("rel", "")
			if rel != "" {
				return
			}

			content := s.Text()
			page, err := strconv.Atoi(content)
			if err != nil {
				return
			}

			if totalPage < page {
				totalPage = page
			}
		})

		return nil
	})
	if err != nil {
		return result, err
	}

	query := &PostalCodeQuery{}
	for i := 1; i <= totalPage; i++ {
		query.Page = i
		data, err := a.postalCode(provinceKey, regionKey, query)
		if err != nil {
			return result, err
		}
		result = append(result, data...)
	}

	return result, err
}
