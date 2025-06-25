package ipostal_api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/api_context"
	"github.com/BTechnopark/ipostal/src/model"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewSearchPostalCode(api PostalCodeApi, cache cache.Cache) ApiMeta {
	return &searchPostalCodeImpl{
		api:   api,
		cache: cache,
	}
}

type searchPostalCodeImpl struct {
	api   PostalCodeApi
	cache cache.Cache
}

type FindPostalCodeQuery struct {
	Q string `json:"q" form:"q" schema:"q"`
}

// Meta implements ApiMeta.
func (f *searchPostalCodeImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &FindPostalCodeQuery{},
		Response:     &ResponseData[[]*model.PostalCode]{},
	}
}

// Handler implements ApiMeta.
func (f *searchPostalCodeImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := FindPostalCodeQuery{}
		result := &ResponseData[model.ListPostalCode]{}

		fullpath := ctx.Request.URL.String()
		hash := md5.Sum([]byte(fullpath))
		cacheKey := hex.EncodeToString(hash[:])

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(f.cache, cacheKey, result).
			Exec(func(seterr func(err error)) {
				var err error

				if query.Q == "" {
					query.Q = "0"
				}

				data, err := f.api.SearchPostalCode(query.Q)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						err = errors.New("third party api timeout")
					}
					seterr(err)
					return
				}

				result.Data = data
			}).
			Finish(result)
	}
}
