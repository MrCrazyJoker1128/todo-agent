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
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
)

// QueryTasksArgs 按时间/标题/内容查任务参数，全可选。
type QueryTasksArgs struct {
	Keyword   string `json:"keyword" jsonschema:"标题/内容关键词，可空"`
	StartDate string `json:"startDate" jsonschema:"开始时间 YYYY-MM-DD HH:MM:SS，可空"`
	EndDate   string `json:"endDate" jsonschema:"结束时间 YYYY-MM-DD HH:MM:SS，可空"`
}

// QueryTasksResult 查询结果，字段对齐 api.TaskItem。
type QueryTasksResult struct {
	Tasks []QueryTaskItem `json:"tasks" jsonschema:"任务列表，最多10条"`
	Count int             `json:"count" jsonschema:"结果数量"`
}

type QueryTaskItem struct {
	ID          string `json:"id" jsonschema:"任务 ID"`
	Title       string `json:"title" jsonschema:"标题"`
	Description string `json:"description" jsonschema:"描述"`
	DueDate     string `json:"dueDate" jsonschema:"截止时间"`
	Priority    string `json:"priority" jsonschema:"优先级"`
	TaskType    string `json:"taskType" jsonschema:"任务类型"`
	Status      string `json:"status" jsonschema:"状态"`
}

func QueryTasks() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "query_tasks",
			Desc: "按时间/标题/内容查询用户任务，最多返回10条。用户表达修改/删除任务意图时先调此工具展示候选；用户只想查看任务时调此工具直接展示结果。空结果返回空列表不报错。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"keyword": {
						Desc:     "标题/内容关键词，可空",
						Type:     "string",
						Required: false,
					},
					"startDate": {
						Desc:     "开始时间 YYYY-MM-DD HH:MM:SS，可空",
						Type:     "string",
						Required: false,
					},
					"endDate": {
						Desc:     "结束时间 YYYY-MM-DD HH:MM:SS，可空",
						Type:     "string",
						Required: false,
					},
				}),
		},
		queryTasks,
	)
}

func queryTasks(ctx context.Context, input *QueryTasksArgs) (string, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return "", errors.New("user_id not found in context")
	}
	uid, err := uuid.Parse(userId)
	if err != nil {
		return "", errors.New("invalid userId")
	}
	tasks, err := dao.DbTasks.SearchTasks(ctx, uid, input.Keyword, input.StartDate, input.EndDate, 10)
	if err != nil {
		glog.Error(ctx, "[queryTasks] search failed: %v", err)
		return "", err
	}
	items := make([]QueryTaskItem, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, QueryTaskItem{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			DueDate:     t.DueDate,
			Priority:    t.Priority,
			TaskType:    t.TaskType,
			Status:      t.Status,
		})
	}
	resultJSON, err := json.Marshal(QueryTasksResult{Tasks: items, Count: len(items)})
	if err != nil {
		return "", err
	}
	return string(resultJSON), nil
}
