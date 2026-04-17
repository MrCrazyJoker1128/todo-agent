import mitt, { type Emitter } from "mitt";

type Events = {
  sseMessage: any;
  sseComplete: void;
  sseError: Error;
  addComplete: void;
};

const emitter: Emitter<Events> = mitt<Events>();

export const eventBus = {
  on: emitter.on,
  off: emitter.off,
  emit: emitter.emit,
  all: emitter.all,
};

export default eventBus;
