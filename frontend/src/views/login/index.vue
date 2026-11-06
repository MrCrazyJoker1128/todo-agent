<template>
  <div
    class="wb-scope min-h-screen w-full flex items-center justify-center relative overflow-hidden bg-surface-1 text-[var(--text-1)]"
  >
    <!-- 背景层：网格 + 辉光 -->
    <div class="login-bg" aria-hidden="true" />
    <div class="login-glow" aria-hidden="true" />

    <div class="w-full max-w-400 flex items-stretch gap-0 login-shell">
      <!-- 左：品牌区 -->
      <div class="login-brand hidden md:flex col-start-start gap-12 p-10">
        <div>
          <div class="flex items-center gap-3">
            <span class="login-mark">◈</span>
            <span class="font-display text-lg font-bold tracking-widest">TODO OPS</span>
          </div>
          <p class="font-mono text-xs tracking-[0.22em] text-[var(--text-3)] mt-3">
            OPS // LIVE
          </p>
        </div>

        <div>
          <h1 class="font-display text-4xl font-extrabold leading-tight">
            进入<br />你的作战室
          </h1>
          <p class="text-sm text-[var(--text-2)] mt-4 leading-relaxed">
            任务、日历、知识库与 AI 副驾，<br />同一块屏幕上并肩作战。
          </p>
          <ul class="login-points font-mono">
            <li><span class="dot"></span>MISSION BOARD · 任务中枢</li>
            <li><span class="dot"></span>AI COPILOT · 副驾待命</li>
            <li><span class="dot"></span>KNOWLEDGE · 弹药库</li>
          </ul>
        </div>

        <p class="font-mono text-xs text-[var(--text-3)] tracking-widest">
          SECURE // ENCRYPTED
        </p>
      </div>

      <!-- 右：表单区 -->
      <div class="login-form-wrap flex-1 flex items-center justify-center p-6 sm:p-10">
        <el-card class="login-card w-full max-w-95">
          <template #header>
            <div class="flex flex-col items-center gap-2 py-2">
              <span class="login-mark sm">◈</span>
              <span class="font-display text-xl font-bold tracking-wider">进入工作台</span>
              <span class="font-mono text-[11px] tracking-[0.2em] text-[var(--text-3)]">
                SIGN IN TO CONSOLE
              </span>
            </div>
          </template>
          <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
            <el-form-item prop="email">
              <el-input
                v-model="form.email"
                placeholder="邮箱"
                prefix-icon="Message"
                clearable
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                placeholder="密码"
                show-password
                prefix-icon="Lock"
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                class="login-button w-full"
                :loading="loginLoadingOption.loading"
                @click="handleLogin"
              >
                进入工作台
              </el-button>
            </el-form-item>
          </el-form>
          <div class="flex items-center justify-center gap-1 text-sm mt-4 text-[var(--text-2)]">
            <span>没有账号？</span>
            <el-link type="primary" @click="showRegister = true">立即注册</el-link>
          </div>
        </el-card>
      </div>
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

const { formRef, form, rules, loginLoadingOption, handleLogin } =
  useLogin(router);
void formRef;
</script>

<style scoped lang="scss">
.login-bg {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(94, 234, 212, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(94, 234, 212, 0.04) 1px, transparent 1px);
  background-size: 28px 28px;
  mask-image: radial-gradient(ellipse 90% 80% at 50% 40%, black 30%, transparent 75%);
}

.login-glow {
  position: absolute;
  top: -20%;
  left: 50%;
  transform: translateX(-50%);
  width: 900px;
  height: 500px;
  background: radial-gradient(
    ellipse,
    rgba(94, 234, 212, 0.09) 0%,
    transparent 65%
  );
  pointer-events: none;
}

.login-shell {
  position: relative;
  z-index: 1;
  background: var(--surface-2);
  border: 1px solid var(--line-soft);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.5);
  min-height: 560px;
}

.login-brand {
  width: 44%;
  background:
    radial-gradient(600px 300px at 20% 0%, rgba(94, 234, 212, 0.1), transparent 60%),
    var(--surface-0);
  border-right: 1px solid var(--line-soft);
}

.login-mark {
  color: var(--brand);
  font-size: 22px;
  line-height: 1;
  text-shadow: 0 0 14px var(--brand-glow);

  &.sm {
    font-size: 30px;
  }
}

.login-points {
  list-style: none;
  margin: 28px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  font-size: 11px;
  letter-spacing: 0.12em;
  color: var(--text-2);

  li {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--brand);
    box-shadow: 0 0 8px var(--brand-glow);
  }
}

.login-form-wrap {
  background:
    radial-gradient(500px 260px at 80% 100%, rgba(96, 165, 250, 0.06), transparent 60%),
    var(--surface-1);
}

.login-card {
  background: transparent;
  border: none;
  box-shadow: none;

  :deep(.el-card__header) {
    border-bottom: none;
    padding: 8px 0 4px;
    background: transparent;
  }

  :deep(.el-card__body) {
    padding: 12px 0 4px;
  }

  :deep(.el-input__wrapper) {
    background: var(--surface-0);
    border: 1px solid var(--line);
    border-radius: 8px;
    padding: 10px 14px;
    box-shadow: none;
    transition: border-color 0.2s, box-shadow 0.2s;

    &:hover {
      border-color: var(--brand-line);
    }

    &.is-focus {
      border-color: var(--brand);
      box-shadow: 0 0 0 2px var(--brand-dim);
    }
  }

  :deep(.el-input__inner) {
    color: var(--text-1);

    &::placeholder {
      color: var(--text-3);
    }
  }
}

.login-button {
  height: 44px;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.1em;
  border-radius: 8px;
}
</style>
