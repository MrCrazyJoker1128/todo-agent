<template>
  <div class="task-summary-card">
    <div class="task-summary-title">{{ titleText }}</div>
    <!-- 维度确认 -->
    <div v-if="payload.kind === 'summary-confirm'" class="task-summary-confirm">
      <div v-if="payload.missing?.includes('period')" class="task-summary-warn">
        AI 未确定时间维度，请选择
      </div>
      <div v-else class="task-summary-hint">
        AI 建议：{{ periodName }}可修改
      </div>
      <el-form label-width="80px" size="small" :disabled="decided">
        <el-form-item label="时间维度">
          <el-radio-group v-model="form.period">
            <el-radio value="day">日</el-radio>
            <el-radio value="week">周</el-radio>
            <el-radio value="month">月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select v-model="form.taskType" placeholder="全部">
            <el-option label="全部" value="" />
            <el-option label="工作" value="work" />
            <el-option label="个人" value="personal" />
            <el-option label="项目" value="project" />
            <el-option label="临时" value="adhoc" />
          </el-select>
        </el-form-item>
        <el-form-item label="紧急程度">
          <el-select v-model="form.priority" placeholder="全部">
            <el-option label="全部" value="" />
            <el-option label="高" value="high" />
            <el-option label="中" value="medium" />
            <el-option label="低" value="low" />
          </el-select>
        </el-form-item>
      </el-form>
      <div class="task-flow-actions">
        <el-button size="small" :disabled="decided" @click="decide('取消生成总结')">取消</el-button>
        <el-button size="small" type="primary" :disabled="decided" @click="confirmGenerate">
          {{ decided ? "已确认" : "确认生成" }}
        </el-button>
      </div>
    </div>
    <!-- 总结结果：只读 -->
    <div v-else class="task-summary-result">
      <div class="task-summary-range">{{ rangeText }}</div>
      <div class="task-summary-tags">
        <el-tag size="small">总数 {{ payload.total ?? 0 }}</el-tag>
        <el-tag size="small" type="success">完成 {{ payload.done ?? 0 }}</el-tag>
        <el-tag size="small" type="warning">进行中 {{ payload.doing ?? 0 }}</el-tag>
        <el-tag size="small" type="info">待办 {{ payload.todo ?? 0 }}</el-tag>
        <el-tag v-if="(payload.cancelled ?? 0) > 0" size="small" type="danger">
          已取消 {{ payload.cancelled }}
        </el-tag>
      </div>
      <pre class="task-summary-content">{{ payload.content ?? "" }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from "element-plus";
import type { SummaryPayload } from "@/components/AIAssistant/hooks";

const props = defineProps<{
  payload: SummaryPayload;
}>();
const emit = defineEmits<{
  confirm: [text: string];
}>();

// decided：用户已确认/取消，该卡锁死防重复提交。历史回放的卡初始即锁。
const decided = ref<boolean>(props.payload.decided ?? false);
function decide(text: string): void {
  if (decided.value) return;
  decided.value = true;
  emit("confirm", text);
}

const periodNameMap: Record<string, string> = {
  day: "日报",
  week: "周报",
  month: "月报",
  custom: "时段总结",
};
const periodName = periodNameMap[props.payload.period] ?? "时段总结";
const titleText =
  props.payload.kind === "summary-confirm" ? "确认总结维度" : (props.payload.title ?? "工作总结");
const rangeText = [props.payload.startDate, props.payload.endDate]
  .filter(Boolean)
  .join(" ~ ");

const form = reactive({
  period: props.payload.period === "custom" ? "" : (props.payload.period ?? ""),
  taskType: "",
  priority: "",
});

function confirmGenerate(): void {
  if (decided.value) return;
  if (!form.period) {
    ElMessage.warning("请选择时间维度");
    return;
  }
  const body = {
    period: form.period,
    startDate: props.payload.startDate ?? "",
    endDate: props.payload.endDate ?? "",
    taskType: form.taskType,
    priority: form.priority,
  };
  decide(`确认生成总结：${JSON.stringify(body)}`);
}
</script>

<style scoped lang="scss">
.task-summary-card {
  min-width: 0;
  max-width: 100%;
}
.task-summary-title {
  font-weight: 600;
  margin-bottom: 8px;
}
.task-summary-warn {
  color: var(--el-color-danger);
  font-size: 12px;
  margin-bottom: 8px;
}
.task-summary-hint {
  color: var(--wb-text-3, #909399);
  font-size: 12px;
  margin-bottom: 8px;
}
.task-summary-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin: 8px 0;
}
.task-summary-range {
  font-size: 12px;
  color: var(--wb-text-3, #909399);
}
.task-summary-content {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 13px;
  line-height: 1.7;
  margin: 4px 0 0;
  font-family: inherit;
}
.task-flow-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}
</style>
