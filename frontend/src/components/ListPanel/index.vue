<template>
  <div class="list-panel h-full">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><i class="i-ep-loading"></i></el-icon>
      <span>加载中...</span>
    </div>

    <template v-else>
      <!-- 模式切换 Tabs -->
      <div class="mode-tabs">
        <button
          :class="['tab-btn', { active: mode === 'today' }]"
          @click="onModeChange('today')"
        >
          <i class="i-ep-calendar"></i>
          <span>今日任务</span>
          <span class="tab-count">{{ todayCount }}</span>
        </button>
        <button
          :class="['tab-btn', { active: mode === 'week' }]"
          @click="onModeChange('week')"
        >
          <i class="i-ep-date"></i>
          <span>本周任务</span>
          <span class="tab-count">{{ weekCount }}</span>
        </button>
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <!-- 优先级筛选 -->
        <div class="priority-filters">
          <span class="filter-label">优先级:</span>
          <button
            v-for="p in priorities"
            :key="p.value"
            :class="[
              'priority-chip',
              p.value,
              { selected: selectedPriorities.includes(p.value) },
            ]"
            @click="togglePriority(p.value)"
          >
            {{ p.label }}
          </button>
        </div>

        <!-- 分类筛选 -->
        <el-select
          v-model="selectedType"
          placeholder="选择分类"
          class="type-select"
          clearable
        >
          <el-option
            v-for="type in taskTypeOptions"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          />
        </el-select>
      </div>

      <!-- 筛选后的任务列表 -->
      <div class="filtered-tasks">
        <div v-if="filteredTasksResult.length === 0" class="empty-state">
          <i class="i-ep-document"></i>
          <span>暂无匹配任务</span>
        </div>
        <div v-else class="task-list">
          <TaskItem
            v-for="task in filteredTasksResult"
            :key="task.id"
            :task="convertToTaskItem(task)"
            @update="handleTaskUpdate"
            @updateTask="() => showUpdateTaskDialog(task)"
          />
        </div>
      </div>

      <!-- 已完成任务（独立区块，不受筛选影响） -->
      <div class="completed-section">
        <div class="section-header" @click="showCompleted = !showCompleted">
          <div class="section-title">
            <i class="i-ep-circle-check section-icon"></i>
            <span>已完成</span>
            <span class="section-count">{{ completedTasks.length }}</span>
          </div>
          <i :class="showCompleted ? 'i-ep-arrow-up' : 'i-ep-arrow-down'"></i>
        </div>
        <div v-show="showCompleted" class="section-content">
          <div v-if="completedTasks.length === 0" class="empty-state small">
            暂无已完成任务
          </div>
          <div v-else class="task-list">
            <TaskItem
              v-for="task in completedTasks"
              :key="task.id"
              :task="convertToTaskItem(task)"
              @update="handleTaskUpdate"
              @updateTask="() => showUpdateTaskDialog(task)"
            />
          </div>
        </div>
      </div>
    </template>
    <UpdateTaskDialog
      v-model:dialogVisible="dialogVisible"
      ref="updateTaskDialogRef"
      @success="() => fetchTasks(mode)"
    />
  </div>
</template>

<script setup lang="ts">
import type { TaskItem } from "@/api/task/type";
import { useTaskFilter, useTaskMode, useTaskList } from "./hooks";
import eventBus from "@/utils/eventbus";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
}

const emit = defineEmits<{
  (e: "today-count", n: number): void;
}>();
const { mode, handleModeChange } = useTaskMode();
const {
  selectedPriorities,
  selectedType,
  filteredTasks,
  togglePriority,
  priorities,
  taskTypeOptions,
} = useTaskFilter();
const {
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
} = useTaskList();

const currentTasks = computed(() => {
  return mode.value === "today" ? todayTasks.value : weekTasks.value;
});

const filteredTasksResult = computed(() => filteredTasks(currentTasks).value);

const onModeChange = (newMode: "today" | "week") => {
  handleModeChange(newMode, () => fetchTasks(newMode));
};
eventBus.on("sseComplete", () => {
  if (mode.value === "today") {
    fetchTasks(mode.value);
  } else if (mode.value === "week") {
    fetchTasks(mode.value);
  }
});
watch(todayCount, (n) => emit("today-count", n), { immediate: true });
onMounted(() => {
  fetchTasks(mode.value);
});
const dialogVisible = ref<boolean>(false);
const updateTaskDialogRef = useTemplateRef("updateTaskDialogRef");
function showUpdateTaskDialog(task: TaskItem) {
  dialogVisible.value = true;
  updateTaskDialogRef.value?.setTaskForm(task);
}
defineExpose({
  refresh: () => fetchTasks(mode.value),
});
</script>

