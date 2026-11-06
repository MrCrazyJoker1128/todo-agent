-- =============================================================================
-- 知识库目录表（用户级）：多个知识库的列表，文件记录只存 ES 不进 PG。
-- 执行方式：docker cp sql/migration_knowledge_bases.sql <容器>:/tmp/mig.sql && docker exec <容器> psql -U postgres -d todo-agent -f /tmp/mig.sql
-- （或本地有 psql 时：psql -d todo-agent -f sql/migration_knowledge_bases.sql）
-- 幂等，可重复执行。
-- db_knowledge_urls 一字不动（旧数据只读兼容，新流程不再写入）。
-- =============================================================================

CREATE TABLE IF NOT EXISTS db_knowledge_bases (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES db_users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE db_knowledge_bases IS '知识库目录表：用户级多个库，文件记录只存 ES';
COMMENT ON COLUMN db_knowledge_bases.id IS '主键 UUID，应用侧赋值';
COMMENT ON COLUMN db_knowledge_bases.user_id IS '归属用户 ID，级联删除';
COMMENT ON COLUMN db_knowledge_bases.name IS '库名';
COMMENT ON COLUMN db_knowledge_bases.description IS '库描述';

CREATE INDEX IF NOT EXISTS idx_kb_user ON db_knowledge_bases(user_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_db_knowledge_bases_updated_at ON db_knowledge_bases;
CREATE TRIGGER trg_db_knowledge_bases_updated_at
    BEFORE UPDATE ON db_knowledge_bases
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();
