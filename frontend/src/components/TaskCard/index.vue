<template>
  <div class="task-card">
    <div class="task-header">
      <span class="task-title">{{ task.title }}</span>
      <el-tag :type="priorityType" size="small">{{ priorityLabel }}</el-tag>
    </div>
    <div class="task-info">
      <div class="info-item">
        <i class="i-ep-time"></i>
        <span>截止时间：{{ task.dueDate }}</span>
      </div>
      <div class="info-item">
        <i class="i-ep-top"></i>
        <span>优先级：</span>
        <span :class="priorityTextClass">{{ priorityLabel }}</span>
      </div>
      <div class="info-item">
        <i class="i-ep-collection-tag"></i>
        <span>分类：{{ taskTypeLabel }}</span>
      </div>
      <div class="info-item">
        <i class="i-ep-coin"></i>
        <span>状态：</span>
        <el-tag :type="statusType" size="small">{{ statusLabel }}</el-tag>
      </div>
    </div>
    <div class="task-actions">
      <!-- 状态变更按钮 -->
      <el-button 
        v-if="showStatusButton"
        size="small" 
        :type="statusButtonType"
        @click="handleStatusChange"
      >
        <i :class="statusButtonIcon"></i>
        {{ statusButtonText }}
      </el-button>
      
      <!-- 取消按钮 -->
      <el-button 
        v-if="showCancelButton"
        size="small" 
        type="danger"
        @click="handleCancel"
      >
        <i class="i-ep-close"></i>
        取消
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { ElMessage } from "element-plus";
import eventBus from "@/utils/eventbus";
import { fetchUpdateTaskStatus } from "@/api/task";

// 任务卡片数据接口（匹配后端实际值）
interface TaskCardItem {
  id: string;
  title: string;
  description?: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  dueDate: string;
  status: "todo" | "doing" | "done" | "cancelled";
}

const props = defineProps<{
  task: TaskCardItem;
}>();

const emit = defineEmits<{
  (e: "update", updatedTask: TaskCardItem): void;
}>();

type TagType = "primary" | "success" | "warning" | "info" | "danger";

// 优先级映射
const priorityMap: Record<string, { label: string; type: TagType; textClass: string }> = {
  high: { label: "高", type: "danger", textClass: "text-danger" },
  medium: { label: "中", type: "warning", textClass: "text-warning" },
  low: { label: "低", type: "success", textClass: "text-success" },
};

const priorityLabel = computed(() => {
  const key = props.task.priority.toLowerCase();
  return priorityMap[key]?.label || props.task.priority;
});

const priorityType = computed(() => {
  const key = props.task.priority.toLowerCase();
  return priorityMap[key]?.type || "info";
});

const priorityTextClass = computed(() => {
  const key = props.task.priority.toLowerCase();
  return priorityMap[key]?.textClass || "";
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

// 状态映射
const statusMap: Record<string, { label: string; type: TagType }> = {
  todo: { label: "待开始", type: "info" },
  doing: { label: "进行中", type: "warning" },
  done: { label: "已完成", type: "success" },
  cancelled: { label: "已取消", type: "info" },
};

const statusLabel = computed(() => {
  return statusMap[props.task.status]?.label || props.task.status;
});

const statusType = computed(() => {
  return statusMap[props.task.status]?.type || "info";
});

// 状态变更按钮逻辑
const showStatusButton = computed(() => {
  return props.task.status === "todo" || props.task.status === "doing";
});

const showCancelButton = computed(() => {
  return props.task.status === "todo" || props.task.status === "doing";
});

const statusButtonText = computed(() => {
  if (props.task.status === "todo") return "开始";
  if (props.task.status === "doing") return "完成";
  return "";
});

const statusButtonType = computed(() => {
  if (props.task.status === "todo") return "primary";
  if (props.task.status === "doing") return "success";
  return "";
});

const statusButtonIcon = computed(() => {
  if (props.task.status === "todo") return "i-ep-video-play";
  if (props.task.status === "doing") return "i-ep-check";
  return "";
});

const handleStatusChange = async () => {
  try {
    const newStatus = props.task.status === "todo" ? "doing" : "done";
    await fetchUpdateTaskStatus({
      id: props.task.id,
      status: newStatus,
    });
    
    emit("update", { ...props.task, status: newStatus });
    eventBus.emit("task-changed");
    ElMessage.success(newStatus === "doing" ? "任务已开始" : "任务已完成");
  } catch (error) {
    ElMessage.error("更新状态失败");
  }
};

const handleCancel = async () => {
  try {
    await fetchUpdateTaskStatus({
      id: props.task.id,
      status: "cancelled",
    });
    
    emit("update", { ...props.task, status: "cancelled" });
    eventBus.emit("task-changed");
    ElMessage.success("任务已取消");
  } catch (error) {
    ElMessage.error("取消任务失败");
  }
};
</script>

<style scoped lang="scss">
.task-card {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 138, 128, 0.2);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;

  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    .task-title {
      font-size: 14px;
      font-weight: 600;
      color: #303133;
    }
  }

  .task-info {
    margin-bottom: 12px;

    .info-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      color: #606266;
      margin-bottom: 6px;

      i {
        font-size: 12px;
        color: rgba(255, 138, 128, 0.8);
      }

      .text-danger {
        color: #f56c6c;
      }

      .text-warning {
        color: #e6a23c;
      }

      .text-success {
        color: #67c23a;
      }
    }
  }

  .task-actions {
    display: flex;
    gap: 8px;

    .el-button {
      flex: 1;
    }
  }
}
</style>
