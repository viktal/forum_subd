package common

import (
	"encoding/json"
)

const (
	EmptyFieldErr = "Обязательные поля не заполнены."
	SessionErr    = "Ошибка авторизации. Попробуйте авторизоваться повторно."
	DataBaseErr   = "Что-то пошло не так. Попробуйте позже."
	UriErrorThread = "Ветка обсуждения отсутсвует в форуме."
)

type Err struct {
	code    int         `json:"code"`
	message string      `json:"message"`
	meta    interface{} `json:"meta"`
}

type RespError struct {
	Err string `json:"error"`
}

func (e Err) Code() int         { return e.code }
func (e Err) Error() string     { return e.message }
func (e Err) Meta() interface{} { return e.meta }

func (e Err) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"code":    e.code,
		"message": e.message,
	}
	if e.meta != nil {
		ret["meta"] = e.meta
	}
	return json.Marshal(ret)
}

func (e Err) String() string {
	data, _ := e.MarshalJSON()
	return string(data)
}

func NewErr(code int, message string, meta interface{}) Err {
	return Err{
		code:    code,
		message: message,
		meta:    meta,
	}
}

var ErrBadRequest = NewErr(400, "Неправильный запрос к серверу", nil)
var ErrInternalServerError = NewErr(500, "Внутренняя ошибка сервера", nil)

var ErrInvalidUpdatePassword = NewErr(1001, "неверный пароль", nil)
