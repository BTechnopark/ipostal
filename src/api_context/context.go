package api_context

import (
	"errors"
	"net/http"
	"time"

	"github.com/BTechnopark/ipostal/pkg/cache"
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
	ForceExec(handler func(seterr func(err error))) ApiContext
	Finish(data Response) ApiContext
	Cache(c cache.Cache, key string, resp any) ApiContext
}

type apiContext struct {
	ctx *gin.Context
	err error

	cache    cache.Cache
	cacheKey string
	useCache bool
}

// Finish implements ApiContext.
func (a *apiContext) Finish(data Response) ApiContext {
	return a.
		Exec(func(seterr func(err error)) {
			if a.cache != nil {
				err := a.cache.Set(a.cacheKey, data, time.Minute*5)
				if err != nil {
					seterr(err)
					return
				}
			}
		}).
		ForceExec(func(seterr func(err error)) {
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

func (a *apiContext) Cache(c cache.Cache, key string, resp any) ApiContext {
	return a.Exec(func(seterr func(err error)) {
		a.cache = c
		a.cacheKey = key

		err := a.cache.Get(key, resp)
		if err != nil {
			if !errors.Is(err, cache.ErrCacheNotFound) {
				seterr(err)
			}
			return
		}

		a.useCache = true
	})
}

// Exec implements ApiContext.
func (a *apiContext) Exec(handler func(seterr func(err error))) ApiContext {
	if a.err != nil {
		return a
	}

	if a.useCache {
		return a
	}

	handler(func(err error) {
		a.err = err
	})

	return a
}

// Exec implements ApiContext.
func (a *apiContext) ForceExec(handler func(seterr func(err error))) ApiContext {
	handler(func(err error) {
		a.err = err
	})

	return a
}
