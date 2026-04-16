package api

import (
	"time"

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

type GetKnowledgeListReq struct {
	g.Meta `method:"GET" path:"/list" summary:"获取知识库列表" tags:"知识库"`
	Title  string `json:"title" description:"知识库标题"`
	Author string `json:"author" description:"知识库作者"`
	Domain string `json:"domain" description:"知识库域名"`
}
type GetKnowledgeListItem struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Domain    string    `json:"domain"`
	Url       string    `json:"url"`
	UploadAt  time.Time `json:"uploadAt"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}
type GetKnowledgeListRes struct {
	Items []GetKnowledgeListItem `json:"items"`
	Total int64                  `json:"total"`
}
