<template>
  <div class="task-item" :class="{ cancelled: task.status === 'cancelled', done: task.status === 'done' }">
    <div class="task-main">
      <div class="task-left">
        <div class="task-info">
          <div class="task-title-row">
            <span class="task-title" :class="{ completed: task.status === 'done' }">{{ task.title }}</span>
            <el-tag :type="priorityType" size="small" effect="plain">{{ priorityLabel }}</el-tag>
            <el-tag size="small" effect="plain">{{ taskTypeLabel }}</el-tag>
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
          <el-button size="small" type="primary" @click="handleStart">
            <i class="i-ep-video-play"></i>
            开始
          </el-button>
          <el-button size="small" type="info" plain @click="handleCancel">
            <i class="i-ep-close"></i>
            取消
          </el-button>
        </template>
        <template v-else-if="task.status === 'doing'">
          <el-button size="small" type="success" @click="handleDone">
            <i class="i-ep-check"></i>
            完成
          </el-button>
          <el-button size="small" type="info" plain @click="handleCancel">
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
</template>

<script setup lang="ts">
import { computed } from "vue";
import { ElMessage } from "element-plus";
import { fetchUpdateTaskStatus } from "@/api/task";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
}

const props = defineProps<{
  task: Task;
}>();

const emit = defineEmits<{
  (e: "update", updatedTask: Task): void;
}>();

type TagType = "primary" | "success" | "warning" | "info" | "danger";

const priorityMap: Record<string, { label: string; type: TagType }> = {
  high: { label: "高", type: "danger" },
  medium: { label: "中", type: "warning" },
  low: { label: "低", type: "success" },
};

const priorityLabel = computed(() => {
  const key = props.task.priority.toLowerCase();
  return priorityMap[key]?.label || props.task.priority;
});
const priorityType = computed(() => {
  const key = props.task.priority.toLowerCase();
  return priorityMap[key]?.type || "info";
});

const taskTypeMap: Record<string, string> = {
  work: "工作",
  personal: "个人",
  project: "项目",
  adhoc: "临时",
};

const taskTypeLabel = computed(() => {
  const key = props.task.taskType?.toLowerCase();
  return taskTypeMap[key] || props.task.taskType || "未分类";
});

const handleStart = async () => {
  try {
    await fetchUpdateTaskStatus({ id: props.task.id, status: "doing" });
    emit("update", { ...props.task, status: "doing" });
    ElMessage.success("任务已开始");
  } catch {
    ElMessage.error("更新状态失败");
  }
};

const handleDone = async () => {
  try {
    await fetchUpdateTaskStatus({ id: props.task.id, status: "done" });
    emit("update", { ...props.task, status: "done" });
    ElMessage.success("任务已完成");
  } catch {
    ElMessage.error("更新状态失败");
  }
};

const handleCancel = async () => {
  try {
    await fetchUpdateTaskStatus({ id: props.task.id, status: "cancelled" });
    emit("update", { ...props.task, status: "cancelled" });
    ElMessage.success("任务已取消");
  } catch {
    ElMessage.error("取消任务失败");
  }
};
</script>

<style scoped lang="scss">
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

.task-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;

  &.completed {
    text-decoration: line-through;
    color: #909399;
  }
}

.priority-tag,
.type-tag,
.status-tag {
  flex-shrink: 0;
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

  .el-button {
    padding: 6px 8px;
    min-width: 28px;
  }
}

.done-icon {
  font-size: 20px;
  color: #67c23a;
}

.cancelled-icon {
  font-size: 20px;
  color: #909399;
}

@media (max-width: 640px) {
  .task-item {
    padding: 12px;
  }

  .task-title-row {
    gap: 6px;
  }

  .priority-tag,
  .type-tag,
  .status-tag {
    font-size: 10px;
    padding: 0 4px;
  }
}
</style>
