// =================================================================================
// 手写（非 gf gen dao 产物）：表 db_knowledge_bases 的读写方法。
// =================================================================================

package dao

import (
	"context"

	"service/internal/dao/internal"
	"service/internal/model/do"
	"service/internal/model/entity"

	"github.com/google/uuid"
)

// dbKnowledgeBasesDao is the data access object for the table db_knowledge_bases.
type dbKnowledgeBasesDao struct {
	*internal.DbKnowledgeBasesDao
}

var (
	// DbKnowledgeBases is a globally accessible object for table db_knowledge_bases operations.
	DbKnowledgeBases = dbKnowledgeBasesDao{internal.NewDbKnowledgeBasesDao()}
)

// InsertBase 插入知识库。
func (dao *dbKnowledgeBasesDao) InsertBase(ctx context.Context, base *entity.DbKnowledgeBases) error {
	_, err := dao.Ctx(ctx).Data(do.DbKnowledgeBases{
		Id:          base.Id,
		UserId:      base.UserId,
		Name:        base.Name,
		Description: base.Description,
	}).Insert()
	return err
}

// ListByUserID 按用户取知识库列表。
func (dao *dbKnowledgeBasesDao) ListByUserID(ctx context.Context, userId uuid.UUID) ([]*entity.DbKnowledgeBases, error) {
	var list []*entity.DbKnowledgeBases
	err := dao.Ctx(ctx).Where(do.DbKnowledgeBases{UserId: userId}).Order("created_at desc").Scan(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteBase 删除知识库（文件靠 ES 侧清理，本函数只删库行）。
func (dao *dbKnowledgeBasesDao) DeleteBase(ctx context.Context, id, userId uuid.UUID) error {
	_, err := dao.Ctx(ctx).Where(do.DbKnowledgeBases{Id: id, UserId: userId}).Delete()
	return err
}
