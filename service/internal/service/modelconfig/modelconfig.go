package modelconfig

import (
	"context"
	"errors"

	"service/internal/dao"
	"service/internal/service/sysconfig"
	"service/internal/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// ModelCredential 内存中的模型凭证。调用方 MUST 不落盘、不打日志。
type ModelCredential struct {
	APIKey       string
	BaseURL      string
	Model        string
	Temperature  *float32
	MaxTokens    *int
	SystemPrompt string
}

// 生成参数内置默认：Resolve 不自动填，由调用方决定是否使用。
const (
	DefaultTemperature = 0.7
	DefaultMaxTokens   = 4096
)

func defaults() (baseURL, model string) {
	return "https://api.minimaxi.com/v1", "MiniMax-M2.7"
}

func resolveUserID(ctx context.Context) (uuid.UUID, error) {
	if v := ctx.Value("user_id"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return uuid.Parse(s)
		}
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		if s := r.GetParam("userId").String(); s != "" {
			return uuid.Parse(s)
		}
	}
	return uuid.Nil, errors.New("user_id not found in context")
}

// Resolve 解密当前用户的模型 Key。无配置返回 ErrNotConfigured，上层转 400 引导去设置页。
var ErrNotConfigured = errors.New("模型 Key 未配置，请先在设置页配置您的模型 Key")

func Resolve(ctx context.Context) (*ModelCredential, error) {
	uid, err := resolveUserID(ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := dao.DbModelConfigs.GetByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, ErrNotConfigured
	}
	dbKEK, err := sysconfig.GetKEK(ctx)
	if err != nil {
		return nil, err
	}
	plain, err := utility.DecryptModelKeyWithKEK(ctx, cfg.ApiKeyCipher, cfg.ApiKeyNonce, dbKEK)
	if err != nil {
		return nil, err
	}
	defBaseURL, defModel := defaults()
	baseURL := defBaseURL
	if cfg.BaseUrl != nil && *cfg.BaseUrl != "" {
		baseURL = *cfg.BaseUrl
	}
	model := defModel
	if cfg.Model != nil && *cfg.Model != "" {
		model = *cfg.Model
	}
	var temperature *float32
	if cfg.Temperature != nil {
		t := float32(*cfg.Temperature)
		temperature = &t
	}
	systemPrompt := ""
	if cfg.SystemPrompt != nil {
		systemPrompt = *cfg.SystemPrompt
	}
	return &ModelCredential{APIKey: plain, BaseURL: baseURL, Model: model, Temperature: temperature, MaxTokens: cfg.MaxTokens, SystemPrompt: systemPrompt}, nil
}
