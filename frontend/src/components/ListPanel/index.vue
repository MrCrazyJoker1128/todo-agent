<template>
  <div class="list-panel">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><i class="i-ep-loading"></i></el-icon>
      <span>加载中...</span>
    </div>
    
    <template v-else>
      <!-- 今日任务 Section -->
      <div class="section today-section">
        <div class="section-header" @click="toggleSection('today')">
          <div class="section-title">
            <i class="i-ep-calendar section-icon"></i>
            <span>今日任务</span>
            <span class="section-count">{{ todayTasks.length }}</span>
          </div>
          <i :class="expandedSections.today ? 'i-ep-arrow-up' : 'i-ep-arrow-down'"></i>
        </div>
        <div v-show="expandedSections.today" class="section-content">
          <div v-if="todayTasks.length === 0" class="empty-state">
            暂无今日任务
          </div>
          <div v-else class="task-list">
            <TaskItem
              v-for="task in todayTasks"
              :key="task.id"
              :task="convertToTaskItem(task)"
            />
          </div>
        </div>
      </div>

      <!-- 本周任务 Section -->
      <div class="section week-section">
        <div class="section-header" @click="toggleSection('week')">
          <div class="section-title">
            <i class="i-ep-date section-icon"></i>
            <span>本周任务</span>
            <span class="section-count">{{ weekTasks.length }}</span>
          </div>
          <i :class="expandedSections.week ? 'i-ep-arrow-up' : 'i-ep-arrow-down'"></i>
        </div>
        <div v-show="expandedSections.week" class="section-content">
          <div v-if="weekTasks.length === 0" class="empty-state">
            暂无本周任务
          </div>
          <div v-else class="task-list">
            <TaskItem
              v-for="task in weekTasks"
              :key="task.id"
              :task="convertToTaskItem(task)"
            />
          </div>
        </div>
      </div>

      <!-- 已完成 Section -->
      <div class="section completed-section">
        <div class="section-header" @click="toggleSection('completed')">
          <div class="section-title">
            <i class="i-ep-circle-check section-icon"></i>
            <span>已完成</span>
            <span class="section-count">{{ completedTasks.length }}</span>
          </div>
          <i :class="expandedSections.completed ? 'i-ep-arrow-up' : 'i-ep-arrow-down'"></i>
        </div>
        <div v-show="expandedSections.completed" class="section-content">
          <div v-if="completedTasks.length === 0" class="empty-state">
            暂无已完成任务
          </div>
          <div v-else class="task-list">
            <TaskItem
              v-for="task in completedTasks"
              :key="task.id"
              :task="convertToTaskItem(task)"
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import TaskItem from "../TaskItem/index.vue";
import { fetchGetTaskList } from "@/api/task";
import type { TaskItem as ApiTaskItem } from "@/api/task/type";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  completed: boolean;
}

interface ExpandedSections {
  today: boolean;
  week: boolean;
  completed: boolean;
}

const tasks = ref<ApiTaskItem[]>([]);
const loading = ref(false);

const expandedSections = ref<ExpandedSections>({
  today: true,
  week: true,
  completed: false,
});

const convertToTaskItem = (task: ApiTaskItem): Task => {
  return {
    id: task.id,
    title: task.title,
    deadline: task.dueDate,
    priority: task.priority as "high" | "medium" | "low",
    completed: task.status === "completed",
  };
};

const isToday = (deadline: string): boolean => {
  if (!deadline) return false;
  const today = new Date();
  const taskDate = new Date(deadline);
  return (
    taskDate.getFullYear() === today.getFullYear() &&
    taskDate.getMonth() === today.getMonth() &&
    taskDate.getDate() === today.getDate()
  );
};

const isThisWeek = (deadline: string): boolean => {
  if (!deadline) return false;
  const today = new Date();
  const taskDate = new Date(deadline);
  const startOfWeek = new Date(today);
  startOfWeek.setDate(today.getDate() - today.getDay());
  startOfWeek.setHours(0, 0, 0, 0);
  const endOfWeek = new Date(startOfWeek);
  endOfWeek.setDate(startOfWeek.getDate() + 6);
  endOfWeek.setHours(23, 59, 59, 999);
  return taskDate >= startOfWeek && taskDate <= endOfWeek && !isToday(deadline);
};

const todayTasks = computed(() => {
  return tasks.value.filter((task) => task.status !== "completed" && isToday(task.dueDate));
});

const weekTasks = computed(() => {
  return tasks.value.filter((task) => task.status !== "completed" && isThisWeek(task.dueDate));
});

const completedTasks = computed(() => {
  return tasks.value.filter((task) => task.status === "completed");
});

const toggleSection = (section: keyof ExpandedSections) => {
  expandedSections.value[section] = !expandedSections.value[section];
};

const fetchTasks = async () => {
  loading.value = true;
  try {
    const res = await fetchGetTaskList();
    tasks.value = res.data.data?.tasks || [];
  } catch (error) {
    ElMessage.error("获取任务列表失败");
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchTasks();
});

defineExpose({
  refresh: fetchTasks,
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

.section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 10px;
  border: 1px solid rgba(255, 138, 128, 0.15);
  overflow: hidden;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 16px rgba(255, 138, 128, 0.12);
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
    background: rgba(255, 138, 128, 0.05);
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
  color: rgba(255, 138, 128, 0.9);
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
  background: rgba(255, 138, 128, 0.9);
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
  padding: 24px 16px;
  text-align: center;
  font-size: 13px;
  color: #909399;
  background: rgba(255, 138, 128, 0.05);
  border-radius: 8px;
  border: 1px dashed rgba(255, 138, 128, 0.2);
}

.today-section {
  .section-icon {
    color: rgba(255, 138, 128, 0.9);
  }
}

.week-section {
  .section-icon {
    color: #409eff;
  }
}

.completed-section {
  .section-icon {
    color: #67c23a;
  }
}
</style>
