package knowledge

import (
	"service/internal/controller/knowledge"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRoutes(group *ghttp.RouterGroup) {
	_ = knowledge.NewKnowledgeController()
	group.Group("/knowledge", func(group *ghttp.RouterGroup) {
		group.POST("/add", knowledge.NewKnowledgeController().AddKnowledge)
		group.GET("/list", knowledge.NewKnowledgeController().GetKnowledgeList)
		group.DELETE("/delete", knowledge.NewKnowledgeController().DeleteKnowledgeUrl)
		group.POST("/bases", knowledge.NewKnowledgeController().CreateBase)
		group.GET("/bases", knowledge.NewKnowledgeController().ListBases)
		group.DELETE("/bases/delete", knowledge.NewKnowledgeController().DeleteBase)
		group.GET("/files", knowledge.NewKnowledgeController().ListFiles)
		group.DELETE("/files/delete", knowledge.NewKnowledgeController().DeleteFile)
		group.POST("/upload", knowledge.NewKnowledgeController().UploadMd)
	})
}
