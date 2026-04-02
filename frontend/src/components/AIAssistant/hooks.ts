import { fetchGetChatList } from "@/api/chat";
import type {
  BubbleListItemProps,
  BubbleListProps,
} from "vue-element-plus-x/types/BubbleList";
import { SSEClient } from "@/utils/sse";
import dayjs from "dayjs";
import eventBus from "@/utils/eventbus";

type ListType = {
  createdAt?: string;
  id?: string;
  taskId?: string;
  title?: string;
};

export const bubbleList = ref<
  BubbleListProps<BubbleListItemProps & ListType>["list"]
>([]);
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
            return {
              placement: item.role === "user" ? "end" : "start",
              loading: false,
              content: item.chatContent,
              role: item.role === "user" ? "assistant" : "user",
              isMarkdown: item.role !== "user",
              createdAt: item.chatCreatedAt,
              id: item.chatId,
              taskId: item.taskId,
              title: item.title,
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
  event: "send_response" | "received" | "complete" | "wait_response";
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
      if (data.event === "received") {
        bubbleList.value.push({
          placement: "end",
          loading: false,
          content: data.message,
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
          bubbleList.value[index].content =
            `${bubbleList.value[index].content}${data.message}`;
        } else {
          bubbleList.value.push({
            placement: "start",
            loading: true,
            content: data.message,
            isMarkdown: true,
            createdAt: dayjs().format("YYYY-MM-DD HH:mm:ss"),
            id: data.chatId,
          });
        }
      } else if (data.event === "complete") {
        sseLoading.value = false;
        const index = bubbleList.value.findIndex(
          (item) => item.id === data.chatId,
        );
        if (index !== -1) {
          bubbleList.value[index].loading = false;
          bubbleList.value[index].content =
            `${bubbleList.value[index].content}${data.message}`;
        }
        eventBus.emit("sseComplete");
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
