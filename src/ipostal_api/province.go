package ipostal_api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/api"
	"github.com/BTechnopark/ipostal/src/api_context"
	"github.com/BTechnopark/ipostal/src/model"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewProvinces(api api.IPostalApi, cache cache.Cache) ApiMeta {
	return &provinceImpl{
		api:   api,
		cache: cache,
	}
}

type provinceImpl struct {
	api   api.IPostalApi
	cache cache.Cache
}

type ProvinceQuery struct {
	Q string `json:"q" form:"q" schema:"q"`
}

// Meta implements ApiMeta.
func (p *provinceImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &ProvinceQuery{},
		Response:     &ResponseData[[]*model.Province]{},
	}
}

// Handler implements ApiMeta.
func (p *provinceImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := ProvinceQuery{}
		result := &ResponseData[[]*model.Province]{}

		fullpath := ctx.Request.URL.String()
		hash := md5.Sum([]byte(fullpath))
		cacheKey := hex.EncodeToString(hash[:])

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(p.cache, cacheKey, result).
			Exec(func(seterr func(err error)) {
				var err error

				data, err := p.api.FindPostalCode("0")
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						err = errors.New("third party api timeout")
					}
					seterr(err)
					return
				}

				result.Data = data.Provinces(query.Q)
			}).
			Finish(result)
	}
}
