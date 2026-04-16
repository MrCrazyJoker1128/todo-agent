package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

type AddKnowledgeReq struct {
	g.Meta `method:"POST" path:"/add" summary:"添加知识库" tags:"知识库"`
	Url    string `json:"url" validate:"required" description:"知识库链接"`
}

type AddKnowledgeRes struct {
	Id uuid.UUID `json:"id"`
}
