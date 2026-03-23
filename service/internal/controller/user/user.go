package user

import (
	"context"
	"errors"
	"service/api"
	"service/internal/dao"
	"service/internal/utility"

	"github.com/gogf/gf/v2/frame/g"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) Login(ctx context.Context, req *api.LoginReq) (res *api.LoginRes, err error) {
	if req.Email == "" || req.Password == "" {
		g.Log().Error(ctx, "email or password is empty")
		return nil, errors.New("email or password is empty")
	}
	result, err := dao.DbUsers.GetUserByEmail(ctx, req.Email)
	if err != nil {
		g.Log().Error(ctx, "get user by email error: %v", err)
		return nil, err
	}
	if result == nil {
		g.Log().Error(ctx, "user not found")
		return nil, errors.New("没有找到用户，请先注册")
	}
	if result.PasswordHash != req.Password {
		g.Log().Error(ctx, "password error")
		return nil, errors.New("邮箱或密码错误，请重新输入")
	}
	// 登录成功，返回token
	token, err := utility.GenerateToken(ctx, result.Id.String())
	utility.RedisSet(ctx, "token:"+result.Id.String(), token, utility.DefaultExpire)
	if err != nil {
		g.Log().Error(ctx, "token生成失败: %v", err)
		return nil, err
	}
	return &api.LoginRes{
		Token: token,
	}, nil
}

func (c *UserController) UserRegister(ctx context.Context, req *api.RegisterReq) (res *api.RegisterRes, err error) {
	if req.Email == "" || req.Password == "" || req.Nickname == "" {
		g.Log().Error(ctx, "email or password or nickname is empty")
		return nil, errors.New("email or password or nickname is empty")
	}
	result, err := dao.DbUsers.GetUserByEmail(ctx, req.Email)
	if err != nil {
		g.Log().Error(ctx, "get user by email error: %v", err)
		return nil, err
	}
	if result != nil {
		g.Log().Error(ctx, "email already registered")
		return nil, errors.New("邮箱已被注册")
	}
	id, err := dao.DbUsers.Register(ctx, req)
	if err != nil {
		g.Log().Error(ctx, "register error: %v", err)
		return nil, err
	}
	g.Log().Info(ctx, "注册成功, id: %v", id)
	return &api.RegisterRes{
		Message: "注册成功",
	}, nil
}
