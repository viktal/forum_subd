package http

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserUseCase user.UseCase
}

type Resp struct {
	User *models.User `json:"user"`
}

func NewRest(router *gin.RouterGroup, useCase user.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup) {
	router.GET("/:nickname/profile", u.GetUserProfile)
	router.POST("/:nickname/create", u.CreateUser)
	router.POST("/:nickname/profile", u.UpdateUser)
}

func (u *UserHandler) GetUserProfile(ctx *gin.Context) {
	var req struct {
		Nickname string `uri:"nickname" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetUserProfile(req.Nickname)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandler) CreateUser(ctx *gin.Context) {
	var req struct {
		Nickname string `uri:"nickname" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user models.User
	user.Nickname = req.Nickname
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userNew, err := u.UserUseCase.CreateUser(user)
	if err != nil {
		if err.Code() == 409 {
			ctx.JSON(http.StatusConflict, userNew)
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return

	}
	ctx.JSON(http.StatusCreated, userNew[0])
}

func (u *UserHandler) UpdateUser(ctx *gin.Context) {
	var req struct {
		Nickname string `uri:"nickname" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var user models.UserUpdate
	user.Nickname = req.Nickname
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userUpdate, err := u.UserUseCase.UpdateUser(user)
	if err != nil {
		if err.Code() == 404 {
			ctx.JSON(http.StatusNotFound, err)
			return
		} else if err.Code() == 409 {
			ctx.JSON(http.StatusConflict, err)
			return
		} else {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

	}
	ctx.JSON(http.StatusOK, userUpdate)
}
