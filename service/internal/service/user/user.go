package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"service/api"
	"service/internal/dao"
	"service/internal/model/entity"
	"service/internal/service/sysconfig"
	"service/internal/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// ModelConfigDefaults 返回对话模型默认值。
func ModelConfigDefaults() (baseURL, model string) {
	return "https://api.minimaxi.com/v1", "MiniMax-M2.7"
}

// UserIDFromString 校验 userId 字符串并转为 uuid。
func UserIDFromString(userID string) (uuid.UUID, error) {
	if userID == "" {
		return uuid.Nil, errors.New("没有找到用户，请先登录")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("用户 ID 非法: %w", err)
	}
	return uid, nil
}

// Login 校验邮箱密码，成功返回 token。
func Login(ctx context.Context, req *api.LoginReq) (*api.LoginRes, error) {
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
	token, err := utility.GenerateToken(ctx, result.Id.String())
	if err != nil {
		g.Log().Error(ctx, "token生成失败: %v", err)
		return nil, err
	}
	if err := utility.RedisSet(ctx, "token:"+result.Id.String(), token, utility.DefaultExpire); err != nil {
		g.Log().Error(ctx, "token缓存失败: %v", err)
		return nil, err
	}
	return &api.LoginRes{Token: token}, nil
}

// UserRegister 注册新用户。
func UserRegister(ctx context.Context, req *api.RegisterReq) (*api.RegisterRes, error) {
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
	return &api.RegisterRes{Message: "注册成功"}, nil
}

// GetUserInfo 按 userId 取用户信息。
func GetUserInfo(ctx context.Context, userID string) (*api.GetUserInfoRes, error) {
	result, err := dao.DbUsers.GetUserById(ctx, userID)
	if err != nil {
		g.Log().Error(ctx, "get user by id error: %v", err)
		return nil, err
	}
	if result == nil {
		g.Log().Error(ctx, "user not found")
		return nil, errors.New("没有找到用户，请先登录")
	}
	return &api.GetUserInfoRes{
		Email:    result.Email,
		Nickname: result.Nickname,
		Id:       result.Id.String(),
	}, nil
}

// SaveModelConfig 加密 ApiKey 后入库。
func SaveModelConfig(ctx context.Context, userID string, req *api.SaveModelConfigReq) (*api.SaveModelConfigRes, error) {
	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return nil, errors.New("temperature 必须在 0~2 之间")
	}
	if req.MaxTokens != nil && *req.MaxTokens <= 0 {
		return nil, errors.New("maxTokens 必须大于 0")
	}
	uid, err := UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	// 对话配置：空=存 NULL（不落库），运行时回退内置默认。
	var baseURL *string
	if trimmed := strings.TrimSpace(req.BaseUrl); trimmed != "" {
		baseURL = &trimmed
	}
	var model *string
	if trimmed := strings.TrimSpace(req.Model); trimmed != "" {
		model = &trimmed
	}
	// 嵌入配置：空=存 NULL（不落库），运行时回退本地 ollama 默认。
	// 默认只出现在 placeholder 和运行时兜底，数据库永远存用户真实意图。
	var embedBaseURL *string
	if trimmed := strings.TrimSpace(req.EmbedBaseUrl); trimmed != "" {
		embedBaseURL = &trimmed
	}
	var embedModel *string
	if trimmed := strings.TrimSpace(req.EmbedModel); trimmed != "" {
		embedModel = &trimmed
	}
	// systemPrompt 去首尾空格后存，空串存 NULL（用 nil 表达“未设置”）。
	var systemPrompt *string
	if trimmed := strings.TrimSpace(req.SystemPrompt); trimmed != "" {
		systemPrompt = &trimmed
	}
	// apiKey 留空=不更新已有 Key（前端文案承诺如此）：只改其它配置时无需重输。
	// 无已有配置时必须给 Key，否则无法鉴权调用模型。
	var cipherB64, nonceB64 string
	if req.ApiKey == "" {
		existing, err := dao.DbModelConfigs.GetByUserID(ctx, uid)
		if err != nil {
			g.Log().Error(ctx, "read existing model config failed: %v", err)
			return nil, err
		}
		if existing == nil {
			return nil, errors.New("首次配置必须填写 apiKey")
		}
		cipherB64, nonceB64 = existing.ApiKeyCipher, existing.ApiKeyNonce
	} else {
		// 加密后入库：内存明文不落盘、不打日志。
		// KEK 来源：db_system_configs 优先（sysconfig），缺行回退 env/config（见 utility.kekFromRaw）。
		dbKEK, err := sysconfig.GetKEK(ctx)
		if err != nil {
			g.Log().Error(ctx, "get db kek failed: %v", err)
			return nil, errors.New("系统配置读取失败，请稍后重试")
		}
		cipherB64, nonceB64, err = utility.EncryptModelKeyWithKEK(ctx, req.ApiKey, dbKEK)
		if err != nil {
			g.Log().Error(ctx, "encrypt model key failed: %v", err)
			return nil, errors.New("模型 Key 加密失败，请先在设置页配置系统 KEK（model_key_enc_key）或服务端 MODEL_KEY_ENC_KEY")
		}
	}
	if err := dao.DbModelConfigs.Upsert(ctx, &entity.DbModelConfigs{
		UserId:       uid,
		ApiKeyCipher: cipherB64,
		ApiKeyNonce:  nonceB64,
		BaseUrl:      baseURL,
		Model:        model,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		SystemPrompt: systemPrompt,
		EmbedBaseUrl: embedBaseURL,
		EmbedModel:   embedModel,
	}); err != nil {
		g.Log().Error(ctx, "save model config failed: %v", err)
		return nil, err
	}
	return &api.SaveModelConfigRes{Configured: true}, nil
}

