import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
} from "axios";
import { ElMessage } from "element-plus";
import type { IResult } from "@/types/type";

const BASE_URL = import.meta.env.VITE_APP_BASE_API || "";

const service: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 50000,
  headers: {
    "Content-Type": "application/json",
  },
}) as AxiosInstance & {
  get<T = any, R = AxiosResponse<IResult<T>>>(
    url: string,
    config?: AxiosRequestConfig,
  ): Promise<R>;
  post<T = any, R = AxiosResponse<IResult<T>>>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<R>;
  put<T = any, R = AxiosResponse<IResult<T>>>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<R>;
  delete<T = any, R = AxiosResponse<IResult<T>>>(
    url: string,
    config?: AxiosRequestConfig,
  ): Promise<R>;
};

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    if (config.method === "get" || config.method === "delete") {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    } else if (config.data instanceof FormData) {
      // 文件上传：不展平 FormData；去掉 JSON 头让 axios 自动带 multipart boundary。
      const headers = config.headers;
      if (headers && "delete" in headers) {
        const del = headers.delete;
        if (typeof del === "function") del.call(headers, "Content-Type");
      }
    } else {
      config.data = {
        ...config.data,
      };
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

service.interceptors.response.use(
  (response: AxiosResponse<IResult<any>>) => {
    const { code, message } = response.data;

    if (code === 200) {
      return response;
    }

    ElMessage.error(message || "请求失败");
    return Promise.reject(new Error(message || "请求失败"));
  },
  (error) => {
    const message =
      error.response?.data?.message || error.message || "网络错误";
    ElMessage.error(message);
    return Promise.reject(error);
  },
);

export const request = service;

export default service;
