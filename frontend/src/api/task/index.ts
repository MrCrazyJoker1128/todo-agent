import type { AxiosResponse } from "axios";
import service from "@/utils/request";
import type {
  CreateTaskReq,
  CreateTaskRes,
  DeleteTaskReq,
  DeleteTaskRes,
  UpdateTaskReq,
  UpdateTaskRes,
  GetTaskListReq,
  GetTaskListRes,
} from "./type";
import type { IResult } from "@/types/type";

/**
 * 创建任务
 * @function fetchCreateTask
 * @param {CreateTaskReq} req - 创建任务请求参数
 * @returns {Promise<AxiosResponse<IResult<CreateTaskRes>>>} 创建任务响应
 */
export function fetchCreateTask(
  req: CreateTaskReq,
): Promise<AxiosResponse<IResult<CreateTaskRes>>> {
  return service.post<CreateTaskRes>("/task/create", req) as unknown as Promise<
    AxiosResponse<IResult<CreateTaskRes>>
  >;
}

/**
 * 删除任务
 * @function fetchDeleteTask
 * @param {DeleteTaskReq} req - 删除任务请求参数
 * @returns {Promise<AxiosResponse<IResult<DeleteTaskRes>>>} 删除任务响应
 */
export function fetchDeleteTask(
  req: DeleteTaskReq,
): Promise<AxiosResponse<IResult<DeleteTaskRes>>> {
  return service.delete<DeleteTaskRes>("/task/delete", { data: req }) as unknown as Promise<
    AxiosResponse<IResult<DeleteTaskRes>>
  >;
}

/**
 * 更新任务
 * @function fetchUpdateTask
 * @param {UpdateTaskReq} req - 更新任务请求参数
 * @returns {Promise<AxiosResponse<IResult<UpdateTaskRes>>>} 更新任务响应
 */
export function fetchUpdateTask(
  req: UpdateTaskReq,
): Promise<AxiosResponse<IResult<UpdateTaskRes>>> {
  return service.put<UpdateTaskRes>("/task/update", req) as unknown as Promise<
    AxiosResponse<IResult<UpdateTaskRes>>
  >;
}

/**
 * 获取任务列表
 * @function fetchGetTaskList
 * @param {GetTaskListReq} [req] - 获取任务列表请求参数（可选）
 * @returns {Promise<AxiosResponse<IResult<GetTaskListRes>>>} 任务列表响应
 */
export function fetchGetTaskList(
  req?: GetTaskListReq,
): Promise<AxiosResponse<IResult<GetTaskListRes>>> {
  return service.get<GetTaskListRes>("/task/list", { params: req }) as unknown as Promise<
    AxiosResponse<IResult<GetTaskListRes>>
  >;
}
