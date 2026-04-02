package middleware

import (
	"context"
	"strings"

	"service/internal/utility"

	"github.com/gogf/gf/v2/net/ghttp"
)

var WhiteList = []string{
	"/api/user/login",
	"/api/user/register",
	"/api.json",
	"/swagger",
}

func MiddlewareAuth(r *ghttp.Request) {
	path := r.URL.Path

	for _, item := range WhiteList {
		if strings.HasPrefix(path, item) {
			r.Middleware.Next()
			return
		}
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "missing token",
		})
		r.Exit()
	}

	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := utility.ParseToken(r.Context(), token)
	if err != nil {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "invalid token",
		})
		r.Exit()
	}

	redisToken, err := utility.RedisGet(r.Context(), "token:"+claims.UserId)
	if err != nil || redisToken == "" {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "token expired",
		})
		r.Exit()
	}
	if strings.TrimSpace(redisToken) != token {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "token mismatch",
		})
		r.Exit()
	}
	r.SetParam("userId", claims.UserId)
	ctx := r.Context()
	ctx = context.WithValue(ctx, "user_id", claims.UserId)
	// 5. 将新 context 设置回 Request
	r.SetCtx(ctx)
	r.Middleware.Next()
}

// 从 Go context 获取 user_id（Tool 层使用）
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}
