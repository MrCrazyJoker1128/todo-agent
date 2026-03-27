# 任务插入流程文档

## 1. 流程概述

插入任务是系统的核心功能，涉及参数验证、业务处理、数据库操作和向量同步等多个环节。

```
┌─────────────────────────────────────────────────────────────────────┐
│                        插入任务完整流程                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐                                                     │
│  │ 1. 请求接收 │ ← POST /api/tasks 或 POST /api/chat               │
│  └──────┬──────┘                                                     │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────┐                                                     │
│  │ 2. 参数验证 │ ← 必填字段、格式校验、用户权限                      │
│  └──────┬──────┘                                                     │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────┐                                                     │
│  │ 3. 业务处理 │ ← 优先级判断、时间解析、类型分类                     │
│  └──────┬──────┘                                                     │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────┐                                                     │
│  │ 4. 数据库事务│ ← PostgreSQL 插入任务                               │
│  └──────┬──────┘                                                     │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────┐                                                     │
│  │ 5. 向量同步 │ ← 异步：生成 Embedding + Upsert Milvus             │
│  └──────┬──────┘                                                     │
│         │                                                             │
│         ▼                                                             │
│  ┌─────────────┐                                                     │
│  │ 6. 响应返回 │ ← 返回任务信息                                      │
│  └──────┬──────┘                                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. 详细步骤

### 2.1 步骤 1：请求接收与解析

**职责**：
- 解析 JSON 请求体
- 从 JWT Token 提取用户 ID
- 记录审计日志

**代码示例**：
```go
func (h *TaskHandler) CreateTask(c *gin.Context) {
    var req dto.CreateTaskRequest
    
    // 1.1 解析请求体
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "参数错误", "details": err.Error()})
        return
    }
    
    // 1.2 获取用户 ID
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(401, gin.H{"error": "未授权"})
        return
    }
    
    // 1.3 记录审计日志
    log.Info("创建任务请求", "user_id", userID, "title", req.Title)
}
```

### 2.2 步骤 2：参数验证

**验证规则**：

| 字段 | 规则 | 错误码 |
|------|------|--------|
| title | 必填，1-255 字符 | 10001 |
| priority_level | HIGH/MEDIUM/LOW 或空 | 10001 |
| task_type | work/personal/project/adhoc 或空 | 10001 |
| due_date | ISO 8601 格式，不能是过去时间 | 10001 |

**代码示例**：
```go
func (r *CreateTaskRequest) Validate() error {
    // 标题验证
    if strings.TrimSpace(r.Title) == "" {
        return errors.New("标题不能为空")
    }
    if len(r.Title) > 255 {
        return errors.New("标题不能超过 255 字符")
    }
    
    // 截止时间验证
    if r.DueDate != nil && r.DueDate.Before(time.Now()) {
        return errors.New("截止时间不能是过去时间")
    }
    
    return nil
}
```

### 2.3 步骤 3：业务逻辑处理

**处理项**：

| 处理项 | 说明 | 默认值 |
|--------|------|--------|
| 优先级自动判断 | 根据关键词/时间推断 | MEDIUM |
| 任务类型分类 | 根据内容分类或默认值 | work |
| 时间标准化 | 转换为 UTC 存储 | - |
| 默认值填充 | status=todo, is_deleted=false | - |

**代码示例**：
```go
func (s *TaskService) PrepareTask(ctx context.Context, userID string, req dto.CreateTaskRequest) (*model.Task, error) {
    task := &model.Task{
        UserID:        userID,
        Title:         strings.TrimSpace(req.Title),
        Description:   strings.TrimSpace(req.Description),
        Status:        "todo",
        IsDeleted:     false,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    // 3.1 优先级处理
    if req.PriorityLevel != "" {
        task.PriorityLevel = req.PriorityLevel
    } else {
        task.PriorityLevel = s.autoDetectPriority(req.Title, req.DueDate)
    }
    
    // 3.2 任务类型处理
    if req.TaskType != "" {
        task.TaskType = req.TaskType
    } else {
        task.TaskType = "work"
    }
    
    // 3.3 时间标准化
    if req.DueDate != nil {
        utcTime := req.DueDate.UTC()
        task.DueDate = &utcTime
    }
    
    // 3.4 生成 Milvus Vector ID
    task.MilvusVectorID = uuid.New()
    
    return task, nil
}

// 自动优先级判断
func (s *TaskService) autoDetectPriority(title string, dueDate *time.Time) string {
    urgentKeywords := []string{"紧急", "马上", "今天", "截止", "故障"}
    for _, keyword := range urgentKeywords {
        if strings.Contains(title, keyword) {
            return "HIGH"
        }
    }
    
    if dueDate != nil {
        hoursUntilDue := time.Until(*dueDate).Hours()
        if hoursUntilDue < 24 {
            return "HIGH"
        } else if hoursUntilDue < 72 {
            return "MEDIUM"
        }
    }
    
    return "MEDIUM"
}
```

### 2.4 步骤 4：数据库事务处理

**数据库操作**：

| 操作 | 表 | 说明 |
|------|-----|------|
| 验证用户 | db_users | 确保用户存在 |
| 插入任务 | db_tasks | 主业务数据 |

**代码示例**：
```go
func (s *TaskService) CreateTask(ctx context.Context, userID string, req dto.CreateTaskRequest) (*model.Task, error) {
    var createdTask *model.Task
    
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 4.1 验证用户存在
        var user model.User
        if err := tx.First(&user, "id = ?", userID).Error; err != nil {
            return errors.New("用户不存在")
        }
        
        // 4.2 准备任务数据
        task, err := s.PrepareTask(ctx, userID, req)
        if err != nil {
            return err
        }
        
        // 4.3 插入任务
        if err := tx.Create(task).Error; err != nil {
            return fmt.Errorf("创建任务失败：%w", err)
        }
        
        createdTask = task
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // 4.4 异步同步向量（事务外）
    go func() {
        bgCtx := context.Background()
        if err := s.vectorService.SyncTaskVector(bgCtx, createdTask); err != nil {
            log.Error("向量同步失败", "task_id", createdTask.ID, "error", err)
        }
    }()
    
    return createdTask, nil
}
```

### 2.5 步骤 5：向量同步（异步）

**同步流程**：
```
任务创建成功
     ↓
启动 goroutine
     ↓
生成 milvus_vector_id (UUID)
     ↓
构建向量化文本
     ↓
调用 MiniMax Embedding API
     ↓
Upsert 到 Milvus (覆盖旧数据)
     ↓
更新任务表的 milvus_vector_id
```

**代码示例**：
```go
func (s *VectorService) SyncTaskVector(ctx context.Context, task *model.Task) error {
    // 5.1 构建向量化文本
    text := s.buildEmbeddingText(task)
    
    // 5.2 生成向量
    vector, err := s.minimaxClient.Embeddings(ctx, text)
    if err != nil {
        return err
    }
    
    // 5.3 Upsert 到 Milvus
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
    
    if err := s.indexer.Upsert(ctx, []Document{doc}); err != nil {
        return err
    }
    
    log.Info("向量同步成功", "task_id", task.ID, "milvus_vector_id", task.MilvusVectorID)
    return nil
}
```

### 2.6 步骤 6：响应返回

**响应字段**：

| 字段 | 说明 | 示例 |
|------|------|------|
| id | 任务 ID | "550e8400-..." |
| title | 任务标题 | "整理会议纪要" |
| priority_level | 优先级 | "HIGH" |
| status | 状态 | "todo" |
| created_at | 创建时间 | "2026-03-24T10:00:00Z" |

**代码示例**：
```go
func (h *TaskHandler) CreateTask(c *gin.Context) {
    // ... 前面步骤 ...
    
    task, err := h.service.CreateTask(c.Request.Context(), userID, req)
    if err != nil {
        log.Error("创建任务失败", "user_id", userID, "error", err)
        c.JSON(500, gin.H{"error": "创建任务失败", "details": err.Error()})
        return
    }
    
    // 返回响应
    c.JSON(201, gin.H{
        "code": 201,
        "message": "任务创建成功",
        "data": task,
    })
}
```

---

## 3. 异常处理

| 异常场景 | 处理策略 |
|----------|----------|
| 用户不存在 | 返回 404，记录安全日志 |
| 参数验证失败 | 返回 400，提示具体字段 |
| 数据库写入失败 | 事务回滚，返回 500 |
| Milvus 不可用 | 向量同步失败，任务仍创建成功 |
| 向量生成失败 | 重试 3 次，仍失败记录日志 |

---

## 4. 检查清单

| 步骤 | 检查项 | 状态 |
|------|--------|------|
| 1 | 请求参数正确解析 | ☐ |
| 2 | 用户权限验证通过 | ☐ |
| 3 | 必填字段验证通过 | ☐ |
| 4 | 业务逻辑处理完成 | ☐ |
| 5 | 数据库事务成功提交 | ☐ |
| 6 | 向量异步同步启动 | ☐ |
| 7 | 响应格式正确 | ☐ |
| 8 | 错误日志记录 | ☐ |
| 9 | 审计日志记录 | ☐ |
