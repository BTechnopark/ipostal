package ipostal_api

import (
	"net/http"

	"github.com/BTechnopark/ipostal/src/api"
	"github.com/BTechnopark/ipostal/src/api_context"
	"github.com/BTechnopark/ipostal/src/model"
	"github.com/BTechnopark/ipostal/src/parser"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewFindPostalCode(api api.IPostalApi) ApiMeta {
	return &findPostalCodeImpl{
		api: api,
	}
}

type findPostalCodeImpl struct {
	api api.IPostalApi
}

type FindPostalCodeQuery struct {
	Q string `json:"q" form:"q" schema:"q"`
}

// Meta implements ApiMeta.
func (f *findPostalCodeImpl) Meta(uri string) *gin_api.ApiData {
	return &gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: uri,
		Query:        &FindPostalCodeQuery{},
		Response:     &ResponseData[[]*model.PostalCode]{},
	}
}

// Handler implements ApiMeta.
func (f *findPostalCodeImpl) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := FindPostalCodeQuery{}
		result := &ResponseData[[]*model.PostalCode]{}

		apiCtx := api_context.NewApiContext(ctx)
		apiCtx.
			BindQuery(&query).
			Exec(func(seterr func(err error)) {
				var err error

				if query.Q == "" {
					query.Q = "0"
				}

				doc, err := f.api.FindPostalCode(query.Q)
				if err != nil {
					seterr(err)
					return
				}

				p := parser.NewIPostalParser(doc)

				data, err := p.PostalCode()
				if err != nil {
					seterr(err)
					return
				}

				result.Data = data
			}).
			Finish(result)
	}
}
