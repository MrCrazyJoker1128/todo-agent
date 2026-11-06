package knowledge

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"service/api"
	ksvc "service/internal/service/knowledge"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

type KnowledgeController struct{}

func NewKnowledgeController() *KnowledgeController {
	return &KnowledgeController{}
}

func (k *KnowledgeController) AddKnowledge(ctx context.Context, req *api.AddKnowledgeReq) (res *api.AddKnowledgeRes, err error) {
	r := g.RequestFromCtx(ctx)
	userIdStr := r.GetParam("userId").String()
	if userIdStr == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	// 异步入库：先写 pending 占位立即返回，后台抓取+向量化（前端 50s 超时不再掐断）。
	if strings.TrimSpace(req.LibraryId) != "" {
		if _, err := ksvc.AddHtmlFileAsync(ctx, userIdStr, req.LibraryId, req.Url); err != nil {
			return nil, err
		}
		return &api.AddKnowledgeRes{Id: uuid.Nil}, nil
	}
	urlId, contentData, hashStr, now, err := ksvc.AddKnowledge(ctx, userIdStr, req)
	if err != nil {
		return nil, err
	}
	go ksvc.ProcessChunksAsync(urlId, userId, contentData, req.Url, hashStr, now)
	return &api.AddKnowledgeRes{Id: urlId}, nil
}

func (k *KnowledgeController) GetKnowledgeList(ctx context.Context, req *api.GetKnowledgeListReq) (res *api.GetKnowledgeListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userIdStr := r.GetParam("userId").String()
	if userIdStr == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	if _, err := uuid.Parse(userIdStr); err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	return ksvc.GetKnowledgeList(ctx, userIdStr, req)
}

func (k *KnowledgeController) DeleteKnowledgeUrl(ctx context.Context, req *api.DeleteKnowledgeUrlReq) (res *api.DeleteKnowledgeUrlRes, err error) {
	urlId, err := uuid.Parse(req.Id)
	if err != nil {
		g.Log().Error(ctx, "parse url_id failed: %v", err)
		return nil, fmt.Errorf("invalid url_id: %w", err)
	}
	return ksvc.DeleteKnowledgeUrl(ctx, urlId)
}

func (k *KnowledgeController) CreateBase(ctx context.Context, req *api.CreateBaseReq) (res *api.CreateBaseRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	base, err := ksvc.CreateBase(ctx, userID, req.LibraryName, req.Description)
	if err != nil {
		return nil, err
	}
	return &api.CreateBaseRes{Id: base.Id}, nil
}

func (k *KnowledgeController) ListBases(ctx context.Context, req *api.ListBasesReq) (res *api.ListBasesRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	bases, err := ksvc.ListBases(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]*api.KnowledgeBaseItem, 0, len(bases))
	for _, b := range bases {
		var created time.Time
		if b.CreatedAt != nil {
			created = b.CreatedAt.Time
		}
		items = append(items, &api.KnowledgeBaseItem{Id: b.Id, Name: b.Name, Description: b.Description, CreatedAt: created})
	}
	return &api.ListBasesRes{Items: items}, nil
}

func (k *KnowledgeController) DeleteBase(ctx context.Context, req *api.DeleteBaseReq) (res *api.DeleteBaseRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	if err := ksvc.DeleteBase(ctx, userID, req.Id); err != nil {
		return nil, err
	}
	return &api.DeleteBaseRes{Message: "删除知识库成功"}, nil
}

func (k *KnowledgeController) ListFiles(ctx context.Context, req *api.ListFilesReq) (res *api.ListFilesRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	files, err := ksvc.ListFiles(ctx, userID, req.LibraryId)
	if err != nil {
		return nil, err
	}
	items := make([]*api.KnowledgeFileItem, 0, len(files))
	for _, f := range files {
		items = append(items, &api.KnowledgeFileItem{Id: f.ID, Title: f.Title, Source: f.Source, ChunkCount: f.ChunkCount, CreatedAt: f.CreatedAt, Status: f.Status, Error: f.Error})
	}
	return &api.ListFilesRes{Items: items}, nil
}

func (k *KnowledgeController) DeleteFile(ctx context.Context, req *api.DeleteFileReq) (res *api.DeleteFileRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	if err := ksvc.DeleteFile(ctx, userID, req.LibraryId, req.Source); err != nil {
		return nil, err
	}
	return &api.DeleteFileRes{Message: "删除文件成功"}, nil
}

func (k *KnowledgeController) UploadMd(ctx context.Context, req *api.UploadMdReq) (res *api.UploadMdRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetParam("userId").String()
	if userID == "" {
		return nil, fmt.Errorf("user_id not found")
	}
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的 md 文件")
	}
	filename := file.Filename
	lower := strings.ToLower(filename)
	if !strings.HasSuffix(lower, ".md") && !strings.HasSuffix(lower, ".markdown") {
		return nil, fmt.Errorf("只支持 .md 文件")
	}
	if file.Size > 2*1024*1024 {
		return nil, fmt.Errorf("文件超过 2MB")
	}
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	content := strings.TrimSpace(string(buf))
	if content == "" {
		return nil, fmt.Errorf("文件内容为空")
	}
	if !utf8.Valid(buf) {
		return nil, fmt.Errorf("文件不是有效 UTF-8 文本")
	}
	title := strings.TrimSuffix(filename, filepath.Ext(filename))
	source, err := ksvc.AddMdFile(ctx, userID, req.LibraryId, title, content)
	if err != nil {
		return nil, err
	}
	return &api.UploadMdRes{Title: title, Source: source}, nil
}
