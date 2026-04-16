package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"service/internal/tools"
	"service/internal/utility"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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

	extractTool := tools.ExtractContentTool()
	extractResult, err := extractTool.InvokableRun(ctx, fmt.Sprintf(`{"url":"%s"}`, req.Url))
	if err != nil {
		g.Log().Error(ctx, "extract content failed: %v", err)
		return nil, fmt.Errorf("extract content failed: %w", err)
	}

	var contentData tools.ExtractContentResult
	if err := json.Unmarshal([]byte(extractResult), &contentData); err != nil {
		g.Log().Error(ctx, "parse extract result failed: %v", err)
		return nil, fmt.Errorf("parse extract result failed: %w", err)
	}

	domain := extractDomain(req.Url)

	urlId := uuid.New()
	now := gtime.Now()

	hashStr := fmt.Sprintf("%s-%d", req.Url, time.Now().UnixNano())

	urlRecord := &entity.DbKnowledgeUrls{
		Id:           urlId,
		UserId:       userId,
		Url:          req.Url,
		Title:        contentData.Title,
		Description:  contentData.Description,
		ContentType:  "article",
		Domain:       domain,
		FaviconUrl:   "",
		Author:       contentData.Author,
		PublishDate:  nil,
		WordCount:    len(contentData.Content),
		ChunkCount:   0,
		ContentHash:  hashStr,
		Status:       "processing",
		ErrorMessage: "",
		IsDeleted:    false,
		DeletedAt:    nil,
		CrawledAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = dao.DbKnowledgeUrls.InsertUrl(ctx, urlRecord)
	if err != nil {
		g.Log().Error(ctx, "insert url record failed: %v", err)
		return nil, fmt.Errorf("insert url record failed: %w", err)
	}

	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   1500,
		OverlapSize: 300,
	})
	if err != nil {
		g.Log().Error(ctx, "create splitter failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		return nil, fmt.Errorf("create splitter failed: %w", err)
	}

	docs := []*schema.Document{
		{
			ID:      urlId.String(),
			Content: contentData.Content,
			MetaData: map[string]any{
				"title":  contentData.Title,
				"author": contentData.Author,
				"url":    req.Url,
			},
		},
	}

	chunks, err := splitter.Transform(ctx, docs)
	if err != nil {
		g.Log().Error(ctx, "split document failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		return nil, fmt.Errorf("split document failed: %w", err)
	}

	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3-embedding:8b",
		Timeout: 120 * time.Second,
	})
	if err != nil {
		g.Log().Error(ctx, "create embedder failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		return nil, fmt.Errorf("create embedder failed: %w", err)
	}

	esClient, err := utility.ES()
	if err != nil {
		g.Log().Error(ctx, "get es client failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		return nil, fmt.Errorf("get es client failed: %w", err)
	}

	esIndexer, err := es8.NewIndexer(ctx, &es8.IndexerConfig{
		Client:    esClient,
		Index:     "knowledge_chunks",
		BatchSize: 1,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (map[string]es8.FieldValue, error) {
			metadata := doc.MetaData
			return map[string]es8.FieldValue{
				"content":     {Value: doc.Content, EmbedKey: "content_vector"},
				"url_id":      {Value: metadata["urlId"]},
				"chunk_index": {Value: metadata["chunkIndex"]},
				"user_id":     {Value: metadata["userId"]},
				"author":      {Value: metadata["author"]},
				"upload_time": {Value: metadata["uploadTime"]},
				"title":       {Value: metadata["title"]},
			}, nil
		},
		Embedding: embedder,
	})
	if err != nil {
		g.Log().Error(ctx, "create es indexer failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		return nil, fmt.Errorf("create es indexer failed: %w", err)
	}

	chunkDocs := make([]*schema.Document, 0, len(chunks))
	for i, chunk := range chunks {
		metadata := map[string]any{
			"urlId":      urlId.String(),
			"userId":     userId.String(),
			"chunkIndex": i,
			"author":     contentData.Author,
			"uploadTime": now.Time.Format("2006-01-02T15:04:05Z"),
			"title":      contentData.Title,
		}
		metadataJSON, _ := json.Marshal(metadata)

		chunkId := uuid.New()
		chunkRecord := &entity.DbKnowledgeChunks{
			Id:               chunkId,
			UrlId:            urlId,
			UserId:           userId,
			ChunkIndex:       i,
			Content:          chunk.Content,
			ContentHash:      fmt.Sprintf("%s-%d", hashStr, i),
			WordCount:        len(chunk.Content),
			MilvusVectorId:   uuid.Nil,
			VectorSyncStatus: "pending",
			Metadata:         string(metadataJSON),
			IsDeleted:        false,
			CreatedAt:        now,
		}
		err = dao.DbKnowledgeChunks.InsertChunk(ctx, chunkRecord)
		if err != nil {
			g.Log().Error(ctx, "insert chunk failed: %v", err)
			continue
		}

		chunkDocs = append(chunkDocs, &schema.Document{
			ID:       chunkId.String(),
			Content:  chunk.Content,
			MetaData: metadata,
		})
	}

	if len(chunkDocs) > 0 {
		_, err = esIndexer.Store(ctx, chunkDocs)
		if err != nil {
			g.Log().Error(ctx, "index chunks to es failed: %v", err)
		}
	}

	updateUrlStatus(ctx, urlId, "completed", "")

	res = &api.AddKnowledgeRes{
		Id: urlId,
	}
	return res, nil
}

func (k *KnowledgeController) GetKnowledgeList(ctx context.Context, req *api.GetKnowledgeListReq) (res *api.GetKnowledgeListRes, err error) {

	return nil, nil
}

func updateUrlStatus(ctx context.Context, urlId uuid.UUID, status, errMsg string) {
	dao.DbKnowledgeUrls.UpdateStatus(ctx, urlId.String(), status, errMsg)
}

func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
