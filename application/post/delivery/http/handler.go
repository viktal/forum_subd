package http

import (
	"forum/application/post"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUseCase post.UseCase
}

func NewRest(router *gin.RouterGroup, useCase post.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup) {
	router.GET("/:id/details", u.GetPostDetails)
	router.POST("/:id/details", u.UpdatePostDetails)
}

func (u *UserHandler) GetPostDetails(ctx *gin.Context) {
	var err error
	var id int
	if id, err = strconv.Atoi(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var ListRelated struct {
		Related []string
	}

	err = ctx.BindJSON(&ListRelated)
	if err != nil && err.Error() != "EOF" {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := u.UserUseCase.GetPostDetails(id, ListRelated.Related)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (u *UserHandler) UpdatePostDetails(ctx *gin.Context) {
	var err error
	var id int
	if id, err = strconv.Atoi(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var body struct {
		Message string
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newPost, err := u.UserUseCase.UpdatePostDetails(id, body.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, newPost)
}
