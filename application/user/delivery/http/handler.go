package http

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/user"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserHandler struct {
	UserUseCase user.UseCase
}

func NewRest(router *fasthttprouter.Router, useCase user.UseCase) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (u *UserHandler) routes(router *fasthttprouter.Router) {
	router.GET("/api/user/:nickname/profile", u.GetUserProfile)
	router.POST("/api/user/:nickname/create", u.CreateUser)
	router.POST("/api/user/:nickname/profile", u.UpdateUser)
}

func (u *UserHandler) GetUserProfile(c *fasthttp.RequestCtx) {
	Nickname := c.UserValue("nickname").(string)

	user, err := u.UserUseCase.GetUserProfile(Nickname)

	c.SetContentType("application/json")
	if err != nil {
		c.SetStatusCode(http.StatusNotFound)
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(user, c)
	}
}

func (u *UserHandler) CreateUser(c *fasthttp.RequestCtx) {
	var user models.User
	_ = user.UnmarshalJSON(c.PostBody())
	user.Nickname = c.UserValue("nickname").(string)

	userNew, err := u.UserUseCase.CreateUser(user)

	c.SetContentType("application/json")
	if err != nil {
		if err.Code() == 409 {
			c.SetStatusCode(http.StatusConflict)
			_, _ = easyjson.MarshalToWriter(models.UserList(userNew), c)
		} else {
			c.SetStatusCode(http.StatusInternalServerError)
			_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
		}
	} else {
		c.SetStatusCode(http.StatusCreated)
		_, _ = easyjson.MarshalToWriter(userNew[0], c)
	}
}

func (u *UserHandler) UpdateUser(c *fasthttp.RequestCtx) {
	var user models.UserUpdate
	_ = user.UnmarshalJSON(c.PostBody())
	user.Nickname = c.UserValue("nickname").(string)

	userUpdate, err := u.UserUseCase.UpdateUser(user)

	c.SetContentType("application/json")
	if err != nil {
		if err.Code() == 404 {
			c.SetStatusCode(http.StatusNotFound)
		} else if err.Code() == 409 {
			c.SetStatusCode(http.StatusConflict)
		} else {
			c.SetStatusCode(http.StatusBadRequest)
		}
		_, _ = easyjson.MarshalToWriter(common.MessageError{Message: err.Error()}, c)
	} else {
		c.SetStatusCode(http.StatusOK)
		_, _ = easyjson.MarshalToWriter(userUpdate, c)
	}
}
