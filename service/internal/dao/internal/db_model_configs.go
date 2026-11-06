// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 手改：columns 追加 temperature / max_tokens / system_prompt（见 sql/migration_model_gen_params.sql）。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DbModelConfigsDao is the data access object for the table db_model_configs.
type DbModelConfigsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  DbModelConfigsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// DbModelConfigsColumns defines and stores column names for the table db_model_configs.
type DbModelConfigsColumns struct {
	UserId       string // 主键/外键：归属用户 ID，级联删除
	ApiKeyCipher string // 加密后的模型 Key（base64），KEK 来自环境变量 MODEL_KEY_ENC_KEY
	ApiKeyNonce  string // GCM 随机数（base64），每次保存重新生成
	BaseUrl      string // 模型 BaseURL，默认 MiniMax OpenAI 兼容地址
	Model        string // 模型名，默认 MiniMax-M2.7
	Temperature  string // 采样温度，NULL=内置默认 0.7
	MaxTokens    string // 最大输出 token 数，NULL=内置默认 4096
	SystemPrompt string // 自定义系统提示词，NULL/空=内置模板
	EmbedBaseUrl string // 嵌入模型服务地址，NULL=本地 ollama 默认
	EmbedModel   string // 嵌入模型名，NULL=用默认
	CreatedAt    string //
	UpdatedAt    string //
}

// dbModelConfigsColumns holds the columns for the table db_model_configs.
var dbModelConfigsColumns = DbModelConfigsColumns{
	UserId:       "user_id",
	ApiKeyCipher: "api_key_cipher",
	ApiKeyNonce:  "api_key_nonce",
	BaseUrl:      "base_url",
	Model:        "model",
	Temperature:  "temperature",
	MaxTokens:    "max_tokens",
	SystemPrompt: "system_prompt",
	EmbedBaseUrl: "embed_base_url",
	EmbedModel:   "embed_model",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewDbModelConfigsDao creates and returns a new DAO object for table data access.
func NewDbModelConfigsDao(handlers ...gdb.ModelHandler) *DbModelConfigsDao {
	return &DbModelConfigsDao{
		group:    "default",
		table:    "db_model_configs",
		columns:  dbModelConfigsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DbModelConfigsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DbModelConfigsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DbModelConfigsDao) Columns() DbModelConfigsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DbModelConfigsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DbModelConfigsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DbModelConfigsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
