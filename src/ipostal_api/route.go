package ipostal_api

import "github.com/muchrief/gin_api"

func RegisterIPostalApi(group gin_api.SdkGroup, api Api) {
	findPostalCode := api.FindPostalCode()
	group.Register(findPostalCode.Meta("search"), findPostalCode.Handler())
}
