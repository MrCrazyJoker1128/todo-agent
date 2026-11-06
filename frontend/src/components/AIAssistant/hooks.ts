import { fetchGetChatList } from "@/api/chat";
import type {
  BubbleListItemProps,
  BubbleListProps,
} from "vue-element-plus-x/types/BubbleList";
import { SSEClient } from "@/utils/sse";
import dayjs from "dayjs";
import eventBus from "@/utils/eventbus";
import { ElMessageBox } from "element-plus";
import { useRouter } from "vue-router";

export interface TaskFlowTask {
  id: string;
  title: string;
  description: string;
  dueDate: string;
  priority: string;
  taskType: string;
  status: string;
}

export interface TaskFlowPayload {
  kind: "candidates" | "edit" | "updated" | "delete-confirm" | "deleted";
  tasks: TaskFlowTask[];
  selectedId?: string;
  mode?: "edit" | "delete";
  // 历史回放的卡默认已决，直接锁死。
  decided?: boolean;
}

export interface SummaryPayload {
  kind: "summary-confirm" | "summary";
  period: string;
  startDate: string;
  endDate: string;
  taskType: string;
  priority: string;
  missing?: string[];
  total?: number;
  done?: number;
  doing?: number;
  todo?: number;
  cancelled?: number;
  rate?: number;
  byType?: Record<string, number>;
  byPriority?: Record<string, number>;
  recordId?: string;
  title?: string;
  content?: string;
  decided?: boolean;
}

type ListType = {
  createdAt?: string;
  id?: string;
  taskId?: string;
  title?: string;
  payload?: TaskFlowPayload;
  summary?: SummaryPayload;
};

export const bubbleList = ref<
  BubbleListProps<BubbleListItemProps & ListType>["list"]
>([]);
const CONFIRM_RE = /^(确认创建任务|确认修改任务|我选删除任务|确认删除任务|取消任务修改|取消删除任务|确认生成总结|取消生成总结)[:：]/;
const THINK_RE = /<think>([\s\S]*?)(?:<\/think>|$)/;
const TASK_FENCE_RE =
  /```task-((?:candidates|edit|updated|delete-confirm|deleted|summary-confirm|summary))\n([\s\S]*?)```/;

// beautifyConfirm 确认指令不展示 JSON：转成人话。非确认消息原样返回。
export function beautifyConfirm(content: string | undefined): string {
  if (!content) return content ?? "";
  const m = CONFIRM_RE.exec(content.trim());
  if (!m) return content;
  const action = m[1];
  let title = "";
  try {
    const raw = JSON.parse(content.slice(m[0].length)) as { title?: string; id?: string };
    title = raw.title ?? (raw.id ? `任务 ${raw.id.slice(0, 8)}` : "");
  } catch {
    // 解析失败也不要漏 JSON
  }
  switch (action) {
    case "确认创建任务":
      return title ? `✅ 确认创建任务「${title}」` : "✅ 确认创建任务";
    case "确认修改任务":
      return title ? `✅ 确认修改任务「${title}」` : "✅ 确认修改任务";
    case "我选删除任务":
      return title ? `🗑️ 选择删除「${title}」待二次确认` : "🗑️ 选择删除任务待二次确认";
    case "确认删除任务":
      return `🗑️ 确认删除${title ? `「${title}」` : "任务"}`;
    case "取消任务修改":
      return "已取消修改";
    case "取消删除任务":
      return "已取消删除";
    case "确认生成总结": {
      const names: Record<string, string> = { day: "日报", week: "周报", month: "月报", custom: "时段总结" };
      let period = "";
      try {
        const raw = JSON.parse(content.slice(m[0].length)) as { period?: string };
        period = names[raw.period ?? ""] ?? "时段总结";
      } catch {
        period = "时段总结";
      }
      return `✅ 确认生成${period}`;
    }
    case "取消生成总结":
      return "已取消生成总结";
    default:
      return "已确认";
  }
}

// foldThink 把 <think> 块收成引用块（markdown 原生支持不依赖 HTML 开关）。
// 流式阶段 think 原文短暂可见可接受，complete 定稿时折叠。
export function foldThink(content: string | undefined): string {
  if (!content) return "";
  const m = THINK_RE.exec(content);
  if (!m) return content;
  const think = m[1].trim();
  const rest = (content || "").replace(m[0], "").trim();
  if (!think) return rest;
  const quoted = think
    .split("\n")
    .map((line) => `> ${line}`)
    .join("\n");
  return `> 💭 思考过程\n>\n${quoted}\n\n---\n\n${rest}`.trim();
}

