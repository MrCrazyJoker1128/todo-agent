<template>
  <aside class="ai-assistant col-start-start">
    <div class="assistant-header">
      <h2>AI COPILOT</h2>
    </div>
    <div class="chat-messages" v-loading="loading">
      <BubbleList :list="bubbleList" ref="bubbleListRef">
        <template #avatar="{ item }">
          <el-icon v-show="item.placement === 'end'"><UserFilled /></el-icon>
          <el-icon v-show="item.placement === 'start'"><Monitor /></el-icon>
        </template>
        <template #footer="{ item }">
          <div>{{ item.createdAt }}</div>
          <TaskSummaryCard
            v-if="item.summary"
            :payload="item.summary"
            @confirm="(text: string) => handleConfirm(text)"
          />
          <TaskFlowCard
            v-else-if="item.payload"
            :payload="item.payload"
            @confirm="(text: string) => handleConfirm(text)"
          />
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
import { UserFilled, Monitor } from "@element-plus/icons-vue";
import { EditorSender, BubbleList } from "vue-element-plus-x";
import { useChatList, useSendChatMessage, bubbleList } from "./hooks";
import TaskFlowCard from "@/components/TaskFlowCard/index.vue";
import TaskSummaryCard from "@/components/TaskSummaryCard/index.vue";
import eventBus from "@/utils/eventbus";
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
      { id: "edit-task", name: "编辑任务" },
      { id: "delete-task", name: "删除任务" },
      { id: "query-task", name: "查询任务" },
      { id: "add-knowledge", name: "添加知识库" },
      { id: "search_knowledge", name: "查询知识库" },
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

// 人机验证卡确认：组装文本走正常 SSE 回环进 agent，由 agent 调工具执行。
function handleConfirm(text: string) {
  sendChatMessage({
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
eventBus.on("open-create-task", () => {
  editorSenderRef.value?.setText("/添加任务");
});
eventBus.on("chatScroll", () => {
  nextTick(() => {
    bubbleListRef.value?.scrollToBottom();
  });
});
onUnmounted(() => {
  eventBus.off("chatScroll");
});
</script>
<style scoped lang="scss">
.ai-assistant {
  width: 360px;
  flex-shrink: 0;
  background: var(--wb-bg-deep);
  border-left: 1px solid var(--wb-line-soft);
  position: relative;
  z-index: 5;
  display: flex;
  flex-direction: column;
  min-height: 0;

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: -1px;
    width: 1px;
    height: 120px;
    background: linear-gradient(180deg, var(--wb-accent), transparent);
    opacity: 0.6;
  }

  .assistant-header {
    padding: 14px 16px 12px;
    border-bottom: 1px solid var(--wb-line-soft);
    display: flex;
    align-items: center;
    gap: 8px;

    &::before {
      content: "";
      width: 7px;
      height: 7px;
      border-radius: 50%;
      background: var(--wb-green);
      box-shadow: 0 0 8px rgba(52, 211, 153, 0.9);
      animation: wb-pulse 2s infinite;
    }

    h2 {
      margin: 0;
      font-size: 12px;
      font-weight: 700;
      letter-spacing: 0.2em;
      color: var(--wb-text-2);
      font-family: var(--wb-font-display);
    }
  }

  .chat-messages {
    padding: 14px;
    overflow-y: auto;
    overflow-x: hidden;
    width: 100%;
    box-sizing: border-box;
    flex: 1;
    min-height: 0;
    min-width: 0;

    &::-webkit-scrollbar {
      width: 6px;
    }
    &::-webkit-scrollbar-thumb {
      background: var(--wb-line);
      border-radius: 3px;
    }

    :deep(.bubble-list) {
      background: transparent;
      max-width: 100%;
    }

    // markdown 表格/代码块等宽内容不把气泡撑出侧栏：文本可断行，
    // 表格/代码块允许在自身内部横向滚动
    :deep(.el-bubble),
    :deep(.el-bubble-content-wrapper),
    :deep(.el-bubble-content) {
      max-width: 100%;
      min-width: 0;
    }
    :deep(.el-bubble-content) {
      overflow-wrap: anywhere;
      word-break: break-word;
    }
    :deep(.el-bubble-content pre),
    :deep(.el-bubble-content table) {
      max-width: 100%;
      overflow-x: auto;
      display: block;
    }
    // 思考过程引用块收紧内边距
    :deep(.el-bubble-content blockquote) {
      margin: 4px 0;
      padding: 6px 0;
      font-size: 12px;
      line-height: 1.6;
    }
    :deep(.el-bubble-content blockquote p) {
      margin: 2px 0;
    }
    :deep(.el-bubble-content hr) {
      margin: 6px 0;
    }
  }

  .chat-input {
    padding: 12px 14px 14px;
    border-top: 1px solid var(--wb-line-soft);
    width: 100%;
    box-sizing: border-box;
    background: rgba(11, 14, 20, 0.6);

    .input-label {
      margin: 0 0 8px;
      font-size: 10px;
      letter-spacing: 0.16em;
      color: var(--wb-text-3);
      font-family: var(--wb-font-mono);
    }
  }
}

@keyframes wb-pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

@media (max-width: 1024px) {
  .ai-assistant {
    width: 300px;
  }
}
@media (max-width: 860px) {
  .ai-assistant {
    display: none;
  }
}
</style>
