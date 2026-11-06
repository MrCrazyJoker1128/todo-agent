-- =============================================================================
-- 系统配置表（key-value）：运维级密钥入库，替代散落的环境变量/配置文件。
-- 执行方式：docker cp sql/migration_system_configs.sql <容器>:/tmp/mig.sql && docker exec <容器> psql -U postgres -d todo-agent -f /tmp/mig.sql
-- （或本地有 psql 时：psql -d todo-agent -f sql/migration_system_configs.sql）
-- 幂等，可重复执行。
-- 首个内置项：model_key_enc_key —— 模型 Key 的 KEK。NULL/缺行 = 回退到
-- 环境变量 MODEL_KEY_ENC_KEY 或 config.yaml modelKey.encKey（见 utility/modelkey.go）。
-- 注意：KEK 本人入库是明文——防的是散落配置，不是防 DBA；生产仍建议环境变量注入。
-- =============================================================================

CREATE TABLE IF NOT EXISTS db_system_configs (
    key        TEXT PRIMARY KEY,
    value      TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE db_system_configs IS '系统配置表：key-value，运维级密钥入库（如模型 Key 的 KEK）';
COMMENT ON COLUMN db_system_configs.key IS '配置键，如 model_key_enc_key';
COMMENT ON COLUMN db_system_configs.value IS '配置值，明文存储';
COMMENT ON COLUMN db_system_configs.updated_at IS '最后更新时间，由触发器维护';

DROP TRIGGER IF EXISTS trg_db_system_configs_updated_at ON db_system_configs;
CREATE TRIGGER trg_db_system_configs_updated_at
    BEFORE UPDATE ON db_system_configs
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();
