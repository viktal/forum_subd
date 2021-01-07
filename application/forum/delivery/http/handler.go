package http

import (
	"fmt"
	"forum/application/forum"
	"forum/application/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type ForumHandler struct {
	UseCaseForum forum.UseCase
}

func NewRest(router *gin.RouterGroup,
	useCaseResume forum.UseCase) *ForumHandler {
	rest := &ForumHandler{
		UseCaseForum: useCaseResume,
	}
	rest.routes(router)
	return rest
}

func (r *ForumHandler) routes(router *gin.RouterGroup) {
	router.POST("/:path1/:path2", r.BuildPath) ///:slug/create CreatePost
	router.POST("/:path1", r.BuildPath) ///create CreateForum
	router.GET("/:slug/details", r.GetForumBySlug)
	router.GET("/:slug/threads", r.GetAllForumTreads)
	router.GET("/:slug/users", r.GetAllForumUsers)
}

func (r *ForumHandler) BuildPath(ctx *gin.Context) {
	path1 := ctx.Param("path1")
	path2 := ctx.Param("path2")

	if path1 == "create" && path2 == "" {
		r.CreateForum(ctx)
	} else if path1 != "" && path2 == "create" {
		slug := path1
		r.CreateThread(ctx, slug)
	} else {
		ctx.Status(http.StatusNotFound)
	}
}

func (r *ForumHandler) CreateForum(ctx *gin.Context) {
	var forum models.ForumCreate
	if err := ctx.ShouldBindJSON(&forum); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.UseCaseForum.CreateForum(forum)
	if err != nil {
		if err.Code() == 404 {
			ctx.JSON(http.StatusNotFound, err)
			return
		} else if err.Code() == 409 {
			ctx.JSON(http.StatusConflict, result)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, result)
}
func (r *ForumHandler) CreateThread(ctx *gin.Context, slugForum string) {
	var template models.Thread
	if err := ctx.ShouldBindBodyWith(&template, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//if err := common.ReqValidation(template); err != nil {
	//	ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
	//	return
	//}

	result, err := r.UseCaseForum.CreateThread(slugForum, template)
	if err != nil {
		if err.Code() == 404 {
			ctx.JSON(http.StatusNotFound, err)
			return
		} else if err.Code() == 409 {
			ctx.JSON(http.StatusConflict, result)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusCreated, result)
}

func (r *ForumHandler) GetForumBySlug(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.UseCaseForum.GetForumBySlug(req.Slug)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if result.ForumID == 0 {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found."))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (r *ForumHandler) GetAllForumTreads(ctx *gin.Context) {
	slugForum := ctx.Param("slug")

	var params models.ForumParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	result, err := r.UseCaseForum.GetAllForumTreads(slugForum, params)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if *result == nil {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("Not found."))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (r *ForumHandler) GetAllForumUsers(ctx *gin.Context) {
	var req struct {
		Slug string `uri:"slug" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var params models.ForumParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	result, err := r.UseCaseForum.GetAllForumUsers(req.Slug, params)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