// parseTaskEnvelope 从消息文本提取人机验证 envelope；失败原样保留文本不抛错。
// fromHistory 为 true 时（历史回放）payload 标记 decided，卡片直接锁死防重复提交。
export function parseTaskEnvelope(
  content: string | undefined,
  fromHistory = false,
): {
  text: string;
  payload?: TaskFlowPayload;
  summary?: SummaryPayload;
} {
  if (!content) return { text: content ?? "" };
  const m = TASK_FENCE_RE.exec(content);
  if (!m) return { text: content };
  try {
    const raw = JSON.parse(m[2]) as {
      tasks?: TaskFlowTask[];
      count?: number;
      title?: string;
      description?: string;
      dueDate?: string;
      priority?: string;
      taskType?: string;
      status?: string;
      success?: boolean;
      taskId?: string;
      id?: string;
    };
    const kind = m[1] as TaskFlowPayload["kind"];
    const decided = fromHistory || undefined;
    if (kind === "candidates") {
      return {
        text: content.replace(m[0], "").trim(),
        payload: { kind, tasks: raw.tasks ?? [], mode: "edit", decided },
      };
    }
    if (kind === "edit") {
      return {
        text: content.replace(m[0], "").trim(),
        payload: {
          kind,
          tasks: [
            {
              id: "",
              title: raw.title ?? "",
              description: raw.description ?? "",
              dueDate: raw.dueDate ?? "",
              priority: raw.priority ?? "",
              taskType: raw.taskType ?? "",
              status: raw.status ?? "",
            },
          ],
          decided,
        },
      };
    }
    if (kind === "updated") {
      return {
        text: content.replace(m[0], "").trim(),
        payload: {
          kind,
          tasks: [
            {
              id: raw.taskId ?? "",
              title: "",
              description: "",
              dueDate: "",
              priority: "",
              taskType: "",
              status: "",
            },
          ],
          selectedId: raw.taskId,
          decided,
        },
      };
    }
    const k = m[1];
    if (k === "summary-confirm" || k === "summary") {
      const sraw = raw as unknown as {
        period?: string; startDate?: string; endDate?: string; taskType?: string; priority?: string;
        missing?: string[]; total?: number; done?: number; doing?: number; todo?: number;
        cancelled?: number; rate?: number; byType?: Record<string, number>; byPriority?: Record<string, number>;
        recordId?: string; title?: string; content?: string; stats?: unknown;
      };
      const summary: SummaryPayload = {
        kind: k as "summary-confirm" | "summary",
        period: sraw.period ?? "",
        startDate: sraw.startDate ?? "",
        endDate: sraw.endDate ?? "",
        taskType: sraw.taskType ?? "",
        priority: sraw.priority ?? "",
        missing: sraw.missing ?? (k === "summary-confirm" ? [] : undefined),
        total: sraw.total, done: sraw.done, doing: sraw.doing, todo: sraw.todo,
        cancelled: sraw.cancelled, rate: sraw.rate,
        byType: sraw.byType, byPriority: sraw.byPriority,
        recordId: sraw.recordId, title: sraw.title, content: sraw.content,
        decided,
      };
      return { text: content.replace(m[0], "").trim(), summary };
    }
    const tid = raw.id ?? "";
    return {
      text: content.replace(m[0], "").trim(),
      payload: {
        kind,
        tasks: [
          {
            id: tid,
            title: raw.title ?? "",
            description: "",
            dueDate: raw.dueDate ?? "",
            priority: "",
            taskType: "",
            status: "",
          },
        ],
        selectedId: tid,
        decided,
      },
    };
  } catch {
    return { text: content };
  }
}
export function useChatList(): {
  getChatList: () => void;
  loading: Ref<boolean>;
} {
  const loading = ref<boolean>(false);
  function getChatList() {
    loading.value = true;
    fetchGetChatList({
      pageIndex: 0,
      pageSize: 0,
    })
      .then((res) => {
        bubbleList.value =
          res.data?.data.list?.map((item) => {
            const isUser = item.role === "user";
            const parsed = parseTaskEnvelope(foldThink(item.chatContent), true);
            return {
              placement: isUser ? "end" : "start",
              loading: false,
              content: isUser ? beautifyConfirm(parsed.text) : parsed.text,
              role: item.role === "user" ? "assistant" : "user",
              isMarkdown: item.role !== "user",
              createdAt: item.chatCreatedAt,
              id: item.chatId,
              taskId: item.taskId,
              title: item.title,
              payload: parsed.payload,
              summary: parsed.summary,
            };
          }) || [];
      })
      .finally(() => {
        loading.value = false;
      });
  }
  return {
    getChatList,
    loading,
  };
}

