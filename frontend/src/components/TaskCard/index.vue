<template>
  <div class="task-card">
    <div class="task-header">
      <span class="task-title">{{ task.title }}</span>
      <el-tag :type="priorityType" size="small">{{ task.priority }}</el-tag>
    </div>
    <div class="task-info">
      <div class="info-item">
        <i class="i-ep-time"></i>
        <span>时间：{{ task.time }}</span>
      </div>
      <div class="info-item">
        <i class="i-ep-rank"></i>
        <span>优先级：</span>
        <i class="i-ep-top" :class="priorityIconClass"></i>
        <span>{{ task.priority }}</span>
      </div>
      <div class="info-item">
        <i class="i-ep-collection-tag"></i>
        <span>标签：{{ task.tags.join(", ") }}</span>
      </div>
    </div>
    <div class="task-actions">
      <el-button size="small" type="success">
        <i class="i-ep-check"></i>
        完成
      </el-button>
      <el-button size="small" type="warning">
        <i class="i-ep-edit"></i>
        编辑
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

interface Task {
  title: string;
  time: string;
  priority: "高" | "中" | "低";
  tags: string[];
}

const props = defineProps<{
  task: Task;
}>();

const priorityType = computed(() => {
  const map: Record<string, any> = {
    高: "danger",
    中: "warning",
    低: "success",
  };
  return map[props.task.priority] || "";
});

const priorityIconClass = computed(() => {
  const map: Record<string, string> = {
    高: "text-danger",
    中: "text-warning",
    低: "text-success",
  };
  return map[props.task.priority] || "";
});
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