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
type SaveModelConfigReq struct {
	g.Meta       `path:"/model-config" method:"post" tags:"User" summary:"保存模型配置（Key加密存储）"`
	ApiKey       string   `json:"apiKey" description:"模型 API Key，明文传输、加密存储"`
	BaseUrl      string   `json:"baseUrl" description:"模型 BaseURL，空则用默认"`
	Model        string   `json:"model" description:"模型名，空则用默认"`
	Temperature  *float64 `json:"temperature" description:"采样温度，nil=用默认 0.7，有效范围[0,2]"`
	MaxTokens    *int     `json:"maxTokens" description:"最大输出 token 数，nil=用默认 4096，必须 > 0"`
	SystemPrompt string   `json:"systemPrompt" description:"自定义系统提示词，空=用内置模板"`
	EmbedBaseUrl string   `json:"embedBaseUrl" description:"嵌入模型服务地址（本地 ollama），空=用默认"`
	EmbedModel   string   `json:"embedModel" description:"嵌入模型名，空=用默认 qwen3-embedding:8b；更换需重建向量"`
}
type SaveModelConfigRes struct {
	Configured bool `json:"configured"`
}
type GetModelConfigReq struct {
	g.Meta `path:"/model-config" method:"get" tags:"User" summary:"获取模型配置（Key脱敏）"`
}
type GetModelConfigRes struct {
	Configured   bool     `json:"configured"`
	ApiKeyMask   string   `json:"apiKeyMask" description:"脱敏后的 Key，如 sk-cp-****9O7p"`
	BaseUrl      *string  `json:"baseUrl" description:"对话 BaseURL，null=用内置默认"`
	Model        *string  `json:"model" description:"对话模型名，null=用内置默认"`
	Temperature  *float64 `json:"temperature" description:"采样温度，null=用默认 0.7"`
	MaxTokens    *int     `json:"maxTokens" description:"最大输出 token 数，null=用默认 4096"`
	SystemPrompt string   `json:"systemPrompt" description:"自定义系统提示词，空=用内置模板"`
	EmbedBaseUrl *string  `json:"embedBaseUrl" description:"嵌入模型服务地址，null=用本地 ollama 默认"`
	EmbedModel   *string  `json:"embedModel" description:"嵌入模型名，null=用默认"`
}
type SaveSystemKEKReq struct {
	g.Meta `path:"/system-kek" method:"post" tags:"User" summary:"设置系统 KEK（模型 Key 加密密钥入库）"`
	Kek    string `json:"kek" description:"KEK 口令，任意长度，服务端派生 32 字节；不回显"`
}
type SaveSystemKEKRes struct {
	Configured bool `json:"configured"`
}
type GetSystemKEKReq struct {
	g.Meta `path:"/system-kek" method:"get" tags:"User" summary:"查询系统 KEK 是否已配置（不回显值）"`
}
type GetSystemKEKRes struct {
	Configured bool `json:"configured" description:"入库 KEK 是否已配置"`
}
type UpdateProfileReq struct {
	g.Meta   `path:"/profile" method:"post" tags:"User" summary:"更新个人资料（昵称/邮箱）"`
	Nickname string `json:"nickname" description:"昵称，空表示不改"`
	Email    string `json:"email" description:"邮箱，空表示不改；被占用时 400"`
}
type UpdateProfileRes struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}
type UpdatePasswordReq struct {
	g.Meta      `path:"/password" method:"post" tags:"User" summary:"修改密码（需旧密码）"`
	OldPassword string `json:"oldPassword" description:"旧密码"`
	NewPassword string `json:"newPassword" description:"新密码"`
}
type UpdatePasswordRes struct {
	Message string `json:"message"`
}
