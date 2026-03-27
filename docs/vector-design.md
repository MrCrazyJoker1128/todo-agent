# 任务向量化文档

## 1. 向量化概述

### 1.1 设计原则

| 原则 | 说明 |
|------|------|
| **文本向量化** | title + description 合并后生成向量 |
| **字段做过滤** | status, priority_level, task_type 等作为 Milvus 过滤条件 |
| **统一 UUID** | PostgreSQL 和 Milvus 使用相同的 UUID |

### 1.2 字段分类

| 字段 | 类型 | 处理方式 |
|------|------|----------|
| title | VARCHAR | ✅ 向量化（合并到文本） |
| description | TEXT | ✅ 向量化（合并到文本） |
| user_id | UUID | 🔒 过滤条件（精确匹配） |
| status | VARCHAR | 🔒 过滤条件（精确匹配） |
| priority_level | VARCHAR | 🔒 过滤条件（精确匹配） |
| task_type | VARCHAR | 🔒 过滤条件（精确匹配） |
| due_date | TIMESTAMPTZ | 🔒 过滤条件（范围查询） |

---

## 2. 向量化文本构建

### 2.1 构建策略

将 `title` + `description` 合并为一个文本字段进行向量化：

```go
func BuildEmbeddingText(task *model.Task) string {
    var parts []string
    
    if task.Title != "" {
        parts = append(parts, fmt.Sprintf("任务标题：%s", task.Title))
    }
    if task.Description != "" {
        parts = append(parts, fmt.Sprintf("任务描述：%s", task.Description))
    }
    
    // 添加上下文信息
    parts = append(parts, fmt.Sprintf("任务类型：%s", task.TaskType))
    parts = append(parts, fmt.Sprintf("优先级：%s", task.PriorityLevel))
    
    return strings.Join(parts, " | ")
}
```

### 2.2 示例

**输入**
```json
{
  "title": "整理会议纪要",
  "description": "汇总本周团队会议内容并发送给所有参会者",
  "task_type": "work",
  "priority_level": "HIGH"
}
```

**向量化文本**
```
任务标题：整理会议纪要 | 任务描述：汇总本周团队会议内容并发送给所有参会者 | 任务类型：work | 优先级：HIGH
```

---

## 3. Milvus Collection 设计

### 3.1 Schema 定义

```python
from pymilvus import CollectionSchema, FieldSchema, DataType

schema = CollectionSchema([
    # 主键字段
    FieldSchema(
        name="id",
        dtype=DataType.VARCHAR,
        max_length=64,
        is_primary=True,
        auto_id=False  # 手动指定 ID（PostgreSQL 的 milvus_vector_id）
    ),
    
    # 向量字段
    FieldSchema(
        name="embedding",
        dtype=DataType.FLOAT_VECTOR,
        dim=1024  # MiniMax Embedding 维度
    ),
    
    # 过滤字段
    FieldSchema(
        name="user_id",
        dtype=DataType.VARCHAR,
        max_length=64
    ),
    FieldSchema(
        name="status",
        dtype=DataType.VARCHAR,
        max_length=20
    ),
    FieldSchema(
        name="priority_level",
        dtype=DataType.VARCHAR,
        max_length=10
    ),
    FieldSchema(
        name="task_type",
        dtype=DataType.VARCHAR,
        max_length=20
    ),
    FieldSchema(
        name="due_date",
        dtype=DataType.VARCHAR,
        max_length=32
    ),
    
    # 返回字段
    FieldSchema(
        name="title",
        dtype=DataType.VARCHAR,
        max_length=255
    ),
    FieldSchema(
        name="content_text",
        dtype=DataType.VARCHAR,
        max_length=2048
    ),
])
```

### 3.2 索引配置

```python
# 向量索引
collection.create_index(
    field_name="embedding",
    index_params={
        "index_type": "IVF_FLAT",
        "metric_type": "COSINE",
        "params": {"nlist": 1024}
    }
)

# 标量索引（加速过滤查询）
collection.create_index(
    field_name="user_id",
    index_params={"index_type": "STL_SORT"}
)
collection.create_index(
    field_name="status",
    index_params={"index_type": "STL_SORT"}
)
```

