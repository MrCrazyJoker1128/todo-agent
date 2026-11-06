package knowledge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"service/internal/service/modelconfig"
	"service/internal/tools"
	"service/internal/utility"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// AddKnowledge extracts content, dedups, inserts the url record and returns
// the pieces the controller needs to launch async chunk processing.
func AddKnowledge(ctx context.Context, userID string, req *api.AddKnowledgeReq) (uuid.UUID, tools.ExtractContentResult, string, *gtime.Time, error) {
	var empty tools.ExtractContentResult
	extractTool := tools.ExtractContentTool()
	extractResult, err := extractTool.InvokableRun(ctx, fmt.Sprintf(`{"url":"%s"}`, req.Url))
	if err != nil {
		g.Log().Error(ctx, "extract content failed: %v", err)
		return uuid.Nil, empty, "", nil, fmt.Errorf("extract content failed: %w", err)
	}
	var contentData tools.ExtractContentResult
	if err := json.Unmarshal([]byte(extractResult), &contentData); err != nil {
		g.Log().Error(ctx, "parse extract result failed: %v", err)
		return uuid.Nil, empty, "", nil, fmt.Errorf("parse extract result failed: %w", err)
	}
	domain := extractDomain(req.Url)
	existingUrl, err := dao.DbKnowledgeUrls.GetUrlByUserAndUrl(ctx, userID, req.Url)
	if err != nil {
		g.Log().Debug(ctx, "check existing url: %v", err)
	}
	if existingUrl != nil {
		return uuid.Nil, empty, "", nil, fmt.Errorf("该知识库已存在，请勿重复添加")
	}
	userId, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, empty, "", nil, fmt.Errorf("invalid user_id: %w", err)
	}
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
		Status:       "pending",
		ErrorMessage: "",
		IsDeleted:    false,
		DeletedAt:    nil,
		CrawledAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := dao.DbKnowledgeUrls.InsertUrl(ctx, urlRecord); err != nil {
		g.Log().Error(ctx, "insert url record failed: %v", err)
		return uuid.Nil, empty, "", nil, fmt.Errorf("insert url record failed: %w", err)
	}
	return urlId, contentData, hashStr, now, nil
}

