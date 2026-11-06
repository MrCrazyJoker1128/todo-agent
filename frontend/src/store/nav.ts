import { defineStore } from "pinia";
import { ref } from "vue";
import { fetchGetTaskList } from "@/api/task";
import { useDateRange } from "@/components/ListPanel/hooks";

export const useNavStore = defineStore("nav", () => {
  /** 今日未完成任务数（菜单 badge 用） */
  const todayCount = ref(0);

  function setTodayCount(n: number) {
    todayCount.value = n;
  }

  /** 拉取今日未完成任务数：today 全量减掉 done/cancelled */
  async function refreshTodayCount() {
    try {
      const { getDateRange } = useDateRange();
      const range = getDateRange("today");
      const res = await fetchGetTaskList({
        startTime: range.startTime,
        endTime: range.endTime,
      });
      const tasks = res.data.data?.tasks ?? [];
      todayCount.value = tasks.filter(
        (t) => t.status !== "done" && t.status !== "cancelled",
      ).length;
    } catch {
      // 静默失败：badge 保持旧值，拦截器已弹错
    }
  }

  return {
    todayCount,
    setTodayCount,
    refreshTodayCount,
  };
});
