<template>
  <main class="task-panel">
    <div class="panel-header flex-col">
      <div class="content-tabs w-full">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="日历" name="calendar" />
          <el-tab-pane label="列表" name="list" />
        </el-tabs>
      </div>
      <div class="flex-1 w-full flex-col">
        <CalendarPanel 
          v-if="activeTab === 'calendar'" 
          @task-click="handleTaskClick"
        />
        <ListPanel 
          v-else 
          ref="listPanelRef"
        />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import CalendarPanel from "../CalendarPanel/index.vue";
import ListPanel from "../ListPanel/index.vue";

const activeTab = ref("list");
const listPanelRef = ref<InstanceType<typeof ListPanel> | null>(null);

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  completed: boolean;
}

const handleTaskClick = (task: Task) => {
  console.log("Task clicked:", task);
};

watch(activeTab, (newTab) => {
  if (newTab === "list" && listPanelRef.value) {
    listPanelRef.value.refresh();
  }
});
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
