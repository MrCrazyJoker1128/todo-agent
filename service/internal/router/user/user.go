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
		group.GET("/userInfo", userController.GetUserInfo)
		group.POST("/model-config", userController.SaveModelConfig)
		group.GET("/model-config", userController.GetModelConfig)
		group.POST("/system-kek", userController.SaveSystemKEK)
		group.GET("/system-kek", userController.GetSystemKEK)
		group.POST("/profile", userController.UpdateProfile)
		group.POST("/password", userController.UpdatePassword)
	})
}
