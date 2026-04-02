package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CreateTaskReq struct {
	g.Meta      `path:"/create" method:"post" tags:"Task" summary:"创建任务"`
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=2048"`
	Priority    string `json:"priority"`
	TaskType    string `json:"taskType"`
	DueDate     string `json:"dueDate"`
}

type CreateTaskRes struct {
	ID string `json:"id"`
}
type DeleteTaskReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"Task" summary:"删除任务"`
	ID     string `json:"id" binding:"required"`
}
type DeleteTaskRes struct {
	Message string `json:"message"`
}
type UpdateTaskReq struct {
	g.Meta      `path:"/update" method:"put" tags:"Task" summary:"更新任务"`
	ID          string `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	TaskType    string `json:"taskType"`
	DueDate     string `json:"dueDate"`
	Status      string `json:"status"`
}

type UpdateTaskRes struct {
	Message string `json:"message"`
}

// 更新任务优先级
type UpdateTaskPriorityReq struct {
	g.Meta   `path:"/update-priority" method:"put" tags:"Task" summary:"更新任务优先级"`
	ID       string `json:"id" binding:"required"`
	Priority string `json:"priority" binding:"required"`
}

type UpdateTaskPriorityRes struct {
	Message string `json:"message"`
}

// 更新任务状态
type UpdateTaskStatusReq struct {
	g.Meta `path:"/update-status" method:"put" tags:"Task" summary:"更新任务状态"`
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type UpdateTaskStatusRes struct {
	Message string `json:"message"`
}

type GetTaskListReq struct {
	g.Meta    `path:"/list" method:"get" tags:"Task" summary:"获取任务列表"`
	StartTime string `json:"startTime" description:"任务开始时间"`
	EndTime   string `json:"endTime" description:"任务结束时间"`
	TaskType  string `json:"taskType" description:"任务类型"`
	Priority  string `json:"priority" description:"任务优先级"`
	Status    string `json:"status" description:"任务状态"`
}

type TaskItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	TaskType    string `json:"taskType"`
	DueDate     string `json:"dueDate"`
	Status      string `json:"status"`
}

type GetTaskListRes struct {
	Tasks []*TaskItem `json:"tasks"`
}
