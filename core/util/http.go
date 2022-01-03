package util

import (
	"context"
	"time"
)

type key int

//CtxKeyUserID defines context key for user id
const CtxKeyUserID = key(1)

var debugToken string = ""

func CurrentUserID(ctx context.Context) int {
	userID := ctx.Value(CtxKeyUserID)
	result := 0
	if userID != nil {
		result = userID.(int)
	}
	return result
}

// Generate a debug token and store it into memory, new one will always override existing one
// After 'duration' it will be removed.
func NewDebugToken(duration time.Duration) string {
	guid := GenerateGUID()
	debugToken = guid
	go clearDebugTokenDelay(duration)
	return debugToken
}

func clearDebugTokenDelay(duration time.Duration) {
	time.Sleep(duration)
	ClearDebugToken()
}

func GetDebugToken() string {
	return debugToken
}

func ClearDebugToken() {
	debugToken = ""
}

//Site visiting anonymous user
func AnonymousUser() int {
	siteUser := GetConfigSectionI("site_settings")["site_user"].(int)
	return siteUser
}

func IsAnonymousUser(ctx context.Context) bool {
	userID := CurrentUserID(ctx)
	if userID == 0 {
		return false
	}
	result := userID == AnonymousUser()
	return result
}
