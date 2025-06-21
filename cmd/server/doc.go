package main

import (
	"fmt"
	"net/http"

	"github.com/BTechnopark/ipostal/config"
	"github.com/gin-gonic/gin"
	"github.com/muchrief/gin_api"
	"github.com/muchrief/go_doc"
)

func RegisterDoc(sdk gin_api.ApiSdk) {
	doc := go_doc.NewGoDocumentation(go_doc.InfoVersion3)
	port := config.GetEnv("PORT", "3000")
	doc.
		AddServer(&go_doc.ServerObject{
			Url:         fmt.Sprintf("http://localhost:%s", port),
			Description: "This API Is Indonesian Postal Code",
		}).
		SetInfo(&go_doc.Info{
			Title:   "Indonesian Postal Code API",
			License: go_doc.LicenceApache,
		})

	sdk.Use(func(a gin_api.Api) {
		doc.RegisterDoc(a)
	})

	r := sdk.GetGinEngine()
	doc.RegisterDataDocumentation("/doc_data", func(method, path string) {
		r.Handle(method, path, func(ctx *gin.Context) {
			ctx.YAML(http.StatusOK, doc)
		})
	})

	doc.RegisterDocumentation("swagger", "/doc_data", "/docs", func(method, path string, template go_doc.TemplateFunc) {
		r.Handle(method, path, func(ctx *gin.Context) {
			temp, err := template()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
				return
			}

			ctx.Data(http.StatusOK, "text/html", []byte(temp))
		})
	})

	doc.RegisterDocumentation("redoc", "/doc_data", "/redoc", func(method, path string, template go_doc.TemplateFunc) {
		r.Handle(method, path, func(ctx *gin.Context) {
			temp, err := template()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
				return
			}

			ctx.Data(http.StatusOK, "text/html", []byte(temp))
		})
	})
}
