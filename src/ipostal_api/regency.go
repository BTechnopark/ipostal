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

func NewRegencies(api api.IPostalApi, cache cache.Cache) ApiMeta {
	return &regencyImpl{
		api:   api,
		cache: cache,
	}
}

type regencyImpl struct {
	api   api.IPostalApi
	cache cache.Cache
}

type RegencyQuery struct {
	Province string `json:"province" form:"province" schema:"province"`
	Q        string `json:"q" form:"q" schema:"q"`
}

// Meta implements ApiMeta.
func (r *regencyImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &RegencyQuery{},
		Response:     &ResponseData[[]*model.Regency]{},
	}
}

// Handler implements ApiMeta.
func (r *regencyImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := RegencyQuery{}
		result := &ResponseData[[]*model.Regency]{}

		fullpath := ctx.Request.URL.String()
		hash := md5.Sum([]byte(fullpath))
		cacheKey := hex.EncodeToString(hash[:])

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(r.cache, cacheKey, result).
			Exec(func(seterr func(err error)) {
				var err error

				data, err := r.api.FindPostalCode("0")
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						err = errors.New("third party api timeout")
					}
					seterr(err)
					return
				}

				result.Data = data.Regencies(query.Province, query.Q)
			}).
			Finish(result)
	}
}