export interface SendChatMessageRequest {
  /**
   * 消息操作
   */
  actions?: string[];
  /**
   * 消息内容
   */
  content?: string;
  /**
   * 消息角色
   */
  role?: string;
}
interface SendChatMessageResponse extends SendChatMessageRequest {
  event:
    | "send_response"
    | "received"
    | "complete"
    | "wait_response"
    | "error:NEED_MODEL_CONFIG"
    | "error:MODEL_AUTH"
    | "error:MODEL_CALL";
  // SSE 事件携带的文本（后端 api.SSEMessage.Message）。
  message?: string;
  chatId?: string;
  taskId?: string;
}
export function useSendChatMessage(): {
  handleSubmit: (data: SendChatMessageRequest) => void;
  sseClient: SSEClient<SendChatMessageRequest>;
  handleStop: () => void;
  sseLoading: Ref<boolean>;
} {
  const sseLoading = ref<boolean>(false);
  const sseClient = new SSEClient<SendChatMessageRequest>(
    `${import.meta.env.VITE_APP_BASE_API || ""}/chat/send`,
  );
  sseClient.setCallbacks({
    onMessage: (res) => {
      const data = JSON.parse(res.data) as SendChatMessageResponse;
      if (
        data.event === "error:NEED_MODEL_CONFIG" ||
        data.message === "NEED_MODEL_CONFIG"
      ) {
        sseLoading.value = false;
        ElMessageBox.confirm(
          "尚未配置模型 Key，聊天不可用。是否前往设置页配置？",
          "需要配置",
          { confirmButtonText: "去配置", cancelButtonText: "稍后" },
        )
          .then(() => {
            useRouter().push("/settings");
          })
          .catch(() => {});
        return;
      }
      if (
        data.event === "error:MODEL_AUTH" ||
        data.message === "MODEL_AUTH"
      ) {
        sseLoading.value = false;
        ElMessageBox.confirm(
          "模型 Key 无效或已过期（服务商返回 401），聊天不可用。是否前往设置页更换？",
          "Key 失效",
          { confirmButtonText: "去更换", cancelButtonText: "稍后" },
        )
          .then(() => {
            useRouter().push("/settings");
          })
          .catch(() => {});
        return;
      }
      if (data.event === "error:MODEL_CALL") {
        sseLoading.value = false;
        bubbleList.value.push({
          placement: "start",
          loading: false,
          content: `⚠️ ${data.message || "模型调用失败请稍后重试"}`,
          isMarkdown: false,
          createdAt: dayjs().format("YYYY-MM-DD HH:mm:ss"),
          id: data.chatId,
        });
        eventBus.emit("sseComplete");
        return;
      }
      if (data.event === "received") {
        bubbleList.value.push({
          placement: "end",
          loading: false,
          content: beautifyConfirm(data.message),
          isMarkdown: false,
          createdAt: dayjs().format("YYYY-MM-DD HH:mm:ss"),
          id: data.chatId,
        });
      } else if (data.event === "send_response") {
        const index = bubbleList.value.findIndex(
          (item) => item.id === data.chatId,
        );
        if (index !== -1) {
          bubbleList.value[index].loading = false;
          bubbleList.value[index].content = foldThink(
            `${bubbleList.value[index].content}${data.message}`,
          );
        } else {
          bubbleList.value.push({
            placement: "start",
            loading: true,
            content: foldThink(data.message),
            isMarkdown: true,
            createdAt: dayjs().format("YYYY-MM-DD HH:mm:ss"),
            id: data.chatId,
          });
        }
        eventBus.emit("chatScroll");
      } else if (data.event === "complete") {
        sseLoading.value = false;
        const index = bubbleList.value.findIndex(
          (item) => item.id === data.chatId,
        );
        if (index !== -1) {
          bubbleList.value[index].loading = false;
          bubbleList.value[index].content =
            `${bubbleList.value[index].content}${data.message}`;
          bubbleList.value[index].content = foldThink(bubbleList.value[index].content);
          const parsed = parseTaskEnvelope(bubbleList.value[index].content);
          bubbleList.value[index].content = parsed.text;
          bubbleList.value[index].payload = parsed.payload;
          bubbleList.value[index].summary = parsed.summary;
        }
        eventBus.emit("sseComplete");
        eventBus.emit("chatScroll");
        // AI 回环的增删改（创建/编辑确认）落库后，通知菜单 badge 刷新
        eventBus.emit("task-changed");
      }
    },
    onClose: () => {
      sseLoading.value = false;
    },
  });
  function handleStop(): void {
    sseClient.stop();
  }
  function handleSubmit(data: SendChatMessageRequest) {
    sseLoading.value = true;
    sseClient.setData(data);
    sseClient.start();
  }
  return {
    handleSubmit,
    sseClient,
    handleStop,
    sseLoading,
  };
}
