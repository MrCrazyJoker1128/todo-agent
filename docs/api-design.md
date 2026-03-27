# 接口设计文档

## 1. 概述

- **认证方式**：JWT Token
- **响应格式**：JSON
- **字符编码**：UTF-8

### 1.1 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 1.2 错误响应格式

```json
{
  "code": 400,
  "message": "参数错误",
  "details": "具体错误信息"
}
```

### 1.3 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

---

## 2. 任务接口 (Task)

### 2.1 创建任务

**POST** `/api/tasks`

**请求头**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**
```json
{
  "title": "整理会议纪要",
  "description": "汇总本周团队会议内容并发送",
  "priority_level": "HIGH",
  "task_type": "work",
  "due_date": "2026-03-25T15:00:00Z",
  "related_message_id": "uuid"  // 可选，聊天创建时关联
}
```

**响应**
```json
{
  "code": 201,
  "message": "任务创建成功",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "整理会议纪要",
    "description": "汇总本周团队会议内容并发送",
    "status": "todo",
    "priority_level": "HIGH",
    "task_type": "work",
    "due_date": "2026-03-25T15:00:00Z",
    "milvus_vector_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2026-03-24T10:00:00Z",
    "updated_at": "2026-03-24T10:00:00Z"
  }
}
```

### 2.2 获取任务列表

**GET** `/api/tasks`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 筛选状态：todo, doing, done, cancelled |
| priority | string | 否 | 筛选优先级：HIGH, MEDIUM, LOW |
| task_type | string | 否 | 筛选类型：work, personal, project, adhoc |
| due_date | string | 否 | 筛选日期：YYYY-MM-DD |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "整理会议纪要",
        "status": "todo",
        "priority_level": "HIGH",
        "due_date": "2026-03-25T15:00:00Z",
        "created_at": "2026-03-24T10:00:00Z"
      }
    ]
  }
}
```

### 2.3 获取任务详情

**GET** `/api/tasks/:id`

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "user-uuid",
    "title": "整理会议纪要",
    "description": "汇总本周团队会议内容并发送",
    "status": "todo",
    "priority": "high",
    "priority_level": "HIGH",
    "task_type": "work",
    "due_date": "2026-03-25T15:00:00Z",
    "milvus_vector_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2026-03-24T10:00:00Z",
    "updated_at": "2026-03-24T10:00:00Z"
  }
}
```

### 2.4 更新任务

**PUT** `/api/tasks/:id`

**请求体**
```json
{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "status": "doing",
  "priority_level": "MEDIUM",
  "due_date": "2026-03-26T15:00:00Z"
}
```

**响应**
```json
{
  "code": 200,
  "message": "任务更新成功",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "更新后的标题",
    "updated_at": "2026-03-24T11:00:00Z"
  }
}
```

### 2.5 完成任务

**PATCH** `/api/tasks/:id/complete`

**响应**
```json
{
  "code": 200,
  "message": "任务已完成",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "done",
    "updated_at": "2026-03-24T11:00:00Z"
  }
}
```

### 2.6 删除任务

**DELETE** `/api/tasks/:id`

**响应**
```json
{
  "code": 200,
  "message": "任务已删除",
  "data": null
}
```

### 2.7 搜索任务（向量搜索）

**POST** `/api/tasks/search`

**请求体**
```json
{
  "query": "找一下关于会议的任务",
  "filters": {
    "status": "todo",
    "priority_level": "HIGH"
  },
  "limit": 10
}
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "results": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "整理会议纪要",
        "priority_level": "HIGH",
        "score": 0.95
      }
    ]
  }
}
```

---

## 3. 聊天接口 (Chat)

### 3.1 发送消息（非流式）

**POST** `/api/chat`

**请求体**
```json
{
  "message": "明天下午 3 点开会"
}
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "reply": "已为您创建任务：团队会议",
    "task": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "团队会议",
      "due_date": "2026-03-25T15:00:00Z",
      "priority_level": "HIGH"
    }
  }
}
```

### 3.2 发送消息（流式 SSE）

**POST** `/api/chat/stream`

**请求体**
```json
{
  "message": "明天下午 3 点开会"
}
```

**响应头**
```
Content-Type: text/event-stream
Cache-Control: no-cache
```

**SSE 事件格式**
```json
event: thinking
data: {"content": "正在分析...", "type": "thinking"}

event: task_created
data: {"task_id": "uuid", "title": "团队会议", "type": "task_created"}

event: done
data: {"content": "已为您创建任务", "type": "done"}
```

### 3.3 获取聊天历史

**GET** `/api/chat/history`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": "msg-uuid",
        "role": "user",
        "content": "明天下午 3 点开会",
        "related_task_id": "task-uuid",
        "created_at": "2026-03-24T10:00:00Z"
      },
      {
        "id": "msg-uuid",
        "role": "assistant",
        "content": "已为您创建任务",
        "related_task_id": "task-uuid",
        "created_at": "2026-03-24T10:00:01Z"
      }
    ]
  }
}
```

---

## 4. 通用查询接口

### 4.1 获取今日任务

**GET** `/api/tasks/today`

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "count": 5,
    "tasks": [...]
  }
}
```

### 4.2 获取本周任务

**GET** `/api/tasks/week`

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "count": 20,
    "tasks": [...]
  }
}
```

### 4.3 获取已完成任务

**GET** `/api/tasks/completed`

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "count": 50,
    "tasks": [...]
  }
}
```

---

## 5. 附录

### 5.1 DTO 结构

```go
// 创建任务请求
type CreateTaskRequest struct {
    Title            string     `json:"title" binding:"required,min=1,max=255"`
    Description      string     `json:"description" binding:"max=2048"`
    PriorityLevel    string     `json:"priority_level" binding:"omitempty,oneof=HIGH MEDIUM LOW"`
    TaskType         string     `json:"task_type" binding:"omitempty,oneof=work personal project adhoc"`
    DueDate          *time.Time `json:"due_date"`
    RelatedMessageID string     `json:"related_message_id"`
}

// 更新任务请求
type UpdateTaskRequest struct {
    Title         *string    `json:"title" binding:"omitempty,min=1,max=255"`
    Description   *string    `json:"description" binding:"omitempty,max=2048"`
    Status        *string    `json:"status" binding:"omitempty,oneof=todo doing done cancelled"`
    PriorityLevel *string    `json:"priority_level" binding:"omitempty,oneof=HIGH MEDIUM LOW"`
    DueDate       *time.Time `json:"due_date"`
}

// 聊天请求
type ChatRequest struct {
    Message string `json:"message" binding:"required"`
}

// 搜索请求
type SearchRequest struct {
    Query   string       `json:"query" binding:"required"`
    Filters SearchFilter `json:"filters"`
    Limit   int          `json:"limit" binding:"omitempty,max=100"`
}

type SearchFilter struct {
    Status        string `json:"status"`
    PriorityLevel string `json:"priority_level"`
    TaskType     string `json:"task_type"`
    DueDate      string `json:"due_date"`
}
```

### 5.2 错误码定义

| 错误码 | 说明 |
|--------|------|
| 10001 | 参数验证失败 |
| 10002 | 用户不存在 |
| 10003 | 任务不存在 |
| 10004 | 权限不足 |
| 20001 | AI 服务错误 |
| 20002 | 向量服务错误 |
| 30001 | 数据库错误 |
| 30002 | 缓存错误 |
