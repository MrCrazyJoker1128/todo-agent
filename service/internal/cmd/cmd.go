package cmd

import (
	"context"

	"service/internal/middleware"
	"service/internal/router/chat"
	"service/internal/router/knowledge"
	"service/internal/router/task"
	"service/internal/router/user"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 前端 dist 同源托管：resource/public/index.html + assets。
			// 单页 fallback：非 /api 路由统一回 index.html（前端 vue-router 接管）。
			s.SetServerRoot("resource/public")
			// SPA fallback：非 /api、非 /swagger、非静态文件一律回 index.html。
			s.BindMiddlewareDefault(func(r *ghttp.Request) {
				r.Middleware.Next()
				if r.Response.Status >= 400 && r.URL.Path != "/api.json" {
					r.Response.ClearBuffer()
					r.Response.ServeFile("resource/public/index.html")
				}
			})
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.MiddlewareAuth)
				group.Middleware(middleware.MiddlewareResponse)
				user.RegisterRoutes(group)
				task.RegisterRoutes(group)
				chat.RegisterRoutes(group)
				knowledge.RegisterRoutes(group)
			})
			s.Run()
			return nil
		},
	}
)
