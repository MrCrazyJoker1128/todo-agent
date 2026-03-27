<template>
  <div class="task-item">
    <div class="task-left">
      <el-checkbox v-model="task.completed" :disabled="disabled" />
      <span class="task-title" :class="{ completed: task.completed }">{{ task.title }}</span>
    </div>
    <div class="task-right">
      <span class="deadline" :class="deadlineClass">
        <i :class="deadlineIcon"></i>
        截止时间 {{ task.deadline }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  completed: boolean;
}

defineProps<{
  task: Task;
  disabled?: boolean;
}>();

const deadlineClass = computed(() => {
  return {
    "deadline-high": true,
    "deadline-medium": true,
    "deadline-low": true,
  };
});

const deadlineIcon = computed(() => {
  return "i-ep-circle-dot";
});
</script>

<style scoped lang="scss">
.task-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  border: 1px solid rgba(255, 138, 128, 0.2);
  transition: all 0.3s;

  &:hover {
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  }

  .task-left {
    display: flex;
    align-items: center;
    gap: 12px;

    .task-title {
      font-size: 14px;
      color: #303133;

      &.completed {
        text-decoration: line-through;
        color: #909399;
      }
    }
  }

  .task-right {
    .deadline {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;

      &.deadline-high {
        color: #f56c6c;

        i {
          color: #f56c6c;
        }
      }

      &.deadline-medium {
        color: #e6a23c;

        i {
          color: #e6a23c;
        }
      }

      &.deadline-low {
        color: #67c23a;

        i {
          color: #67c23a;
        }
      }
    }
  }
}
</style>