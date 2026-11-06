<template>
  <el-dialog
    v-model="visible"
    title="注册"
    width="420px"
    :close-on-click-modal="false"
    class="register-dialog"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
      <el-form-item prop="nickname">
        <el-input
          v-model="form.nickname"
          placeholder="请输入用户名"
          prefix-icon="User"
        />
      </el-form-item>
      <el-form-item prop="email">
        <el-input
          v-model="form.email"
          placeholder="请输入邮箱"
          prefix-icon="Message"
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
      <el-form-item prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          show-password
          prefix-icon="Lock"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          :loading="registerLoadingOption.loading"
          @click="handleRegister"
          >注册</el-button
        >
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { useRegister } from "./hook";
defineOptions({
  name: "Register",
});

const visible = defineModel<boolean>("visible", {
  default: false,
});
const { formRef, form, rules, registerLoadingOption, handleRegister } =
  useRegister(visible);

// 告诉 TypeScript 该变量是有意保留用于模板引用的
void formRef;
function handleClose(): void {
  formRef.value?.resetFields();
  visible.value = false;
}
</script>

<style scoped lang="scss">
.register-dialog {
  :deep(.el-dialog) {
    border-radius: var(--radius);
    background: var(--surface-3);
    border: 1px solid var(--line);
  }

  :deep(.el-dialog__header) {
    padding: 20px 24px;
    border-bottom: 1px solid var(--line-soft);
  }

  :deep(.el-dialog__title) {
    font-size: 17px;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: var(--text-1);
    font-family: var(--font-display);
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid var(--line-soft);
  }

  :deep(.el-input__wrapper) {
    background: var(--surface-0);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: none;

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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
