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
	Id           uuid.UUID   `json:"id"           orm:"id"            description:"主键：用户唯一标识，使用 UUID 避免泄露业务量"`  // 主键：用户唯一标识，使用 UUID 避免泄露业务量
	Email        string      `json:"email"        orm:"email"         description:"登录邮箱：唯一索引，用于身份验证"`           // 登录邮箱：唯一索引，用于身份验证
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:"密码哈希：存储 bcrypt 加密后的密码，不存明文"` // 密码哈希：存储 bcrypt 加密后的密码，不存明文
	Nickname     string      `json:"nickname"     orm:"nickname"      description:"用户昵称：前端显示用，可为空"`             // 用户昵称：前端显示用，可为空
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"账户创建时间"`                     // 账户创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"账户最后更新时间"`                   // 账户最后更新时间
}
