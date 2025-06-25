package ipostal_api

import (
	"net/http"

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

				result.Data = data
			}).
			Finish(result)
	}
}
