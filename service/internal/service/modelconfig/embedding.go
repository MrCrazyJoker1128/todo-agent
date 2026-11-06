package modelconfig

import (
	"context"
	"time"

	"service/internal/dao"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/gogf/gf/v2/frame/g"
)

// EmbeddingDefaults 返回嵌入模型默认值（本地 ollama）。
func EmbeddingDefaults() (baseURL, model string) {
	return "http://localhost:11434", "qwen3-embedding:8b"
}

// EmbeddingCredential 嵌入模型配置。本地 ollama 无密钥，故只含地址与模型名。
type EmbeddingCredential struct {
	BaseURL string
	Model   string
}

// ResolveEmbedding 取当前用户的嵌入配置。
// 无 user_id（后台/系统任务）或该用户未建配置时回退默认，不报错：
// 嵌入是知识库链路的基础能力，缺配置不该让后台任务直接失败。
func ResolveEmbedding(ctx context.Context) *EmbeddingCredential {
	defBaseURL, defModel := EmbeddingDefaults()
	cred := &EmbeddingCredential{BaseURL: defBaseURL, Model: defModel}

	uid, err := resolveUserID(ctx)
	if err != nil {
		return cred
	}
	cfg, err := dao.DbModelConfigs.GetByUserID(ctx, uid)
	if err != nil {
		g.Log().Warning(ctx, "read embedding config failed, fallback to default: %v", err)
		return cred
	}
	if cfg == nil {
		return cred
	}
	if cfg.EmbedBaseUrl != nil && *cfg.EmbedBaseUrl != "" {
		cred.BaseURL = *cfg.EmbedBaseUrl
	}
	if cfg.EmbedModel != nil && *cfg.EmbedModel != "" {
		cred.Model = *cfg.EmbedModel
	}
	return cred
}

// NewEmbedder 按当前用户的嵌入配置创建 ollama embedder。
// timeout 为 0 表示不设超时（与 eino 默认一致）。
func NewEmbedder(ctx context.Context, timeout time.Duration) (*ollama.Embedder, error) {
	cred := ResolveEmbedding(ctx)
	return ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: cred.BaseURL,
		Model:   cred.Model,
		Timeout: timeout,
	})
}
