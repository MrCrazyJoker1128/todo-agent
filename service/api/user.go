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
