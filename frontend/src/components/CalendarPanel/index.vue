<template>
  <div class="calendar-panel">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><i class="i-ep-loading"></i></el-icon>
      <span>加载中...</span>
    </div>
    <el-calendar v-else v-model="currentMonth">
      <template #date-cell="{ data }">
        <div
          class="calendar-cell"
          :class="{
            'is-selected': data.day === selectedDate,
            'is-today': data.day === todayStr,
            'other-month':
              data.type === 'prev-month' || data.type === 'next-month',
          }"
          @click="handleDateClick(data.day)"
        >
          <span class="day-number">{{
            data.day.split("-").slice(2).join("")
          }}</span>
          <div class="task-list">
            <template v-if="getFirstTaskForDate(data.day)">
              <div
                class="task-chip"
                @click.stop="handleTaskClick(getFirstTaskForDate(data.day)!)"
              >
                <span class="task-name">{{
                  getFirstTaskForDate(data.day)!.title
                }}</span>
              </div>
            </template>
          </div>
        </div>
      </template>
    </el-calendar>

    <div class="tasks-section" ref="tasksSectionRef">
      <div class="tasks-header" ref="tasksHeaderRef">
        <span class="tasks-title">{{ formatSelectedDate(selectedDate) }}</span>
        <span class="tasks-count">{{ selectedDateTasks.length }} 个任务</span>
      </div>
      <div v-if="selectedDateTasks.length === 0" class="empty-state">
        暂无任务
      </div>
      <el-scrollbar
        v-else
        class="tasks-scrollbar"
        :height="scrollbarHeight + 'px'"
        max-height="200px"
      >
        <div
          v-for="task in selectedDateTasks"
          :key="task.id"
          class="task-item"
          :class="{
            done: task.status === 'done',
            cancelled: task.status === 'cancelled',
          }"
        >
          <div class="task-main">
            <div class="task-left">
              <div class="task-info">
                <div class="task-title-row">
                  <span
                    class="task-item-title"
                    :class="{ completed: task.status === 'done' }"
                    >{{ task.title }}</span
                  >
                  <el-tag
                    :type="getPriorityTagType(task.priority)"
                    size="small"
                    effect="plain"
                    >{{ getPriorityLabel(task.priority) }}</el-tag
                  >
                  <el-tag size="small" effect="plain">{{
                    getTaskTypeLabel(task.taskType)
                  }}</el-tag>
                </div>
                <div class="task-meta">
                  <span class="deadline">
                    <i class="i-ep-time"></i>
                    {{ task.deadline }}
                  </span>
                </div>
              </div>
            </div>
            <div class="task-actions">
              <template v-if="task.status === 'todo'">
                <el-button
                  size="small"
                  type="primary"
                  @click="handleStart(task)"
                >
                  <i class="i-ep-video-play"></i>
                  开始
                </el-button>
                <el-button
                  size="small"
                  type="info"
                  plain
                  @click="handleCancel(task)"
                >
                  <i class="i-ep-close"></i>
                  取消
                </el-button>
              </template>
              <template v-else-if="task.status === 'doing'">
                <el-button
                  size="small"
                  type="success"
                  @click="handleDone(task)"
                >
                  <i class="i-ep-check"></i>
                  完成
                </el-button>
                <el-button
                  size="small"
                  type="info"
                  plain
                  @click="handleCancel(task)"
                >
                  <i class="i-ep-close"></i>
                  取消
                </el-button>
              </template>
              <template v-else-if="task.status === 'done'">
                <el-tag type="success" size="small" effect="plain">
                  <i class="i-ep-circle-check"></i>
                  已完成
                </el-tag>
              </template>
              <template v-else-if="task.status === 'cancelled'">
                <el-tag type="info" size="small" effect="plain">
                  <i class="i-ep-remove"></i>
                  已取消
                </el-tag>
              </template>
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch, onMounted } from "vue";
import {
  useTask,
  useTaskDisplay,
  useScrollbarHeight,
  useCalendar,
  type Task,
} from "./hooks";
import eventBus from "@/utils/eventbus";

const emit = defineEmits<{
  (e: "task-click", task: Task): void;
}>();

const { currentMonth, selectedDate, todayStr } = useCalendar();

