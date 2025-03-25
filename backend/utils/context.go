package utils

import (
	"context"
	"errors"
)

type ContextKey int

const (
	AuthCtxKey ContextKey = iota
)

func WithUserId(ctx *context.Context, userId int) {
	*ctx = context.WithValue(*ctx, AuthCtxKey, userId)
}

func ExtractUserId(ctx *context.Context) (int, error) {
	id, ok := (*ctx).Value(AuthCtxKey).(int)
	if !ok {
		return 0, errors.New("failed to extract user id")
	}
	return id, nil
}
