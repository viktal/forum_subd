package http

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/thread"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUseCase thread.UseCase
}

type Request struct {
	ID   int32 `uri:"identification" binding:"required"`
	Slug string `uri:"identification" binding:"string"`
}

func NewRest(router *fasthttprouter.Router, useCase thread.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *fasthttprouter.Router) {
	router.POST("/api/thread/:slug_or_id/create", u.CreatePost)
	router.GET("/api/thread/:slug_or_id/details", u.GetThreadDetails)
	router.POST("/api/thread/:slug_or_id/details", u.UpdateThread)
	router.GET("/api/thread/:slug_or_id/posts", u.GetPostsThread)
	router.POST("/api/thread/:slug_or_id/vote", u.VoteOnThread)
}

func (u *UserHandler) GetThreadDetails(c *fasthttp.RequestCtx) {
	slugOrID := c.UserValue("slug_or_id").(string)
	thread, err := u.UserUseCase.GetThreadDetails(slugOrID)

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(thread, c)
	}
}

func (u *UserHandler) CreatePost(c *fasthttp.RequestCtx) {
	slugOrID := c.UserValue("slug_or_id").(string)

	posts := new(models.ListPosts)
	_ = posts.UnmarshalJSON(c.PostBody())

	thr, err := u.UserUseCase.GetThreadDetails(slugOrID)

	c.SetContentType("application/json")
	if err != nil || thr.ThreadID == 0 {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
		return
	} else if len(*posts) == 0 {
		c.SetStatusCode(http.StatusCreated)
		_, _ = easyjson.MarshalToWriter(posts, c)
		return
	}

	newPosts, err := u.UserUseCase.CreatePosts(slugOrID, *posts)

	if err != nil {
		if err.Error() == "Parent post was created in another thread" {
			c.SetStatusCode(http.StatusConflict)
		} else if err.Error() == "Not found" {
			c.SetStatusCode(http.StatusNotFound)
		} else {
			c.SetStatusCode(http.StatusInternalServerError)
		}
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
	} else {
		c.SetStatusCode(http.StatusCreated)
		_, _ = easyjson.MarshalToWriter(newPosts, c)
	}
}

func (u *UserHandler) UpdateThread(c *fasthttp.RequestCtx) {
	slugOrID := c.UserValue("slug_or_id").(string)

	var thread models.ThreadUpdate
	_ = thread.UnmarshalJSON(c.PostBody())

	c.SetContentType("application/json")

	if thread.Message == nil && thread.Title == nil {
		oldThread, err := u.UserUseCase.GetThreadDetails(slugOrID)
		if err != nil {
			c.SetStatusCode(http.StatusNotFound)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
			return
		} else {
			c.SetStatusCode(http.StatusOK)
			_, _ = easyjson.MarshalToWriter(oldThread, c)
			return
		}
	}

	newThread, err := u.UserUseCase.UpdateThread(slugOrID, thread)
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(newThread, c)
	}
}

func (u *UserHandler) GetPostsThread(c *fasthttp.RequestCtx) {
	slugOrID := c.UserValue("slug_or_id").(string)


	limit, _ := strconv.Atoi(string(c.URI().QueryArgs().Peek("limit")))
	desc, _ := strconv.ParseBool(string(c.URI().QueryArgs().Peek("desc")))
	since, _ := strconv.Atoi(string(c.URI().QueryArgs().Peek("since")))
	sort := string(c.URI().QueryArgs().Peek("sort"))
	params := models.PostParams{
		Limit: uint(limit),
		Since: &since,
		Sort: sort,
		Desc: desc,
	}
	if *params.Since == 0 {
		params.Since = nil
	}

	if params.Sort == "" {
		params.Sort = common.Flat
	}

	posts, err := u.UserUseCase.GetPostsThread(slugOrID, params)

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else if posts == nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
	} else if *posts == nil {
		c.SetStatusCode(http.StatusOK)
		c.Write([]byte("[]"))
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(models.ListPosts(*posts), c)
	}
}

func (u *UserHandler) VoteOnThread(c *fasthttp.RequestCtx) {
	slugOrID := c.UserValue("slug_or_id").(string)

	var vote models.Vote
	_ = vote.UnmarshalJSON(c.PostBody())


	c.SetContentType("application/json")
	if vote.Voice != 1 && vote.Voice != -1 {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Wrong voice"}, c)
		return
	}

	thread, err := u.UserUseCase.VoteOnThread(slugOrID, vote)
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(thread, c)
	}
}
