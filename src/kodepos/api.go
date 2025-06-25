package kodepos

import (
	"github.com/BTechnopark/ipostal/src/client"
	"github.com/BTechnopark/ipostal/src/model"
)

func NewKodePosApi(config KodePosConfig, api client.Api) KodePos {
	return &kodePosImpl{
		config: config,
		api:    api,
	}
}

type KodePos interface {
	Province() ([]*Province, error)
	Region(provinceKey string) ([]*Region, error)
	PostalCode(provinceKey, regionKey string) ([]*KodePosData, error)
	SearchPostalCode(q string) (model.ListPostalCode, error)
}

type KodePosConfig interface {
	GetBaseUrl() string
}

type kodePosImpl struct {
	config KodePosConfig
	api    client.Api
}
