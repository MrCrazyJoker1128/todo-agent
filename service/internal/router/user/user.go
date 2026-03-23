package user

import (
	"service/internal/controller/user"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRoutes(group *ghttp.RouterGroup) {
	userController := user.NewUserController()
	group.Group("/user", func(group *ghttp.RouterGroup) {
		group.POST("/login", userController.Login)
		group.POST("/register", userController.UserRegister)
	})
}
