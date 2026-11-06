import type { AxiosResponse } from "axios";
import service from "@/utils/request";
import type {
  LoginReq,
  LoginRes,
  ModelConfigInfo,
  RegisterReq,
  RegisterRes,
  SaveModelConfigReq,
  UpdatePasswordReq,
  UpdateProfileReq,
  UserInfo,
} from "./type";
import type { IResult } from "@/types/type";

/**
 * 用户登录
 * @function fetchLogin
 * @param {LoginReq} req - 登录请求参数
 * @returns {Promise<AxiosResponse<IResult<LoginRes>>>} 登录响应
 */
export function fetchLogin(
  req: LoginReq,
): Promise<AxiosResponse<IResult<LoginRes>>> {
  return service.post<LoginRes>("/user/login", req) as unknown as Promise<
    AxiosResponse<IResult<LoginRes>>
  >;
}

/**
 * 用户注册
 * @function fetchRegister
 * @param {RegisterReq} req - 注册请求参数
 * @returns {Promise<AxiosResponse<IResult<RegisterRes>>>} 注册响应
 */
export function fetchRegister(
  req: RegisterReq,
): Promise<AxiosResponse<IResult<RegisterRes>>> {
  return service.post<RegisterRes>("/user/register", req) as unknown as Promise<
    AxiosResponse<IResult<RegisterRes>>
  >;
}

/**
 * 获取用户信息
 * @function fetchGetUserInfo
 * @returns {Promise<AxiosResponse<IResult<UserInfo>>>} 用户信息响应
 */
export function fetchGetUserInfo(): Promise<AxiosResponse<IResult<UserInfo>>> {
  return service.get<UserInfo>("/user/userInfo") as unknown as Promise<
    AxiosResponse<IResult<UserInfo>>
  >;
}
/**
 * 保存模型配置（Key 服务端加密存储）
 * @function fetchSaveModelConfig
 */
export function fetchSaveModelConfig(
  req: SaveModelConfigReq,
): Promise<AxiosResponse<IResult<{ configured: boolean }>>> {
  return service.post<{ configured: boolean }>(
    "/user/model-config",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<{ configured: boolean }>>>;
}

/**
 * 获取模型配置（Key 脱敏回显）
 * @function fetchGetModelConfig
 */
export function fetchGetModelConfig(): Promise<
  AxiosResponse<IResult<ModelConfigInfo>>
> {
  return service.get<ModelConfigInfo>(
    "/user/model-config",
  ) as unknown as Promise<AxiosResponse<IResult<ModelConfigInfo>>>;
}

/**
 * 设置系统 KEK（入库，不回显）
 * @function fetchSaveSystemKEK
 */
export function fetchSaveSystemKEK(req: {
  kek: string;
}): Promise<AxiosResponse<IResult<{ configured: boolean }>>> {
  return service.post<{ configured: boolean }>(
    "/user/system-kek",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<{ configured: boolean }>>>;
}

/**
 * 查询系统 KEK 是否已配置（不回显值）
 * @function fetchGetSystemKEK
 */
export function fetchGetSystemKEK(): Promise<
  AxiosResponse<IResult<{ configured: boolean }>>
> {
  return service.get<{ configured: boolean }>(
    "/user/system-kek",
  ) as unknown as Promise<AxiosResponse<IResult<{ configured: boolean }>>>;
}

/**
 * 更新个人资料（昵称/邮箱）
 * @function fetchUpdateProfile
 */
export function fetchUpdateProfile(
  req: UpdateProfileReq,
): Promise<AxiosResponse<IResult<UserInfo>>> {
  return service.post<UserInfo>(
    "/user/profile",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<UserInfo>>>;
}

/**
 * 修改密码（需旧密码）
 * @function fetchUpdatePassword
 */
export function fetchUpdatePassword(
  req: UpdatePasswordReq,
): Promise<AxiosResponse<IResult<{ message: string }>>> {
  return service.post<{ message: string }>(
    "/user/password",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<{ message: string }>>>;
}
