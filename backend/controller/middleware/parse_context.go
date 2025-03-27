package middleware

import (
	"errors"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

const GinContextKeyUser string = "user"

func GetUserFromContext(ctx *gin.Context) (*store.User, error) {
	userAny, exist := ctx.Get(GinContextKeyUser)
	if !exist {
		return nil, errors.New("no user found in the gin context")
	}
	user, ok := userAny.(store.User)
	if !ok {
		return nil, errors.New("failed to convert any to store.User")
	}
	return &user, nil
}
