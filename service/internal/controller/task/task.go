package task

import (
	"context"
	"errors"
	"service/api"
	tasksvc "service/internal/service/task"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

type TaskController struct {
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (c *TaskController) CreateTask(ctx context.Context, req *api.CreateTaskReq) (res *api.CreateTaskRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetParam("userId").String()
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("无效的用户ID")
	}
	return tasksvc.CreateTask(ctx, userID, req)
}

func (c *TaskController) DeleteTask(ctx context.Context, req *api.DeleteTaskReq) (res *api.DeleteTaskRes, err error) {
	return tasksvc.DeleteTask(ctx, req)
}

func (c *TaskController) UpdateTask(ctx context.Context, req *api.UpdateTaskReq) (res *api.UpdateTaskRes, err error) {
	return tasksvc.UpdateTask(ctx, req)
}

func (c *TaskController) UpdateTaskPriority(ctx context.Context, req *api.UpdateTaskPriorityReq) (res *api.UpdateTaskPriorityRes, err error) {
	return tasksvc.UpdateTaskPriority(ctx, req)
}

func (c *TaskController) UpdateTaskStatus(ctx context.Context, req *api.UpdateTaskStatusReq) (res *api.UpdateTaskStatusRes, err error) {
	return tasksvc.UpdateTaskStatus(ctx, req)
}

func (c *TaskController) GetTaskList(ctx context.Context, req *api.GetTaskListReq) (res *api.GetTaskListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetParam("userId").String()
	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("无效的用户ID")
	}
	return tasksvc.GetTaskList(ctx, userID, req)
}