---

## 4. Golang 实现

### 4.1 Milvus Indexer

```go
// pkg/vector/milvus_indexer.go
type MilvusIndexer struct {
    client     client.Client
    collection string
    useUpsert  bool  // 是否启用 upsert
}

func (m *MilvusIndexer) Upsert(ctx context.Context, docs []Document) error {
    // 1. 生成向量
    embeddings, err := m.generateEmbeddings(ctx, docs)
    if err != nil {
        return err
    }
    
    // 2. 准备列数据
    var ids []string
    var userIDs []string
    var vectors []entity.Vector
    var titles []string
    var statuses []string
    var priorities []string
    
    for i, doc := range docs {
        ids = append(ids, doc.ID)
        userIDs = append(userIDs, doc.Metadata["user_id"])
        vectors = append(vectors, entity.FloatVector(embeddings[i]))
        titles = append(titles, doc.Metadata["title"])
        statuses = append(statuses, doc.Metadata["status"])
        priorities = append(priorities, doc.Metadata["priority_level"])
    }
    
    // 3. 创建列
    columns := []entity.Column{
        entity.NewColumnVarChar("id", ids),
        entity.NewColumnVarChar("user_id", userIDs),
        entity.NewColumnFloatVector("embedding", 1024, vectors),
        entity.NewColumnVarChar("title", titles),
        entity.NewColumnVarChar("status", statuses),
        entity.NewColumnVarChar("priority_level", priorities),
    }
    
    // 4. 执行 Upsert（覆盖旧数据）
    _, err = m.client.Upsert(ctx, m.collection, "", columns...)
    return err
}
```

### 4.2 向量同步服务

```go
// internal/service/vector.go
type VectorService struct {
    indexer *MilvusIndexer
    minimax *minimax.Client
    db      *gorm.DB
}

// 同步任务向量
func (s *VectorService) SyncTaskVector(ctx context.Context, task *model.Task) error {
    // 1. 构建向量化文本
    text := s.buildEmbeddingText(task)
    
    // 2. 生成向量
    vector, err := s.minimax.Embeddings(ctx, text)
    if err != nil {
        return err
    }
    
    // 3. 准备文档
    doc := Document{
        ID:      task.MilvusVectorID.String(),
        Content: text,
        Metadata: map[string]interface{}{
            "user_id":        task.UserID.String(),
            "source_id":      task.ID.String(),
            "status":         task.Status,
            "priority_level": task.PriorityLevel,
            "task_type":      task.TaskType,
            "title":          task.Title,
        },
        Embedding: vector,
    }
    
    // 4. 调用 Upsert
    return s.indexer.Upsert(ctx, []Document{doc})
}

// 异步同步
func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (*model.Task, error) {
    // ... 数据库操作 ...
    
    // 异步同步向量
    go func() {
        bgCtx := context.Background()
        if err := s.vectorService.SyncTaskVector(bgCtx, task); err != nil {
            log.Error("向量同步失败", "task_id", task.ID, "error", err)
        }
    }()
    
    return task, nil
}
```

---

## 5. 向量搜索

### 5.1 搜索实现

```go
func (s *VectorService) SearchTasks(ctx context.Context, req SearchRequest) ([]model.Task, error) {
    // 1. 生成查询向量
    queryVector, err := s.minimax.Embeddings(ctx, req.Query)
    if err != nil {
        return nil, err
    }
    
    // 2. 构建过滤表达式
    expr := fmt.Sprintf("user_id == '%s'", req.UserID)
    if req.Filters.Status != "" {
        expr += fmt.Sprintf(" and status == '%s'", req.Filters.Status)
    }
    if req.Filters.PriorityLevel != "" {
        expr += fmt.Sprintf(" and priority_level == '%s'", req.Filters.PriorityLevel)
    }
    
    // 3. 执行搜索
    results, err := s.client.Search(
        ctx,
        "tasks",
        []string{"title", "priority_level", "due_date"},
        expr,
        []entity.Vector{entity.FloatVector(queryVector)},
        "embedding",
        entity.COSINE,
        int64(req.Limit),
    )
    if err != nil {
        return nil, err
    }
    
    // 4. 解析结果
    var tasks []model.Task
    for _, hit := range results[0].Hits {
        task := model.Task{
            Title:        hit.Fields["title"].(string),
            PriorityLevel: hit.Fields["priority_level"].(string),
        }
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}
```

