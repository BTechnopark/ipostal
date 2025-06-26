package ipostal_api

import (
	"context"
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

type SearchPostalCodeQuery struct {
	Q     string `json:"q" form:"q" schema:"q"`
	Page  int    `json:"page" form:"page" schema:"page"`
	Limit int    `json:"limit" form:"limit" schema:"limit"`
}

// Meta implements ApiMeta.
func (f *searchPostalCodeImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &SearchPostalCodeQuery{},
		Response:     &ResponseData[[]*model.PostalCode]{},
	}
}

// Handler implements ApiMeta.
func (f *searchPostalCodeImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := SearchPostalCodeQuery{
			Page:  1,
			Limit: 20,
		}
		result := &ResponseData[model.ListPostalCode]{
			Data:     model.ListPostalCode{},
			PageInfo: &PageInfo{},
		}

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Cache(f.cache, "", result).
			Exec(func(seterr func(err error)) {
				if query.Limit%10 != 0 {
					seterr(errors.New("limit must be multiple of 10"))
					return
				}
			}).
			Exec(func(seterr func(err error)) {
				var err error

				data, hasMore, err := f.api.SearchPostalCode(query.Q, query.Page, query.Limit)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						err = errors.New("third party api timeout")
					}
					seterr(err)
					return
				}

				result.Data = data
				result.PageInfo = &PageInfo{
					CurrentPage: query.Page,
					HasMore:     hasMore,
				}

			}).
			Finish(result)
	}
}
