export interface GetChatListReq {
  /**
   * 页码
   */
  pageIndex?: number;
  /**
   * 每页数量
   */
  pageSize?: number;
}
export interface ServiceapiChatMessageWithTask {
  chatContent?: string;
  chatCreatedAt?: string;
  chatId?: string;
  description?: string;
  dueDate?: string;
  priority?: string;
  priorityLevel?: string;
  relatedTaskId?: string;
  role?: string;
  status?: string;
  taskCreatedAt?: string;
  taskId?: string;
  taskType?: string;
  taskUpdatedAt?: string;
  title?: string;
}
export interface GetChatListResp {
  list?: ServiceapiChatMessageWithTask[];
  total?: number;
}
