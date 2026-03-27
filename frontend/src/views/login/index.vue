<template>
  <div class="login-container">
    <div class="login-content col-center-center">
      <el-card class="login-card">
        <template #header>
          <div class="card-header row-center-center">
            <i class="i-ep-check mr-[8px]"></i>
            <span>{{ title }}</span>
          </div>
        </template>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              placeholder="请输入邮箱"
              prefix-icon="User"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              show-password
              prefix-icon="Lock"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              class="login-button"
              :loading="loginLoadingOption.loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        <div class="footer row-center-center">
          <span>没有账号？</span>
          <el-link type="primary" @click="showRegister = true"
            >立即注册</el-link
          >
        </div>
      </el-card>
    </div>
    <Register v-model:visible="showRegister" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import Register from "./components/Register/index.vue";
import { useLogin } from "./hook";

const router = useRouter();
const showRegister = ref(false);
const title = import.meta.env.VITE_TITLE;

const { formRef, form, rules, loginLoadingOption, handleLogin } =
  useLogin(router);
void formRef;
</script>

<style scoped lang="scss">
.login-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: linear-gradient(135deg, #fff5f5 0%, #fff0ed 40%, #ffe8e6 100%);
  position: relative;

  &::before {
    content: "";
    position: absolute;
    top: -30%;
    left: -30%;
    width: 160%;
    height: 160%;
    background: radial-gradient(
      circle,
      rgba(255, 138, 128, 0.15) 0%,
      rgba(255, 138, 128, 0.05) 40%,
      transparent 70%
    );
    animation: float 15s ease-in-out infinite;
  }
}

@keyframes float {
  0%,
  100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -30px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

.login-content {
  width: 100%;
  height: 100%;
}

.login-card {
  width: 420px;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.3);
  z-index: 1;

  :deep(.el-card__header) {
    padding: 24px;
    border-bottom: none;
    background: transparent;
  }

  :deep(.el-card__body) {
    padding: 0 32px 32px;
  }
}

.card-header {
  flex-direction: column;
  gap: 8px;

  span {
    font-size: 24px;
    font-weight: 600;
    color: #2c3e50;
  }

  i {
    font-size: 48px;
    color: #ff8a80;
  }
}

.el-form-item {
  margin-bottom: 20px;

  :deep(.el-input__wrapper) {
    border-radius: 12px;
    padding: 12px 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    }

    &.is-focus {
      box-shadow: 0 4px 16px rgba(255, 138, 128, 0.2);
    }
  }
}

.login-button {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(135deg, #ff8a80 0%, #ff5252 100%);
  border: none;
  box-shadow: 0 4px 16px rgba(255, 138, 128, 0.4);
  transition: all 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(255, 138, 128, 0.5);
  }

  &:active {
    transform: translateY(0);
  }
}

.footer {
  margin-top: 16px;
  gap: 4px;
  color: #666;
}
</style>
