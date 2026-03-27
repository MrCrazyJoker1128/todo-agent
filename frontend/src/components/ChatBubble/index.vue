<template>
  <div class="chat-bubble" :class="type">
    <div class="bubble-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  type?: "user" | "ai";
}>();
</script>

<style scoped lang="scss">
.chat-bubble {
  margin-bottom: 12px;

  &.user {
    display: flex;
    justify-content: flex-end;

    .bubble-content {
      background: #409eff;
      color: #fff;
      border-radius: 12px 12px 0 12px;
      padding: 10px 14px;
      max-width: 80%;
    }
  }

  &.ai {
    display: flex;
    justify-content: flex-start;

    .bubble-content {
      background: #f4f4f5;
      color: #303133;
      border-radius: 12px 12px 12px 0;
      padding: 10px 14px;
      max-width: 80%;
    }
  }
}
</style>