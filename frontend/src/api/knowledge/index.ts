import { type AxiosResponse } from "axios";
import { type IResult } from "./../../types/type";
import type {
  CreateKnowledgeUrlRequest,
  CreateKnowledgeUrlResponse,
  GetKnowledgeUrlsListRequest,
  GetKnowledgeUrlsListResponse,
} from "./type";
import service from "@/utils/request";
export function fetchGetKnowledgeUrlsList(req: GetKnowledgeUrlsListRequest) {
  return service.get<GetKnowledgeUrlsListResponse>("/knowledge/list", {
    params: req,
  }) as unknown as Promise<
    AxiosResponse<IResult<GetKnowledgeUrlsListResponse>>
  >;
}

export function fetchCreateKnowledgeUrl(req: CreateKnowledgeUrlRequest) {
  return service.post<CreateKnowledgeUrlResponse>(
    "/knowledge/add",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<CreateKnowledgeUrlResponse>>>;
}
