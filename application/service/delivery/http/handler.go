package http

import (
	"forum/application/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserUseCase service.UseCase
}

func NewRest(router *gin.RouterGroup, useCase service.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup) {
	router.POST("/clear", u.ClearDB)
	router.GET("/status", u.StatusDB)
}

func (u *UserHandler) StatusDB(ctx *gin.Context) {

	result, err := u.UserUseCase.GetStatusDB()

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (u *UserHandler) ClearDB(ctx *gin.Context) {
	err := u.UserUseCase.ClearDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
