import { fetchEventSource } from "@microsoft/fetch-event-source";

export class SSEClient<T> {
  private controller: AbortController = new AbortController();
  private url: string;
  private data?: T;
  private onMessage?: (msg: any) => void;
  private onClose?: () => void;
  private onError?: (err: Error) => void;
  constructor(url: string) {
    this.url = url;
  }
  public setData(data: T) {
    this.data = data;
  }
  // 设置回调函数
  public setCallbacks(options?: {
    onMessage?: (msg: any) => void;
    onClose?: () => void;
    onError?: (err: Error) => void;
  }) {
    this.onMessage = options?.onMessage;
    this.onClose = options?.onClose;
    this.onError = options?.onError;
  }

  // 开始连接
  public start() {
    const token = localStorage.getItem("token");
    fetchEventSource(this.url, {
      method: "POST",
      signal: this.controller.signal,
      headers: {
        "Content-Type": "text/event-stream",
        "Cache-Control": "no-cache",
        Connection: "keep-alive",
        Authorization: token ?? "",
      },
      body: JSON.stringify(this.data),
      async onopen() {
        // 建立连接的回调
      },
      onmessage: (msg) => {
        // 接收一次数据段时回调，因为是流式返回，所以这个回调会被调用多次
        if (msg.event === "") {
          // 进行连接正常的操作
          if (this.onMessage) {
            this.onMessage(msg);
          }
        } else if (msg.event === "close") {
          this.controller.abort();
        }
      },
      onclose: () => {
        // 正常结束的回调
        this.controller.abort(); // 关闭连接
        if (this.onClose) {
          this.onClose();
        }
      },
      onerror: (err) => {
        // 连接出现异常回调
        if (this.onError) {
          this.onError(err);
        }
        throw err; // 必须抛出错误才会停止
      },
    });
  }

  // 关闭连接
  public stop() {
    this.controller.abort();
  }
}
