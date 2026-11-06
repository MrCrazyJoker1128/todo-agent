// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 手改：追加 Temperature / MaxTokens / SystemPrompt（any 表可空，见 sql/migration_model_gen_params.sql）、
// EmbedBaseUrl / EmbedModel（见 sql/migration_embed_config.sql）。
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbModelConfigs is the golang structure of table db_model_configs for DAO operations like Where/Data.
type DbModelConfigs struct {
	g.Meta       `orm:"table:db_model_configs, do:true"`
	UserId       any         // 主键/外键：归属用户 ID，级联删除
	ApiKeyCipher any         // 加密后的模型 Key（base64），KEK 来自环境变量 MODEL_KEY_ENC_KEY
	ApiKeyNonce  any         // GCM 随机数（base64），每次保存重新生成
	BaseUrl      any         // 模型 BaseURL，默认 MiniMax OpenAI 兼容地址
	Model        any         // 模型名，默认 MiniMax-M2.7
	Temperature  any         // 采样温度，NULL=内置默认 0.7
	MaxTokens    any         // 最大输出 token 数，NULL=内置默认 4096
	SystemPrompt any         // 自定义系统提示词，NULL/空=内置模板
	EmbedBaseUrl any         // 嵌入模型服务地址，NULL=本地 ollama 默认
	EmbedModel   any         // 嵌入模型名，NULL=用默认
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
}
