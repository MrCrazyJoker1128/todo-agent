import type { AxiosResponse } from "axios";
import service from "@/utils/request";
import type { GetChatListReq, GetChatListResp } from "./type";
import type { IResult } from "@/types/type";

export function fetchGetChatList(
  req: GetChatListReq,
): Promise<AxiosResponse<IResult<GetChatListResp>>> {
  return service.get<GetChatListResp>("/chat/list", {
    params: req,
  }) as unknown as Promise<AxiosResponse<IResult<GetChatListResp>>>;
}
