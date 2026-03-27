<template>
  <header class="layout-header">
    <div class="header-left">
      <i class="i-ep-user header-icon"></i>
    </div>
    <div class="header-center">
      <h1 class="app-title">{{ title }}</h1>
    </div>
    <div class="header-right">
      <i class="i-ep-setting header-icon"></i>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  title?: string;
}>();
</script>

<style scoped lang="scss">
.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid rgba(255, 138, 128, 0.2);
  backdrop-filter: blur(10px);
  position: relative;
  z-index: 2;

  .header-left,
  .header-right {
    width: 40px;
    display: flex;
    align-items: center;
    justify-content: center;

    .header-icon {
      font-size: 20px;
      color: #606266;
      cursor: pointer;
      transition: color 0.3s;

      &:hover {
        color: #ff8a80;
      }
    }
  }

  .header-center {
    flex: 1;
    text-align: center;

    .app-title {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }
}
</style>
