import { ref, computed } from "vue";
import type { Ref, ComputedRef } from "vue";
import { ElMessage } from "element-plus";
import { fetchGetTaskList } from "@/api/task";
import type { TaskItem as ApiTaskItem, Priority } from "@/api/task/type";

/**
 * 组件内部任务类型
 * @typedef {Object} Task
 * @property {string} id 任务ID
 * @property {string} title 任务标题
 * @property {string} deadline 截止日期
 * @property {"high" | "medium" | "low"} priority 优先级
 * @property {string} taskType 任务类型
 * @property {"todo" | "doing" | "done" | "cancelled"} status 任务状态
 */
interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
}

/**
 * 时间模式类型
 * @typedef {"today" | "week"} Mode
 */
type Mode = "today" | "week";

/**
 * 日期范围工具函数
 * @function useDateRange
 * @returns {Object} 包含日期处理相关的方法
 * @returns {function(Date): string} returns.formatDate - 将Date对象格式化为YYYY-MM-DD HH:mm:ss字符串
 * @returns {function(Mode): {startTime: string, endTime: string}} returns.getDateRange - 根据模式获取日期范围
 */
export function useDateRange() {
  /**
   * 将Date对象格式化为"YYYY-MM-DD HH:mm:ss"字符串
   * @param {Date} date 日期对象
   * @returns {string} 格式化后的日期字符串
   */
  const formatDate = (date: Date): string => {
    const pad = (n: number) => n.toString().padStart(2, "0");
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
  };

  /**
   * 根据模式获取日期范围
   * - today: 当天 00:00:00 ~ 23:59:59
   * - week: 本周一 00:00:00 ~ 本周日 23:59:59
   * @param {Mode} modeType 时间模式
   * @returns {{startTime: string, endTime: string}} 日期范围
   */
  const getDateRange = (modeType: Mode): { startTime: string; endTime: string } => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (modeType === "today") {
      const endOfDay = new Date(today);
      endOfDay.setHours(23, 59, 59, 999);
      return {
        startTime: formatDate(today),
        endTime: formatDate(endOfDay),
      };
    }

    const dayOfWeek = today.getDay();
    const diffToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
    const monday = new Date(today);
    monday.setDate(today.getDate() + diffToMonday);

    const sunday = new Date(monday);
    sunday.setDate(monday.getDate() + 6);
    sunday.setHours(23, 59, 59, 999);

    return {
      startTime: formatDate(monday),
      endTime: formatDate(sunday),
    };
  };

  return {
    getDateRange,
    formatDate,
  };
}

/**
 * 优先级选项配置
 * @type {Array<{value: Priority, label: string}>}
 */
const priorities = [
  { value: "high" as Priority, label: "高" },
  { value: "medium" as Priority, label: "中" },
  { value: "low" as Priority, label: "低" },
];

/**
 * 任务类型选项配置
 * @type {Array<{value: string, label: string}>}
 */
const taskTypeOptions = [
  { value: "work", label: "工作" },
  { value: "personal", label: "个人" },
  { value: "project", label: "项目" },
  { value: "adhoc", label: "临时" },
];

/**
 * 任务筛选逻辑 Hook
 * @function useTaskFilter
 * @returns {Object} 筛选相关状态和方法
 * @returns {Ref<Priority[]>} returns.selectedPriorities 当前选中的优先级列表
 * @returns {Ref<string>} returns.selectedType 当前选中的任务类型
 * @returns {function(Ref<ApiTaskItem[]>): ComputedRef<ApiTaskItem[]>} returns.filteredTasks - 返回筛选后的任务计算属性
 * @returns {function(Priority): void} returns.togglePriority - 切换优先级选中状态
 * @returns {typeof priorities} returns.priorities 优先级选项列表
 * @returns {typeof taskTypeOptions} returns.taskTypeOptions 任务类型选项列表
 */
export function useTaskFilter() {
  /** @type {Ref<Priority[]>} */
  const selectedPriorities: Ref<Priority[]> = ref([]);
  /** @type {Ref<string>} */
  const selectedType: Ref<string> = ref("");

  /**
   * 根据选中的优先级和类型筛选任务
   * @param {Ref<ApiTaskItem[]>} tasks 任务列表Ref
   * @returns {ComputedRef<ApiTaskItem[]>} 筛选后的任务列表
   */
  const filteredTasks = (tasks: Ref<ApiTaskItem[]>): ComputedRef<ApiTaskItem[]> => {
    return computed(() => {
      return tasks.value.filter((task) => {
        if (selectedPriorities.value.length > 0) {
          const taskPriorityLower = task.priority.toLowerCase();
          const matchPriority = selectedPriorities.value.some(
            (p) => p.toLowerCase() === taskPriorityLower
          );
          if (!matchPriority) return false;
        }

        if (selectedType.value && task.taskType.toLowerCase() !== selectedType.value) {
          return false;
        }

        return true;
      });
    });
  };

  /**
   * 切换优先级选中状态
   * - 如果已选中则取消选中
   * - 如果未选中则添加到选中列表
   * @param {Priority} priority 优先级
   */
  const togglePriority = (priority: Priority) => {
    const index = selectedPriorities.value.indexOf(priority);
    if (index === -1) {
      selectedPriorities.value.push(priority);
    } else {
      selectedPriorities.value.splice(index, 1);
    }
  };

  return {
    selectedPriorities,
    selectedType,
    filteredTasks,
    togglePriority,
    priorities,
    taskTypeOptions,
  };
}

/**
 * 时间模式切换 Hook
 * @function useTaskMode
 * @returns {Object} 模式切换相关状态和方法
 * @returns {Ref<Mode>} returns.mode 当前模式 ("today" | "week")
 * @returns {function(Mode, function()?: void): void} returns.handleModeChange - 切换模式并执行回调
 */
