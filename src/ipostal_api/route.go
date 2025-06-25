package ipostal_api

import "github.com/muchrief/gin_api"

func RegisterIPostalApi(group gin_api.SdkGroup, api IPostalApi) {
	searchPostalCode := api.SearchPostalCode()
	group.Register(searchPostalCode.Meta("search"), searchPostalCode.Handler())

	province := api.Province()
	group.Register(province.Meta("province"), province.Handler())

	region := api.Region()
	group.Register(region.Meta("region"), region.Handler())
}
