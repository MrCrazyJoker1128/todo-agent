<template>
  <main class="wb-taskzone">
    <div class="wb-zone-head">
      <div class="wb-zone-title">
        <span class="wb-zone-eyebrow wb-mono">MISSION BOARD</span>
        <div class="wb-zone-title-row">
          <h2 class="wb-font-display">{{ zoneTitle }}</h2>
          <div class="wb-view-switch">
            <button
              v-for="tab in tabs"
              :key="tab.name"
              :class="['wb-view-btn', { active: activeTab === tab.name }]"
              @click="activeTab = tab.name"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>
      </div>
      <div class="wb-zone-meta">
        <span class="wb-count-chip wb-mono">
          TODAY <b>{{ todayCount }}</b>
        </span>
        <el-button
          v-if="activeTab === 'list'"
          type="primary"
          size="small"
          @click="openCreate"
        >
          <i class="i-ep-plus"></i> 新建任务
        </el-button>
      </div>
    </div>

    <div class="wb-zone-body">
      <component
        :is="tabComponent"
        @task-click="handleTaskClick"
        @today-count="forwardTodayCount"
      />
    </div>
  </main>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import CalendarPanel from "../CalendarPanel/index.vue";
import ListPanel from "../ListPanel/index.vue";
import { useTaskPanelStore, type TabName } from "@/store/taskPanel";
import eventBus from "@/utils/eventbus";
defineOptions({
  name: "TaskPanel",
  components: {
    CalendarPanel,
    ListPanel,
  },
});
const taskPanelStore = useTaskPanelStore();
const { activeTab } = storeToRefs(taskPanelStore);

const emit = defineEmits<{
  (e: "today-count", n: number): void;
}>();

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
}
interface TabItem {
  label: string;
  name: TabName;
  component: string;
}
const tabs: TabItem[] = [
  { label: "日历", name: "calendar", component: "CalendarPanel" },
  { label: "列表", name: "list", component: "ListPanel" },
];
const tabComponent = computed(() => {
  const tab: TabItem | undefined = tabs.find(
    (item) => item.name === activeTab.value,
  );
  return tab?.component;
});
const zoneTitle = computed(() => {
  const tab = tabs.find((item) => item.name === activeTab.value);
  return tab?.label ?? "";
});
const todayCount = ref(0);
function forwardTodayCount(n: number) {
  todayCount.value = n;
  emit("today-count", n);
}
const handleTaskClick = (task: Task) => {
  console.log("Task clicked:", task);
};
function openCreate() {
  eventBus.emit("open-create-task");
}
</script>

<style scoped lang="scss">
.wb-taskzone {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 16px 20px;
  overflow: hidden;
}

.wb-zone-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--wb-line-soft);
  margin-bottom: 12px;
  flex-shrink: 0;
}

.wb-zone-eyebrow {
  font-size: 10px;
  letter-spacing: 0.22em;
  color: var(--wb-text-3);
}

.wb-zone-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 2px;
}

.wb-zone-title-row h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 0.04em;
  color: var(--wb-text);
}

.wb-view-switch {
  display: inline-flex;
  border: 1px solid var(--wb-line);
  border-radius: 8px;
  overflow: hidden;
  background: var(--wb-panel);
}

.wb-view-btn {
  border: 0;
  background: transparent;
  color: var(--wb-text-2);
  font-size: 12px;
  padding: 4px 12px;
  cursor: pointer;

  &.active {
    background: var(--wb-accent-dim);
    color: var(--wb-text);
  }
}

.wb-zone-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.wb-count-chip {
  font-size: 11px;
  letter-spacing: 0.1em;
  color: var(--wb-text-2);
  border: 1px solid var(--wb-line);
  background: var(--wb-panel);
  border-radius: 6px;
  padding: 5px 10px;

  b {
    color: var(--wb-amber);
    font-size: 13px;
    margin-left: 4px;
  }
}

.wb-zone-body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding-right: 4px;

  &::-webkit-scrollbar {
    width: 8px;
  }
  &::-webkit-scrollbar-thumb {
    background: var(--wb-line);
    border-radius: 4px;
  }
  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

@media (max-width: 720px) {
  .wb-taskzone {
    padding: 12px;
  }
  .wb-zone-title h2 {
    font-size: 18px;
  }
}
</style>
