# 数据库设计文档

## 1. 设计概述

- **数据库类型**：PostgreSQL
- **向量数据库**：Milvus（独立部署）
- **设计原则**：一级任务结构（无父子任务），Milvus Vector ID 存储在任务表

---

## 2. ER 关系图

```
┌─────────────┐
│  db_users   │ ← 用户基础信息
└──────┬──────┘
       │ (user_id 外键)
       ▼
┌─────────────┐      ┌─────────────┐
│  db_tasks   │ ←──→ │   Milvus    │
└─────────────┘      └─────────────┘
       │
       │ (related_task_id 外键)
       ▼
┌───────────────────┐
│ db_chat_messages  │
└───────────────────┘
```

---

## 3. 表结构

### 3.1 用户表 (db_users)

```sql
CREATE TABLE db_users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    nickname        VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键，用户唯一标识 |
| email | VARCHAR(255) | 登录邮箱，唯一索引 |
| password_hash | VARCHAR(255) | 密码哈希 (bcrypt) |
| nickname | VARCHAR(100) | 用户昵称，可选 |
| created_at | TIMESTAMPTZ | 账户创建时间 |
| updated_at | TIMESTAMPTZ | 最后更新时间 |

### 3.2 任务表 (db_tasks)

```sql
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
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键，任务唯一标识 |
| user_id | UUID | 外键，归属用户 ID |
| title | VARCHAR(255) | 任务标题，1-255 字符 |
| description | TEXT | 任务详细描述 |
| status | VARCHAR(20) | 状态：todo, doing, done, cancelled |
| priority | VARCHAR(20) | 优先级：low, medium, high |
| due_date | TIMESTAMPTZ | 截止时间 |
| **milvus_vector_id** | **UUID** | Milvus 向量 ID（与 Milvus 主键对应） |
| task_type | VARCHAR(20) | 类型：work, personal, project, adhoc |
| priority_level | VARCHAR(10) | 优先级级别：HIGH, MEDIUM, LOW |
| is_deleted | BOOLEAN | 软删除标记 |
| created_at | TIMESTAMPTZ | 创建时间 |
| updated_at | TIMESTAMPTZ | 更新时间 |

### 3.3 聊天消息表 (db_chat_messages)

```sql
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
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | UUID | 主键，消息唯一标识 |
| user_id | UUID | 外键，归属用户 ID |
| role | VARCHAR(20) | 角色：user, assistant |
| content | TEXT | 消息内容 |
| related_task_id | UUID | 关联任务 ID（如果消息涉及任务） |
| created_at | TIMESTAMPTZ | 创建时间 |

---

## 4. 索引设计

### 4.1 任务表索引

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| idx_db_tasks_user_id | user_id | B-tree | 按用户查询任务（主索引） |
| idx_db_tasks_status | status | B-tree | 按状态筛选任务 |
| idx_db_tasks_priority | priority_level | B-tree | 按优先级排序 |
| idx_db_tasks_due_date | due_date | B-tree | 按截止时间排序 |
| idx_db_tasks_deleted | is_deleted | 部分索引 | 查询未删除任务 |
| idx_db_tasks_milvus_id | milvus_vector_id | B-tree | 按 Milvus 向量 ID 查询 |
| idx_db_tasks_type | task_type | B-tree | 按任务类型筛选 |

### 4.2 消息表索引

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| idx_db_messages_user_id | user_id | B-tree | 按用户查询消息 |
| idx_db_messages_task_id | related_task_id | B-tree | 按关联任务查询消息 |
| idx_db_messages_created_at | created_at | B-tree | 按时间排序消息 |

---

## 5. 触发器

### 5.1 更新时间戳触发器

```sql
CREATE OR REPLACE FUNCTION db_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 用户表触发器
CREATE TRIGGER trg_db_users_updated_at
    BEFORE UPDATE ON db_users
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();

-- 任务表触发器
CREATE TRIGGER trg_db_tasks_updated_at
    BEFORE UPDATE ON db_tasks
    FOR EACH ROW
    EXECUTE FUNCTION db_update_updated_at();
