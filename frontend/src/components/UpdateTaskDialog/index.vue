<template>
  <el-dialog v-model="dialogVisible" title="编辑任务" width="500">
    <el-form :model="taskForm">
      <el-form-item label="任务标题" prop="title">
        <el-input v-model="taskForm.title" placeholder="请输入任务标题" />
      </el-form-item>
      <el-form-item label="任务描述" prop="description">
        <el-input v-model="taskForm.description" placeholder="请输入任务描述" />
      </el-form-item>
      <el-form-item label="任务截止时间" prop="dueDate">
        <el-date-picker
          v-model="taskForm.dueDate"
          type="datetime"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
          placeholder="请选择任务截止时间"
        />
      </el-form-item>
      <el-form-item label="任务优先级" prop="priority">
        <el-radio-group v-model="taskForm.priority">
          <el-radio label="LOW">低</el-radio>
          <el-radio label="MEDIUM">中</el-radio>
          <el-radio label="HIGH">高</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="任务类型" prop="taskType">
        <el-radio-group v-model="taskForm.taskType">
          <el-radio label="work">工作</el-radio>
          <el-radio label="personal">个人</el-radio>
          <el-radio label="project">项目</el-radio>
          <el-radio label="adhoc">临时</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="copyLoadingOption.loading"
        >
          确认
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import {
  useTaskFormDialog,
  useTaskFormSubmit,
} from "@/components/UpdateTaskDialog/hooks";
defineOptions({
  name: "UpdateTaskDialog",
});
const dialogVisible = defineModel<boolean>("dialogVisible", {
  default: false,
});
const emit = defineEmits(["success"]);
const { taskForm, setTaskForm, resetTaskForm } = useTaskFormDialog();
const { submitForm, copyLoadingOption } = useTaskFormSubmit(handleSuccess);

function handleClose(): void {
  dialogVisible.value = false;
  resetTaskForm();
}
function handleSuccess(): void {
  handleClose();
  emit("success");
}
function handleSubmit(): void {
  submitForm(taskForm);
}
defineExpose({
  setTaskForm,
});
</script>

<style lang="sass" scoped></style>
