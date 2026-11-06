package tools

import (
	"context"
	"encoding/json"
	"errors"

	"service/internal/dao"
	"service/internal/middleware"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
)

// UpdateTaskArgs 按需更新任务参数：taskId 必填其余只传改动字段。
// description/dueDate 清空走弹窗本工具不支持（空=不改）。
type UpdateTaskArgs struct {
	TaskID      string `json:"taskId" jsonschema:"任务 ID，必填"`
	Title       string `json:"title" jsonschema:"标题只传改动"`
	Description string `json:"description" jsonschema:"描述只传改动清空不支持"`
	Priority    string `json:"priority" jsonschema:"优先级只传改动"`
	TaskType    string `json:"taskType" jsonschema:"任务类型只传改动"`
	DueDate     string `json:"dueDate" jsonschema:"截止时间只传改动清空不支持"`
	Status      string `json:"status" jsonschema:"状态只传改动"`
}

func UpdateTaskTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "update_task",
			Desc: "用户在人机验证卡上确认后按需更新任务。taskId 必填其余只传改动字段空=不改；description/dueDate 清空不支持走弹窗。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"taskId": {
						Desc:     "任务 ID，必填",
						Type:     "string",
						Required: true,
					},
					"title": {
						Desc:     "标题，只传改动",
						Type:     "string",
						Required: false,
					},
					"description": {
						Desc:     "描述只传改动",
						Type:     "string",
						Required: false,
					},
					"priority": {
						Desc:     "优先级只传改动",
						Type:     "string",
						Required: false,
					},
					"taskType": {
						Desc:     "任务类型只传改动",
						Type:     "string",
						Required: false,
					},
					"dueDate": {
						Desc:     "截止时间只传改动",
						Type:     "string",
						Required: false,
					},
					"status": {
						Desc:     "状态只传改动",
						Type:     "string",
						Required: false,
					},
				}),
		},
		updateTask,
	)
}

func updateTask(ctx context.Context, input *UpdateTaskArgs) (string, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return "", errors.New("user_id not found in context")
	}
	uid, err := uuid.Parse(userId)
	if err != nil {
		return "", errors.New("invalid userId")
	}
	if _, err := uuid.Parse(input.TaskID); err != nil {
		return "", errors.New("invalid taskId")
	}
	// 已完成/已取消的任务不允许编辑（候选已过滤，此处防绕过）。
	existing, err := dao.DbTasks.GetTask(ctx, input.TaskID)
	if err != nil {
		return "", err
	}
	if existing == nil || existing.Status == "done" || existing.Status == "cancelled" {
		return "", errors.New("该任务已完成或已取消不能编辑")
	}
	data := g.Map{}
	if input.Title != "" {
		data["title"] = input.Title
	}
	if input.Description != "" {
		data["description"] = input.Description
	}
	if input.Priority != "" {
		data["priority"] = input.Priority
	}
	if input.TaskType != "" {
		data["task_type"] = input.TaskType
	}
	if input.DueDate != "" {
		data["due_date"] = input.DueDate
	}
	if input.Status != "" {
		data["status"] = input.Status
	}
	if len(data) == 0 {
		return "", errors.New("没有改动字段")
	}
	if err := dao.DbTasks.PatchTask(ctx, uid, input.TaskID, data); err != nil {
		glog.Error(ctx, "[updateTask] patch failed: %v", err)
		return "", err
	}
	resultJSON, err := json.Marshal(map[string]any{"success": true, "taskId": input.TaskID})
	if err != nil {
		return "", err
	}
	return string(resultJSON), nil
}
