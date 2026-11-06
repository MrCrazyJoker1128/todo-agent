import { type AxiosResponse } from "axios";
import { type IResult } from "./../../types/type";
import type {
  CreateKnowledgeUrlRequest,
  CreateKnowledgeUrlResponse,
  GetKnowledgeUrlsListRequest,
  GetKnowledgeUrlsListResponse,
  KnowledgeFile,
  KnowledgeLibrary,
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
export function fetchDeleteKnowledgeUrl(req: { id: string }) {
  return service.delete<{ message: string }>("/knowledge/delete", {
    params: req,
  }) as unknown as Promise<AxiosResponse<IResult<{ message: string }>>>;
}

export function fetchListBases() {
  return service.get<{ items: KnowledgeLibrary[] }>(
    "/knowledge/bases",
  ) as unknown as Promise<
    AxiosResponse<IResult<{ items: KnowledgeLibrary[] }>>
  >;
}

export function fetchCreateBase(req: { libraryName: string; description?: string }) {
  return service.post<{ id: string }>(
    "/knowledge/bases",
    req,
  ) as unknown as Promise<AxiosResponse<IResult<{ id: string }>>>;
}

export function fetchDeleteBase(req: { id: string }) {
  return service.delete<{ message: string }>("/knowledge/bases/delete", {
    params: req,
  }) as unknown as Promise<AxiosResponse<IResult<{ message: string }>>>;
}

export function fetchGetBaseFiles(req: { libraryId: string }) {
  return service.get<{ items: KnowledgeFile[] }>(
    "/knowledge/files",
    { params: req },
  ) as unknown as Promise<
    AxiosResponse<IResult<{ items: KnowledgeFile[] }>>
  >;
}

export function fetchDeleteFile(req: { libraryId: string; source: string }) {
  return service.delete<{ message: string }>("/knowledge/files/delete", {
    params: req,
  }) as unknown as Promise<AxiosResponse<IResult<{ message: string }>>>;
}

export function fetchUploadMd(form: FormData) {
  return service.post<{ title: string; source: string }>(
    "/knowledge/upload",
    form,
  ) as unknown as Promise<
    AxiosResponse<IResult<{ title: string; source: string }>>
  >;
}
