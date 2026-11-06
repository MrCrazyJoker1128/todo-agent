<template>
  <div
    class="task-item"
    :class="{
      cancelled: task.status === 'cancelled',
      done: task.status === 'done',
    }"
  >
    <div class="task-main">
      <div class="task-left">
        <div class="task-info">
          <div class="task-title-row">
            <span
              class="task-title"
              :class="{ completed: task.status === 'done' }"
              >{{ task.title }}</span
            >
            <el-tag :type="priorityType" size="small" effect="plain">{{
              priorityLabel
            }}</el-tag>
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
          <el-button size="small" type="primary" @click="handleUpdate">
            <i class="i-ep-video-play"></i>
            编辑
          </el-button>
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
import eventBus from "@/utils/eventbus";
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
  (e: "updateTask"): void;
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
    eventBus.emit("task-changed");
    ElMessage.success("任务已完成");
  } catch {
    ElMessage.error("更新状态失败");
  }
};

const handleCancel = async () => {
  try {
    await fetchUpdateTaskStatus({ id: props.task.id, status: "cancelled" });
    emit("update", { ...props.task, status: "cancelled" });
    eventBus.emit("task-changed");
    ElMessage.success("任务已取消");
  } catch {
    ElMessage.error("取消任务失败");
  }
};
function handleUpdate() {
  emit("updateTask");
}
</script>

<style scoped lang="scss">
.task-item {
  position: relative;
  padding: 10px 12px 10px 34px;
  background: var(--wb-panel);
  border-radius: 8px;
  border: 1px solid var(--wb-line-soft);
  border-left: 2px solid var(--wb-blue);
  transition: border-color 0.2s, background 0.2s;

  &:hover {
    border-color: var(--wb-line);
    border-left-color: var(--wb-accent);
    background: var(--wb-panel-2);
  }

  &::before {
    content: "";
    position: absolute;
    left: 15px;
    top: 50%;
    width: 7px;
    height: 7px;
    margin: -3.5px 0 0 -3.5px;
    border-radius: 50%;
    background: var(--wb-blue);
  }

  &.done {
    opacity: 0.55;
    border-left-color: var(--wb-green);
    &::before {
      background: var(--wb-green);
    }
  }
  &.cancelled {
    opacity: 0.45;
    border-left-color: var(--wb-text-3);
    &::before {
      background: var(--wb-text-3);
    }
  }
}

.task-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.task-left {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
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
  font-size: 13.5px;
  font-weight: 500;
  color: var(--wb-text);
  &.completed {
    text-decoration: line-through;
    color: var(--wb-text-3);
  }
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

:deep(.el-tag) {
  font-family: var(--wb-font-mono);
  font-size: 10px;
  letter-spacing: 0.06em;
  border-radius: 4px;
  background: var(--wb-panel-3);
  border-color: var(--wb-line);
  color: var(--wb-text-2);
}
:deep(.el-tag--danger) {
  background: rgba(248, 113, 113, 0.1);
  border-color: rgba(248, 113, 113, 0.3);
  color: var(--wb-red);
}
:deep(.el-tag--warning) {
  background: var(--wb-amber-dim);
  border-color: rgba(251, 191, 36, 0.3);
  color: var(--wb-amber);
}
:deep(.el-tag--success) {
  background: rgba(52, 211, 153, 0.1);
  border-color: rgba(52, 211, 153, 0.3);
  color: var(--wb-green);
}
:deep(.el-tag--primary) {
  background: var(--wb-accent-dim);
  border-color: var(--wb-accent-line);
  color: var(--wb-accent);
}
.deadline {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-family: var(--wb-font-mono);
  color: var(--wb-text-3);
  i {
    font-size: 12px;
  }
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s;
}
.task-item:hover .task-actions {
  opacity: 1;
}

@media (max-width: 640px) {
  .task-actions {
    opacity: 1;
  }
  .task-item {
    padding: 10px 10px 10px 30px;
  }
}
</style>
