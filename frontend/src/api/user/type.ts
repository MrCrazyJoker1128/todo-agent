/**
 * 用户相关类型定义
 * @module user/type
 */

/**
 * 登录请求参数
 * @interface LoginReq
 */
export interface LoginReq {
  /**
   * 邮箱
   * @type {string}
   */
  email: string;

  /**
   * 密码
   * @type {string}
   */
  password: string;
}

/**
 * 登录响应数据
 * @interface LoginRes
 */
export interface LoginRes {
  /**
   * JWT Token
   * @type {string}
   */
  token: string;
}

/**
 * 注册请求参数
 * @interface RegisterReq
 */
export interface RegisterReq {
  /**
   * 用户名
   * @type {string}
   */
  nickname: string;

  /**
   * 邮箱
   * @type {string}
   */
  email: string;

  /**
   * 密码
   * @type {string}
   */
  password: string;
}

/**
 * 注册响应数据
 * @interface RegisterRes
 */
export interface RegisterRes {
  /**
   * 响应消息
   * @type {string}
   */
  message: string;
}
/**
 * 用户信息响应数据
 * @interface UserInfo
 * @property {string} email - 邮箱
 * @property {string} id - 用户ID
 * @property {string} nickname - 用户名
 */
export interface UserInfo {
  /**
   * 邮箱
   * @type {string}
   */
  email: string;
  /**
   * 用户ID
   * @type {string}
   */
  id: string;
  /**
   * 用户名
   * @type {string}
   */
  nickname: string;
}
/**
 * 保存模型配置请求
 * @interface SaveModelConfigReq
 */
export interface SaveModelConfigReq {
  /** 模型 API Key（明文传输、服务端加密存储） */
  apiKey: string;
  /** 模型 BaseURL，空则用服务端默认 */
  baseUrl?: string;
  /** 模型名，空则用服务端默认 */
  model?: string;
  /** 采样温度，null/undefined = 用默认 0.7，有效范围 [0, 2] */
  temperature?: number | null;
  /** 最大输出 token 数，null/undefined = 用默认 4096 */
  maxTokens?: number | null;
  /** 自定义系统提示词，空 = 用内置模板 */
  systemPrompt?: string;
  /** 嵌入模型服务地址（本地 ollama），空 = 用默认 */
  embedBaseUrl?: string;
  /** 嵌入模型名，空 = 用默认 qwen3-embedding:8b；更换需重建知识库向量 */
  embedModel?: string;
}

/**
 * 模型配置响应（Key 脱敏，仅用于展示配置状态）
 * @interface ModelConfigInfo
 */
export interface ModelConfigInfo {
  /** 是否已配置 */
  configured: boolean;
  /** 脱敏后的 Key，如 sk-cp-****9O7p */
  apiKeyMask: string;
  /** 对话 BaseURL，null = 用内置默认 */
  baseUrl: string | null;
  /** 对话模型名，null = 用内置默认 */
  model: string | null;
  /** 采样温度，null = 用默认 0.7 */
  temperature: number | null;
  /** 最大输出 token 数，null = 用默认 4096 */
  maxTokens: number | null;
  /** 自定义系统提示词，空 = 用内置模板 */
  systemPrompt: string;
  /** 嵌入模型服务地址，null = 用本地 ollama 默认 */
  embedBaseUrl: string | null;
  /** 嵌入模型名，null = 用默认 */
  embedModel: string | null;
}

/**
 * 更新个人资料请求
 * @interface UpdateProfileReq
 */
export interface UpdateProfileReq {
  /** 昵称，空表示不改 */
  nickname?: string;
  /** 邮箱，空表示不改 */
  email?: string;
}

/**
 * 修改密码请求
 * @interface UpdatePasswordReq
 */
export interface UpdatePasswordReq {
  /** 旧密码 */
  oldPassword: string;
  /** 新密码 */
  newPassword: string;
}
