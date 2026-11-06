<template>
  <div class="task-flow-card">
    <div class="task-flow-title">{{ titleText }}</div>
    <!-- 候选列表：修改模式选中行展开内联编辑，删除模式只单选 -->
    <div v-if="payload.kind === 'candidates'" class="task-flow-list">
      <div
        v-for="t in payload.tasks"
        :key="t.id"
        :class="['task-flow-row', { selected: selectedId === t.id, locked: decided }]"
        @click="!decided && (selectedId = t.id)"
      >
        <el-radio :model-value="selectedId" :value="t.id" @click.stop @change="!decided && (selectedId = t.id)">
          <span class="task-flow-row-title">{{ t.title }}</span>
        </el-radio>
        <div v-if="t.dueDate || t.priority" class="task-flow-meta">
          <span v-if="t.dueDate" class="task-flow-due">{{ t.dueDate }}</span>
          <el-tag v-if="t.priority" size="small">{{ t.priority }}</el-tag>
        </div>
      </div>
      <!-- 内联编辑：仅 edit 模式 -->
      <div v-if="payload.mode !== 'delete' && selectedTask" class="task-flow-edit">
        <el-form label-width="70px" size="small" :disabled="decided">
          <el-form-item label="任务标题">
            <el-input v-model="editForm.title" placeholder="请输入任务标题" />
          </el-form-item>
          <el-form-item label="任务描述">
            <el-input v-model="editForm.description" placeholder="请输入任务描述" />
          </el-form-item>
          <el-form-item label="截止时间">
            <el-date-picker
              v-model="editForm.dueDate"
              type="datetime"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              placeholder="请选择任务截止时间"
            />
          </el-form-item>
          <el-form-item label="优先级">
            <el-radio-group v-model="editForm.priority">
              <el-radio label="LOW">低</el-radio>
              <el-radio label="MEDIUM">中</el-radio>
              <el-radio label="HIGH">高</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="任务类型">
            <el-radio-group v-model="editForm.taskType">
              <el-radio label="work">工作</el-radio>
              <el-radio label="personal">个人</el-radio>
              <el-radio label="project">项目</el-radio>
              <el-radio label="adhoc">临时</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
      <div class="task-flow-actions">
        <el-button size="small" :disabled="decided" @click="decide(cancelText)">取消</el-button>
        <el-button size="small" type="primary" :disabled="decided" @click="confirmCandidates">{{ decided ? "已确认" : "确认" }}</el-button>
      </div>
    </div>
    <!-- 新增预览：空表单 -->
    <div v-else-if="payload.kind === 'edit'" class="task-flow-edit">
      <el-form label-width="70px" size="small" :disabled="decided">
        <el-form-item label="任务标题">
          <el-input v-model="editForm.title" placeholder="请输入任务标题" />
        </el-form-item>
        <el-form-item label="任务描述">
          <el-input v-model="editForm.description" placeholder="请输入任务描述" />
        </el-form-item>
        <el-form-item label="截止时间">
          <el-date-picker
            v-model="editForm.dueDate"
            type="datetime"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择任务截止时间"
          />
        </el-form-item>
        <el-form-item label="优先级">
          <el-radio-group v-model="editForm.priority">
            <el-radio label="LOW">低</el-radio>
            <el-radio label="MEDIUM">中</el-radio>
            <el-radio label="HIGH">高</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-radio-group v-model="editForm.taskType">
            <el-radio label="work">工作</el-radio>
            <el-radio label="personal">个人</el-radio>
            <el-radio label="project">项目</el-radio>
            <el-radio label="adhoc">临时</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div class="task-flow-actions">
        <el-button size="small" :disabled="decided" @click="decide('取消任务修改')">取消</el-button>
        <el-button size="small" type="primary" :disabled="decided" @click="confirmCreate">{{ decided ? "已创建" : "确认创建" }}</el-button>
      </div>
    </div>
    <!-- 修改结果：只读 -->
    <div v-else-if="payload.kind === 'updated'" class="task-flow-result">
      <el-tag type="success" size="small">已更新</el-tag>
      <span class="task-flow-row-title">{{ resultTitle }}</span>
    </div>
    <!-- 二次确认：红色警告 -->
    <div v-else-if="payload.kind === 'delete-confirm'" class="task-flow-danger">
      <el-alert type="error" :closable="false" show-icon :title="`确认删除「${selectedTask?.title ?? ''}」吗？此操作可在任务列表中找回（软删除）。`" />
      <div v-if="selectedTask?.dueDate" class="task-flow-due">截止：{{ selectedTask.dueDate }}</div>
      <div class="task-flow-actions">
        <el-button size="small" :disabled="decided" @click="decide('取消删除任务')">取消</el-button>
        <el-button size="small" type="danger" :disabled="decided" @click="confirmDelete">{{ decided ? "已删除" : "确认删除" }}</el-button>
      </div>
    </div>
    <!-- 删除结果：只读 -->
    <div v-else-if="payload.kind === 'deleted'" class="task-flow-result">
      <el-tag type="danger" size="small">已删除</el-tag>
      <span class="task-flow-row-title">{{ selectedTask?.title ?? "" }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus";
import type { TaskFlowPayload, TaskFlowTask } from "@/components/AIAssistant/hooks";

const props = defineProps<{
  payload: TaskFlowPayload;
}>();
const emit = defineEmits<{
  confirm: [text: string];
}>();

const isDeleteMode = props.payload.mode === "delete";
const titleText =
  props.payload.kind === "candidates"
    ? isDeleteMode
      ? "选择要删除的任务"
      : "选择要修改的任务"
    : props.payload.kind === "edit"
      ? "确认创建以下任务"
      : "";
const cancelText = isDeleteMode ? "取消删除任务" : "取消任务修改";
const resultTitle =
  props.payload.tasks.find((t) => t.id === props.payload.selectedId)?.title ?? "";

const selectedId = ref<string>(props.payload.selectedId ?? "");
const selectedTask = computed<TaskFlowTask | undefined>(() =>
  props.payload.tasks.find((t) => t.id === selectedId.value),
);
// decided：用户已确认/取消，该卡锁死防重复提交。历史回放的卡初始即锁。
const decided = ref<boolean>(props.payload.decided ?? false);
function decide(text: string): void {
  if (decided.value) return;
  decided.value = true;
  emit("confirm", text);
}
const editForm = reactive({
  title: "",
  description: "",
  dueDate: "",
  priority: "",
  taskType: "",
  status: "",
});

// 选中行变化时把原值灌入编辑表单；新增预览用 tasks[0] 空表单初值
watch(
  selectedTask,
  (t) => {
    if (!t) return;
    editForm.title = t.title ?? "";
    editForm.description = t.description ?? "";
    editForm.dueDate = t.dueDate ?? "";
    editForm.priority = t.priority ?? "";
    editForm.taskType = t.taskType ?? "";
    editForm.status = t.status ?? "";
  },
  { immediate: true },
);

function diffFields(): Record<string, string> {
  const t = selectedTask.value;
  const out: Record<string, string> = {};
  if (!t) return out;
  (["title", "description", "dueDate", "priority", "taskType", "status"] as const).forEach(
    (k) => {
      if (editForm[k] !== (t[k] ?? "") && editForm[k] !== "") out[k] = editForm[k];
    },
  );
  return out;
}

function confirmCandidates(): void {
  if (decided.value) return;
  if (!selectedId.value) {
    ElMessage.warning("请先选择一条任务");
    return;
  }
  if (isDeleteMode) {
    const t = selectedTask.value;
    decide(`我选删除任务：${JSON.stringify({ id: t?.id, title: t?.title })}`);
    return;
  }
  const diff = diffFields();
  if (Object.keys(diff).length === 0) {
    ElMessage.warning("请至少修改一个字段");
    return;
  }
  decide(`确认修改任务：${JSON.stringify({ id: selectedId.value, ...diff })}`);
}

function confirmCreate(): void {
  if (decided.value) return;
  if (!editForm.title) {
    ElMessage.warning("标题不能为空");
    return;
  }
  const body: Record<string, string> = { title: editForm.title };
  (["description", "dueDate", "priority", "taskType", "status"] as const).forEach((k) => {
    if (editForm[k]) body[k] = editForm[k];
  });
  decide(`确认创建任务：${JSON.stringify(body)}`);
}

function confirmDelete(): void {
  if (decided.value) return;
  if (!selectedId.value) {
    ElMessage.warning("请先选择一条任务");
    return;
  }
  decide(`确认删除任务：${JSON.stringify({ id: selectedId.value })}`);
}
</script>

<style scoped lang="scss">
.task-flow-card {
  min-width: 0;
  max-width: 100%;
}
.task-flow-title {
  font-weight: 600;
  margin-bottom: 8px;
}
.task-flow-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.task-flow-row {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 4px;
  padding: 6px 8px;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  &.selected {
    border-color: var(--el-color-primary);
    background: var(--el-color-primary-light-9);
  }
}
.task-flow-row-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.task-flow-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 22px;
}
.task-flow-due {
  font-size: 12px;
  color: var(--wb-text-3, #909399);
  white-space: nowrap;
}
.task-flow-edit {
  margin-top: 8px;
}
.task-flow-row.locked {
  cursor: default;
  opacity: 0.8;
}
.task-flow-result {
  display: flex;
  align-items: center;
  gap: 8px;
}
.task-flow-danger {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
