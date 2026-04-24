/**
 * 任务相关类型定义
 * @module task/type
 */

/**
 * 任务优先级
 * @type {"high" | "medium" | "low"}
 */
export type Priority = "high" | "medium" | "low";

/**
 * 任务状态
 * @type {"todo" | "doing" | "done" | "cancelled"}
 */
export type TaskStatus = "todo" | "doing" | "done" | "cancelled";

/**
 * 任务类型
 * @type {string}
 */
export type TaskType = "work" | "personal" | "project" | "adhoc";

/**
 * 创建任务请求参数
 * @interface CreateTaskReq
 */
export interface CreateTaskReq {
  /**
   * 任务标题
   * @type {string}
   */
  title: string;

  /**
   * 任务描述
   * @type {string}
   */
  description: string;

  /**
   * 优先级
   * @type {Priority}
   */
  priority: Priority;

  /**
   * 任务类型
   * @type {TaskType}
   */
  taskType: TaskType;

  /**
   * 截止日期
   * @type {string}
   */
  dueDate: string;
}

/**
 * 创建任务响应数据
 * @interface CreateTaskRes
 */
export interface CreateTaskRes {
  /**
   * 任务ID
   * @type {string}
   */
  id: string;
}

/**
 * 删除任务请求参数
 * @interface DeleteTaskReq
 */
export interface DeleteTaskReq {
  /**
   * 任务ID
   * @type {string}
   */
  id: string;
}

/**
 * 删除任务响应数据
 * @interface DeleteTaskRes
 */
export interface DeleteTaskRes {
  /**
   * 响应消息
   * @type {string}
   */
  message: string;
}

/**
 * 更新任务请求参数
 * @interface UpdateTaskReq
 */
export interface UpdateTaskReq {
  /**
   * 任务ID
   * @type {string}
   */
  id: string;

  /**
   * 任务标题
   * @type {string}
   */
  title?: string;

  /**
   * 任务描述
   * @type {string}
   */
  description?: string;

  /**
   * 优先级
   * @type {Priority}
   */
  priority?: Priority;

  /**
   * 任务类型
   * @type {TaskType}
   */
  taskType?: TaskType;

  /**
   * 截止日期
   * @type {string}
   */
  dueDate?: string;
}

/**
 * 更新任务响应数据
 * @interface UpdateTaskRes
 */
export interface UpdateTaskRes {
  /**
   * 响应消息
   * @type {string}
   */
  message: string;
}

/**
 * 获取任务列表请求参数
 * @interface GetTaskListReq
 */
export interface GetTaskListReq {
  /**
   * 开始时间
   * @type {string}
   */
  startTime?: string;

  /**
   * 结束时间
   * @type {string}
   */
  endTime?: string;

  /**
   * 任务类型
   * @type {TaskType}
   */
  taskType?: TaskType;

  /**
   * 优先级
   * @type {Priority}
   */
  priority?: Priority;

  /**
   * 状态
   * @type {TaskStatus}
   */
  status?: TaskStatus;
}

/**
 * 任务项
 * @interface TaskItem
 */
export interface TaskItem {
  /**
   * 任务ID
   * @type {string}
   */
  id: string;

  /**
   * 任务标题
   * @type {string}
   */
  title: string;

  /**
   * 任务描述
   * @type {string}
   */
  description: string;

  /**
   * 优先级
   * @type {Priority}
   */
  priority: Priority;

  /**
   * 任务类型
   * @type {TaskType}
   */
  taskType: TaskType;

  /**
   * 截止日期
   * @type {string}
   */
  dueDate: string;

  /**
   * 状态
   * @type {TaskStatus}
   */
  status: TaskStatus;
}

/**
 * 获取任务列表响应数据
 * @interface GetTaskListRes
 */
export interface GetTaskListRes {
  tasks: TaskItem[];
}

// 更新任务优先级
export interface UpdateTaskPriorityReq {
  id: string;
  priority: Priority;
}

export interface UpdateTaskPriorityRes {
  message: string;
}

// 更新任务状态
export interface UpdateTaskStatusReq {
  id: string;
  status: TaskStatus;
}

export interface UpdateTaskStatusRes {
  message: string;
}
