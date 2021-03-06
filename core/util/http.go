package util

import (
	"context"
)

type key int

//CtxKeyUserID defines context key for user id
const CtxKeyUserID = key(1)

func CurrentUserID(ctx context.Context) int {
	userID := ctx.Value(CtxKeyUserID)
	result := 0
	if userID != nil {
		result = userID.(int)
	}
	return result
}

func AnonymousUser() int {
	return 2 //todo: get from config
}

func IsAnonymousUser(ctx context.Context) bool {
	userID := CurrentUserID(ctx)
	if userID == 0 {
		return false
	}
	result := userID == AnonymousUser()
	return result
}
