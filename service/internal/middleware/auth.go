package middleware

import (
	"fmt"
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

	redisToken, err := utility.RedisGet(r.Context(), fmt.Sprintf("token:%s", claims.UserId))
	if err != nil || redisToken == nil {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "token expired",
		})
		r.Exit()
	}

	redisTokenStr, ok := redisToken.(string)
	if !ok || redisTokenStr != token {
		r.Response.WriteJson(Response{
			Code:    401,
			Message: "token mismatch",
		})
		r.Exit()
	}

	r.SetParam("userId", claims.UserId)
	r.Middleware.Next()
}
