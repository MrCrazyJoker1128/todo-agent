import type { IResult, LoadingOption, PageReq, PageResp } from "@/types/type";
import { useThrottleFn } from "@vueuse/core";
import type { AxiosResponse } from "axios";
import type { Ref } from "vue";

/**
 * 定义请求钩子的接口
 */
interface RequestHook {
  copyLoadingOption: LoadingOption;
  func<T>(data: T): void;
}

/**
 * 定义请求选项的接口
 */
interface RequestOption<T> {
  loadingText?: string;
  requestFunc: Function;
  success: (res: IResult<T>) => void;
  fail?: (err: any) => void;
  finally?: () => void;
  time?: number;
}

/**
 * 自定义钩子，用于处理请求
 * @param options 请求选项，包括加载文本、请求函数、成功回调、失败回调、最终回调和延迟时间
 * @returns {
 * 	copyLoadingOption: 只读的加载选项副本，包含加载状态和加载文本
 * 	func: 处理请求的函数，接受请求数据作为参数
 * }
 */
export function useRequestHook<U>(options: RequestOption<U>): RequestHook {
  // 初始化加载选项
  let loadingOption = reactive<LoadingOption>({
    loading: false,
    loadingText: options.loadingText || "提交",
  });

  /**
   * 处理请求的函数
   * @param data 请求的数据
   */
  function hanldeFunc<T>(data: T) {
    // 请求开始时，设置加载状态为true，加载文本为"提交中"
    loadingOption.loading = true;
    loadingOption.loadingText = "提交中";

    // 执行请求函数，并处理成功、失败和最终状态
    options
      .requestFunc(data)
      .then((res: AxiosResponse<IResult<U>>) => {
        options.success(res.data);
      })
      .catch((err: any) => {
        // 如果有失败回调，则执行失败回调
        if (options.fail) {
          options.fail(err);
        }
      })
      .finally(() => {
        // 请求结束后，恢复加载状态为false，加载文本恢复初始值
        loadingOption.loading = false;
        loadingOption.loadingText = options.loadingText || "提交";

        // 如果有最终回调，则执行最终回调
        if (options.finally) {
          options.finally();
        }
      });
  }

  // 创建只读的加载选项副本
  const copyLoadingOption = readonly<LoadingOption>(loadingOption);

  // 使用节流函数包装请求处理函数，以限制请求频率
  const func = useThrottleFn((data) => {
    hanldeFunc(data);
  }, options.time ?? 1000);

  // 返回请求钩子
  return { copyLoadingOption, func };
}

/**
 * 定义分页配置接口，继承自PageBase
 * @extends PageBase
 */
export interface pageOptionType extends PageReq {
  total: number; // 总数据数
  loading: boolean; // 加载状态
}

/**
 * 定义页面钩子选项接口
 */
interface PageHookOption {
  pageOption: pageOptionType; // 分页配置
  func: Function; // 请求函数
}

/**
 * 定义页面选项接口
 * @template T - 期望的响应数据类型
 */
interface PageOptions<T> {
  per_page?: number; // 每页数据量，默认为10
  requestFunc: Function; // 请求数据的函数
  success: (res: IResult<T>) => void; // 成功时的回调函数
  fail?: (err: any) => void; // 失败时的回调函数
  finally?: () => void; // 无论成功或失败都会执行的回调函数
}

/**
 * 导出一个分页钩子函数
 * @template U - 期望的响应数据类型
 * @param {PageOptions<U>} options - 分页选项
 * @returns {PageHookOption} - 页面钩子选项，包含分页配置和请求函数
 */
