package http

import (
	"fmt"
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
	router.POST("/:slug_or_id/create", u.CreatePost) //+
	router.GET("/:slug_or_id/details", u.GetThreadDetails) //+
	router.POST("/:slug_or_id/details", u.UpdateThread) //+
	router.GET("/:slug_or_id/posts", u.GetPostsThread) //+
	router.POST("/:slug_or_id/vote", u.VoteOnThread)
}

func (u *UserHandler) GetThreadDetails(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")
	thread, err := u.UserUseCase.GetThreadDetails(slugOrID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found thread"))
		return
	}

	ctx.JSON(http.StatusOK, thread)
}

func (u *UserHandler) CreatePost(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var posts models.ListPosts
	if err := ctx.ShouldBindJSON(&posts); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	thr, err := u.UserUseCase.GetThreadDetails(slugOrID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found thread"))
		return
	}
	if thr.ThreadID == 0 {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found thread"))
		return
	}

	if len(posts) == 0 {
		ctx.JSON(http.StatusCreated, posts)
		return
	}

	newPosts, err := u.UserUseCase.CreatePosts(slugOrID, posts)

	if err != nil {
		if err.Error() == "Parent post was created in another thread" {
			ctx.JSON(http.StatusConflict, err)
			return
		} else if err.Error() == "Not found" {
			ctx.JSON(http.StatusNotFound, err)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, newPosts)
}

func (u *UserHandler) UpdateThread(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var thread models.ThreadUpdate
	if err := ctx.ShouldBindJSON(&thread); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if thread.Message == nil && thread.Title == nil {
		oldThread, err := u.UserUseCase.GetThreadDetails(slugOrID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found thread"))
			return
		} else {
			ctx.JSON(http.StatusOK, oldThread)
			return
		}
	}

	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	newThread, err := u.UserUseCase.UpdateThread(slugOrID, thread)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, newThread)
}

func (u *UserHandler) GetPostsThread(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var params models.PostParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if params.Sort == "" {
		params.Sort = common.Flat
	}

	posts, err := u.UserUseCase.GetPostsThread(slugOrID, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if *posts == nil {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found."))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (u *UserHandler) VoteOnThread(ctx *gin.Context) {
	slugOrID := ctx.Param("slug_or_id")

	var vote models.Vote
	if err := ctx.ShouldBindJSON(&vote); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//if err := common.ReqValidation(&req); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	if vote.Voice != 1 && vote.Voice != -1 {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: "Wrong voice"})
		return
	}

	threads, err := u.UserUseCase.VoteOnThread(slugOrID, vote)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, threads)
}