package knowledge

import (
	"service/internal/controller/knowledge"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRoutes(group *ghttp.RouterGroup) {
	_ = knowledge.NewKnowledgeController()
	group.Group("/knowledge", func(group *ghttp.RouterGroup) {
		group.POST("/add", knowledge.NewKnowledgeController().AddKnowledge)
	})
}
