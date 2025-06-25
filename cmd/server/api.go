package main

import (
	"net/http"
	"time"

	"github.com/BTechnopark/ipostal/config"
	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/client"
	"github.com/BTechnopark/ipostal/src/ipostal_api"
	"github.com/BTechnopark/ipostal/src/kodepos"
	"github.com/BTechnopark/ipostal/src/session"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
)

func SetUpSdk() gin_api.ApiSdk {
	r := gin.Default()
	sdk := gin_api.NewGinApiSdk(r)

	return sdk
}

func CreateApi(sdk gin_api.ApiSdk) error {
	DevMode := config.GetEnv("DEV_MODE", "") != ""
	if DevMode {
		RegisterDoc(sdk)
	}
	CheckHealth(sdk)

	cacheDuration := config.GetEnv("CACHE_DURATION", "5m")
	d, err := time.ParseDuration(cacheDuration)
	if err != nil {
		return err
	}
	cache := cache.NewCache(d)

	v1 := sdk.Group("v1")

	session := session.NewSession("")
	baseApi := client.NewApi(session)

	// indoPosconfig := &Config{
	// 	BaseUrl: config.GetEnv("POSTAL_URI", "https://kodepos.posindonesia.co.id"),
	// }
	// indonesianPostalCodeApi := api.NewPosIndonesiaApi(indoPosconfig, baseApi, cache)

	kodeposConfig := &Config{
		BaseUrl: config.GetEnv("KODEDEPOS_URI", "https://kodepos.co.id"),
	}
	kodePosApi := kodepos.NewKodePosApi(kodeposConfig, baseApi)

	api := ipostal_api.NewApi(kodePosApi, cache)

	ipostal_api.RegisterIPostalApi(v1, api)

	return nil
}

func CheckHealth(sdk gin_api.ApiSdk) {
	sdk.Register(&gin_api.ApiData{
		Method:       http.MethodGet,
		RelativePath: "/health",
		Response:     &ipostal_api.ResponseData[any]{},
	}, func(ctx *gin.Context) {
		response := &ipostal_api.ResponseData[any]{
			Message: "OK",
		}
		ctx.JSON(http.StatusOK, &response)
	})
}
