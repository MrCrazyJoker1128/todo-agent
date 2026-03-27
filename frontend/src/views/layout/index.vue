<template>
  <DefaultLayout>
    <template #task-list>
      <TaskItem v-for="task in taskList" :key="task.id" :task="task" />
    </template>

    <template #chat-messages>
      <ChatBubble type="ai">已为您创建任务</ChatBubble>
      <TaskCard :task="currentTask" />
    </template>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref } from "vue";
import DefaultLayout from "@/layout/DefaultLayout/index.vue";
import TaskItem from "@/components/TaskItem/index.vue";
import ChatBubble from "@/components/ChatBubble/index.vue";
import TaskCard from "@/components/TaskCard/index.vue";

interface Task {
  id: string;
  title: string;
  deadline: string;
  priority: "high" | "medium" | "low";
  completed: boolean;
}

interface TaskDetail {
  title: string;
  time: string;
  priority: "高" | "中" | "低";
  tags: string[];
}

const taskList = ref<Task[]>([
  {
    id: "1",
    title: "整理会议纪要",
    deadline: "今天 18:00",
    priority: "high",
    completed: false,
  },
  {
    id: "2",
    title: "整理会议纪要",
    deadline: "今天 18:00",
    priority: "medium",
    completed: false,
  },
  {
    id: "3",
    title: "整理会议纪要",
    deadline: "今天 18:00",
    priority: "low",
    completed: false,
  },
  {
    id: "4",
    title: "整理会议纪要",
    deadline: "今天 18:00",
    priority: "medium",
    completed: false,
  },
]);

const currentTask = ref<TaskDetail>({
  title: "整理会议纪要",
  time: "明天 10:00",
  priority: "高",
  tags: ["工作", "会议"],
});
</script>
