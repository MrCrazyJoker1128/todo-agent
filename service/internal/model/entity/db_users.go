// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// DbUsers is the golang structure for table db_users.
type DbUsers struct {
	Id           uuid.UUID   `json:"id"           orm:"id"            description:"主键：用户唯一标识 (UUID)"` // 主键：用户唯一标识 (UUID)
	Email        string      `json:"email"        orm:"email"         description:"登录邮箱：唯一索引"`        // 登录邮箱：唯一索引
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:"密码哈希：bcrypt 加密"`   // 密码哈希：bcrypt 加密
	Nickname     string      `json:"nickname"     orm:"nickname"      description:"用户昵称：可选"`          // 用户昵称：可选
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"账户创建时间"`           // 账户创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"账户最后更新时间"`         // 账户最后更新时间
}
