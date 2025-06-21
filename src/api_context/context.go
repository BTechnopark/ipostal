package api_context

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApiContext(ctx *gin.Context) ApiContext {
	return &apiContext{
		ctx: ctx,
	}
}

type Response interface {
	SetError(err error)
}

type ApiContext interface {
	BindQuery(query interface{}) ApiContext
	BindPayload(payload interface{}) ApiContext
	Exec(func(seterr func(err error))) ApiContext
	Finish(data Response) ApiContext
}

type apiContext struct {
	ctx *gin.Context
	err error
}

// Finish implements ApiContext.
func (a *apiContext) Finish(data Response) ApiContext {
	return a.Exec(func(seterr func(err error)) {
		if a.err != nil {
			data.SetError(a.err)
			a.ctx.JSON(http.StatusInternalServerError, data)
			return
		}

		a.ctx.JSON(200, data)
	})
}

// BindPayload implements ApiContext.
func (a *apiContext) BindPayload(payload interface{}) ApiContext {
	return a.Exec(func(seterr func(err error)) {
		err := a.ctx.BindJSON(payload)
		if err != nil {
			seterr(err)
		}
	})
}

// BindQuery implements ApiContext.
func (a *apiContext) BindQuery(query interface{}) ApiContext {
	return a.Exec(func(seterr func(err error)) {
		err := a.ctx.BindQuery(query)
		if err != nil {
			seterr(err)
		}
	})
}

// Exec implements ApiContext.
func (a *apiContext) Exec(handler func(seterr func(err error))) ApiContext {
	if a.err != nil {
		return a
	}

	handler(func(err error) {
		a.err = err
	})

	return a
}