// GetModelConfig 取模型配置，Key 脱敏回显。
func GetModelConfig(ctx context.Context, userID string) (*api.GetModelConfigRes, error) {
	uid, err := UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	cfg, err := dao.DbModelConfigs.GetByUserID(ctx, uid)
	if err != nil {
		g.Log().Error(ctx, "get model config failed: %v", err)
		return nil, err
	}
	if cfg == nil {
		// 未配置：对话/嵌入回包保持 nil，默认只由前端 placeholder 与运行时兜底展示。
		return &api.GetModelConfigRes{Configured: false}, nil
	}
	systemPrompt := ""
	if cfg.SystemPrompt != nil {
		systemPrompt = *cfg.SystemPrompt
	}
	// 脱敏回显：仅用于判断是否已配置，前端不展示明文。
	dbKEK, err := sysconfig.GetKEK(ctx)
	if err != nil {
		g.Log().Error(ctx, "get db kek failed: %v", err)
		return nil, errors.New("系统配置读取失败，请稍后重试")
	}
	plain, err := utility.DecryptModelKeyWithKEK(ctx, cfg.ApiKeyCipher, cfg.ApiKeyNonce, dbKEK)
	if err != nil {
		g.Log().Error(ctx, "decrypt model key failed: %v", err)
		return &api.GetModelConfigRes{Configured: true, ApiKeyMask: "****", BaseUrl: cfg.BaseUrl, Model: cfg.Model, Temperature: cfg.Temperature, MaxTokens: cfg.MaxTokens, SystemPrompt: systemPrompt, EmbedBaseUrl: cfg.EmbedBaseUrl, EmbedModel: cfg.EmbedModel}, nil
	}
	return &api.GetModelConfigRes{
		Configured:   true,
		ApiKeyMask:   utility.MaskAPIKey(plain),
		BaseUrl:      cfg.BaseUrl,
		Model:        cfg.Model,
		Temperature:  cfg.Temperature,
		MaxTokens:    cfg.MaxTokens,
		SystemPrompt: systemPrompt,
		EmbedBaseUrl: cfg.EmbedBaseUrl,
		EmbedModel:   cfg.EmbedModel,
	}, nil
}

// SaveSystemKEK 写入/轮换系统 KEK（入库）。空串拒绝——清 KEK 会锁死所有用户 Key。
func SaveSystemKEK(ctx context.Context, req *api.SaveSystemKEKReq) (*api.SaveSystemKEKRes, error) {
	if err := sysconfig.SetKEK(ctx, req.Kek); err != nil {
		g.Log().Error(ctx, "save system kek failed: %v", err)
		return nil, err
	}
	return &api.SaveSystemKEKRes{Configured: true}, nil
}

// GetSystemKEK 查询入库 KEK 是否已配置（不回显值）。
func GetSystemKEK(ctx context.Context) (*api.GetSystemKEKRes, error) {
	ok, err := sysconfig.HasKEK(ctx)
	if err != nil {
		g.Log().Error(ctx, "get system kek failed: %v", err)
		return nil, errors.New("系统配置读取失败，请稍后重试")
	}
	return &api.GetSystemKEKRes{Configured: ok}, nil
}

// UpdateProfile 更新昵称/邮箱：两者都空直接拒绝；邮箱被占用透出 400。
// token 只带 userId，改邮箱不影响已发 token。
func UpdateProfile(ctx context.Context, userID string, req *api.UpdateProfileReq) (*api.UpdateProfileRes, error) {
	nickname := strings.TrimSpace(req.Nickname)
	email := strings.TrimSpace(req.Email)
	if nickname == "" && email == "" {
		return nil, errors.New("昵称和邮箱不能都为空")
	}
	if email != "" && !strings.Contains(email, "@") {
		return nil, errors.New("邮箱格式不正确")
	}
	uid, err := UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	if err := dao.DbUsers.UpdateProfile(ctx, uid.String(), nickname, email); err != nil {
		if errors.Is(err, dao.ErrEmailTaken) {
			return nil, err
		}
		g.Log().Error(ctx, "update profile failed: %v", err)
		return nil, err
	}
	updated, err := dao.DbUsers.GetUserById(ctx, uid.String())
	if err != nil {
		return nil, err
	}
	return &api.UpdateProfileRes{Email: updated.Email, Nickname: updated.Nickname}, nil
}

// UpdatePassword 校验旧密码后改新密码。
func UpdatePassword(ctx context.Context, userID string, req *api.UpdatePasswordReq) (*api.UpdatePasswordRes, error) {
	if req.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}
	if len(req.NewPassword) < 6 {
		return nil, errors.New("新密码至少 6 位")
	}
	uid, err := UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	if err := dao.DbUsers.UpdatePassword(ctx, uid.String(), req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, dao.ErrPasswordMismatch) {
			return nil, err
		}
		g.Log().Error(ctx, "update password failed: %v", err)
		return nil, err
	}
	return &api.UpdatePasswordRes{Message: "密码修改成功"}, nil
}