// ProcessChunksAsync splits content, persists chunks and indexes them to ES.
func ProcessChunksAsync(urlId uuid.UUID, userId uuid.UUID, contentData tools.ExtractContentResult, pageUrl, hashStr string, now *gtime.Time) {
	ctx := context.Background()
	updateUrlStatus(ctx, urlId, "processing", "")
	g.Log().Info(ctx, "processChunksAsync: started processing urlId=%s", urlId.String())
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{ChunkSize: 1500, OverlapSize: 300})
	if err != nil {
		g.Log().Error(ctx, "create splitter failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		g.Log().Info(ctx, "processChunksAsync: failed at splitter urlId=%s, error=%v", urlId.String(), err)
		return
	}
	docs := []*schema.Document{{
		ID:      urlId.String(),
		Content: contentData.Content,
		MetaData: map[string]any{
			"title":  contentData.Title,
			"author": contentData.Author,
			"url":    pageUrl,
		},
	}}
	chunks, err := splitter.Transform(ctx, docs)
	if err != nil {
		g.Log().Error(ctx, "split document failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		g.Log().Info(ctx, "processChunksAsync: failed at split urlId=%s, error=%v", urlId.String(), err)
		return
	}
	embedder, err := modelconfig.NewEmbedder(ctx, 10*time.Minute)
	if err != nil {
		g.Log().Error(ctx, "create embedder failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		g.Log().Info(ctx, "processChunksAsync: failed at embedder urlId=%s, error=%v", urlId.String(), err)
		return
	}
	esClient, err := utility.ES()
	if err != nil {
		g.Log().Error(ctx, "get es client failed: %v", err)
		updateUrlStatus(ctx, urlId, "failed", err.Error())
		g.Log().Info(ctx, "processChunksAsync: failed at ES client urlId=%s, error=%v", urlId.String(), err)
		return
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
		g.Log().Info(ctx, "processChunksAsync: failed at ES indexer urlId=%s, error=%v", urlId.String(), err)
		return
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
		if err := dao.DbKnowledgeChunks.InsertChunk(ctx, chunkRecord); err != nil {
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
		if _, err := esIndexer.Store(ctx, chunkDocs); err != nil {
			g.Log().Error(ctx, "index chunks to es failed: %v", err)
			updateUrlStatus(ctx, urlId, "failed", err.Error())
			g.Log().Info(ctx, "processChunksAsync: failed urlId=%s, error=%v", urlId.String(), err)
			return
		}
	}
	updateUrlStatus(ctx, urlId, "completed", "")
	g.Log().Info(ctx, "processChunksAsync: completed urlId=%s", urlId.String())
}

// GetKnowledgeList returns the paged knowledge url list for a user.
func GetKnowledgeList(ctx context.Context, userID string, req *api.GetKnowledgeListReq) (*api.GetKnowledgeListRes, error) {
	res, err := dao.DbKnowledgeUrls.GetKnowledgeList(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("get knowledge list failed: %w", err)
	}
	return res, nil
}

// DeleteKnowledgeUrl soft-deletes a knowledge url record.
func DeleteKnowledgeUrl(ctx context.Context, urlID uuid.UUID) (*api.DeleteKnowledgeUrlRes, error) {
	if err := dao.DbKnowledgeUrls.DeleteKnowledgeUrl(ctx, urlID); err != nil {
		g.Log().Error(ctx, "delete knowledge url failed: %v", err)
		return nil, fmt.Errorf("delete knowledge url failed: %w", err)
	}
	return &api.DeleteKnowledgeUrlRes{Message: "删除知识库成功"}, nil
}

// CreateBase 新建知识库。
func CreateBase(ctx context.Context, userID, name, description string) (*entity.DbKnowledgeBases, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("库名不能为空")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	base := &entity.DbKnowledgeBases{
		Id:          uuid.New(),
		UserId:      uid,
		Name:        name,
		Description: strings.TrimSpace(description),
	}
	if err := dao.DbKnowledgeBases.InsertBase(ctx, base); err != nil {
		return nil, err
	}
	return base, nil
}

// ListBases 按用户取知识库列表。
func ListBases(ctx context.Context, userID string) ([]*entity.DbKnowledgeBases, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}
	return dao.DbKnowledgeBases.ListByUserID(ctx, uid)
}

// DeleteBase 删库：先清该库全部 ES 文档再删 PG 库行，ES 失败则不删库。
func DeleteBase(ctx context.Context, userID, id string) error {
	baseId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %w", err)
	}
	// 库归属校验：无此库直接报错。
	bases, err := dao.DbKnowledgeBases.ListByUserID(ctx, uid)
	if err != nil {
		return err
	}
	found := false
	for _, b := range bases {
		if b.Id == baseId {
			found = true
			break
		}
	}
	if !found {
		return errors.New("知识库不存在")
	}
	if err := deleteLibraryDocs(ctx, userID, baseId.String()); err != nil {
		return fmt.Errorf("清理库文件失败: %w", err)
	}
	return dao.DbKnowledgeBases.DeleteBase(ctx, baseId, uid)
}

func updateUrlStatus(ctx context.Context, urlId uuid.UUID, status, errMsg string) {
	dao.DbKnowledgeUrls.UpdateStatus(ctx, urlId.String(), status, errMsg)
}

// FileItem 库文件列表项。
type FileItem struct {
	ID         string
	Title      string
	Source     string
	ChunkCount int
	CreatedAt  string
	Status     string
	Error      string
}

// checkLibrary 校验库归属，返回库 id。
func checkLibrary(ctx context.Context, userID, libraryID string) (uuid.UUID, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id: %w", err)
	}
	lid, err := uuid.Parse(libraryID)
	if err != nil {
		return uuid.Nil, errors.New("库ID非法")
	}
	bases, err := dao.DbKnowledgeBases.ListByUserID(ctx, uid)
	if err != nil {
		return uuid.Nil, err
	}
	for _, b := range bases {
		if b.Id == lid {
			return lid, nil
		}
	}
	return uuid.Nil, errors.New("知识库不存在")
}

// mdSource md 文件身份标识。
func mdSource(libraryID, content string) string {
	sum := sha256.Sum256([]byte(content))
	return "md://" + libraryID + "/" + hex.EncodeToString(sum[:])[:32]
}

