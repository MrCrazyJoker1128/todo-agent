-- =============================================================================
-- 对话模型配置列改为可空：未填写存 NULL（不落库），运行时回退内置默认。
-- 执行方式：psql -d todo-agent -f sql/migration_model_nullable.sql
-- 幂等，可重复执行。
-- 存量清洗：把“从未手工填写、由旧逻辑自动填入”的 MiniMax 默认值改回 NULL。
-- （运行时对 NULL 与该默认值行为完全一致，故对真实 MiniMax 用户无影响。）
-- =============================================================================

ALTER TABLE db_model_configs ALTER COLUMN base_url DROP NOT NULL;
ALTER TABLE db_model_configs ALTER COLUMN base_url DROP DEFAULT;
ALTER TABLE db_model_configs ALTER COLUMN model DROP NOT NULL;
ALTER TABLE db_model_configs ALTER COLUMN model DROP DEFAULT;

UPDATE db_model_configs SET base_url = NULL WHERE base_url = 'https://api.minimaxi.com/v1';
UPDATE db_model_configs SET model = NULL WHERE model = 'MiniMax-M2.7';

COMMENT ON COLUMN db_model_configs.base_url IS '对话模型 BaseURL，NULL=运行时用内置默认；不存默认字符串';
COMMENT ON COLUMN db_model_configs.model IS '对话模型名，NULL=运行时用内置默认；不存默认字符串';
