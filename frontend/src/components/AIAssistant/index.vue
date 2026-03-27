<template>
  <aside class="ai-assistant">
    <div class="assistant-header">
      <h2>AI 任务助手</h2>
    </div>

    <div class="chat-messages">
      <slot name="chat-messages"></slot>
    </div>

    <div class="chat-input">
      <p class="input-label">还需要我帮您做什么？</p>
      <div class="input-wrapper">
        <el-input placeholder="输入指令或提问..." v-model="inputMessage">
          <template #append>
            <el-button><i class="i-ep-arrow-right"></i></el-button>
          </template>
        </el-input>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
const inputMessage = ref("");
</script>

<style scoped lang="scss">
.ai-assistant {
  width: 320px;
  background: transparent;
  border-left: 1px solid rgba(255, 138, 128, 0.1);
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;

  .assistant-header {
    padding: 16px;
    border-bottom: 1px solid rgba(255, 138, 128, 0.1);

    h2 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
    }
  }

  .chat-messages {
    flex: 1;
    padding: 16px;
    overflow-y: auto;
  }

  .chat-input {
    padding: 16px;
    border-top: 1px solid rgba(255, 138, 128, 0.1);

    .input-label {
      margin: 0 0 8px;
      font-size: 12px;
      color: #909399;
    }

    .input-wrapper {
      :deep(.el-input-group__append) {
        padding: 0 12px;
      }
    }
  }
}
</style>
