<template>
  <div class="list-panel">
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
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import TaskItem from "../TaskItem/index.vue";
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
onMounted(() => {
  fetchTasks(mode.value);
});

defineExpose({
  refresh: () => fetchTasks(mode.value),
});
</script>

<style scoped lang="scss">
.list-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(255, 138, 128, 0.08);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  gap: 12px;
  color: #909399;

  .el-icon {
    font-size: 24px;
  }
}

/* Mode Tabs */
.mode-tabs {
  display: flex;
  gap: 8px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 10px;
  border: 1px solid rgba(255, 138, 128, 0.12);
}

.tab-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.25s ease;

  i {
    font-size: 16px;
  }

  &:hover {
    color: rgba(255, 138, 128, 0.9);
    background: rgba(255, 138, 128, 0.06);
  }

  &.active {
    color: #fff;
    background: rgba(255, 138, 128, 0.9);
    box-shadow: 0 2px 8px rgba(255, 138, 128, 0.35);

    .tab-count {
      background: rgba(255, 255, 255, 0.25);
      color: #fff;
    }
  }
}

.tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  font-size: 11px;
  font-weight: 600;
  color: rgba(255, 138, 128, 0.9);
  background: rgba(255, 138, 128, 0.12);
  border-radius: 10px;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 10px;
  border: 1px solid rgba(255, 138, 128, 0.1);
}

.priority-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 13px;
  color: #909399;
  margin-right: 4px;
}

.priority-chip {
  padding: 5px 12px;
  font-size: 12px;
  font-weight: 500;
  border: 1px solid transparent;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s ease;

  &.high {
    color: #f56c6c;
    background: rgba(245, 108, 108, 0.1);
    border-color: rgba(245, 108, 108, 0.2);

    &.selected {
      background: #f56c6c;
      color: #fff;
      border-color: #f56c6c;
      box-shadow: 0 2px 6px rgba(245, 108, 108, 0.4);
    }
  }

  &.medium {
    color: #e6a23c;
    background: rgba(230, 162, 60, 0.1);
    border-color: rgba(230, 162, 60, 0.2);

    &.selected {
      background: #e6a23c;
      color: #fff;
      border-color: #e6a23c;
      box-shadow: 0 2px 6px rgba(230, 162, 60, 0.4);
    }
  }

  &.low {
    color: #67c23a;
    background: rgba(103, 194, 58, 0.1);
    border-color: rgba(103, 194, 58, 0.2);

    &.selected {
      background: #67c23a;
      color: #fff;
      border-color: #67c23a;
      box-shadow: 0 2px 6px rgba(103, 194, 58, 0.4);
    }
  }

  &:hover:not(.selected) {
    opacity: 0.8;
    transform: translateY(-1px);
  }
}

.type-select {
  width: 160px;

  :deep(.el-select__wrapper) {
    background: rgba(255, 255, 255, 0.8);
    border-color: rgba(255, 138, 128, 0.15);
    box-shadow: none;

    &:hover {
      border-color: rgba(255, 138, 128, 0.3);
    }

    &.is-focused {
      border-color: rgba(255, 138, 128, 0.5);
      box-shadow: 0 0 0 2px rgba(255, 138, 128, 0.1);
    }
  }
}

/* Filtered Tasks */
.filtered-tasks {
  min-height: 80px;
}

/* Completed Section */
.completed-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 10px;
  border: 1px solid rgba(103, 194, 58, 0.2);
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 16px rgba(103, 194, 58, 0.1);
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.2s ease;

  &:hover {
    background: rgba(103, 194, 58, 0.05);
  }
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.section-icon {
  font-size: 18px;
  color: #67c23a;
}

.section-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 8px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  background: #67c23a;
  border-radius: 11px;
}

.section-content {
  padding: 0 16px 16px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 32px 16px;
  text-align: center;
  font-size: 13px;
  color: #909399;
  background: rgba(255, 138, 128, 0.05);
  border-radius: 8px;
  border: 1px dashed rgba(255, 138, 128, 0.2);

  i {
    font-size: 28px;
    color: rgba(255, 138, 128, 0.3);
  }

  &.small {
    padding: 20px 16px;

    i {
      font-size: 22px;
    }
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
