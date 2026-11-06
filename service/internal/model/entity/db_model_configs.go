// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 手改：追加 temperature / max_tokens / system_prompt 可空列（见 sql/migration_model_gen_params.sql）、
// embed_base_url / embed_model（见 sql/migration_embed_config.sql）。
// 未跑 `gf gen dao`（plan 模式禁装 CLI）；列增减是标准场景，若后续重新生成，以生成结果为准重做本手改。
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbModelConfigs is the golang structure for table db_model_configs.
type DbModelConfigs struct {
	UserId       uuid.UUID   `json:"userId"       orm:"user_id"        description:"主键/外键：归属用户 ID，级联删除"`                              // 主键/外键：归属用户 ID，级联删除
	ApiKeyCipher string      `json:"apiKeyCipher" orm:"api_key_cipher" description:"加密后的模型 Key（base64），KEK 来自环境变量 MODEL_KEY_ENC_KEY"` // 加密后的模型 Key（base64），KEK 来自环境变量 MODEL_KEY_ENC_KEY
	ApiKeyNonce  string      `json:"apiKeyNonce"  orm:"api_key_nonce"  description:"GCM 随机数（base64），每次保存重新生成"`                        // GCM 随机数（base64），每次保存重新生成
	BaseUrl      *string     `json:"baseUrl"      orm:"base_url"       description:"模型 BaseURL，NULL=运行时用内置默认"`                        // 模型 BaseURL，NULL=运行时用内置默认
	Model        *string     `json:"model"        orm:"model"        description:"模型名，NULL=运行时用内置默认"`                                 // 模型名，NULL=运行时用内置默认
	Temperature  *float64    `json:"temperature"  orm:"temperature"    description:"采样温度，NULL=内置默认 0.7"`                              // 采样温度，NULL=内置默认 0.7
	MaxTokens    *int        `json:"maxTokens"    orm:"max_tokens"     description:"最大输出 token 数，NULL=内置默认 4096"`                     // 最大输出 token 数，NULL=内置默认 4096
	SystemPrompt *string     `json:"systemPrompt" orm:"system_prompt"  description:"自定义系统提示词，NULL/空=内置模板"`                            // 自定义系统提示词，NULL/空=内置模板
	EmbedBaseUrl *string     `json:"embedBaseUrl" orm:"embed_base_url" description:"嵌入模型服务地址，NULL=本地 ollama 默认"`                      // 嵌入模型服务地址，默认本地 ollama
	EmbedModel   *string     `json:"embedModel"   orm:"embed_model"    description:"嵌入模型名，NULL=用默认；更换需重建向量"`                          // 嵌入模型名，默认 qwen3-embedding:8b；更换需重建向量
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     description:""`                                                //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     description:""`                                                //
}