// storeChunksToES 切分→向量→写 ES（复用 ProcessChunksAsync 三件套：splitter/embedder/Store，不另起）。
func storeChunksToES(ctx context.Context, userID, libraryID, title, source, content string) error {
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{ChunkSize: 1500, OverlapSize: 300})
	if err != nil {
		return fmt.Errorf("create splitter failed: %w", err)
	}
	docs := []*schema.Document{{
		ID:      uuid.NewString(),
		Content: content,
		MetaData: map[string]any{
			"title":     title,
			"source":    source,
			"libraryId": libraryID,
			"userId":    userID,
		},
	}}
	chunks, err := splitter.Transform(ctx, docs)
	if err != nil {
		return fmt.Errorf("split document failed: %w", err)
	}
	embedder, err := modelconfig.NewEmbedder(ctx, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("create embedder failed: %w", err)
	}
	esClient, err := utility.ES()
	if err != nil {
		return fmt.Errorf("get es client failed: %w", err)
	}
	esIndexer, err := es8.NewIndexer(ctx, &es8.IndexerConfig{
		Client:    esClient,
		Index:     "knowledge_chunks",
		BatchSize: 1,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (map[string]es8.FieldValue, error) {
			md := doc.MetaData
			return map[string]es8.FieldValue{
				"content":     {Value: doc.Content, EmbedKey: "content_vector"},
				"library_id":  {Value: md["libraryId"]},
				"user_id":     {Value: md["userId"]},
				"title":       {Value: md["title"]},
				"source":      {Value: md["source"]},
				"chunk_index": {Value: md["chunkIndex"]},
			}, nil
		},
		Embedding: embedder,
	})
	if err != nil {
		return fmt.Errorf("create es indexer failed: %w", err)
	}
	chunkDocs := make([]*schema.Document, 0, len(chunks))
	for i, chunk := range chunks {
		chunkDocs = append(chunkDocs, &schema.Document{
			ID:      uuid.NewString(),
			Content: chunk.Content,
			MetaData: map[string]any{
				"libraryId":  libraryID,
				"userId":     userID,
				"title":      title,
				"source":     source,
				"chunkIndex": i,
			},
		})
	}
	if len(chunkDocs) > 0 {
		if _, err := esIndexer.Store(ctx, chunkDocs); err != nil {
			return fmt.Errorf("index chunks to es failed: %w", err)
		}
	}
	return nil
}

// AddMdFile 入库 md 全文到指定库（不经抓取与门控）。异步：先写 pending 占位立即返回，后台切分+向量化。
func AddMdFile(ctx context.Context, userID, libraryID, title, content string) (string, error) {
	if _, err := checkLibrary(ctx, userID, libraryID); err != nil {
		return "", err
	}
	source := mdSource(libraryID, content)
	writeStatusDoc(context.Background(), userID, libraryID, title, source, "pending", "")
	go func() {
		bg := context.WithValue(context.Background(), "user_id", userID)
		if err := storeChunksToES(bg, userID, libraryID, title, source, content); err != nil {
			g.Log().Error(bg, "async md store failed: %v", err)
			writeStatusDoc(bg, userID, libraryID, title, source, "failed", err.Error())
			return
		}
		writeStatusDoc(bg, userID, libraryID, title, source, "completed", "")
	}()
	return source, nil
}

// AddHtmlFileAsync 异步入库 html：先写 pending 占位立即返回，后台抓取+切分+向量化。
// ctx 必须是请求 ctx（取 userID/libraryID 用）；后台用 context.Background() detach，不受前端超时取消影响。
func AddHtmlFileAsync(ctx context.Context, userID, libraryID, pageURL string) (string, error) {
	if _, err := checkLibrary(ctx, userID, libraryID); err != nil {
		return "", err
	}
	title := pageURL
	writeStatusDoc(context.Background(), userID, libraryID, title, pageURL, "pending", "")
	go func() {
		// 后台 ctx 必须带 user_id（modelconfig.Resolve 与工具链从 ctx 取用户）；detach 不受前端超时取消影响。
		bg := context.WithValue(context.Background(), "user_id", userID)
		extractTool := tools.ExtractContentTool()
		result, err := extractTool.InvokableRun(bg, fmt.Sprintf(`{"url":%q}`, pageURL))
		if err != nil {
			g.Log().Error(bg, "async extract failed: %v", err)
			writeStatusDoc(bg, userID, libraryID, title, pageURL, "failed", err.Error())
			return
		}
		var data tools.ExtractContentResult
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			g.Log().Error(bg, "async parse failed: %v", err)
			writeStatusDoc(bg, userID, libraryID, title, pageURL, "failed", err.Error())
			return
		}
		if data.Title != "" {
			title = data.Title
		}
		if err := storeChunksToES(bg, userID, libraryID, title, pageURL, data.Content); err != nil {
			g.Log().Error(bg, "async store failed: %v", err)
			writeStatusDoc(bg, userID, libraryID, title, pageURL, "failed", err.Error())
			return
		}
		writeStatusDoc(bg, userID, libraryID, title, pageURL, "completed", "")
	}()
	return pageURL, nil
}

// writeStatusDoc 写/更新文件状态占位 doc（id 固定，保证幂等覆盖）。
func writeStatusDoc(ctx context.Context, userID, libraryID, title, source, status, errMsg string) {
	id := "status-" + hex.EncodeToString([]byte(source))[:32]
	doc := map[string]any{
		"library_id":  libraryID,
		"user_id":     userID,
		"title":       title,
		"source":      source,
		"chunk_index": -1,
		"content":     "",
		"status":      status,
		"error":       errMsg,
	}
	if err := utility.ESIndex(ctx, "knowledge_chunks", id, doc); err != nil {
		g.Log().Error(ctx, "write status doc failed: %v", err)
	}
}