### 5.2 搜索示例

**场景 1：语义搜索**
```
输入: "找一下关于会议的任务"
输出: 匹配所有包含"会议"语义的的任务
```

**场景 2：带过滤的语义搜索**
```
输入: "本周高优先级的客户相关任务"
过滤: priority_level == 'HIGH', due_date >= '2026-03-24', due_date <= '2026-03-31'
输出: 高优先级且本周截止的客户相关任务
```

---

## 6. 数据覆盖策略

### 6.1 覆盖机制

| 操作 | Milvus 接口 | 说明 |
|------|-----------|------|
| 新建任务 | Insert/Upsert | 插入新向量 |
| 更新任务 | Upsert | 根据 ID 覆盖旧向量 |
| 删除任务 | Delete | 从 Milvus 删除向量 |

### 6.2 Upsert 前提条件

1. **Milvus 版本**：需要 2.3+
2. **主键配置**：`auto_id = false`
3. **主键类型**：VARCHAR（支持 UUID 字符串）

### 6.3 回退方案

如果 Milvus 版本不支持 Upsert，使用 Delete + Insert：

```go
func (s *MilvusIndexer) DeleteAndInsert(ctx context.Context, doc Document) error {
    // 1. 删除旧数据
    expr := fmt.Sprintf("id == '%s'", doc.ID)
    s.client.Delete(ctx, s.collection, "", expr)
    
    // 2. 插入新数据
    return s.Insert(ctx, []Document{doc})
}
```

---

## 7. 同步状态管理

### 7.1 同步状态字段

| 状态 | 说明 |
|------|------|
| pending | 待同步 |
| syncing | 同步中 |
| synced | 已同步 |
| failed | 同步失败 |

### 7.2 监控视图

```sql
CREATE VIEW v_db_pending_vector_sync AS
SELECT 
    id, 
    user_id, 
    title, 
    priority_level, 
    milvus_vector_id, 
    updated_at
FROM db_tasks
WHERE is_deleted = FALSE 
  AND status = 'todo';
```

---

## 8. 性能优化

### 8.1 批量同步

```go
func (s *VectorService) BatchSync(ctx context.Context, tasks []model.Task) error {
    // 1. 批量生成向量
    texts := make([]string, len(tasks))
    for i, task := range tasks {
        texts[i] = s.buildEmbeddingText(&task)
    }
    vectors, err := s.minimax.BatchEmbeddings(ctx, texts)
    if err != nil {
        return err
    }
    
    // 2. 批量 Upsert
    docs := make([]Document, len(tasks))
    for i, task := range tasks {
        docs[i] = Document{
            ID:        task.MilvusVectorID.String(),
            Metadata:  map[string]interface{}{"user_id": task.UserID.String()},
            Embedding: vectors[i],
        }
    }
    
    return s.indexer.BatchUpsert(ctx, docs)
}
```

### 8.2 缓存策略

| 缓存类型 | TTL | 说明 |
|---------|-----|------|
| 相似查询缓存 | 5 分钟 | 避免重复调用 Embedding API |
| 热门任务缓存 | 30 分钟 | 高频任务向量预加载 |

---

## 9. 错误处理

| 错误类型 | 处理策略 |
|---------|---------|
| Embedding API 超时 | 重试 3 次，间隔 1s |
| Milvus 不可用 | 降级为数据库 LIKE 查询 |
| 向量维度不匹配 | 记录错误，标记同步失败 |
| Upsert 不支持 | 回退到 Delete + Insert |
