package pos_indonesia

import (
	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/client"
	"github.com/BTechnopark/ipostal/src/kodepos"
	"github.com/BTechnopark/ipostal/src/model"
)

func NewPosIndonesiaApi(config PostalConfig, api client.Api, cache cache.Cache) PosIndonesiaApi {
	return &iPostalApiImpl{
		config: config,
		api:    api,
		cache:  cache,
	}
}

type PosIndonesiaApi interface {
	Province() ([]*kodepos.Province, error)
	Region(provinceKey string) ([]*kodepos.Region, error)
	PostalCode(provinceKey, regionKey string) ([]*kodepos.KodePosData, error)
	SearchPostalCode(q string) (model.ListPostalCode, error)
}

type PostalConfig interface {
	GetBaseUrl() string
}

type iPostalApiImpl struct {
	config PostalConfig
	api    client.Api
	cache  cache.Cache
}

// PostalCode implements ipostal_api.PostalCodeApi.
func (p *iPostalApiImpl) PostalCode(provinceKey string, regionKey string) ([]*kodepos.KodePosData, error) {
	panic("unimplemented")
}

// Province implements ipostal_api.PostalCodeApi.
func (p *iPostalApiImpl) Province() ([]*kodepos.Province, error) {
	panic("unimplemented")
}

// Region implements ipostal_api.PostalCodeApi.
func (p *iPostalApiImpl) Region(provinceKey string) ([]*kodepos.Region, error) {
	panic("unimplemented")
}
