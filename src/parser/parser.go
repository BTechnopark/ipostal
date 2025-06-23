package parser

import (
	"strconv"

	"github.com/BTechnopark/ipostal/src/model"
	"github.com/PuerkitoBio/goquery"
)

func NewIPostalParser(doc *goquery.Document) IPostalParser {
	return &iPostalParserImpl{
		doc: doc,
	}
}

type IPostalParser interface {
	Headers() ([]string, error)
	PostalCode() (model.ListPostalCode, error)
}

type iPostalParserImpl struct {
	doc *goquery.Document
}

// PostalCode implements IPostalParser.
func (p *iPostalParserImpl) PostalCode() (model.ListPostalCode, error) {
	var err error
	results := model.ListPostalCode{}

	rows := p.GetBodyRows()
	rows.Each(func(i int, s *goquery.Selection) {
		if err != nil {
			return
		}

		postalCode := model.PostalCode{}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			content := s.Text()

			switch i {
			case 0:
				num, err := strconv.Atoi(content)
				if err != nil {
					return
				}
				postalCode.ID = num
			case 1:
				postalCode.PostalCode = content
			case 2:
				postalCode.Village = content
			case 3:
				postalCode.District = content
			case 4:
				postalCode.Regency = content
			case 5:
				postalCode.Province = content
			}
		})

		results = append(results, &postalCode)
	})

	return results, err
}

// Headers implements IPostalParser.
func (p *iPostalParserImpl) Headers() ([]string, error) {
	var err error
	result := []string{}

	headers := p.GetHeader()
	headers.Find("th").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		result = append(result, content)
	})

	return result, err
}

func (p *iPostalParserImpl) GetTable() *goquery.Selection {
	return p.doc.Find("table")
}

func (p *iPostalParserImpl) GetHeader() *goquery.Selection {
	return p.GetTable().Find("thead")
}

func (p *iPostalParserImpl) GetBody() *goquery.Selection {
	return p.GetTable().Find("tbody")
}

func (p *iPostalParserImpl) GetBodyRows() *goquery.Selection {
	return p.GetBody().Find("tr")
}
