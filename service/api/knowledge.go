package api

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

type AddKnowledgeReq struct {
	g.Meta    `method:"POST" path:"/add" summary:"添加知识库" tags:"知识库"`
	Url       string `json:"url" validate:"required" description:"知识库链接"`
	LibraryId string `json:"libraryId" validate:"required" description:"归属库ID，新UI必传"`
}

type AddKnowledgeRes struct {
	Id uuid.UUID `json:"id"`
}

type GetKnowledgeListReq struct {
	g.Meta      `method:"GET" path:"list" summary:"获取知识库列表" tags:"知识库"`
	Title       string `json:"title" description:"知识库标题（中文分词搜索）"`
	Author      string `json:"author" description:"知识库作者（中文分词搜索）"`
	Domain      string `json:"domain" description:"知识库域名（模糊搜索）"`
	Status      string `json:"status" description:"知识库状态"`
	CurrentPage int    `json:"current_page" description:"当前页码"`
	PerPage     int    `json:"per_page" description:"每页数量"`
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
	Items []*GetKnowledgeListItem `json:"items"`
	Total int64                   `json:"total"`
}

type DeleteKnowledgeUrlReq struct {
	g.Meta `method:"DELETE" path:"/delete" summary:"删除知识库" tags:"知识库"`
	Id     string `json:"id" validate:"required" description:"知识库ID"`
}
type DeleteKnowledgeUrlRes struct {
	Message string `json:"message"`
}
type CreateBaseReq struct {
	g.Meta      `method:"POST" path:"/bases" summary:"新建知识库" tags:"知识库"`
	LibraryName string `json:"libraryName" validate:"required" description:"库名"`
	Description string `json:"description" description:"库描述"`
}
type CreateBaseRes struct {
	Id uuid.UUID `json:"id"`
}
type KnowledgeBaseItem struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
type ListBasesRes struct {
	Items []*KnowledgeBaseItem `json:"items"`
}
type ListBasesReq struct {
	g.Meta `method:"GET" path:"/bases" summary:"知识库列表" tags:"知识库"`
}
type DeleteBaseReq struct {
	g.Meta `method:"DELETE" path:"/bases/delete" summary:"删除知识库" tags:"知识库"`
	Id     string `json:"id" validate:"required" description:"库ID"`
}
type DeleteBaseRes struct {
	Message string `json:"message"`
}
type KnowledgeFileItem struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Source     string `json:"source"`
	ChunkCount int    `json:"chunkCount"`
	CreatedAt  string `json:"createdAt"`
	Status     string `json:"status"`
	Error      string `json:"error"`
}
type ListFilesRes struct {
	Items []*KnowledgeFileItem `json:"items"`
}
type ListFilesReq struct {
	g.Meta    `method:"GET" path:"/files" summary:"某库文件列表" tags:"知识库"`
	LibraryId string `json:"libraryId" validate:"required" description:"库ID"`
}
type DeleteFileReq struct {
	g.Meta    `method:"DELETE" path:"/files/delete" summary:"删除库内文件" tags:"知识库"`
	LibraryId string `json:"libraryId" validate:"required" description:"库ID"`
	Source    string `json:"source" validate:"required" description:"文件标识：html地址原文或 md://…"`
}
type DeleteFileRes struct {
	Message string `json:"message"`
}
type UploadMdReq struct {
	g.Meta    `method:"POST" path:"/upload" summary:"上传md文件到库" tags:"知识库"`
	LibraryId string `json:"libraryId" validate:"required" description:"库ID"`
}
type UploadMdRes struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}
