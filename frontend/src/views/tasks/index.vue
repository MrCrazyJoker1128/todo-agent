<template>
  <TaskPanel @today-count="navStore.setTodayCount" />
</template>

<script setup lang="ts">
import TaskPanel from "@/components/TaskPanel/index.vue";
import { useNavStore } from "@/store/nav";

defineOptions({
  name: "TasksView",
});

const navStore = useNavStore();
</script>
