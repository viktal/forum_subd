package http

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/post"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	UserUseCase post.UseCase
}

func NewRest(router *fasthttprouter.Router, useCase post.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *fasthttprouter.Router) {
	router.GET("/api/post/:id/details", u.GetPostDetails)
	router.POST("/api/post/:id/details", u.UpdatePostDetails)
}

func (u *UserHandler) GetPostDetails(c *fasthttp.RequestCtx) {
	related := string(c.URI().QueryArgs().Peek("related"))
	id, _ := strconv.ParseInt(c.UserValue("id").(string), 10, 64)

	result, err := u.UserUseCase.GetPostDetails(int(id), strings.Split(related, ","))

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(result, c)
	}
}

func (u *UserHandler) UpdatePostDetails(c *fasthttp.RequestCtx) {
	id, _ := strconv.ParseInt(c.UserValue("id").(string), 10, 64)

	var message common.MessageError
	_ = message.UnmarshalJSON(c.PostBody())

	var newPost *models.Post
	var err error
	if message.Message == "" {
		newPost, err = u.UserUseCase.GetPostByID(int(id))
	} else {
		newPost, err = u.UserUseCase.UpdatePostDetails(int(id), message.Message)
	}

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(newPost, c)
	}
}
