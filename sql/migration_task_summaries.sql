-- =============================================================================
-- 工作总结记录表（用户级）：agent 生成的日/周/月总结落库，可按周期检索回看。
-- 执行方式：docker cp sql/migration_task_summaries.sql <容器>:/tmp/mig.sql && docker exec <容器> psql -U postgres -d todo-agent -f /tmp/mig.sql
-- （或本地有 psql 时：psql -d todo-agent -f sql/migration_task_summaries.sql）
-- 幂等，可重复执行。
-- =============================================================================

CREATE TABLE IF NOT EXISTS db_task_summaries (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES db_users(id) ON DELETE CASCADE,
    period     TEXT NOT NULL,
    start_date TIMESTAMPTZ,
    end_date   TIMESTAMPTZ,
    task_type  TEXT NOT NULL DEFAULT '',
    priority   TEXT NOT NULL DEFAULT '',
    title      TEXT NOT NULL DEFAULT '',
    content    TEXT NOT NULL,
    stats      TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE db_task_summaries IS '工作总结记录表：agent 生成的日/周/月总结落库';
COMMENT ON COLUMN db_task_summaries.id IS '主键 UUID，应用侧赋值';
COMMENT ON COLUMN db_task_summaries.user_id IS '归属用户 ID，级联删除';
COMMENT ON COLUMN db_task_summaries.period IS '周期：day|week|month|custom';
COMMENT ON COLUMN db_task_summaries.start_date IS '统计窗口起点（TIMESTAMPTZ）';
COMMENT ON COLUMN db_task_summaries.end_date IS '统计窗口终点（TIMESTAMPTZ）';
COMMENT ON COLUMN db_task_summaries.task_type IS '任务类型过滤，空=全部';
COMMENT ON COLUMN db_task_summaries.priority IS '紧急程度过滤（priority 小写版），空=全部';
COMMENT ON COLUMN db_task_summaries.title IS '总结标题，如 周总结 11.02~11.08';
COMMENT ON COLUMN db_task_summaries.content IS '服务端拼好的总结正文 markdown';
COMMENT ON COLUMN db_task_summaries.stats IS '统计 JSON 字符串（不用 JSONB），如 total/done/rate/byType';

CREATE INDEX IF NOT EXISTS idx_task_summaries_user ON db_task_summaries(user_id, created_at DESC);

DROP TRIGGER IF EXISTS trg_db_task_summaries_updated_at ON db_task_summaries;
CREATE TRIGGER trg_db_task_summaries_updated_at
    BEFORE UPDATE ON db_task_summaries
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();
