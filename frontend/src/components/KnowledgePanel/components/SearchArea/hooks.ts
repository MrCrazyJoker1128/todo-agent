import { fetchCreateKnowledgeUrl } from "@/api/knowledge";
import type { CreateKnowledgeUrlRequest } from "@/api/knowledge/type";
import type { LoadingOption } from "@/types/type";
import eventBus from "@/utils/eventbus";
import { useRequestHook } from "@/utils/hooks";
import { ElMessage } from "element-plus";

export function useCreateKnowledgeUrlHooks(): {
  loadingOption: LoadingOption;
  knowledgeUrlObj: CreateKnowledgeUrlRequest;
  handleAdd: () => void;
} {
  const knowledgeUrlObj = reactive<CreateKnowledgeUrlRequest>({ url: "" });
  const { copyLoadingOption: loadingOption, func: createKnowledgeUrlFunc } =
    useRequestHook({
      requestFunc: fetchCreateKnowledgeUrl,
      success: (res) => {
        if (res.code === 200) {
          knowledgeUrlObj.url = "";
          ElMessage.success("添加成功");
          eventBus.emit("addComplete");
        }
      },
    });
  function handleAdd() {
    createKnowledgeUrlFunc(knowledgeUrlObj);
  }
  return {
    loadingOption,
    knowledgeUrlObj,
    handleAdd,
  };
}
