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

func NewRegionApi(api PostalCodeApi, cache cache.Cache) ApiMeta {
	return &regionImpl{
		api:   api,
		cache: cache,
	}
}

type regionImpl struct {
	api   PostalCodeApi
	cache cache.Cache
}

type RegionQuery struct {
	ProvinceKey string `json:"province_key" form:"province_key" schema:"province_key" binding:"required"`
	Q           string `json:"q" form:"q" schema:"q"`
	Page        int    `json:"page" form:"page" schema:"page"`
	Limit       int    `json:"limit" form:"limit" schema:"limit"`
}

// Meta implements ApiMeta.
func (p *regionImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &RegionQuery{},
		Response:     &ResponseData[[]*kodepos.Region]{},
	}
}

// Handler implements ApiMeta.
func (p *regionImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := RegionQuery{}
		result := &ResponseData[[]*kodepos.Region]{}

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(p.cache, query.ProvinceKey, result).
			Exec(func(seterr func(err error)) {
				data, err := p.api.Region(query.ProvinceKey)
				if err != nil {
					seterr(err)
					return
				}

				if query.Q != "" {
					newData := []*kodepos.Region{}

					key := strings.ToLower(query.Q)
					for _, item := range data {
						regionName := strings.ToLower(item.Name)
						if strings.Contains(regionName, key) {
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
