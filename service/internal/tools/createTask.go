package tools

import (
	"context"
	"errors"
	"service/internal/dao"
	"service/internal/middleware"
	"service/internal/model/entity"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// 创建任务参数
type CreateTaskArgs struct {
	Title       string  `json:"title" jsonschema:"任务标题，必填"`
	Description string  `json:"description" jsonschema:"任务详细描述，可选"`
	DueDate     *string `json:"dueDate" jsonschema:"截止时间，ISO 8601 格式，例如 2026-03-25T10:00:00Z"`
	Priority    string `json:"priority" jsonschema:"优先级：HIGH|MEDIUM|LOW，默认 MEDIUM"`
	TaskType    string `json:"taskType" jsonschema:"任务类型：work|personal|project|adhoc，默认 work"`
}

func CreateTask() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "create_task",
			Desc: "创建新的待办任务。当用户表达添加任务、安排会议、设置提醒等意图时调用。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"title": {
						Desc:     "任务标题，简短概括任务内容",
						Type:     "string",
						Required: true,
					},
					"description": {
						Desc:     "任务详细描述",
						Type:     "string",
						Required: false,
					},
					"dueDate": {
						Desc:     "截止时间，ISO 8601 格式，例如 2026-03-25T10:00:00Z",
						Type:     "string",
						Required: false,
					},
					"priority": {
						Desc:     "优先级",
						Type:     "string",
						Required: false,
						Enum:     []string{"HIGH", "MEDIUM", "LOW"},
					},
					"taskType": {
						Desc:     "任务类型",
						Type:     "string",
						Required: false,
						Enum:     []string{"work", "personal", "project", "adhoc"},
					},
				}),
		},
		createTask,
	)
}

func createTask(ctx context.Context, input *CreateTaskArgs) (uuid.UUID, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, errors.New("无效的用户ID")
	}

	if input.Title == "" {
		return uuid.Nil, errors.New("任务标题不能为空")
	}

	var dueDate *gtime.Time
	if input.DueDate != nil && *input.DueDate != "" {
		formats := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"}
		for _, format := range formats {
			if parsed, err := time.Parse(format, *input.DueDate); err == nil {
				dueDate = gtime.NewFromTime(parsed.UTC())
				break
			}
		}
		if dueDate == nil {
			g.Log().Warning(ctx, "无法解析日期格式: ", *input.DueDate)
		}
	}
	g.Log().Info(ctx, "dueDate: ", dueDate)
	if input.Priority == "" {
		input.Priority = "MEDIUM"
	}
	if input.TaskType == "" {
		input.TaskType = "work"
	}

	taskId := uuid.New()

	tx, err := dao.DbTasks.InsertTask(ctx, &entity.DbTasks{
		Id:          taskId,
		UserId:      userIdUUID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
		Priority:    input.Priority,
		TaskType:    input.TaskType,
	})
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return uuid.Nil, err
	}

	tx.Commit()
	return taskId, nil
}