<style scoped lang="scss">
.list-panel {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  gap: 12px;
  color: var(--wb-text-3);

  .el-icon {
    font-size: 24px;
    color: var(--wb-accent);
  }
}

/* Mode Tabs */
.mode-tabs {
  display: flex;
  gap: 6px;
  padding: 4px;
  background: var(--wb-panel);
  border-radius: 8px;
  border: 1px solid var(--wb-line-soft);
}

.tab-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 14px;
  font-size: 13px;
  font-weight: 500;
  color: var(--wb-text-3);
  background: transparent;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;

  i {
    font-size: 15px;
  }

  &:hover {
    color: var(--wb-text-2);
    background: var(--wb-panel-2);
  }

  &.active {
    color: var(--wb-accent);
    background: var(--wb-accent-dim);
    border-color: var(--wb-accent-line);

    .tab-count {
      background: rgba(94, 234, 212, 0.2);
      color: var(--wb-accent);
    }
  }
}

.tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 18px;
  padding: 0 6px;
  font-size: 11px;
  font-weight: 600;
  font-family: var(--wb-font-mono);
  color: var(--wb-text-3);
  background: var(--wb-panel-3);
  border-radius: 9px;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 12px;
  background: var(--wb-panel);
  border-radius: 8px;
  border: 1px solid var(--wb-line-soft);
}

.priority-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 11px;
  letter-spacing: 0.12em;
  color: var(--wb-text-3);
  margin-right: 2px;
  font-family: var(--wb-font-mono);
}

.priority-chip {
  padding: 5px 12px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.04em;
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: var(--wb-font-mono);

  &.high {
    color: var(--wb-red);
    background: rgba(248, 113, 113, 0.1);
    border-color: rgba(248, 113, 113, 0.3);
    &.selected {
      background: var(--wb-red);
      color: #0b0e14;
      border-color: var(--wb-red);
    }
  }
  &.medium {
    color: var(--wb-amber);
    background: var(--wb-amber-dim);
    border-color: rgba(251, 191, 36, 0.3);
    &.selected {
      background: var(--wb-amber);
      color: #0b0e14;
      border-color: var(--wb-amber);
    }
  }
  &.low {
    color: var(--wb-green);
    background: rgba(52, 211, 153, 0.1);
    border-color: rgba(52, 211, 153, 0.3);
    &.selected {
      background: var(--wb-green);
      color: #0b0e14;
      border-color: var(--wb-green);
    }
  }
  &:hover:not(.selected) {
    filter: brightness(1.2);
  }
}

.type-select {
  width: 160px;
}

/* Filtered Tasks */
.filtered-tasks {
  min-height: 80px;
}

/* Completed Section */
.completed-section {
  background: var(--wb-panel);
  border-radius: var(--wb-radius);
  border: 1px solid var(--wb-line-soft);
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.2s;
  &:hover {
    background: var(--wb-panel-2);
  }
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.06em;
  color: var(--wb-text-2);
}

.section-icon {
  font-size: 16px;
  color: var(--wb-green);
}

.section-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 18px;
  padding: 0 7px;
  font-size: 11px;
  font-weight: 600;
  font-family: var(--wb-font-mono);
  color: var(--wb-green);
  background: rgba(52, 211, 153, 0.12);
  border: 1px solid rgba(52, 211, 153, 0.3);
  border-radius: 9px;
}

.section-content {
  padding: 0 12px 12px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;

  &::before {
    content: "";
    position: absolute;
    left: 15px;
    top: 12px;
    bottom: 12px;
    width: 1px;
    background: var(--wb-line-soft);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 28px 16px;
  text-align: center;
  font-size: 12px;
  color: var(--wb-text-3);
  background: var(--wb-panel);
  border-radius: 8px;
  border: 1px dashed var(--wb-line);
  font-family: var(--wb-font-mono);
  letter-spacing: 0.06em;

  i {
    font-size: 24px;
    color: var(--wb-line);
  }
  &.small {
    padding: 18px 14px;
  }
}

/* Responsive */
@media (max-width: 640px) {
  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
  }
  .type-select {
    width: 100%;
  }
  .mode-tabs {
    flex-direction: column;
  }
  .tab-btn {
    justify-content: flex-start;
  }
}
</style>
