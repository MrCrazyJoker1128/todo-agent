-- =============================================================================
-- 嵌入模型配置列改为可空：未配置存 NULL（不落库），运行时回退本地 ollama 默认。
-- 执行方式：psql -d todo-agent -f sql/migration_embed_nullable.sql
-- 幂等，可重复执行。
-- 注意：已存的默认值字符串保持不动（视为用户显式配置）；新保存空值写 NULL。
-- =============================================================================

ALTER TABLE db_model_configs ALTER COLUMN embed_base_url DROP NOT NULL;
ALTER TABLE db_model_configs ALTER COLUMN embed_base_url DROP DEFAULT;
ALTER TABLE db_model_configs ALTER COLUMN embed_model DROP NOT NULL;
ALTER TABLE db_model_configs ALTER COLUMN embed_model DROP DEFAULT;

COMMENT ON COLUMN db_model_configs.embed_base_url IS '嵌入模型服务地址，NULL=运行时用本地 ollama 默认；不存默认字符串';
COMMENT ON COLUMN db_model_configs.embed_model IS '嵌入模型名，NULL=运行时用默认；更换需重建知识库向量';
