package chat

import (
	"service/internal/controller/chat"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRoutes(group *ghttp.RouterGroup) {
	chatController := chat.NewChatController()
	group.Group("/chat", func(group *ghttp.RouterGroup) {
		group.GET("/list", chatController.GetChatList)
		group.POST("/insert", chatController.InsertChatMessage)
		group.POST("/send", chatController.SendChatMessage)
	})
}
