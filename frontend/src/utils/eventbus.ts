import mitt, { type Emitter } from "mitt";

type Events = {
  sseMessage: any;
  sseComplete: void;
  sseError: Error;
  addComplete: void;
  chatScroll: void;
  /** 任务增删改：创建/完成/取消后触发，菜单 badge 刷新 */
  "task-changed": void;
  "open-create-task": void;
};

const emitter: Emitter<Events> = mitt<Events>();

export const eventBus = {
  on: emitter.on,
  off: emitter.off,
  emit: emitter.emit,
  all: emitter.all,
};

export default eventBus;
