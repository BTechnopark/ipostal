package ipostal_api

import (
	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/api"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewApi(api api.IPostalApi, cache cache.Cache) Api {
	return &apiImpl{
		api:   api,
		cache: cache,
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
	api   api.IPostalApi
	cache cache.Cache
}

// FindPostalCode implements Api.
func (a *apiImpl) FindPostalCode() ApiMeta {
	return NewFindPostalCode(a.api, a.cache)
}
