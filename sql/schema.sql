-- =============================================================================
-- AI TODO Agent 数据库初始化脚本
-- 版本：v1.0
-- 说明：一级任务结构，milvus_vector_id 为 UUID 类型
-- 警告：此脚本会删除现有表，请先备份数据！
-- =============================================================================

-- 启用扩展
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- =============================================================================
-- 1. 删除现有表（级联删除）
-- =============================================================================
DROP TABLE IF EXISTS db_chat_messages CASCADE;
DROP TABLE IF EXISTS db_tasks CASCADE;
DROP TABLE IF EXISTS db_users CASCADE;

-- 删除触发器函数（如果存在）
DROP FUNCTION IF EXISTS db_update_updated_at() CASCADE;

-- 删除视图
DROP VIEW IF EXISTS v_db_active_tasks CASCADE;
DROP VIEW IF EXISTS v_db_pending_vector_sync CASCADE;
DROP VIEW IF EXISTS v_db_today_tasks CASCADE;

-- =============================================================================
-- 2. 用户表
-- =============================================================================
CREATE TABLE db_users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    nickname        VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE db_users IS '用户表：存储系统用户基础信息';
COMMENT ON COLUMN db_users.id IS '主键：用户唯一标识 (UUID)';
COMMENT ON COLUMN db_users.email IS '登录邮箱：唯一索引';
COMMENT ON COLUMN db_users.password_hash IS '密码哈希：bcrypt 加密';
COMMENT ON COLUMN db_users.nickname IS '用户昵称：可选';
COMMENT ON COLUMN db_users.created_at IS '账户创建时间';
COMMENT ON COLUMN db_users.updated_at IS '账户最后更新时间';

-- =============================================================================
-- 3. 任务表（一级结构，milvus_vector_id 为 UUID）
-- =============================================================================
CREATE TABLE db_tasks (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL,
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'todo',
    priority            VARCHAR(20) DEFAULT 'medium',
    due_date            TIMESTAMPTZ,
    milvus_vector_id    UUID DEFAULT gen_random_uuid(),
    task_type           VARCHAR(20) DEFAULT 'work',
    priority_level      VARCHAR(10) DEFAULT 'MEDIUM',
    is_deleted          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_tasks_user FOREIGN KEY (user_id) REFERENCES db_users(id) ON DELETE CASCADE
);

COMMENT ON TABLE db_tasks IS '任务表：存储用户待办事项（一级结构）';
COMMENT ON COLUMN db_tasks.id IS '主键：任务唯一标识 (UUID)';
COMMENT ON COLUMN db_tasks.user_id IS '外键：归属用户 ID';
COMMENT ON COLUMN db_tasks.title IS '任务标题：1-255 字符';
COMMENT ON COLUMN db_tasks.description IS '任务详细描述：可选';
COMMENT ON COLUMN db_tasks.status IS '任务状态：todo|doing|done|cancelled';
COMMENT ON COLUMN db_tasks.priority IS '优先级：low|medium|high';
COMMENT ON COLUMN db_tasks.due_date IS '截止时间：带时区';
COMMENT ON COLUMN db_tasks.milvus_vector_id IS 'Milvus 向量 ID (UUID 类型，与 Milvus 主键对应)';
COMMENT ON COLUMN db_tasks.task_type IS '任务类型：work|personal|project|adhoc';
COMMENT ON COLUMN db_tasks.priority_level IS '优先级级别：HIGH|MEDIUM|LOW';
COMMENT ON COLUMN db_tasks.is_deleted IS '软删除标记：true 表示已删除';
COMMENT ON COLUMN db_tasks.created_at IS '任务创建时间';
COMMENT ON COLUMN db_tasks.updated_at IS '任务最后更新时间';

-- =============================================================================
-- 4. 聊天消息表
-- =============================================================================
CREATE TABLE db_chat_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    role            VARCHAR(20) NOT NULL,
    content         TEXT NOT NULL,
    related_task_id UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_messages_user FOREIGN KEY (user_id) REFERENCES db_users(id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_task FOREIGN KEY (related_task_id) REFERENCES db_tasks(id) ON DELETE SET NULL
);

COMMENT ON TABLE db_chat_messages IS '聊天消息表：存储用户与 AI 的对话记录';
COMMENT ON COLUMN db_chat_messages.id IS '主键：消息唯一标识 (UUID)';
COMMENT ON COLUMN db_chat_messages.user_id IS '外键：归属用户 ID';
COMMENT ON COLUMN db_chat_messages.role IS '消息角色：user|assistant';
COMMENT ON COLUMN db_chat_messages.content IS '消息内容：文本';
COMMENT ON COLUMN db_chat_messages.related_task_id IS '关联任务 ID：如果消息涉及任务';
COMMENT ON COLUMN db_chat_messages.created_at IS '消息创建时间';

