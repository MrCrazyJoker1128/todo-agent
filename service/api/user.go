package api

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"User" summary:"用户登录"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRes struct {
	Token string `json:"token"`
}
type RegisterReq struct {
	g.Meta   `path:"/register" method:"post" tags:"User" summary:"用户注册"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
type RegisterRes struct {
	Message string `json:"message"`
}
type GetUserInfoReq struct {
	g.Meta `path:"/userInfo" method:"get" tags:"User" summary:"通过token获取用户信息"`
}
type GetUserInfoRes struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Id       string `json:"id"`
}
