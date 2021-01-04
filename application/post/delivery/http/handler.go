package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/post"
	"net/http"
)

type UserHandler struct {
	UserUseCase post.UseCase
}

func NewRest(router *gin.RouterGroup, useCase post.UseCase, AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/:id/details", u.GetPostDetails)
	router.POST("/:id/details", u.UpdatePostDetails)
}

func (u *UserHandler) GetPostDetails(ctx *gin.Context) {
	var req struct {
		ID int `uri:"id" binding:"required,int"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var relatad []string
	if err := ctx.ShouldBindQuery(&relatad); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}

	result, err := u.UserUseCase.GetPostDetails(req.ID, relatad)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (u *UserHandler) UpdatePostDetails(ctx *gin.Context) {
	var req struct {
		ID int `uri:"id" binding:"required,int"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var body struct {
		NewMessage string `json:"id" binding:"required,string"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newPost, err := u.UserUseCase.UpdatePostDetails(req.ID, body.NewMessage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, newPost)
}
