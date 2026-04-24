/**
 * @file 日历面板相关 hooks
 * @module components/CalendarPanel/hooks
 */
import { ref, watch, onMounted, onUnmounted, nextTick } from "vue";
import { ElMessage } from "element-plus";
import { fetchGetTaskList, fetchUpdateTaskStatus } from "@/api/task";
import type { TaskItem as ApiTaskItem, TaskType } from "@/api/task/type";

/**
 * 优先级颜色配置
 */
export const PriorityColor = {
  high: "#f56c6c",
  medium: "#e6a23c",
  low: "#67c23a",
} as const;

/**
 * 优先级类型
 * @remarks 从 PriorityColor 的键名推导
 */
export type Priority = keyof typeof PriorityColor;

/**
 * 任务对象
 */
export interface Task {
  /** 任务ID */
  id: string;
  /** 任务标题 */
  title: string;
  /** 截止时间 */
  dueDate: string;
  /** 优先级 */
  priority: Priority;
  /** 任务类型 */
  taskType: TaskType;
  /** 任务状态 */
  status: "todo" | "doing" | "done" | "cancelled";
  /** 任务描述 */
  description: string;
}

/**
 * Element Plus 标签类型
 */
export type TagType = "primary" | "success" | "warning" | "info" | "danger";

const taskTypeLabels: Record<string, string> = {
  work: "工作",
  personal: "个人",
  project: "项目",
  adhoc: "临时",
};

const priorityMap: Record<string, { label: string; type: TagType }> = {
  high: { label: "高", type: "danger" },
  medium: { label: "中", type: "warning" },
  low: { label: "低", type: "success" },
};

/**
 * 日历相关状态和操作
 * @returns 日历状态和辅助方法
 */
