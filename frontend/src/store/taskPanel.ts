import { defineStore } from "pinia";
import { ref, watch } from "vue";

export type TabName = "calendar" | "list" | "knowledge";

const STORAGE_KEY = "task-panel-active-tab";

function loadFromStorage(): TabName {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored !== null) {
    return stored as TabName;
  }
  return "calendar";
}

export const useTaskPanelStore = defineStore("taskPanel", () => {
  const activeTab = ref<TabName>(loadFromStorage());

  watch(
    activeTab,
    (newTab) => {
      localStorage.setItem(STORAGE_KEY, newTab);
    },
    { immediate: true },
  );

  return {
    activeTab,
  };
});
