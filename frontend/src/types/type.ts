/**
 * 通用 API 响应格式
 * @interface IResult
 * @template T - data 的类型
 */
export interface IResult<T = any> {
  /**
   * 状态码
   * @type {number}
   * @example 200
   */
  code: number;

  /**
   * 响应数据
   * @type {T}
   */
  data: T;

  /**
   * 响应消息
   * @type {string}
   * @example "success"
   */
  message: string;
}
/**
 * @description 定义加载选项的接口
 * @interface LoadingOption
 * @property {boolean} loading - 是否正在加载
 * @property {string} loadingText - 加载文本
 */
export interface LoadingOption {
  /**
   * @description 是否正在加载
   * @type {boolean}
   * @example true
   */
  loading: boolean;
  /**
   * @description 加载文本
   * @type {string}
   * @example "提交"
   */
  loadingText: string;
}
/**
 * @description 分页请求参数
 * @interface PageReq
 * @property {number} current_page - 当前页码
 * @property {number} per_page - 每页数据量
 */
export interface PageReq {
  /**
   * @description 当前页码
   * @type {number}
   * @example 1
   */
  current_page?: number;
  /**
   * @description 每页数据量
   * @type {number}
   * @example 10
   */
  per_page?: number;
}
export interface PageResp {
  total?: number; // 总数据量
}
export interface ListBaseResp {
  page: PageResp;
}
/**
 * TableColumn接口定义了表格表头的基本属性。
 */
export interface TableColumn {
  minWidth?: string | number;
  width?: string | number;
  label: string;
  prop: string;
  className?: string;
  type?: string;
  isSorted?: boolean;
  align?: string;
  fixed?: TableFixed | boolean;
  slot?: string;
}
