import { fetchGetKnowledgeUrlsList } from "@/api/knowledge";
import type {
  GetKnowledgeUrlsListRequest,
  GetKnowledgeUrlsListResponse,
  ServiceapiGetKnowledgeListItem,
  Status,
} from "@/api/knowledge/type";
import type { TableColumn } from "@/types/type";
import { usePageHook, type pageOptionType } from "@/utils/hooks";
import dayjs from "dayjs";

export function useGetKnowledgeListHooks(): {
  tableList: Ref<ServiceapiGetKnowledgeListItem[]>;
  initTableList: (searchReq?: GetKnowledgeUrlsListRequest) => void;
  getKnowledgeUrlsList: Function;
  pageOption: pageOptionType;
  tableHeadConfig: Ref<TableColumn[]>;
} {
  const tableHeadConfig: Ref<TableColumn[]> = ref([
    {
      label: "标题",
      prop: "title",
      minWidth: "150px",
    },
    {
      label: "作者",
      prop: "author",
    },
    {
      label: "状态",
      prop: "status",
      slot: "status",
    },
    {
      label: "域名",
      prop: "domain",
    },
    {
      label: "地址",
      prop: "url",
      slot: "url",
    },
    {
      label: "创建时间",
      prop: "createdAt",
    },
    {
      label: "操作",
      prop: "operation",
      slot: "operation",
    },
  ]);
  const tableList = ref<ServiceapiGetKnowledgeListItem[]>([]);
  const { func: getKnowledgeUrlsList, pageOption } =
    usePageHook<GetKnowledgeUrlsListResponse>({
      requestFunc: fetchGetKnowledgeUrlsList,
      success: (res) => {
        tableList.value =
          res.data.items?.map((item) => ({
            ...item,
            createdAt: item.createdAt
              ? dayjs(item.createdAt).format("YYYY-MM-DD HH:mm:ss")
              : "",
          })) || [];
      },
    });
  function initTableList(searchReq?: GetKnowledgeUrlsListRequest) {
    tableList.value = [];
    pageOption.current_page = 1;
    getKnowledgeUrlsList(searchReq || {});
  }
  return {
    tableList,
    initTableList,
    getKnowledgeUrlsList,
    pageOption,
    tableHeadConfig: tableHeadConfig,
  };
}

export function useStatusOptions(): {
  statusOptions: { label: string; value: Status }[];
  getStatusOptionsLabel: (key: Status) => string;
} {
  //pending|processing|completed|failed
  const statusOptions: { label: string; value: Status }[] = [
    { label: "等待", value: "pending" },
    { label: "处理中", value: "processing" },
    { label: "已完成", value: "completed" },
    { label: "处理失败", value: "failed" },
  ];
  function getStatusOptionsLabel(key: Status) {
    return statusOptions.find((item) => item.value === key)?.label || "";
  }
  return { statusOptions, getStatusOptionsLabel };
}
