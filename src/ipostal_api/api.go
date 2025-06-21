package ipostal_api

import (
	"github.com/BTechnopark/ipostal/src/api"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewApi(api api.IPostalApi) Api {
	return &apiImpl{
		api: api,
	}
}

type Api interface {
	FindPostalCode() ApiMeta
}

type ApiMeta interface {
	Handler() gin.HandlerFunc
	Meta(uri string) *gin_api.ApiData
}

type apiImpl struct {
	api api.IPostalApi
}

// FindPostalCode implements Api.
func (a *apiImpl) FindPostalCode() ApiMeta {
	return NewFindPostalCode(a.api)
}