export function usePageHook<U>(options: PageOptions<U>): PageHookOption {
  // 初始化分页配置
  const pageOption = reactive<pageOptionType>({
    loading: false, // 初始加载状态为false
    current_page: 1, // 初始页码为1
    per_page: options.per_page ?? 10, // 每页数据量，如果options中未提供则默认为10
    total: 0, // 初始总数据数为0
  });

  /**
   * 请求函数，用于获取分页数据
   * @template T - 期望的请求数据类型
   * @param {T} data - 请求数据
   */
  function func<T extends Object>(data?: T) {
    pageOption.loading = true; // 请求开始时设置加载状态为true
    // 执行请求数据的函数，传入请求数据、当前页码和每页数据量
    options
      .requestFunc({
        ...data,
        current_page: pageOption.current_page,
        per_page: pageOption.per_page,
      })
      .then((res: AxiosResponse<IResult<U & PageResp>>) => {
        pageOption.total = res.data.data.total ?? 0; // 更新总数据数
        options.success(res.data); // 调用成功时的回调函数，传入响应数据
      })
      .catch((err: any) => {
        // 如果请求抛出错误
        if (options.fail) {
          // 如果提供了失败时的回调函数
          options.fail(err); // 调用失败时的回调函数，传入错误数据
        }
      })
      .finally(() => {
        // 无论成功或失败都会执行
        pageOption.loading = false; // 请求结束时设置加载状态为false
        if (options.finally) {
          // 如果提供了无论成功或失败都会执行的回调函数
          options.finally(); // 调用该回调函数
        }
      });
  }

  return { pageOption, func }; // 返回分页配置和请求函数
}
/**
 * 使用防抖动钩子函数。
 *
 * 此函数为Vue的自定义hook，用于实现输入值的防抖动处理。在某些场景下，比如表单输入，我们可能不希望每次输入都触发某些操作，
 * 而是希望用户输入“稳定”一段时间后（比如200毫秒）再触发。这就是防抖动的作用。
 *
 * @param value 输入的值。
 * @param delay 延迟执行的时间，单位为毫秒，默认为200毫秒。
 * @returns 返回一个自定义的ref对象，该对象的get方法用于获取当前的值，set方法用于更新值并应用防抖动逻辑。
 */
export function useDebouncedHook<T>(value: T, delay: number = 200): Ref<T> {
  let timeout: ReturnType<typeof setTimeout>;
  return customRef((track, trigger) => {
    return {
      get() {
        track();
        return value;
      },
      set(newValue) {
        clearTimeout(timeout);
        timeout = setTimeout((): void => {
          value = newValue;
          trigger();
        }, delay);
      },
    };
  });
}
/**
 * 自定义钩子用于管理一个响应式状态对象，并提供重置状态的功能
 * @param initialState 初始化状态对象，作为响应式状态的基线
 * @returns 返回一个对象，包含响应式状态和重置状态的函数
 */
export function useReactiveReset<T extends object>(
  initialState: T,
): { state: T; resetState: (newState?: Partial<T>) => void } {
  // 创建响应式对象
  const state = reactive<T>({ ...initialState }) as T;

  /**
   * 重置 state 到初始状态或指定状态
   * @param newState 可选的部分状态更新，不传则完全重置到初始状态
   */
  function resetState(newState: Partial<T> = {}): void {
    Object.assign(state, { ...initialState, ...newState });
  }

  // 返回响应式状态和重置状态的函数
  return {
    state,
    resetState,
  };
}
/**
 * 导出一个用于计算表格高度的钩子函数
 * @returns {Object} - 包含表格高度的响应式变量和搜索区域的引用
 */
export function useCalcTableHeight() {
  // 初始化表格高度为 100%，使用 ref 创建响应式变量
  const tableHeight = ref<string>("100%");
  // 获取搜索区域的引用，用于后续计算表格高度，使用 useTemplateRef 函数
  const searchAreaRef = useTemplateRef<HTMLDivElement>("searchAreaRef");
  /**
   * 在组件挂载完成后，调整表格高度
   * 根据搜索区域的实际高度来计算表格的高度，以优化布局
   */
  useResizeObserver(searchAreaRef, (entries) => {
    const entry = entries[0];
    const { height } = entry.contentRect;
    tableHeight.value = `calc(100% - ${height ?? 0}px)`;
  });
  // 返回包含表格高度和搜索区域引用的对象
  return {
    tableHeight,
    searchAreaRef,
  };
}
