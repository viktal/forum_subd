package http

import (
	"forum/application/common"
	"forum/application/forum"
	"forum/application/models"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"time"
)

type ForumHandler struct {
	UseCaseForum forum.UseCase
}

func NewRest(router *fasthttprouter.Router,
	useCaseResume forum.UseCase) *ForumHandler {
	rest := &ForumHandler{
		UseCaseForum: useCaseResume,
	}
	rest.routes(router)
	return rest
}

func (r *ForumHandler) routes(router *fasthttprouter.Router) {


	router.POST("/api/forum/:slug", r.CreateForum)
	router.POST("/api/forum/:slug/create", r.CreateThread)
	router.GET("/api/forum/:slug/details", r.GetForumBySlug)
	router.GET("/api/forum/:slug/threads", r.GetAllForumTreads)
	router.GET("/api/forum/:slug/users", r.GetAllForumUsers)
}

func (r *ForumHandler) CreateForum(c *fasthttp.RequestCtx) {
	var forum models.ForumCreate
	_ = forum.UnmarshalJSON(c.PostBody())


	result, err := r.UseCaseForum.CreateForum(forum)

	c.SetContentType("application/json")
	if err != nil {
		if err.Code() == 404 {
			c.SetStatusCode(http.StatusNotFound)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
		} else if err.Code() == 409 {
			c.SetStatusCode(http.StatusConflict)
			_, _ = easyjson.MarshalToWriter(result, c)
		} else {
			c.SetStatusCode(http.StatusInternalServerError)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
		}
	} else {
		c.SetStatusCode(http.StatusCreated)
		_, _ = easyjson.MarshalToWriter(result, c)
	}
}

func (r *ForumHandler) CreateThread(c *fasthttp.RequestCtx) {
	var template models.Thread
	template.UnmarshalJSON(c.PostBody())
	template.Forum = c.UserValue("slug").(string)


	result, err := r.UseCaseForum.CreateThread(template.Forum, template)

	c.SetContentType("application/json")
	if err != nil {
		if err.Code() == 404 {
			c.SetStatusCode(http.StatusNotFound)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
		} else if err.Code() == 409 {
			c.SetStatusCode(http.StatusConflict)
			_, _ = easyjson.MarshalToWriter(result, c)
		} else {
			c.SetStatusCode(http.StatusNotFound)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
		}
	} else {
		c.SetStatusCode(http.StatusCreated)
		_, _ = easyjson.MarshalToWriter(result, c)
	}
}

func (r *ForumHandler) GetForumBySlug(c *fasthttp.RequestCtx) {
	slug := c.UserValue("slug").(string)
	result, err := r.UseCaseForum.GetForumBySlug(slug)


	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else if result.ForumID == 0 {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found."}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(result, c)
	}
}

func (r *ForumHandler) GetAllForumTreads(c *fasthttp.RequestCtx) {
	slugForum := c.UserValue("slug").(string)

	limit, _ := strconv.Atoi(string(c.URI().QueryArgs().Peek("limit")))
	desc, _ := strconv.ParseBool(string(c.URI().QueryArgs().Peek("desc")))
	since := string(c.URI().QueryArgs().Peek("since"))

	params := models.ForumParams{
		Limit: uint(limit),
		Desc:  desc,
	}

	if since != "" {
		params.Since, _ = time.Parse(time.RFC3339, since)
	}

	result, err := r.UseCaseForum.GetAllForumTreads(slugForum, params)

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else if *result == nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(models.ListThread(*result), c)
	}
}

func (r *ForumHandler) GetAllForumUsers(c *fasthttp.RequestCtx) {
	slugForum := c.UserValue("slug").(string)


	limit, _ := strconv.Atoi(string(c.URI().QueryArgs().Peek("limit")))
	desc, _ := strconv.ParseBool(string(c.URI().QueryArgs().Peek("desc")))
	since := string(c.URI().QueryArgs().Peek("since"))


	params := models.UserParams{
		Since: &since,
		Limit: uint(limit),
		Desc:  desc,
	}
	if since == "" {
		params.Since = nil
	}


	result, err := r.UseCaseForum.GetAllForumUsers(slugForum, params)


	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else if result == nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: "Not found"}, c)
	} else if *result == nil {
		c.SetStatusCode(http.StatusOK)
		_, _ = c.Write([]byte("[]"))
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(models.UserList(*result), c)
	}
}

