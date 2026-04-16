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
