package http

import (
	"forum/application/models"
	"forum/application/post"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
		Related string `form:"related"`
	}

	_ = ctx.ShouldBindQuery(&ListRelated)
	err = nil

	result, err := u.UserUseCase.GetPostDetails(id, strings.Split(ListRelated.Related, ","))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
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

	var newPost *models.Post
	if body.Message == "" {
		newPost, err = u.UserUseCase.GetPostByID(id)
	} else {
		newPost, err = u.UserUseCase.UpdatePostDetails(id, body.Message)
	}

	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, newPost)
}
