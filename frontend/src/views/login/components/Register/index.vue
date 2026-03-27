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
    border-radius: 16px;
  }

  :deep(.el-dialog__header) {
    padding: 20px 24px;
    border-bottom: 1px solid #f0f0f0;
  }

  :deep(.el-dialog__title) {
    font-size: 18px;
    font-weight: 600;
    color: #2c3e50;
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid #f0f0f0;
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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;

  .el-button {
    border-radius: 8px;
  }

  .el-button--primary {
    background: linear-gradient(135deg, #ff8a80 0%, #ff5252 100%);
    border: none;
    box-shadow: 0 4px 16px rgba(255, 138, 128, 0.4);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(255, 138, 128, 0.5);
    }
  }
}
</style>
