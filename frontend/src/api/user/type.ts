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
