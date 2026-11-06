<template>
  <div class="wb-scope wb-home" v-loading="loading">
    <div class="wb-home-head">
      <span class="wb-mono wb-eyebrow">OVERVIEW</span>
      <h2 class="wb-font-display">首页</h2>
    </div>

    <div class="wb-home-grid">
      <el-card class="home-card">
        <template #header>
          <span>API 状态</span>
        </template>

        <template v-if="configInfo?.configured">
          <el-tag type="success" class="home-tag">已配置</el-tag>
          <el-descriptions :column="1" border class="home-status">
            <el-descriptions-item label="API Key">
              {{ configInfo.apiKeyMask }}
            </el-descriptions-item>
            <el-descriptions-item label="Base URL">
              {{ configInfo.baseUrl ?? "内置默认" }}
            </el-descriptions-item>
            <el-descriptions-item label="模型">
              {{ configInfo.model ?? "内置默认" }}
            </el-descriptions-item>
          </el-descriptions>
        </template>
        <template v-else>
          <el-alert
            type="warning"
            :closable="false"
            show-icon
            title="尚未配置模型 Key，聊天与知识库提取不可用。"
            class="home-status"
          />
          <el-button type="primary" @click="goSettings">前往配置</el-button>
        </template>
      </el-card>

      <el-card class="home-card">
        <template #header>
          <span>任务概况</span>
        </template>

        <div class="home-stats">
          <div class="home-stat">
            <span class="home-stat-label">今日待办</span>
            <span class="home-stat-num wb-mono">{{ todayPending }}</span>
          </div>
          <div class="home-stat">
            <span class="home-stat-label">本周任务</span>
            <span class="home-stat-num wb-mono">{{ weekTotal }}</span>
          </div>
          <div class="home-stat">
            <span class="home-stat-label">本周已完成</span>
            <span class="home-stat-num wb-mono">{{ weekDone }}</span>
          </div>
        </div>
        <div class="home-rate">
          <span class="home-stat-label">本周完成率</span>
          <el-progress :percentage="weekRate" />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useHomeDashboard } from "./hook";

defineOptions({
  name: "HomeView",
});

const router = useRouter();
const {
  loading,
  configInfo,
  todayPending,
  weekTotal,
  weekDone,
  weekRate,
} = useHomeDashboard();

function goSettings() {
  router.push("/settings");
}
</script>

<style scoped lang="scss">
.wb-home {
  padding: 16px 20px;
  overflow-y: auto;
}

.wb-home-head {
  margin-bottom: 16px;
}

.wb-eyebrow {
  font-size: 10px;
  letter-spacing: 0.22em;
  color: var(--wb-text-3);
}

.wb-home-head h2 {
  margin: 4px 0 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 0.04em;
  color: var(--wb-text);
}

.wb-home-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.home-tag {
  margin-bottom: 12px;
}

.home-status {
  margin-bottom: 16px;
}

.home-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
}

.home-stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.home-stat-label {
  font-size: 12px;
  color: var(--wb-text-2);
}

.home-stat-num {
  font-size: 28px;
  font-weight: 700;
  color: var(--wb-text);
}

.home-rate {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
