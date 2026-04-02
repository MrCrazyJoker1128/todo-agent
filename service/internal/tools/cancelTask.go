package tools

import (
	"context"
	"errors"
	"service/internal/dao"
	"service/internal/middleware"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

type CancelTaskTool struct {
	TaskId string `json:"taskId"`
	ChatId string `json:"chatId"`
}

func CancelTask() tool.InvokableTool {
	return utils.NewTool(&schema.ToolInfo{
		Name: "cancel_task",
		Desc: "取消/删除任务。支持三种方式：1)直接指定taskId 2)通过chatId查找关联任务 3)取消用户最近创建的任务",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"taskId": {
					Desc:     "任务 ID，直接指定要取消的任务",
					Type:     "string",
					Required: false,
				},
				"chatId": {
					Desc:     "会话 ID，通过会话查找关联的任务",
					Type:     "string",
					Required: false,
				},
			}),
	}, cancelTask)
}

func cancelTask(ctx context.Context, input *CancelTaskTool) (string, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return "", errors.New("user_id not found in context")
	}

	var taskIdToCancel string

	if input.TaskId != "" {
		taskIdToCancel = input.TaskId
	} else if input.ChatId != "" {
		chatId, err := uuid.Parse(input.ChatId)
		if err != nil {
			return "", errors.New("invalid chatId format")
		}
		chatMsg, err := dao.DbChatMessages.GetChatMessageById(ctx, chatId)
		if err != nil {
			return "", errors.New("chat message not found")
		}
		if chatMsg.RelatedTaskId == uuid.Nil {
			return "", errors.New("该会话没有关联的任务")
		}
		taskIdToCancel = chatMsg.RelatedTaskId.String()
	} else {
		userIdUUID, err := uuid.Parse(userId)
		if err != nil {
			return "", errors.New("invalid userId")
		}
		task, err := dao.DbTasks.GetLatestTaskByUserId(ctx, userIdUUID)
		if err != nil {
			return "", errors.New("没有找到可取消的任务")
		}
		taskIdToCancel = task.Id.String()
	}

	err := dao.DbTasks.DeleteTask(ctx, taskIdToCancel)
	if err != nil {
		return "", err
	}
	return "任务已取消", nil
}
