package http

import (
	"forum/application/common"
	"forum/application/service"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserHandler struct {
	UserUseCase service.UseCase
}

func NewRest(router *fasthttprouter.Router, useCase service.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *fasthttprouter.Router) {
	router.POST("/api/service/clear", u.ClearDB)
	router.GET("/api/service/status", u.StatusDB)
}

func (u *UserHandler) StatusDB(c *fasthttp.RequestCtx) {

	result, err := u.UserUseCase.GetStatusDB()

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(result, c)
	}
}

func (u *UserHandler) ClearDB(c *fasthttp.RequestCtx) {
	err := u.UserUseCase.ClearDB()

	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
	}
}