```

---

## 6. 视图

### 6.1 活跃任务视图

```sql
CREATE VIEW v_db_active_tasks AS
SELECT 
    id, user_id, title, description, status, priority, 
    priority_level, task_type, due_date, milvus_vector_id, created_at
FROM db_tasks
WHERE is_deleted = FALSE 
  AND status IN ('todo', 'doing');
```

### 6.2 今日任务视图

```sql
CREATE VIEW v_db_today_tasks AS
SELECT 
    id, user_id, title, priority_level, due_date
FROM db_tasks
WHERE is_deleted = FALSE 
  AND status = 'todo'
  AND (due_date IS NULL OR due_date::date = CURRENT_DATE);
```

---

## 7. Milvus Collection 设计

### 7.1 字段定义

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | VARCHAR(64) | 主键，使用 PostgreSQL 中的 milvus_vector_id |
| user_id | VARCHAR(64) | 用户 ID，数据隔离 |
| embedding | FLOAT_VECTOR(1024) | 任务文本的向量表示 |
| title | VARCHAR(255) | 任务标题（原始文本，用于返回） |
| content_text | VARCHAR(2048) | 完整文本内容：标题 + 描述 |
| status | VARCHAR(20) | 任务状态：todo, doing, done |
| priority_level | VARCHAR(10) | 优先级：HIGH, MEDIUM, LOW |
| task_type | VARCHAR(20) | 任务类型：work, personal, project |
| due_date | VARCHAR(32) | 截止时间（ISO 格式字符串） |

### 7.2 向量化策略

**合并文本向量化**：将 `title` + `description` 合并为一个文本字段进行向量化

```python
def build_embedding_text(task: Task) -> str:
    parts = []
    if task.title:
        parts.append(f"任务标题：{task.title}")
    if task.description:
        parts.append(f"任务描述：{task.description}")
    parts.append(f"任务类型：{task.task_type}")
    parts.append(f"优先级：{task.priority_level}")
    return " | ".join(parts)
```

### 7.3 向量搜索示例

```python
# 语义搜索任务
query_vector = embedding_model.encode("关于会议的任务")

results = collection.search(
    data=[query_vector],
    anns_field="embedding",
    limit=10,
    expr="user_id == 'xxx' and status == 'todo'",
    output_fields=["title", "priority_level", "due_date"]
)
```

---

## 8. 同步机制

### 8.1 同步流程

```
任务创建/更新
     ↓
PostgreSQL 插入/更新任务
     ↓
更新 milvus_vector_id
     ↓
异步：生成 Embedding + 插入 Milvus
     ↓
（使用 Upsert 覆盖旧数据）
```

### 8.2 数据覆盖策略

| 场景 | 策略 | 说明 |
|------|------|------|
| 新建任务 | Insert | 首次插入 Milvus |
| 更新任务 | Upsert | 根据 milvus_vector_id 覆盖 |
| 删除任务 | Delete | 从 Milvus 删除向量 |

**前提条件**：
- Milvus Collection 主键使用 `VARCHAR` 类型
- `auto_id = false`（手动指定 ID）
- PostgreSQL 和 Milvus 使用相同 UUID 作为 ID

---

## 9. 枚举值说明

### 9.1 任务状态 (status)

| 值 | 说明 |
|----|------|
| todo | 待办 |
| doing | 进行中 |
| done | 已完成 |
| cancelled | 已取消 |

### 9.2 优先级级别 (priority_level)

| 值 | 说明 |
|----|------|
| HIGH | 高优先级 |
| MEDIUM | 中优先级 |
| LOW | 低优先级 |

### 9.3 任务类型 (task_type)

| 值 | 说明 |
|----|------|
| work | 工作 |
| personal | 个人 |
| project | 项目 |
| adhoc | 临时 |

---

## 10. SQL 脚本

详见：[schema.sql](../sql/schema.sql)
