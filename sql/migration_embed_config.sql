-- =============================================================================
-- 嵌入模型配置（用户级）：embed_base_url / embed_model
-- 执行方式：psql -d todo-agent -f sql/migration_embed_config.sql
-- 前置：db_model_configs 已存在（migration_model_configs.sql）。
-- 幂等，可重复执行。
-- 默认值为本地 ollama；插入/更新时由应用层兜底空串，故此处设 NOT NULL DEFAULT。
-- 注意：更换嵌入模型会改变向量维度与语义空间，已入库向量必须清空重建，
--      否则检索要么报维度不匹配，要么召回语义错乱。
-- =============================================================================

ALTER TABLE db_model_configs
  ADD COLUMN IF NOT EXISTS embed_base_url VARCHAR(255) NOT NULL DEFAULT 'http://localhost:11434';
ALTER TABLE db_model_configs
  ADD COLUMN IF NOT EXISTS embed_model VARCHAR(100) NOT NULL DEFAULT 'qwen3-embedding:8b';

COMMENT ON COLUMN db_model_configs.embed_base_url IS '嵌入模型服务地址，默认本地 ollama http://localhost:11434';
COMMENT ON COLUMN db_model_configs.embed_model IS '嵌入模型名，默认 qwen3-embedding:8b；更换后需重建知识库向量';
