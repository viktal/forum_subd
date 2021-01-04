package http

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/thread"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserUseCase thread.UseCase
}

type Request struct {
	ID   int32 `uri:"identification" binding:"required"`
	Slug string `uri:"identification" binding:"string"`
}

func NewRest(router *gin.RouterGroup, useCase thread.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup) {
	router.POST("/:slug_or_id/create", u.CreatePost)
	router.GET("/:slug_or_id/details", u.GetThreadDetails) //+
	router.POST("/:slug_or_id/details", u.UpdateThread) //+
	router.GET("/:slug_or_id/posts", u.GetPostsThread)
	router.GET("/:slug_or_id/vote", u.VoteOnThread)
}

func (u *UserHandler) GetThreadDetails(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")
	thread, err := u.UserUseCase.GetThreadDetails(slugOrID)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, thread)
}

func (u *UserHandler) CreatePost(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var posts models.ListPosts
	//var newThreads models.ListThread
	if err := ctx.ShouldBindJSON(&posts); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newPosts, err := u.UserUseCase.CreatePosts(slugOrID, posts)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, newPosts)
}

func (u *UserHandler) UpdateThread(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var thread models.Thread
	if err := ctx.ShouldBindJSON(&thread); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newThread, err := u.UserUseCase.UpdateThread(slugOrID, thread)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, newThread)
}

func (u *UserHandler) GetPostsThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var query models.ThreadParams
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		query.SlugOrID = req.Slug
	} else if req.ID != 0 {
		query.SlugOrID = string(req.ID)
	}
	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	threads, err := u.UserUseCase.GetPostsThread(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, threads)
}

func (u *UserHandler) VoteOnThread(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if (Request{} == req) {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.UriErrorThread})
		return
	}

	var vote models.Vote
	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if req.Slug != "" {
		vote.SlugOrID = req.Slug
	} else if req.ID != 0 {
		vote.SlugOrID = string(req.ID)
	}

	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	threads, err := u.UserUseCase.VoteOnThread(vote)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, threads)
}