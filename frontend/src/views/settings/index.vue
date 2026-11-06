<template>
  <div class="wb-scope wb-settings">
    <div class="wb-settings-head">
      <span class="wb-mono wb-eyebrow">CONFIG</span>
      <h2 class="wb-font-display">设置</h2>
    </div>
    <el-collapse v-model="activePanels" class="settings-collapse">
      <el-collapse-item name="model" class="settings-panel">
        <template #title>
          <span class="settings-panel-title">模型配置</span>
          <el-tag v-if="configInfo?.configured" type="success" size="small" class="settings-panel-tag">已配置</el-tag>
          <el-tag v-else type="warning" size="small" class="settings-panel-tag">未配置</el-tag>
        </template>
        <div v-loading="loadingConfig" class="settings-panel-body">

      <el-alert
        type="info"
        :closable="false"
        show-icon
        title="您的模型 Key 仅为您个人使用，服务端 AES-256-GCM 加密存储，页面仅展示脱敏结果。"
        class="settings-tip"
      />

      <el-descriptions
        v-if="configInfo?.configured"
        :column="1"
        border
        class="settings-status"
      >
        <el-descriptions-item label="配置状态">
          <el-tag type="success">已配置</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="API Key">
          {{ configInfo.apiKeyMask }}
        </el-descriptions-item>
        <el-descriptions-item label="Base URL">
          {{ configInfo.baseUrl ?? "内置默认" }}
        </el-descriptions-item>
        <el-descriptions-item label="模型">
          {{ configInfo.model ?? "内置默认" }}
        </el-descriptions-item>
        <el-descriptions-item label="嵌入模型">
          {{ (configInfo.embedModel ?? "默认") }} @ {{ (configInfo.embedBaseUrl ?? "本地 ollama 默认") }}
        </el-descriptions-item>
      </el-descriptions>
      <el-alert
        v-else
        type="warning"
        :closable="false"
        show-icon
        title="尚未配置模型 Key，聊天与知识库提取不可用。"
        class="settings-status"
      />

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="90px"
        class="settings-form"
      >
        <el-form-item label="API Key" prop="apiKey">
          <el-input
            v-model="form.apiKey"
            type="password"
            show-password
            placeholder="请输入您的模型 API Key，留空则不更新"
            clearable
          />
        </el-form-item>
        <el-form-item label="Base URL" prop="baseUrl">
          <el-input
            v-model="form.baseUrl"
            placeholder="例如：https://api.minimaxi.com/v1"
            clearable
          />
        </el-form-item>
        <el-form-item label="模型" prop="model">
          <el-input
            v-model="form.model"
            placeholder="例如：MiniMax-M2.7"
            clearable
          />
        </el-form-item>
        <el-divider content-position="left">生成参数</el-divider>
        <el-form-item label="Temperature" prop="temperature">
          <el-input-number
            v-model="form.temperature"
            :min="0"
            :max="2"
            :step="0.1"
            placeholder="默认 0.7"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="Max Tokens" prop="maxTokens">
          <el-input-number
            v-model="form.maxTokens"
            :min="1"
            :step="1"
            placeholder="默认 4096"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="系统提示词" prop="systemPrompt">
          <el-input
            v-model="form.systemPrompt"
            type="textarea"
            :rows="6"
            placeholder="空 = 用内置模板；填写后日期时间上下文块自动追加其后"
            clearable
          />
        </el-form-item>
        <el-divider content-position="left">嵌入模型（知识库向量化）</el-divider>
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          title="更换嵌入模型会改变向量维度与语义空间，已有知识库数据需清空重传，否则检索失效。"
          class="settings-tip"
        />
        <el-form-item label="Base URL" prop="embedBaseUrl">
          <el-input
            v-model="form.embedBaseUrl"
            placeholder="本地默认：http://localhost:11434（安装 ollama 后可用）"
            clearable
          />
        </el-form-item>
        <el-form-item label="嵌入模型" prop="embedModel">
          <el-input
            v-model="form.embedModel"
            placeholder="例如：qwen3-embedding:8b、nomic-embed-text"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="saveLoadingOption.loading"
            @click="handleSave"
          >
            保存配置
          </el-button>
        </el-form-item>
      </el-form>
        </div>
      </el-collapse-item>
      <el-collapse-item name="profile" class="settings-panel">
        <template #title>
          <span class="settings-panel-title">个人资料</span>
          <span v-if="profileForm.nickname" class="settings-panel-sub">{{ profileForm.nickname }}</span>
        </template>
        <div v-loading="loadingProfile" class="settings-panel-body">
      <el-form
        ref="profileRef"
        :model="profileForm"
        :rules="profileRules"
        label-width="90px"
        class="settings-form"
      >
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="profileForm.nickname" placeholder="显示名称" clearable />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" placeholder="登录邮箱，被占用时保存失败" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saveProfileLoading.loading" @click="handleSaveProfile">
            保存资料
          </el-button>
        </el-form-item>
      </el-form>
      <el-divider content-position="left">修改密码</el-divider>
      <el-form
        ref="passwordRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="90px"
        class="settings-form"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password clearable />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="至少 6 位" clearable />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="savePasswordLoading.loading" @click="handleSavePassword">
            修改密码
          </el-button>
        </el-form-item>
      </el-form>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup lang="ts">
import { useModelConfig, useProfile } from "./hook";

// 默认展开模型配置；个人资料收起（允许多项同时展开）。
const activePanels = ref<string[]>(["model"]);

const {
  formRef,
  form,
  rules,
  configInfo,
  loadingConfig,
  saveLoadingOption,
  handleSave,
} = useModelConfig();
void formRef;
const {
  profileRef,
  passwordRef,
  loadingProfile,
  profileForm,
  passwordForm,
  profileRules,
  passwordRules,
  saveProfileLoading,
  savePasswordLoading,
  handleSaveProfile,
  handleSavePassword,
} = useProfile();
void profileRef;
void passwordRef;
</script>

<style scoped lang="scss">
.wb-settings {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 32px 24px 48px;
  max-width: 680px;
  width: 100%;
  margin: 0 auto;
  background:
    radial-gradient(800px 300px at 50% -10%, rgba(94, 234, 212, 0.07), transparent 60%),
    var(--wb-bg);
  color: var(--wb-text);
  overflow-wrap: anywhere;

  &::-webkit-scrollbar {
    width: 8px;
  }
  &::-webkit-scrollbar-thumb {
    background: var(--wb-line);
    border-radius: 4px;
  }
  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

.wb-settings-head {
  margin-bottom: 20px;
}

.wb-eyebrow {
  font-size: 10px;
  letter-spacing: 0.22em;
  color: var(--wb-text-3);
}

.wb-settings-head h2 {
  margin: 4px 0 0;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.settings-tip {
  margin-bottom: 16px;
}

.settings-status {
  margin-bottom: 16px;
}

.settings-collapse {
  --el-collapse-header-bg-color: transparent;
  --el-collapse-content-bg-color: transparent;
  --el-collapse-header-text-color: var(--wb-text);
  --el-collapse-content-text-color: var(--wb-text);
  --el-collapse-border-color: var(--wb-line-soft);
  border-top: 1px solid var(--wb-line-soft);
  border-bottom: 1px solid var(--wb-line-soft);
}

.settings-panel-title {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.04em;
}

.settings-panel-tag {
  margin-left: 10px;
}

.settings-panel-sub {
  margin-left: 10px;
  font-size: 12px;
  color: var(--wb-text-3);
}

.settings-panel-body {
  padding: 4px 2px 12px;
}
</style>
