package task

import (
	"context"
	"errors"
	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type TaskController struct {
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (c *TaskController) CreateTask(ctx context.Context, req *api.CreateTaskReq) (res *api.CreateTaskRes, err error) {
	g.Log().Info(ctx, "req: %+v", req)
	if req.Title == "" {
		g.Log().Error(ctx, "请输入标题")
		return nil, errors.New("请输入标题")
	}
	if req.Description == "" {
		g.Log().Error(ctx, "请输入描述")
		return nil, errors.New("请输入描述")
	}
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}
	taskId := uuid.New()
	milvusVectorId := uuid.New()
	var dueDate *gtime.Time
	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
		if err == nil {
			dueDate = gtime.NewFromTime(t)
		}
	}
	tx, err := dao.DbTasks.InsertTask(ctx, &entity.DbTasks{
		Id:             taskId,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		TaskType:       req.TaskType,
		UserId:         userIdUUID,
		DueDate:        dueDate,
		MilvusVectorId: milvusVectorId,
	})
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		g.Log().Error(ctx, err)
		return nil, err
	}
	docs := []*schema.Document{
		{
			ID:      milvusVectorId.String(),
			Content: req.Title + "-" + req.Description,
			MetaData: map[string]any{
				"userId": userId,
				"date":   dueDate,
			},
		},
	}
	_, err = StoreIndex(ctx, docs)
	if err != nil {
		tx.Rollback()
		g.Log().Error(ctx, err)
		return nil, err
	}
	tx.Commit()
	return &api.CreateTaskRes{
		ID: taskId.String(),
	}, nil
}

func (c *TaskController) DeleteTask(ctx context.Context, req *api.DeleteTaskReq) (res *api.DeleteTaskRes, err error) {
	if req.ID == "" {
		g.Log().Error(ctx, "请输入任务ID")
		return nil, errors.New("请输入任务ID")
	}
	err = dao.DbTasks.DeleteTask(ctx, req.ID)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return &api.DeleteTaskRes{
		Message: "任务删除成功",
	}, nil
}
func (c *TaskController) UpdateTask(ctx context.Context, req *api.UpdateTaskReq) (res *api.UpdateTaskRes, err error) {
	g.Log().Info(ctx, "req: %+v", req)
	if req.ID == "" {
		g.Log().Error(ctx, "任务 ID 为空")
		return nil, errors.New("无效的任务 ID")
	}
	Id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.New("无效的任务 ID")
	}
	milvusVectorId := uuid.New()
	var dueDate *gtime.Time
	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
		if err == nil {
			dueDate = gtime.NewFromTime(t)
		}
	}
	tx, err := dao.DbTasks.UpdateTask(ctx, &entity.DbTasks{
		Id:             Id,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		TaskType:       req.TaskType,
		DueDate:        dueDate,
		MilvusVectorId: milvusVectorId,
	})
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		g.Log().Error(ctx, err)
		return nil, err
	}
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	docs := []*schema.Document{
		{
			ID:      milvusVectorId.String(),
			Content: req.Title + "-" + req.Description,
			MetaData: map[string]any{
				"userId": userId,
				"date":   dueDate,
			},
		},
	}
	_, err = StoreIndex(ctx, docs)
	if err != nil {
		tx.Rollback()
		g.Log().Error(ctx, err)
		return nil, err
	}
	tx.Commit()
	return &api.UpdateTaskRes{
		Message: "任务更新成功",
	}, nil
}

func (c *TaskController) UpdateTaskPriority(ctx context.Context, req *api.UpdateTaskPriorityReq) (res *api.UpdateTaskPriorityRes, err error) {
	if req.ID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	if req.Priority == "" {
		return nil, errors.New("优先级不能为空")
	}

	err = dao.DbTasks.UpdateTaskPriority(ctx, req.ID, req.Priority)
	if err != nil {
		g.Log().Error(ctx, "更新任务优先级失败: ", err)
		return nil, err
	}

	return &api.UpdateTaskPriorityRes{
		Message: "优先级更新成功",
	}, nil
}

func (c *TaskController) UpdateTaskStatus(ctx context.Context, req *api.UpdateTaskStatusReq) (res *api.UpdateTaskStatusRes, err error) {
	if req.ID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	if req.Status == "" {
		return nil, errors.New("状态不能为空")
	}

	err = dao.DbTasks.UpdateTaskStatus(ctx, req.ID, req.Status)
	if err != nil {
		g.Log().Error(ctx, "更新任务状态失败: ", err)
		return nil, err
	}

	return &api.UpdateTaskStatusRes{
		Message: "状态更新成功",
	}, nil
}

func (c *TaskController) GetTaskList(ctx context.Context, req *api.GetTaskListReq) (res *api.GetTaskListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := r.GetParam("userId").String()
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}

	tasks, err := dao.DbTasks.GetTaskList(ctx, req, userIdUUID)
	if err != nil {
		g.Log().Error(ctx, "获取任务列表失败: ", err)
		return nil, err
	}

	return &api.GetTaskListRes{
		Tasks: tasks,
	}, nil
}
func StoreIndex(ctx context.Context, docs []*schema.Document) ([]string, error) {
	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3-embedding:8b",
		Timeout: 30 * time.Second,
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: "localhost:19530",
			DBName:  "AwesomeEino",
		},
		Collection: "test",
		Vector: &milvus2.VectorConfig{
			Dimension:    4096,
			MetricType:   milvus2.COSINE,
			IndexBuilder: milvus2.NewHNSWIndexBuilder().WithM(16).WithEfConstruction(200),
		},
		Embedding: embedder,
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	return ids, nil
}
