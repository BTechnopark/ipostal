package ipostal_api

import (
	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/kodepos"
	"github.com/BTechnopark/ipostal/src/model"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func NewApi(api PostalCodeApi, cache cache.Cache) IPostalApi {
	return &apiImpl{
		api:   api,
		cache: cache,
	}
}

type IPostalApi interface {
	SearchPostalCode() ApiMeta
	Province() ApiMeta
	Region() ApiMeta
}

type PostalCodeApi interface {
	Province() ([]*kodepos.Province, error)
	Region(provinceKey string) ([]*kodepos.Region, error)
	PostalCode(provinceKey, regionKey string) ([]*kodepos.KodePosData, error)
	SearchPostalCode(q string, page, limit int) (model.ListPostalCode, bool, error)
}

type ApiMeta interface {
	Meta(uri string) *gin_api.ApiData
	Handler() gin.HandlerFunc
}

type apiImpl struct {
	api   PostalCodeApi
	cache cache.Cache
}

// Province implements IPostalApi.
func (a *apiImpl) Province() ApiMeta {
	return NewProvinceApi(a.api, a.cache)
}

// Region implements IPostalApi.
func (a *apiImpl) Region() ApiMeta {
	return NewRegionApi(a.api, a.cache)
}

// FindPostalCode implements Api.
func (a *apiImpl) SearchPostalCode() ApiMeta {
	return NewSearchPostalCode(a.api, a.cache)
}