export function useCalendar() {
  /** 当前查看的月份 */
  const currentMonth = ref(new Date());
  /** 选中的日期 */
  const selectedDate = ref<string>("");

  const today = new Date();
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, "0")}-${String(today.getDate()).padStart(2, "0")}`;

  return {
    currentMonth,
    selectedDate,
    todayStr,
  };
}

/**
 * 任务数据管理和API调用
 * @returns 任务状态和相关操作方法
 */
export function useTask() {
  /** 当前月份所有任务 */
  const tasks = ref<Task[]>([]);
  /** 选中日期的任务列表 */
  const selectedDateTasks = ref<Task[]>([]);
  /** 加载状态 */
  const loading = ref(false);

  /**
   * 格式化日期为 YYYY-MM-DD HH:mm:ss
   * @param date - Date对象
   */
  const formatDate = (date: Date): string => {
    const pad = (n: number) => n.toString().padStart(2, "0");
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
  };

  /**
   * 获取月份的时间范围
   * @param date - 日期对象
   * @returns 包含 startTime 和 endTime 的对象
   */
  const getMonthRange = (
    date: Date,
  ): { startTime: string; endTime: string } => {
    const year = date.getFullYear();
    const month = date.getMonth();
    const startDate = new Date(year, month, 1);
    const endDate = new Date(year, month + 1, 0, 23, 59, 59);
    return {
      startTime: formatDate(startDate),
      endTime: formatDate(endDate),
    };
  };

  /**
   * 获取日期的时间范围
   * @param dateStr - 日期字符串 YYYY-MM-DD
   * @returns 包含 startTime 和 endTime 的对象
   */
  const getDayRange = (
    dateStr: string,
  ): { startTime: string; endTime: string } => {
    return {
      startTime: `${dateStr} 00:00:00`,
      endTime: `${dateStr} 23:59:59`,
    };
  };

  /**
   * 筛选指定日期的任务
   * @param date - 日期字符串
   * @returns 匹配的任务数组
   */
  const getTasksForDate = (date: string): Task[] => {
    return tasks.value.filter((t) => t.dueDate.startsWith(date));
  };

  /**
   * 获取指定日期的第一个任务
   * @param date - 日期字符串
   * @returns 第一个任务对象，无任务时返回 null
   */
  const getFirstTaskForDate = (date: string): Task | null => {
    const taskList = getTasksForDate(date);
    return taskList.length > 0 ? taskList[0] : null;
  };

  /**
   * 获取当前月份的任务列表
   * @param currentMonth - 当前月份日期对象
   */
  const fetchTasks = async (currentMonth: Date) => {
    loading.value = true;
    try {
      const range = getMonthRange(currentMonth);
      const res = await fetchGetTaskList({
        startTime: range.startTime,
        endTime: range.endTime,
      });
      const apiTasks = res.data.data?.tasks || [];
      tasks.value = apiTasks
        .filter(
          (t: ApiTaskItem) => t.status !== "done" && t.status !== "cancelled",
        )
        .map((t: ApiTaskItem) => ({
          id: t.id,
          title: t.title,
          dueDate: t.dueDate,
          priority: t.priority as Priority,
          taskType: t.taskType,
          status: t.status as Task["status"],
          description: t.description,
        }));
    } catch (error) {
      ElMessage.error("获取任务列表失败");
    } finally {
      loading.value = false;
    }
  };

  /**
   * 获取指定日期的任务列表
   * @param date - 日期字符串
   */
  const fetchTasksForDate = async (date: string) => {
    loading.value = true;
    try {
      const range = getDayRange(date);
      const res = await fetchGetTaskList({
        startTime: range.startTime,
        endTime: range.endTime,
      });
      const apiTasks = res.data.data?.tasks || [];
      selectedDateTasks.value = apiTasks
        .filter(
          (t: ApiTaskItem) => t.status !== "done" && t.status !== "cancelled",
        )
        .map((t: ApiTaskItem) => ({
          id: t.id,
          title: t.title,
          dueDate: t.dueDate,
          priority: t.priority as Priority,
          taskType: t.taskType,
          status: t.status as Task["status"],
          description: t.description,
        }));
    } catch (error) {
      ElMessage.error("获取任务列表失败");
    } finally {
      loading.value = false;
    }
  };

  /**
   * 开始任务
   * @param task - 任务对象
   */
  const handleStart = async (task: Task) => {
    try {
      await fetchUpdateTaskStatus({ id: task.id, status: "doing" });
      const index = selectedDateTasks.value.findIndex((t) => t.id === task.id);
      if (index !== -1) {
        selectedDateTasks.value[index] = { ...task, status: "doing" };
      }
      ElMessage.success("任务已开始");
    } catch {
      ElMessage.error("更新状态失败");
    }
  };

  /**
   * 完成任务
   * @param task - 任务对象
   */
  const handleDone = async (task: Task) => {
    try {
      await fetchUpdateTaskStatus({ id: task.id, status: "done" });
      selectedDateTasks.value = selectedDateTasks.value.filter(
        (t) => t.id !== task.id,
      );
      ElMessage.success("任务已完成");
    } catch {
      ElMessage.error("更新状态失败");
    }
  };

  /**
   * 取消任务
   * @param task - 任务对象
   */
  const handleCancel = async (task: Task) => {
    try {
      await fetchUpdateTaskStatus({ id: task.id, status: "cancelled" });
      selectedDateTasks.value = selectedDateTasks.value.filter(
        (t) => t.id !== task.id,
      );
      ElMessage.success("任务已取消");
    } catch {
      ElMessage.error("取消任务失败");
    }
  };

  return {
    tasks,
    selectedDateTasks,
    loading,
    getTasksForDate,
    getFirstTaskForDate,
    getMonthRange,
    getDayRange,
    fetchTasks,
    fetchTasksForDate,
    handleStart,
    handleDone,
    handleCancel,
  };
}

/**
 * 任务显示相关格式化方法
 * @returns 任务显示辅助函数
 */
export function useTaskDisplay() {
  /**
   * 获取任务类型标签
   * @param type - 任务类型标识
   * @returns 任务类型中文标签
   */
  const getTaskTypeLabel = (type: string): string => {
    return taskTypeLabels[type] || type;
  };

  /**
   * 获取优先级标签
   * @param priority - 优先级
   * @returns 优先级中文标签
   */
  const getPriorityLabel = (priority: Priority): string => {
    return priorityMap[priority.toLowerCase()]?.label || priority;
  };

  /**
   * 获取优先级对应的 Element Plus 标签类型
   * @param priority - 优先级
   * @returns Element Plus Tag 类型
   */
  const getPriorityTagType = (priority: Priority): TagType => {
    return priorityMap[priority.toLowerCase()]?.type || "info";
  };

  /**
   * 格式化选中日期显示
   * @param dateStr - 日期字符串 YYYY-MM-DD
   * @returns 格式化后的文本，如：2024年1月1日
   */
  const formatSelectedDate = (dateStr: string): string => {
    if (!dateStr) return "";
    const [year, month, day] = dateStr.split("-");
    return `${year}年${parseInt(month)}月${parseInt(day)}日`;
  };

  return {
    getTaskTypeLabel,
    getPriorityLabel,
    getPriorityTagType,
    formatSelectedDate,
  };
}

/**
 * 滚动条高度动态计算
 * @description 根据任务列表容器和头部高度自动计算滚动条所需高度
 * @returns 滚动条高度相关状态和方法
 */
export function useScrollbarHeight() {
  /** 任务列表容器引用 */
  const tasksSectionRef = ref<HTMLElement>();
  /** 任务列表头部引用 */
  const tasksHeaderRef = ref<HTMLElement>();
  /** 滚动条高度 */
  const scrollbarHeight = ref(200);

  /**
   * 计算滚动条高度
   * @description 滚动条高度 = 容器高度 - 头部高度 - 头部底部间距(12px)
   */
  const calculateScrollbarHeight = () => {
    nextTick(() => {
      if (tasksSectionRef.value && tasksHeaderRef.value) {
        const sectionHeight = tasksSectionRef.value.clientHeight;
        const headerHeight = tasksHeaderRef.value.clientHeight;
        scrollbarHeight.value = Math.max(
          100,
          sectionHeight - headerHeight - 12,
        );
      }
    });
  };

  watch(tasksSectionRef, (el) => {
    if (el) {
      calculateScrollbarHeight();
    }
  });

  /** 初始化窗口 resize 监听 */
  const initResizeListener = () => {
    window.addEventListener("resize", calculateScrollbarHeight);
  };

  /** 销毁窗口 resize 监听 */
  const destroyResizeListener = () => {
    window.removeEventListener("resize", calculateScrollbarHeight);
  };

  onMounted(() => {
    initResizeListener();
  });

  onUnmounted(() => {
    destroyResizeListener();
  });

  return {
    tasksSectionRef,
    tasksHeaderRef,
    scrollbarHeight,
    calculateScrollbarHeight,
  };
}