// AddHtmlFile 抓取 html 后入库到指定库（同步版，保留给直接调用方）。
func AddHtmlFile(ctx context.Context, userID, libraryID, pageURL string) (string, error) {
	if _, err := checkLibrary(ctx, userID, libraryID); err != nil {
		return "", err
	}
	extractTool := tools.ExtractContentTool()
	result, err := extractTool.InvokableRun(ctx, fmt.Sprintf(`{"url":%q}`, pageURL))
	if err != nil {
		return "", fmt.Errorf("抓取失败: %w", err)
	}
	var data tools.ExtractContentResult
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return "", fmt.Errorf("解析抓取结果失败: %w", err)
	}
	title := data.Title
	if title == "" {
		title = pageURL
	}
	if err := storeChunksToES(ctx, userID, libraryID, title, pageURL, data.Content); err != nil {
		return "", err
	}
	return pageURL, nil
}

// ListFiles 按库列文件（ES size 500 拉回内存按 source 去重）。
func ListFiles(ctx context.Context, userID, libraryID string) ([]*FileItem, error) {
	if _, err := checkLibrary(ctx, userID, libraryID); err != nil {
		return nil, err
	}
	query := map[string]any{
		"size": 500,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{"term": map[string]any{"user_id.keyword": userID}},
					{"term": map[string]any{"library_id.keyword": libraryID}},
				},
			},
		},
		"_source": []string{"title", "source", "status", "error", "chunk_index"},
		"sort":    []map[string]any{{"chunk_index": map[string]any{"order": "asc"}}},
	}
	raw, err := utility.ESSearch(ctx, "knowledge_chunks", query)
	if err != nil {
		return nil, fmt.Errorf("ES搜索失败: %w", err)
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				Source struct {
					Title      string `json:"title"`
					Source     string `json:"source"`
					ChunkIndex *int   `json:"chunk_index"`
					Status     string `json:"status"`
					Error      string `json:"error"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("解析ES响应失败: %w", err)
	}
	seen := map[string]*FileItem{}
	order := []*FileItem{}
	for _, h := range resp.Hits.Hits {
		if h.Source.Source == "" {
			continue
		}
		// 状态占位 doc（chunk_index=-1 或缺失）：只记状态不计分片。
		if h.Source.ChunkIndex == nil || *h.Source.ChunkIndex == -1 {
			it, ok := seen[h.Source.Source]
			if !ok {
				it = &FileItem{ID: h.Source.Source, Title: h.Source.Title, Source: h.Source.Source, Status: h.Source.Status}
				seen[h.Source.Source] = it
				order = append(order, it)
			} else {
				it.Status = h.Source.Status
				if h.Source.Title != "" {
					it.Title = h.Source.Title
				}
			}
			if h.Source.Error != "" {
				it.Error = h.Source.Error
			}
			continue
		}
		it, ok := seen[h.Source.Source]
		if !ok {
			it = &FileItem{ID: h.Source.Source, Title: h.Source.Title, Source: h.Source.Source}
			seen[h.Source.Source] = it
			order = append(order, it)
		}
		it.ChunkCount++
	}
	return order, nil
}

// DeleteFile 删除库内一个文件（source）的全部 chunk。
func DeleteFile(ctx context.Context, userID, libraryID, source string) error {
	if _, err := checkLibrary(ctx, userID, libraryID); err != nil {
		return err
	}
	query := map[string]any{
		"size": 500,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{"term": map[string]any{"user_id.keyword": userID}},
					{"term": map[string]any{"library_id.keyword": libraryID}},
					{"term": map[string]any{"source.keyword": source}},
				},
			},
		},
		"_source": []string{"title"},
	}
	raw, err := utility.ESSearch(ctx, "knowledge_chunks", query)
	if err != nil {
		return fmt.Errorf("ES搜索失败: %w", err)
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return fmt.Errorf("解析ES响应失败: %w", err)
	}
	for _, h := range resp.Hits.Hits {
		if err := utility.ESDelete(ctx, "knowledge_chunks", h.ID); err != nil {
			return fmt.Errorf("删除文档失败: %w", err)
		}
	}
	return nil
}

// deleteLibraryDocs 清该库全部 ES 文档（删库前调用）。
func deleteLibraryDocs(ctx context.Context, userID, libraryID string) error {
	files, err := ListFiles(ctx, userID, libraryID)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := DeleteFile(ctx, userID, libraryID, f.Source); err != nil {
			return err
		}
	}
	return nil
}

func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
