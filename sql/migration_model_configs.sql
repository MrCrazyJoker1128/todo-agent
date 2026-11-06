-- =============================================================================
-- 模型 Key 配置表（用户级，AES-256-GCM 密文存储）
-- 执行方式：psql -d todo-agent -f sql/migration_model_configs.sql
-- 前置：依赖 schema.sql 中的 db_update_updated_at() 函数（新库先跑 schema.sql）。
-- 幂等，可重复执行。
-- =============================================================================

CREATE TABLE IF NOT EXISTS db_model_configs (
    user_id        UUID PRIMARY KEY REFERENCES db_users(id) ON DELETE CASCADE,
    api_key_cipher TEXT NOT NULL,
    api_key_nonce  TEXT NOT NULL,
    base_url       VARCHAR(255) NOT NULL DEFAULT 'https://api.minimaxi.com/v1',
    model          VARCHAR(100) NOT NULL DEFAULT 'MiniMax-M2.7',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE db_model_configs IS '模型配置表：每用户一行，api_key 为 AES-256-GCM 密文（base64），nonce 为随机数（base64）';
COMMENT ON COLUMN db_model_configs.user_id IS '主键/外键：归属用户 ID，级联删除';
COMMENT ON COLUMN db_model_configs.api_key_cipher IS '加密后的模型 Key（base64），KEK 来自环境变量 MODEL_KEY_ENC_KEY';
COMMENT ON COLUMN db_model_configs.api_key_nonce IS 'GCM 随机数（base64），每次保存重新生成';
COMMENT ON COLUMN db_model_configs.base_url IS '模型 BaseURL，默认 MiniMax OpenAI 兼容地址';
COMMENT ON COLUMN db_model_configs.model IS '模型名，默认 MiniMax-M2.7';

DROP TRIGGER IF EXISTS trg_db_model_configs_updated_at ON db_model_configs;
CREATE TRIGGER trg_db_model_configs_updated_at
    BEFORE UPDATE ON db_model_configs
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();
