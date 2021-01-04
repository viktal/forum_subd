package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/service"
	"net/http"
)

type UserHandler struct {
	UserUseCase service.UseCase
}

func NewRest(router *gin.RouterGroup, useCase service.UseCase, AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.POST("/clear", u.ClearDB)
	router.GET("/status", u.StatusDB)
}

func (u *UserHandler) StatusDB(ctx *gin.Context) {

	result, err := u.UserUseCase.GetStatusDB()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (u *UserHandler) ClearDB(ctx *gin.Context) {
	err := u.UserUseCase.ClearDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.Status(http.StatusOK)
}