export function useTaskMode() {
  /** @type {Ref<Mode>} */
  const mode: Ref<Mode> = ref("today");

  /**
   * 切换时间模式
   * @param {Mode} newMode 新模式
   * @param {function(): void} [callback] 切换后执行的回调函数
   */
  const handleModeChange = (newMode: Mode, callback?: () => void) => {
    mode.value = newMode;
    callback?.();
  };

  return {
    mode,
    handleModeChange,
  };
}

/**
 * 任务列表管理 Hook
 * @function useTaskList
 * @returns {Object} 任务列表相关状态和方法
 * @returns {Ref<ApiTaskItem[]>} returns.todayTasks 今日任务列表
 * @returns {Ref<ApiTaskItem[]>} returns.weekTasks 本周任务列表
 * @returns {Ref<ApiTaskItem[]>} returns.completedTasks 已完成任务列表
 * @returns {Ref<boolean>} returns.loading 加载状态
 * @returns {Ref<boolean>} returns.showCompleted 是否显示已完成区块
 * @returns {ComputedRef<number>} returns.todayCount 今日任务数量
 * @returns {ComputedRef<number>} returns.weekCount 本周任务数量
 * @returns {function(Mode?): Promise<void>} returns.fetchTasks - 获取任务列表
 * @returns {function(Task): void} returns.handleTaskUpdate - 处理任务更新
 * @returns {function(ApiTaskItem): Task} returns.convertToTaskItem - 转换任务格式
 */
export function useTaskList() {
  /** @type {Ref<ApiTaskItem[]>} */
  const todayTasks: Ref<ApiTaskItem[]> = ref([]);
  /** @type {Ref<ApiTaskItem[]>} */
  const weekTasks: Ref<ApiTaskItem[]> = ref([]);
  /** @type {Ref<ApiTaskItem[]>} */
  const completedTasks: Ref<ApiTaskItem[]> = ref([]);
  /** @type {Ref<boolean>} */
  const loading: Ref<boolean> = ref(false);
  /** @type {Ref<boolean>} */
  const showCompleted: Ref<boolean> = ref(false);

  /** 今日任务数量 */
  const todayCount = computed(() => todayTasks.value.length);
  /** 本周任务数量 */
  const weekCount = computed(() => weekTasks.value.length);

  const { getDateRange } = useDateRange();

  /**
   * 将API任务对象转换为组件内部任务格式
   * @param {ApiTaskItem} task API返回的任务对象
   * @returns {Task} 组件内部任务格式
   */
  const convertToTaskItem = (task: ApiTaskItem): Task => {
    return {
      id: task.id,
      title: task.title,
      deadline: task.dueDate,
      priority: task.priority as "high" | "medium" | "low",
      taskType: task.taskType,
      status: task.status as "todo" | "doing" | "done" | "cancelled",
    };
  };

  /**
   * 获取任务列表
   * 并行请求今日、本周和已完成任务
   * @param {Mode} [mode="today"] 当前模式，用于确定已完成任务的时间范围
   */
  const fetchTasks = async (mode: Mode = "today") => {
    loading.value = true;
    try {
      const todayRange = getDateRange("today");
      const weekRange = getDateRange("week");
      const currentRange = getDateRange(mode);

      const [todayRes, weekRes, completedRes] = await Promise.all([
        fetchGetTaskList({ startTime: todayRange.startTime, endTime: todayRange.endTime }),
        fetchGetTaskList({ startTime: weekRange.startTime, endTime: weekRange.endTime }),
        fetchGetTaskList({
          startTime: currentRange.startTime,
          endTime: currentRange.endTime,
          status: "done",
        }),
      ]);

      todayTasks.value = todayRes.data.data?.tasks || [];
      weekTasks.value = weekRes.data.data?.tasks || [];
      completedTasks.value = completedRes.data.data?.tasks || [];
    } catch (error) {
      ElMessage.error("获取任务列表失败");
    } finally {
      loading.value = false;
    }
  };

  /**
   * 处理任务状态更新
   * - 更新对应列表中的任务状态
   * - 已完成任务移到completedTasks
   * - 取消任务从列表移除
   * @param {Task} updatedTask 更新后的任务
   */
  const handleTaskUpdate = (updatedTask: Task) => {
    const todayIndex = todayTasks.value.findIndex((t) => t.id === updatedTask.id);
    if (todayIndex !== -1) {
      todayTasks.value[todayIndex] = {
        ...todayTasks.value[todayIndex],
        status: updatedTask.status,
      };
    }

    const weekIndex = weekTasks.value.findIndex((t) => t.id === updatedTask.id);
    if (weekIndex !== -1) {
      weekTasks.value[weekIndex] = {
        ...weekTasks.value[weekIndex],
        status: updatedTask.status,
      };
    }

    if (updatedTask.status === "done") {
      if (todayIndex !== -1) {
        const [movedTask] = todayTasks.value.splice(todayIndex, 1);
        completedTasks.value.unshift(movedTask);
      }
      if (weekIndex !== -1) {
        const [movedTask] = weekTasks.value.splice(weekIndex, 1);
        if (!completedTasks.value.find((t) => t.id === movedTask.id)) {
          completedTasks.value.unshift(movedTask);
        }
      }
    }

    if (updatedTask.status === "cancelled") {
      if (todayIndex !== -1) {
        todayTasks.value.splice(todayIndex, 1);
      }
      if (weekIndex !== -1) {
        weekTasks.value.splice(weekIndex, 1);
      }
    }
  };

  return {
    todayTasks,
    weekTasks,
    completedTasks,
    loading,
    showCompleted,
    todayCount,
    weekCount,
    fetchTasks,
    handleTaskUpdate,
    convertToTaskItem,
  };
}