-- =============================================================================
-- 5. 触发器函数
-- =============================================================================
CREATE OR REPLACE FUNCTION db_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION db_update_updated_at() IS '通用触发器函数：自动更新 updated_at 字段为当前时间';

-- =============================================================================
-- 6. 创建触发器
-- =============================================================================
CREATE TRIGGER trg_db_users_updated_at
    BEFORE UPDATE ON db_users
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();

CREATE TRIGGER trg_db_tasks_updated_at
    BEFORE UPDATE ON db_tasks
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();

COMMENT ON TRIGGER trg_db_users_updated_at ON db_users IS '用户表更新时间触发器';
COMMENT ON TRIGGER trg_db_tasks_updated_at ON db_tasks IS '任务表更新时间触发器';

-- =============================================================================
-- 7. 创建索引
-- =============================================================================
-- 任务表索引
CREATE INDEX idx_db_tasks_user_id ON db_tasks(user_id);
CREATE INDEX idx_db_tasks_status ON db_tasks(status);
CREATE INDEX idx_db_tasks_priority ON db_tasks(priority_level);
CREATE INDEX idx_db_tasks_due_date ON db_tasks(due_date);
CREATE INDEX idx_db_tasks_deleted ON db_tasks(is_deleted) WHERE is_deleted = FALSE;
CREATE INDEX idx_db_tasks_milvus_id ON db_tasks(milvus_vector_id);
CREATE INDEX idx_db_tasks_type ON db_tasks(task_type);
CREATE INDEX idx_db_tasks_title_trgm ON db_tasks USING gin (title gin_trgm_ops);
CREATE INDEX idx_db_tasks_desc_trgm ON db_tasks USING gin (description gin_trgm_ops);

-- 聊天消息表索引
CREATE INDEX idx_db_messages_user_id ON db_chat_messages(user_id);
CREATE INDEX idx_db_messages_task_id ON db_chat_messages(related_task_id);
CREATE INDEX idx_db_messages_created_at ON db_chat_messages(created_at);

-- 索引备注
COMMENT ON INDEX idx_db_tasks_user_id IS '按用户查询任务（主索引）';
COMMENT ON INDEX idx_db_tasks_status IS '按状态筛选任务';
COMMENT ON INDEX idx_db_tasks_priority IS '按优先级排序任务';
COMMENT ON INDEX idx_db_tasks_due_date IS '按截止时间排序';
COMMENT ON INDEX idx_db_tasks_deleted IS '查询未删除任务（部分索引，性能优化）';
COMMENT ON INDEX idx_db_tasks_milvus_id IS '按 Milvus 向量 ID 查询（UUID 类型）';
COMMENT ON INDEX idx_db_tasks_type IS '按任务类型筛选';
COMMENT ON INDEX idx_db_tasks_title_trgm IS '任务标题模糊搜索索引（pg_trgm）';
COMMENT ON INDEX idx_db_tasks_desc_trgm IS '任务描述模糊搜索索引（pg_trgm）';
COMMENT ON INDEX idx_db_messages_user_id IS '按用户查询消息';
COMMENT ON INDEX idx_db_messages_task_id IS '按关联任务查询消息';
COMMENT ON INDEX idx_db_messages_created_at IS '按时间排序消息';

-- =============================================================================
-- 8. 创建视图
-- =============================================================================
-- 活跃任务视图（未删除且未完成的任务）
CREATE VIEW v_db_active_tasks AS
SELECT 
    id, 
    user_id, 
    title, 
    description, 
    status, 
    priority, 
    priority_level, 
    task_type, 
    due_date, 
    milvus_vector_id, 
    created_at
FROM db_tasks
WHERE is_deleted = FALSE 
  AND status IN ('todo', 'doing');

COMMENT ON VIEW v_db_active_tasks IS '活跃任务视图：仅展示未删除的待办和进行中任务';

-- 今日任务视图
CREATE VIEW v_db_today_tasks AS
SELECT 
    id, 
    user_id, 
    title, 
    priority_level, 
    due_date
FROM db_tasks
WHERE is_deleted = FALSE 
  AND status = 'todo'
  AND (due_date IS NULL OR due_date::date = CURRENT_DATE);

COMMENT ON VIEW v_db_today_tasks IS '今日任务视图：展示今天需要完成的任务';

-- =============================================================================
-- 脚本执行完成
-- =============================================================================
