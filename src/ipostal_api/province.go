package ipostal_api

import (
	"math"
	"net/http"
	"strings"

	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/api_context"
	"github.com/BTechnopark/ipostal/src/kodepos"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewProvinceApi(api PostalCodeApi, cache cache.Cache) ApiMeta {
	return &provinceImpl{
		api:   api,
		cache: cache,
	}
}

type provinceImpl struct {
	api   PostalCodeApi
	cache cache.Cache
}

type ProvinceQuery struct {
	Q     string `json:"q" form:"q" schema:"q"`
	Page  int    `json:"page" form:"page" schema:"page"`
	Limit int    `json:"limit" form:"limit" schema:"limit"`
}

// Meta implements ApiMeta.
func (p *provinceImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &ProvinceQuery{},
		Response:     &ResponseData[[]*kodepos.Province]{},
	}
}

// Handler implements ApiMeta.
func (p *provinceImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := ProvinceQuery{
			Page:  1,
			Limit: 20,
		}
		result := &ResponseData[[]*kodepos.Province]{
			Data:     []*kodepos.Province{},
			PageInfo: &PageInfo{},
		}

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(p.cache, "", result).
			Exec(func(seterr func(err error)) {
				data, err := p.api.Province()
				if err != nil {
					seterr(err)
					return
				}

				if query.Q != "" {
					newData := []*kodepos.Province{}

					key := strings.ToLower(query.Q)
					for _, item := range data {
						provinceName := strings.ToLower(item.Name)

						if strings.Contains(provinceName, key) {
							newData = append(newData, item)
						}
					}

					data = newData
				}

				result.PageInfo.TotalItems = len(data)
				result.PageInfo.CurrentPage = query.Page
				result.PageInfo.TotalPages = int(math.Ceil(float64(len(data)) / float64(query.Limit)))

				start := (query.Page - 1) * query.Limit
				end := query.Page * query.Limit
				if end > len(data) {
					end = len(data)
				}
				if start > end {
					return
				}

				result.Data = data[start:end]
			}).
			Finish(result)
	}
}
