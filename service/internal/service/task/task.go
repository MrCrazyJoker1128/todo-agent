package task

import (
	"context"
	"errors"
	"fmt"
	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"service/internal/service/modelconfig"
	"time"

	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func parseDueDate(ctx context.Context, raw string) *gtime.Time {
	if raw == "" {
		return nil
	}
	formats := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"}
	for _, format := range formats {
		if parsed, err := time.Parse(format, raw); err == nil {
			return gtime.NewFromTime(parsed.UTC())
		}
	}
	g.Log().Warning(ctx, "无法解析日期格式: ", raw)
	return nil
}

func CreateTask(ctx context.Context, userID string, req *api.CreateTaskReq) (*api.CreateTaskRes, error) {
	g.Log().Info(ctx, "req: %+v", req)
	if req.Title == "" {
		g.Log().Error(ctx, "请输入标题")
		return nil, errors.New("请输入标题")
	}
	if req.Description == "" {
		g.Log().Error(ctx, "请输入描述")
		return nil, errors.New("请输入描述")
	}
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}
	taskID := uuid.New()
	milvusVectorID := uuid.New()
	dueDate := parseDueDate(ctx, req.DueDate)
	tx, err := dao.DbTasks.InsertTask(ctx, &entity.DbTasks{
		Id:             taskID,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		TaskType:       req.TaskType,
		UserId:         userIDUUID,
		DueDate:        dueDate,
		MilvusVectorId: milvusVectorID,
	})
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}
	docs := []*schema.Document{
		{
			ID:      milvusVectorID.String(),
			Content: req.Title + "-" + req.Description,
			MetaData: map[string]any{
				"userId": userID,
				"date":   dueDate,
			},
		},
	}
	if _, err = StoreIndex(ctx, docs); err != nil {
		tx.Rollback()
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("任务向量索引失败: %w", err)
	}
	tx.Commit()
	return &api.CreateTaskRes{ID: taskID.String()}, nil
}

func DeleteTask(ctx context.Context, req *api.DeleteTaskReq) (*api.DeleteTaskRes, error) {
	if req.ID == "" {
		g.Log().Error(ctx, "请输入任务ID")
		return nil, errors.New("请输入任务ID")
	}
	if err := dao.DbTasks.DeleteTask(ctx, req.ID); err != nil {
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("删除任务失败: %w", err)
	}
	return &api.DeleteTaskRes{Message: "任务删除成功"}, nil
}

func UpdateTask(ctx context.Context, req *api.UpdateTaskReq) (*api.UpdateTaskRes, error) {
	g.Log().Info(ctx, "req: %+v", req)
	if req.ID == "" {
		g.Log().Error(ctx, "任务 ID 为空")
		return nil, errors.New("无效的任务 ID")
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.New("无效的任务 ID")
	}
	milvusVectorID := uuid.New()
	dueDate := parseDueDate(ctx, req.DueDate)
	g.Log().Info(ctx, "dueDate: %+v", dueDate)
	tx, err := dao.DbTasks.UpdateTask(ctx, &entity.DbTasks{
		Id:             id,
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		TaskType:       req.TaskType,
		DueDate:        dueDate,
		MilvusVectorId: milvusVectorID,
	})
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("更新任务失败: %w", err)
	}
	tx.Commit()
	return &api.UpdateTaskRes{Message: "任务更新成功"}, nil
}

func UpdateTaskPriority(ctx context.Context, req *api.UpdateTaskPriorityReq) (*api.UpdateTaskPriorityRes, error) {
	if req.ID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	if req.Priority == "" {
		return nil, errors.New("优先级不能为空")
	}
	if err := dao.DbTasks.UpdateTaskPriority(ctx, req.ID, req.Priority); err != nil {
		g.Log().Error(ctx, "更新任务优先级失败: ", err)
		return nil, fmt.Errorf("更新任务优先级失败: %w", err)
	}
	return &api.UpdateTaskPriorityRes{Message: "优先级更新成功"}, nil
}

func UpdateTaskStatus(ctx context.Context, req *api.UpdateTaskStatusReq) (*api.UpdateTaskStatusRes, error) {
	if req.ID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	if req.Status == "" {
		return nil, errors.New("状态不能为空")
	}
	if err := dao.DbTasks.UpdateTaskStatus(ctx, req.ID, req.Status); err != nil {
		g.Log().Error(ctx, "更新任务状态失败: ", err)
		return nil, fmt.Errorf("更新任务状态失败: %w", err)
	}
	return &api.UpdateTaskStatusRes{Message: "状态更新成功"}, nil
}

func GetTaskList(ctx context.Context, userID string, req *api.GetTaskListReq) (*api.GetTaskListRes, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}
	tasks, err := dao.DbTasks.GetTaskList(ctx, req, userIDUUID)
	if err != nil {
		g.Log().Error(ctx, "获取任务列表失败: ", err)
		return nil, fmt.Errorf("获取任务列表失败: %w", err)
	}
	return &api.GetTaskListRes{Tasks: tasks}, nil
}

func StoreIndex(ctx context.Context, docs []*schema.Document) ([]string, error) {
	embedder, err := modelconfig.NewEmbedder(ctx, 30*time.Second)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("创建 embedder 失败: %w", err)
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
		return nil, fmt.Errorf("创建 indexer 失败: %w", err)
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, fmt.Errorf("向量索引存储失败: %w", err)
	}
	return ids, nil
}
