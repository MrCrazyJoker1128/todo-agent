import { fetchUpdateTask } from "@/api/task";
import type { UpdateTaskReq } from "@/api/task/type";
import { useReactiveReset, useRequestHook } from "@/utils/hooks";
import { ElMessage } from "element-plus";

export function useTaskFormDialog() {
  const { state: taskForm, resetState } = useReactiveReset<UpdateTaskReq>({
    id: "",
    description: "",
    title: "",
    dueDate: "",
    priority: undefined,
    taskType: undefined,
  });
  function setTaskForm(task: UpdateTaskReq) {
    resetState({
      id: task.id,
      description: task.description,
      title: task.title,
      dueDate: task.dueDate,
      priority: task.priority,
      taskType: task.taskType,
    });
  }
  return {
    taskForm,
    setTaskForm,
    resetTaskForm: resetState,
  };
}
export function useTaskFormSubmit(callback: () => void) {
  const { func, copyLoadingOption } = useRequestHook({
    requestFunc: fetchUpdateTask,
    success: (res) => {
      if (res.code === 200) {
        ElMessage.success("更新成功");
        callback();
      }
    },
  });
  return {
    submitForm: func,
    copyLoadingOption,
  };
}
