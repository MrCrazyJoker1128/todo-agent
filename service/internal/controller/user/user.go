package user

import (
	"context"
	"service/api"
	svc "service/internal/service/user"

	"github.com/gogf/gf/v2/frame/g"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) Login(ctx context.Context, req *api.LoginReq) (res *api.LoginRes, err error) {
	return svc.Login(ctx, req)
}

func (c *UserController) UserRegister(ctx context.Context, req *api.RegisterReq) (res *api.RegisterRes, err error) {
	return svc.UserRegister(ctx, req)
}

func (c *UserController) GetUserInfo(ctx context.Context, req *api.GetUserInfoReq) (res *api.GetUserInfoRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	return svc.GetUserInfo(ctx, userID)
}

func (c *UserController) SaveModelConfig(ctx context.Context, req *api.SaveModelConfigReq) (res *api.SaveModelConfigRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	return svc.SaveModelConfig(ctx, userID, req)
}

func (c *UserController) GetModelConfig(ctx context.Context, req *api.GetModelConfigReq) (res *api.GetModelConfigRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	return svc.GetModelConfig(ctx, userID)
}

func (c *UserController) SaveSystemKEK(ctx context.Context, req *api.SaveSystemKEKReq) (res *api.SaveSystemKEKRes, err error) {
	return svc.SaveSystemKEK(ctx, req)
}

func (c *UserController) GetSystemKEK(ctx context.Context, req *api.GetSystemKEKReq) (res *api.GetSystemKEKRes, err error) {
	return svc.GetSystemKEK(ctx)
}

func (c *UserController) UpdateProfile(ctx context.Context, req *api.UpdateProfileReq) (res *api.UpdateProfileRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	return svc.UpdateProfile(ctx, userID, req)
}

func (c *UserController) UpdatePassword(ctx context.Context, req *api.UpdatePasswordReq) (res *api.UpdatePasswordRes, err error) {
	userID := g.RequestFromCtx(ctx).GetParam("userId").String()
	return svc.UpdatePassword(ctx, userID, req)
}
