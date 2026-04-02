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
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
}

interface TaskDetail {
  id: string;
  title: string;
  dueDate: string;
  priority: "high" | "medium" | "low";
  taskType: string;
  status: "todo" | "doing" | "done" | "cancelled";
  description?: string;
}

const taskList = ref<Task[]>([
  {
    id: "1",
    title: "整理会议纪要",
    deadline: "今天 18:00",
    priority: "high",
    taskType: "工作",
    status: "todo",
  },
  {
    id: "2",
    title: "回复邮件",
    deadline: "今天 18:00",
    priority: "medium",
    taskType: "日常",
    status: "doing",
  },
  {
    id: "3",
    title: "阅读文档",
    deadline: "今天 18:00",
    priority: "low",
    taskType: "学习",
    status: "todo",
  },
  {
    id: "4",
    title: "代码审查",
    deadline: "今天 18:00",
    priority: "medium",
    taskType: "工作",
    status: "done",
  },
]);

const currentTask = ref<TaskDetail>({
  id: "test-1",
  title: "整理会议纪要",
  dueDate: "2026-03-30 18:00",
  priority: "high",
  taskType: "工作",
  status: "todo",
  description: "会议结束后整理并发送给相关人员",
});
</script>
