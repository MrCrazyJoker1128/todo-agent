<template>
  <main class="task-panel">
    <div class="panel-header flex-col">
      <div class="content-tabs w-full">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="日历" name="calendar" />
          <el-tab-pane label="列表" name="list" />
        </el-tabs>
      </div>
      <div class="flex-1 w-full flex-col">
        <div class="task-list-container">
          <div class="content-filters">
            <el-button
              :type="filter === 'today' ? 'primary' : 'default'"
              @click="filter = 'today'"
            >
              <i class="i-ep-calendar filter-icon"></i>
              今日任务
            </el-button>
            <el-button
              :type="filter === 'week' ? 'primary' : 'default'"
              @click="filter = 'week'"
            >
              <i class="i-ep-date filter-icon"></i>
              本周任务
            </el-button>
            <el-button
              :type="filter === 'completed' ? 'primary' : 'default'"
              @click="filter = 'completed'"
            >
              <i class="i-ep-circle-check filter-icon"></i>
              已完成
            </el-button>
          </div>
        </div>
        <div class="task-list flex-1">
          <slot name="task-list"></slot>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
const activeTab = ref("list");
const filter = ref("today");
</script>

<style scoped lang="scss">
.task-panel {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  background: transparent;
  margin-right: 4px;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
}

.panel-header {
  flex: 1;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(255, 138, 128, 0.08);
  overflow: hidden;
}

.content-tabs {
  margin-bottom: 16px;

  :deep(.el-tabs__header) {
    margin: 0;
  }

  :deep(.el-tabs__nav-wrap::after) {
    height: 1px;
    background-color: rgba(255, 138, 128, 0.2);
  }
}

.task-list-container {
  margin-top: 16px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(255, 138, 128, 0.08);
}

.content-filters {
  display: flex;
  gap: 10px;

  .el-button {
    flex: 1;
    height: 36px;
    border-radius: 8px;
    font-size: 13px;

    .filter-icon {
      margin-right: 4px;
      font-size: 14px;
    }
  }
}

.task-list {
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
}

.task-calendar {
  min-height: 300px;
}
</style>
