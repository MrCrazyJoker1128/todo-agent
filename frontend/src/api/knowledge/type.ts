import type { PageReq } from "@/types/type";

export type Status = "pending" | "processing" | "completed" | "failed" | "";
/**
 * 获取知识库列表请求参数
 * @interface
 */
export interface GetKnowledgeUrlsSearchOptions {
  /**
   * 知识库作者（中文分词搜索）
   * 支持 zhparser 中文全文检索
   * @example '大知闲闲i'
   */
  author?: string;

  /**
   * 知识库域名（模糊搜索）
   * 用于按网站分组展示
   * @example 'juejin.cn'
   */
  domain?: string;

  /**
   * 知识库标题（中文分词搜索 + 单字匹配）
   * 使用 zhparser 全文检索，支持词和单字匹配
   * @example '面试'
   */
  title?: string;

  /**
   * 处理状态
   * @enum pending|processing|completed|failed
   * @example 'completed'
   */
  status?: Status;
}

export type GetKnowledgeUrlsListRequest = PageReq &
  GetKnowledgeUrlsSearchOptions;

/**
 * 知识库列表项
 * @interface
 */
export interface ServiceapiGetKnowledgeListItem {
  /**
   * 文章作者
   * @example '大知闲闲i'
   */
  author?: string;

  /**
   * 创建时间 (ISO 8601 格式)
   * @format date-time
   * @example '2026-04-16T06:54:24+08:00'
   */
  createdAt?: string;

  /**
   * 域名
   * 用于按网站分组展示
   * @example 'juejin.cn'
   */
  domain?: string;

  /**
   * 知识库记录唯一标识 (UUID)
   */
  id?: string;

  /**
   * 处理状态
   * @enum pending|processing|completed|failed
   * @example 'completed'
   */
  status?: string;

  /**
   * 文章标题（带搜索高亮标签）
   * 当存在搜索关键词时，返回带 <em> 标签高亮的内容
   * @example '<em>面试</em>官：前端批量请求失败 Toast 重复弹窗怎么解决？'
   */
  title?: string;

  /**
   * 用户收藏时间 (ISO 8601 格式)
   * 用于列表排序
   * @format date-time
   * @example '2026-04-16T06:54:24+08:00'
   */
  uploadAt?: string;

  /**
   * 网页完整地址
   * @format uri
   * @example 'https://juejin.cn/post/7504247651577479180'
   */
  url?: string;
}

/**
 * 获取知识库列表响应
 * @interface
 */
export interface GetKnowledgeUrlsListResponse {
  /**
   * 知识库列表
   */
  items?: ServiceapiGetKnowledgeListItem[];

  /**
   * 符合条件的总数
   * 用于分页计算
   * @example 100
   */
  total?: number;
}

/**
 * 创建知识库请求参数
 * @interface
 */
export interface CreateKnowledgeUrlRequest {
  /**
   * 知识库链接
   * 支持 http/https 协议的网页 URL
   * @format uri
   * @example 'https://juejin.cn/post/7504247651577479180'
   */
  url?: string;
}

/**
 * 创建知识库响应
 * @interface
 */
export interface CreateKnowledgeUrlResponse {
  /**
   * 创建成功的知识库记录 ID (UUID)
   */
  id?: string;
}
