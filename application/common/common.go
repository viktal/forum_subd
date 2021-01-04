package common

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CandID = "cand_id"
	EmplID = "empl_id"
	UserID = "user_id"
	Candidate = "candidate"
	Employer = "employer"
)


func HandlerGetCurrentUserID(ctx *gin.Context, user string) (id uuid.UUID, err error) {
	session := sessions.Default(ctx)
	userIDStr := session.Get(user)
	if userIDStr == nil {
		return uuid.Nil, nil
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
