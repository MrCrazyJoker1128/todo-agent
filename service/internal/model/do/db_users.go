// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DbUsers is the golang structure of table db_users for DAO operations like Where/Data.
type DbUsers struct {
	g.Meta       `orm:"table:db_users, do:true"`
	Id           any         // 主键：用户唯一标识，使用 UUID 避免泄露业务量
	Email        any         // 登录邮箱：唯一索引，用于身份验证
	PasswordHash any         // 密码哈希：存储 bcrypt 加密后的密码，不存明文
	Nickname     any         // 用户昵称：前端显示用，可为空
	CreatedAt    *gtime.Time // 账户创建时间
	UpdatedAt    *gtime.Time // 账户最后更新时间
}