const {
  selectedDateTasks,
  loading,
  getFirstTaskForDate,
  fetchTasks,
  fetchTasksForDate,
  handleStart,
  handleDone,
  handleCancel,
} = useTask();

const {
  getTaskTypeLabel,
  getPriorityLabel,
  getPriorityTagType,
  formatSelectedDate,
} = useTaskDisplay();

const { scrollbarHeight, calculateScrollbarHeight } = useScrollbarHeight();

const handleDateClick = async (date: string) => {
  selectedDate.value = date;
  await fetchTasksForDate(date);
  calculateScrollbarHeight();
};

const handleTaskClick = (task: Task) => {
  emit("task-click", task);
};

watch(
  selectedDateTasks,
  () => {
    calculateScrollbarHeight();
  },
  { flush: "post" },
);

onMounted(() => {
  fetchTasks(currentMonth.value);
  selectedDate.value = todayStr;
  handleDateClick(todayStr);
  calculateScrollbarHeight();
});
eventBus.on("sseComplete", () => {
  refresh();
});
const refresh = () => {
  fetchTasks(currentMonth.value);
};

defineExpose({
  refresh,
});
</script>

<style scoped lang="scss">
.calendar-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(255, 138, 128, 0.08);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120px;
  gap: 12px;
  color: #909399;

  .el-icon {
    font-size: 24px;
  }
}

:deep(.el-calendar) {
  border: none;

  .el-calendar__header {
    padding: 12px 16px;
    border-bottom: none;
  }

  .el-calendar__body {
    padding: 8px;
  }

  .el-calendar-table__header {
    th {
      padding: 8px 0;
      color: #909399;
      font-weight: 500;
      border: none;
    }
  }

  .el-calendar-table__body {
    td {
      border: 1px solid rgba(255, 138, 128, 0.08);
      transition: all 0.2s ease;
      cursor: pointer;
      vertical-align: top;

      &.prev-month,
      &.next-month {
        .day-number {
          color: #c0c4cc;
        }
      }

      &.today {
        .day-number {
          color: rgba(255, 138, 128, 1);
          font-weight: 700;
        }
      }

      &.is-selected {
        background: rgba(255, 138, 128, 0.15);
      }

      &:hover {
        background: rgba(255, 255, 255, 0.8);
      }
    }
  }
}

.calendar-cell {
  min-height: 70px;
  padding: 6px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.5);
  transition: all 0.2s ease;

  &.other-month {
    background: rgba(255, 255, 255, 0.25);
  }

  &.is-today {
    background: rgba(255, 138, 128, 0.08);
  }

  .day-number {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #303133;
    margin-bottom: 4px;
  }
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.task-chip {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.8);

  .task-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
    font-weight: 500;
  }

  &:hover {
    background: rgba(255, 255, 255, 1);
  }
}

.more-tasks {
  font-size: 9px;
  color: #909399;
  padding: 2px 6px;
  text-align: center;
}

.tasks-section {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tasks-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-shrink: 0;

  .tasks-title {
    font-size: 15px;
    font-weight: 600;
    color: #303133;
  }

  .tasks-count {
    font-size: 12px;
    color: #909399;
    background: rgba(255, 138, 128, 0.1);
    padding: 4px 10px;
    border-radius: 12px;
  }
}

.empty-state {
  text-align: center;
  padding: 24px;
  color: #909399;
  font-size: 13px;
}

.tasks-scrollbar {
  flex: 1;
  min-height: 0;
}

.task-item {
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  border: 1px solid rgba(255, 138, 128, 0.2);
  transition: all 0.3s;

  &:hover {
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  }

  &.done {
    opacity: 0.7;
    background: rgba(103, 194, 58, 0.05);
    border-color: rgba(103, 194, 58, 0.2);
  }

  &.cancelled {
    opacity: 0.6;
    background: rgba(144, 147, 153, 0.05);
    border-color: rgba(144, 147, 153, 0.2);
  }
}

.task-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.task-left {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
  flex: 1;
}

.task-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.task-item-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;

  &.completed {
    text-decoration: line-through;
    color: #909399;
  }
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.deadline {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;

  i {
    font-size: 12px;
  }
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;

  :deep(.el-button) {
    padding: 6px 8px;
    min-width: 28px;
  }
}
</style>
