import type { AxiosResponse } from "axios";
import service from "@/utils/request";
import type {
  LoginReq,
  LoginRes,
  RegisterReq,
  RegisterRes,
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
