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
}

type RespError struct {
	Err string `json:"error"`
}

func (e Err) Code() int         { return e.code }
func (e Err) Error() string     { return e.message }

func (e Err) MarshalJSON() ([]byte, error) {
	ret := map[string]interface{}{
		"code":    e.code,
		"message": e.message,
	}
	return json.Marshal(ret)
}

func (e Err) String() string {
	data, _ := e.MarshalJSON()
	return string(data)
}

func NewErr(code int, message string) Err {
	return Err{
		code:    code,
		message: message,
	}
}
