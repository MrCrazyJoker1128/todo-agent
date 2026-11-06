import { computed, onMounted, onUnmounted, ref } from "vue";
import { fetchGetTaskList } from "@/api/task";
import { fetchGetModelConfig } from "@/api/user";
import type { TaskItem } from "@/api/task/type";
import type { ModelConfigInfo } from "@/api/user/type";
import { useDateRange } from "@/components/ListPanel/hooks";
import eventBus from "@/utils/eventbus";

export function useHomeDashboard() {
  const { getDateRange } = useDateRange();

  const loading = ref(false);
  const configInfo = ref<ModelConfigInfo | null>(null);
  const todayTasks = ref<TaskItem[]>([]);
  const weekTasks = ref<TaskItem[]>([]);
  const weekDoneCount = ref(0);

  const todayPending = computed(() => todayTasks.value.length);
  const weekTotal = computed(() => weekTasks.value.length);
  const weekDone = computed(() => weekDoneCount.value);
  const weekRate = computed(() => {
    const total = weekTasks.value.length + weekDoneCount.value;
    if (total === 0) return 0;
    return Math.round((weekDoneCount.value / total) * 100);
  });

  async function load() {
    loading.value = true;
    try {
      const today = getDateRange("today");
      const week = getDateRange("week");
      const [configRes, todayRes, weekRes, doneRes] = await Promise.all([
        fetchGetModelConfig(),
        fetchGetTaskList({
          startTime: today.startTime,
          endTime: today.endTime,
        }),
        fetchGetTaskList({
          startTime: week.startTime,
          endTime: week.endTime,
        }),
        fetchGetTaskList({
          startTime: week.startTime,
          endTime: week.endTime,
          status: "done",
        }),
      ]);
      configInfo.value = configRes.data.data;
      todayTasks.value = todayRes.data.data?.tasks ?? [];
      weekTasks.value = weekRes.data.data?.tasks ?? [];
      weekDoneCount.value = doneRes.data.data?.tasks?.length ?? 0;
    } catch {
      // 单路失败不掀桌：保留已有数据，拦截器已弹过错误提示
    } finally {
      loading.value = false;
    }
  }

  function handleSseComplete() {
    load();
  }

  onMounted(() => {
    load();
    eventBus.on("sseComplete", handleSseComplete);
  });

  onUnmounted(() => {
    eventBus.off("sseComplete", handleSseComplete);
  });

  return {
    loading,
    configInfo,
    todayPending,
    weekTotal,
    weekDone,
    weekRate,
    load,
  };
}
