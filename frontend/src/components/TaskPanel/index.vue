<template>
  <main class="task-panel">
    <div class="panel-header flex-col">
      <div class="content-tabs w-full">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="日历" name="calendar" />
          <el-tab-pane label="列表" name="list" />
          <el-tab-pane label="知识库" name="knowledge" />
        </el-tabs>
      </div>
      <div class="flex-1 w-full flex-col">
        <component
          ref="listPanelRef"
          :is="tabComponent"
          @task-click="handleTaskClick"
        />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { storeToRefs } from "pinia";
import CalendarPanel from "../CalendarPanel/index.vue";
import ListPanel from "../ListPanel/index.vue";
import KnowledgePanel from "../KnowledgePanel/index.vue";
import { useTaskPanelStore } from "@/store/taskPanel";
defineOptions({
  name: "TaskPanel",
  components: {
    CalendarPanel,
    ListPanel,
    KnowledgePanel,
  },
});
const taskPanelStore = useTaskPanelStore();
const { activeTab } = storeToRefs(taskPanelStore);
const listPanelRef = ref<InstanceType<typeof ListPanel> | null>(null);

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
  name: string;
  component: string;
}
const tabs: TabItem[] = [
  { label: "日历", name: "calendar", component: "CalendarPanel" },
  { label: "列表", name: "list", component: "ListPanel" },
  { label: "知识库", name: "knowledge", component: "KnowledgePanel" },
];
const tabComponent = computed(() => {
  const tab = tabs.find((item) => item.name === activeTab.value);
  return tab.component;
});
const handleTaskClick = (task: Task) => {
  console.log("Task clicked:", task);
};
</script>

<style scoped lang="scss">
.task-panel {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  background: transparent;
  margin-right: 4px;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
}

.panel-header {
  flex: 1;
  background: transparent;
  border-radius: 12px;
  padding: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.content-tabs {
  margin-bottom: 16px;

  :deep(.el-tabs__header) {
    margin: 0;
  }

  :deep(.el-tabs__nav-wrap::after) {
    height: 1px;
    background-color: rgba(255, 138, 128, 0.2);
  }
}
</style>
