package kodepos

import (
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/BTechnopark/ipostal/src/model"
	"github.com/PuerkitoBio/goquery"
)

type searchKodePostQuery struct {
	Q     string `json:"q" schema:"q"`
	Page  int    `json:"page" schema:"page"`
	Limit int    `json:"limit" schema:"limit"`
}

func (a *kodePosImpl) searchPostalCode(query searchKodePostQuery) (model.ListPostalCode, error) {
	var err error
	result := model.ListPostalCode{}

	baseUrl := a.config.GetBaseUrl()
	endpoint := "/cari"
	uri, err := url.JoinPath(baseUrl, endpoint)
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

		doc.Find("table tbody tr").Each(func(rowIndex int, s *goquery.Selection) {
			postalCode := model.PostalCode{}

			s.Find("td").Each(func(colIndex int, s *goquery.Selection) {
				content := s.Text()

				switch colIndex {
				case 0:
					postalCode.PostalCode = content
				case 1:
					postalCode.Village = content
				case 2:
					postalCode.District = content
				case 3:
					postalCode.Region = content
				case 4:
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

	return result, err
}

func (a *kodePosImpl) SearchPostalCode(q string, page, limit int) (model.ListPostalCode, bool, error) {
	var err error
	result := model.ListPostalCode{}

	baseUrl := a.config.GetBaseUrl()
	endpoint := "/cari"
	uri, err := url.JoinPath(baseUrl, endpoint)
	if err != nil {
		return result, false, err
	}

	query := searchKodePostQuery{
		Q: q,
	}

	req, err := a.api.NewRequest(http.MethodGet, uri, query, nil, nil)
	if err != nil {
		return result, false, err
	}

	totalPage := 0
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

		return err
	})
	if err != nil {
		return result, false, err
	}

	baseLimit := 10
	groupLimit := limit / baseLimit
	groupPage := int(math.Ceil(float64(totalPage) / float64(groupLimit)))

	startPage := 1
	if page > 1 {
		startPage = (page * groupLimit) - 1
	}

	endPage := min((startPage+groupLimit)-1, groupPage)

	for i := startPage; i <= endPage; i++ {
		query.Page = i
		data, err := a.searchPostalCode(query)
		if err != nil {
			return result, false, err
		}

		result = append(result, data...)
	}

	hasMore := groupPage > page

	return result, hasMore, err
}
