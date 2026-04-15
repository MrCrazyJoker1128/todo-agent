<template>
  <aside class="ai-assistant col-start-start">
    <div class="assistant-header">
      <h2>AI 任务助手</h2>
    </div>

    <div class="chat-messages flex-1" v-loading="loading">
      <BubbleList :list="bubbleList" ref="bubbleListRef">
        <template #avatar="{ item }">
          <el-icon v-show="item.placement === 'end'"><UserFilled /></el-icon>
          <el-icon v-show="item.placement === 'start'"><Monitor /></el-icon>
        </template>
        <template #footer="{ item }">
          <div>{{ item.createdAt }}</div>
          <div
            class="row-start-center mt-[6px]"
            v-if="item.placement === 'start'"
          >
            <el-button
              type="success"
              :icon="Check"
              circle
              size="small"
            ></el-button>
            <el-button
              type="danger"
              :icon="Delete"
              circle
              size="small"
              @click="
                handleDelete({ taskId: item.taskId, id: item.id as string })
              "
            />
          </div>
        </template>
      </BubbleList>
    </div>

    <div class="chat-input">
      <p class="input-label">还需要我帮您做什么？</p>
      <div class="input-wrapper">
        <EditorSender
          ref="editorSenderRef"
          :customTrigger="customTagsArr"
          placeholder="输入‘/’获取更多功能。"
          :loading="sseLoading"
          @cancel="handleStop"
          @submit="handleSubmit"
        />
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { Delete, Check, UserFilled, Monitor } from "@element-plus/icons-vue";
import { EditorSender, BubbleList } from "vue-element-plus-x";
import { useChatList, useSendChatMessage, bubbleList } from "./hooks";
import type {
  CustomTag,
  SubmitResult,
} from "vue-element-plus-x/types/EditorSender";
defineOptions({
  name: "AIAssistant",
  components: {
    EditorSender,
    BubbleList,
  },
});
const editorSenderRef = useTemplateRef("editorSenderRef");
const bubbleListRef = useTemplateRef("bubbleListRef");
const customTagsArr = ref<CustomTag[]>([
  {
    dialogTitle: "任务标签",
    prefix: "/",
    tagList: [
      { id: "add-task", name: "添加任务" },
      { id: "delete-task", name: "删除任务" },
      { id: "query-task", name: "查询任务" },
      { id: "add-knowledge", name: "添加知识库" },
    ],
  },
]);
const {
  handleSubmit: sendChatMessage,
  sseLoading,
  handleStop,
} = useSendChatMessage();
function handleSubmit(payload: SubmitResult) {
  const text = payload.text;
  const customTags = payload.customTags || {};

  const prefix = "/";
  const actions: string[] = [];

  customTags[prefix]?.forEach((tag: { name: string; id: string }) => {
    if (text.includes(tag.name)) {
      actions.push(tag.id);
    }
  });
  editorSenderRef.value?.clear();
  sendChatMessage({
    actions,
    content: text,
    role: "user",
  });
}

const { getChatList, loading } = useChatList();

onMounted(() => {
  getChatList();
  nextTick(() => {
    bubbleListRef.value?.scrollToBottom();
  });
});
function handleDelete(item: { taskId?: string; id: string }) {
  editorSenderRef.value?.setCustomTag("/", "delete-task");
  if (item.taskId) {
    editorSenderRef.value?.setText(`取消任务${item.taskId}`);
    return;
  }
  if (item.id) {
    editorSenderRef.value?.setText(`取消会话${item.id}的相关任务`);
  }
}
</script>

<style scoped lang="scss">
.ai-assistant {
  width: 320px;
  background: transparent;
  border-left: 1px solid rgba(255, 138, 128, 0.1);
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
    padding: 16px;
    overflow-y: auto;
    width: 100%;
    box-sizing: border-box;
    background: #fff;
    border-radius: 10px;
  }

  .chat-input {
    padding: 16px;
    border-top: 1px solid rgba(255, 138, 128, 0.1);
    width: 100%;
    box-sizing: border-box;

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
