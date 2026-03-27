<template>
  <div class="calendar-panel">
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading"><i class="i-ep-loading"></i></el-icon>
      <span>加载中...</span>
    </div>
    
    <template v-else>
      <div class="calendar-header">
        <div class="calendar-nav">
          <el-button text @click="prevMonth">
            <i class="i-ep-arrow-left"></i>
          </el-button>
          <span class="current-date">{{ currentYear }}年 {{ currentMonth + 1 }}月</span>
          <el-button text @click="nextMonth">
            <i class="i-ep-arrow-right"></i>
          </el-button>
        </div>
        <el-button size="small" @click="goToToday">今天</el-button>
      </div>

      <div class="calendar-weekdays">
        <div v-for="day in weekdays" :key="day" class="weekday">{{ day }}</div>
      </div>

      <div class="calendar-grid">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          class="calendar-cell"
          :class="{
            'other-month': !day.isCurrentMonth,
            'today': day.isToday,
            'has-tasks': day.tasks.length > 0
          }"
        >
          <span class="day-number">{{ day.date }}</span>
          <div class="task-list">
            <div
              v-for="task in day.tasks.slice(0, 3)"
              :key="task.id"
              class="task-chip"
              :class="[`priority-${task.priority}`, { completed: task.completed }]"
              @click="handleTaskClick(task)"
            >
              <span class="task-chip-title">{{ task.title }}</span>
            </div>
            <div v-if="day.tasks.length > 3" class="more-tasks">
              +{{ day.tasks.length - 3 }} 更多
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { fetchGetTaskList } from "@/api/task";
import type { TaskItem as ApiTaskItem } from "@/api/task/type";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  completed: boolean;
}

interface CalendarDay {
  date: number;
  isCurrentMonth: boolean;
  isToday: boolean;
  tasks: Task[];
}

const emit = defineEmits<{
  (e: "task-click", task: Task): void;
}>();

const tasks = ref<Task[]>([]);
const loading = ref(false);

const weekdays = ["日", "一", "二", "三", "四", "五", "六"];

const today = new Date();
const currentYear = ref(today.getFullYear());
const currentMonth = ref(today.getMonth());

const prevMonth = () => {
  if (currentMonth.value === 0) {
    currentMonth.value = 11;
    currentYear.value--;
  } else {
    currentMonth.value--;
  }
};

const nextMonth = () => {
  if (currentMonth.value === 11) {
    currentMonth.value = 0;
    currentYear.value++;
  } else {
    currentMonth.value++;
  }
};

const goToToday = () => {
  currentYear.value = today.getFullYear();
  currentMonth.value = today.getMonth();
};

const getDaysInMonth = (year: number, month: number) => {
  return new Date(year, month + 1, 0).getDate();
};

const getFirstDayOfMonth = (year: number, month: number) => {
  return new Date(year, month, 1).getDay();
};

const isSameDate = (d1: Date, d2: Date) => {
  return (
    d1.getFullYear() === d2.getFullYear() &&
    d1.getMonth() === d2.getMonth() &&
    d1.getDate() === d2.getDate()
  );
};

const calendarDays = computed<CalendarDay[]>(() => {
  const daysInMonth = getDaysInMonth(currentYear.value, currentMonth.value);
  const firstDay = getFirstDayOfMonth(currentYear.value, currentMonth.value);
  const daysInPrevMonth = getDaysInMonth(
    currentMonth.value === 0 ? currentYear.value - 1 : currentYear.value,
    currentMonth.value === 0 ? 11 : currentMonth.value - 1
  );

  const days: CalendarDay[] = [];

  for (let i = firstDay - 1; i >= 0; i--) {
    const date = daysInPrevMonth - i;
    const prevMonthVal = currentMonth.value === 0 ? 11 : currentMonth.value - 1;
    const prevYear = currentMonth.value === 0 ? currentYear.value - 1 : currentYear.value;
    const fullDate = `${prevYear}-${String(prevMonthVal + 1).padStart(2, "0")}-${String(date).padStart(2, "0")}`;
    days.push({
      date,
      isCurrentMonth: false,
      isToday: false,
      tasks: tasks.value.filter((t) => t.deadline === fullDate),
    });
  }

  for (let date = 1; date <= daysInMonth; date++) {
    const fullDate = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, "0")}-${String(date).padStart(2, "0")}`;
    days.push({
      date,
      isCurrentMonth: true,
      isToday: isSameDate(
        new Date(currentYear.value, currentMonth.value, date),
        today
      ),
      tasks: tasks.value.filter((t) => t.deadline === fullDate),
    });
  }

  const remainingDays = 42 - days.length;
  for (let date = 1; date <= remainingDays; date++) {
    const nextMonthVal = currentMonth.value === 11 ? 0 : currentMonth.value + 1;
    const nextYear = currentMonth.value === 11 ? currentYear.value + 1 : currentYear.value;
    const fullDate = `${nextYear}-${String(nextMonthVal + 1).padStart(2, "0")}-${String(date).padStart(2, "0")}`;
    days.push({
      date,
      isCurrentMonth: false,
      isToday: false,
      tasks: tasks.value.filter((t) => t.deadline === fullDate),
    });
  }

  return days;
});

const handleTaskClick = (task: Task) => {
  emit("task-click", task);
};

const fetchTasks = async () => {
  loading.value = true;
  try {
    const res = await fetchGetTaskList();
    const apiTasks = res.data.data?.tasks || [];
    tasks.value = apiTasks.map((t: ApiTaskItem) => ({
      id: t.id,
      title: t.title,
      deadline: t.dueDate,
      priority: t.priority as "high" | "medium" | "low",
      completed: t.status === "completed",
    }));
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
.calendar-panel {
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

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .calendar-nav {
    display: flex;
    align-items: center;
    gap: 12px;

    .current-date {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      min-width: 140px;
      text-align: center;
    }
  }
}

.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 8px;

  .weekday {
    text-align: center;
    font-size: 12px;
    font-weight: 500;
    color: #909399;
    padding: 8px 0;
  }
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-template-rows: repeat(6, 1fr);
  gap: 4px;
}

.calendar-cell {
  min-height: 80px;
  padding: 6px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  border: 1px solid rgba(255, 138, 128, 0.1);
  transition: all 0.2s ease;
  overflow: hidden;

  &:hover {
    background: rgba(255, 255, 255, 0.8);
    box-shadow: 0 2px 8px rgba(255, 138, 128, 0.12);
  }

  &.other-month {
    background: rgba(255, 255, 255, 0.25);

    .day-number {
      color: #c0c4cc;
    }
  }

  &.today {
    border-color: rgba(255, 138, 128, 0.5);
    background: rgba(255, 138, 128, 0.08);

    .day-number {
      color: rgba(255, 138, 128, 1);
      font-weight: 700;
    }
  }

  &.has-tasks {
    background: rgba(255, 255, 255, 0.9);
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
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  &.priority-high {
    background: rgba(245, 108, 108, 0.15);
    color: #f56c6c;
    border-left: 2px solid #f56c6c;
  }

  &.priority-medium {
    background: rgba(230, 162, 60, 0.15);
    color: #e6a23c;
    border-left: 2px solid #e6a23c;
  }

  &.priority-low {
    background: rgba(103, 194, 58, 0.15);
    color: #67c23a;
    border-left: 2px solid #67c23a;
  }

  &.completed {
    opacity: 0.5;
    text-decoration: line-through;
  }

  &:hover {
    transform: translateX(2px);
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  }

  .task-chip-title {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.more-tasks {
  font-size: 9px;
  color: #909399;
  padding: 2px 6px;
  text-align: center;
}
</style>
