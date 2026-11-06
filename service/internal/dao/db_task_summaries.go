// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_task_summaries 的读写方法。
// =================================================================================

package dao

import (
	"context"

	"service/internal/dao/internal"
	"service/internal/model/do"
	"service/internal/model/entity"

	"github.com/google/uuid"
)

// dbTaskSummariesDao is the data access object for the table db_task_summaries.
type dbTaskSummariesDao struct {
	*internal.DbTaskSummariesDao
}

var (
	// DbTaskSummaries is a globally accessible object for table db_task_summaries operations.
	DbTaskSummaries = dbTaskSummariesDao{internal.NewDbTaskSummariesDao()}
)

// InsertSummary 插入一条工作总结记录。
func (dao *dbTaskSummariesDao) InsertSummary(ctx context.Context, summary *entity.DbTaskSummaries) error {
	_, err := dao.Ctx(ctx).Data(do.DbTaskSummaries{
		Id:        summary.Id,
		UserId:    summary.UserId,
		Period:    summary.Period,
		StartDate: summary.StartDate,
		EndDate:   summary.EndDate,
		TaskType:  summary.TaskType,
		Priority:  summary.Priority,
		Title:     summary.Title,
		Content:   summary.Content,
		Stats:     summary.Stats,
	}).Insert()
	return err
}

// ListByUserID 按用户取总结记录，limit<=0 则 20。
func (dao *dbTaskSummariesDao) ListByUserID(ctx context.Context, userId uuid.UUID, limit int) ([]*entity.DbTaskSummaries, error) {
	if limit <= 0 {
		limit = 20
	}
	var list []*entity.DbTaskSummaries
	err := dao.Ctx(ctx).Where(do.DbTaskSummaries{UserId: userId}).Order("created_at desc").Limit(limit).Scan(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
