-- =============================================================================
-- 模型生成参数（用户级）：temperature / max_tokens / system_prompt
-- 执行方式：psql -d todo-agent -f sql/migration_model_gen_params.sql
-- 幂等，可重复执行。
-- NULL = 用内置默认（temperature 默认 0.7 / max_tokens 默认 4096 / system_prompt 空用内置模板），
-- 与 base_url/model 的 NOT NULL DEFAULT 风格区分：此处不设 DEFAULT 值。
-- =============================================================================

ALTER TABLE db_model_configs ADD COLUMN IF NOT EXISTS temperature DOUBLE PRECISION NULL;
ALTER TABLE db_model_configs ADD COLUMN IF NOT EXISTS max_tokens INTEGER NULL;
ALTER TABLE db_model_configs ADD COLUMN IF NOT EXISTS system_prompt TEXT NULL;

COMMENT ON COLUMN db_model_configs.temperature IS '采样温度，NULL=用内置默认 0.7，有效范围 [0,2]，0 为有效值（判 nil 而非判 0）';
COMMENT ON COLUMN db_model_configs.max_tokens IS '最大输出 token 数（Chat Completion 语义），NULL=用内置默认 4096；Responses API 模型下被 SDK 忽略';
COMMENT ON COLUMN db_model_configs.system_prompt IS '自定义系统提示词，NULL/空=用内置模板；非空时日期时间上下文块追加其后';
