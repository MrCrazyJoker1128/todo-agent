import {
  fetchCreateBase,
  fetchDeleteBase,
  fetchDeleteFile,
  fetchGetBaseFiles,
  fetchListBases,
  fetchUploadMd,
  fetchCreateKnowledgeUrl,
} from "@/api/knowledge";
import type { KnowledgeFile, KnowledgeLibrary } from "@/api/knowledge/type";
import { ElMessage, ElMessageBox } from "element-plus";
import { useRequestHook } from "@/utils/hooks";
import eventBus from "@/utils/eventbus";

export function useKnowledgeBase() {
  const libs = ref<KnowledgeLibrary[]>([]);
  const libsLoading = ref(false);
  const activeId = ref("");
  const files = ref<KnowledgeFile[]>([]);
  const filesLoading = ref(false);
  const keyword = ref("");
  const page = ref(1);
  const pageSize = ref(10);
  const filteredFiles = computed(() => {
    const kw = keyword.value.trim();
    if (!kw) return files.value;
    return files.value.filter((f) => f.title.includes(kw));
  });
  const pageFiles = computed(() => {
    const start = (page.value - 1) * pageSize.value;
    return filteredFiles.value.slice(start, start + pageSize.value);
  });
  function resetPage() {
    page.value = 1;
  }

  const showCreate = ref(false);
  const showAdd = ref(false);
  const createForm = reactive({ name: "", description: "" });
  const createLoading = ref(false);
  const addTab = ref("url");
  const urlForm = reactive({ url: "" });
  const addUrlLoading = ref(false);
  const mdFile = ref<File | null>(null);
  const mdName = ref("");
  const dragOver = ref(false);
  const uploadLoading = ref(false);

  function openMdPicker() {
    document.getElementById("mdFileInput")?.click();
  }
  async function loadLibs() {
    libsLoading.value = true;
    try {
      const res = await fetchListBases();
      if (res.data.code === 200) {
        libs.value = res.data.data.items ?? [];
        if (!libs.value.find((l) => l.id === activeId.value)) {
          activeId.value = libs.value[0]?.id ?? "";
        }
        if (activeId.value) loadFiles();
        else files.value = [];
      }
    } finally {
      libsLoading.value = false;
    }
  }

  let pollTimer: number | undefined;
  function stopPoll() {
    if (pollTimer !== undefined) {
      clearInterval(pollTimer);
      pollTimer = undefined;
    }
  }
  function maybePoll() {
    stopPoll();
    if (files.value.some((f) => !f.status || f.status === "pending")) {
      pollTimer = window.setInterval(() => loadFiles(true), 3000);
    }
  }

  async function loadFiles(silent = false) {
    if (!activeId.value) return;
    if (!silent) filesLoading.value = true;
    try {
      const res = await fetchGetBaseFiles({ libraryId: activeId.value });
      if (res.data.code === 200) {
        files.value = res.data.data.items ?? [];
        resetPage();
        maybePoll();
      }
    } finally {
      if (!silent) filesLoading.value = false;
    }
  }

  function selectLib(id: string) {
    stopPoll();
    activeId.value = id;
    keyword.value = "";
    resetPage();
    loadFiles();
  }

  const { func: createFunc } = useRequestHook<{ libraryName: string; description?: string }>({
    requestFunc: fetchCreateBase,
    success: () => {
      ElMessage.success("知识库已创建");
      showCreate.value = false;
      createForm.name = "";
      createForm.description = "";
      loadLibs();
    },
  });

  function handleCreate() {
    if (!createForm.name.trim()) {
      ElMessage.warning("请输入库名");
      return;
    }
    createLoading.value = true;
    Promise.resolve(
      createFunc({ libraryName: createForm.name.trim(), description: createForm.description.trim() }),
    ).finally(() => {
      createLoading.value = false;
    });
  }

  function confirmDelete(lib: KnowledgeLibrary) {
    ElMessageBox.confirm(`删除知识库「${lib.name}」将同时删除其下所有文件，继续？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    }).then(() => {
      fetchDeleteBase({ id: lib.id }).then((res) => {
        if (res.data.code === 200) {
          ElMessage.success("已删除");
          if (activeId.value === lib.id) activeId.value = "";
          loadLibs();
        }
      });
    }).catch(() => {});
  }

  const { func: addUrlFunc } = useRequestHook<{
    url: string;
    libraryId: string;
  }>({
    requestFunc: fetchCreateKnowledgeUrl,
    success: () => {
      ElMessage.success("已添加，处理中");
      urlForm.url = "";
      showAdd.value = false;
      loadFiles();
      eventBus.emit("addComplete");
    },
  });

  function handleAddUrl() {
    if (!urlForm.url.trim()) {
      ElMessage.warning("请输入 html 页面地址");
      return;
    }
    if (!activeId.value) {
      ElMessage.warning("请先选择知识库");
      return;
    }
    addUrlLoading.value = true;
    Promise.resolve(
      addUrlFunc({ url: urlForm.url.trim(), libraryId: activeId.value }),
    ).finally(() => {
      addUrlLoading.value = false;
    });
  }
  function pickMdFile(file: File | null | undefined) {
    if (!file) return;
    if (!/\.md$|\.markdown$/i.test(file.name)) {
      ElMessage.warning("仅支持 md 文件");
      return;
    }
    if (file.size > 2 * 1024 * 1024) {
      ElMessage.warning("文件超过 2MB");
      return;
    }
    mdFile.value = file;
    mdName.value = file.name;
  }

  function handleMdInput(e: Event) {
    const input = e.target as HTMLInputElement;
    pickMdFile(input.files?.[0]);
    input.value = "";
  }

  function handleMdDrop(e: DragEvent) {
    dragOver.value = false;
    pickMdFile(e.dataTransfer?.files?.[0]);
  }

  function handleUploadMd() {
    if (!mdFile.value) {
      ElMessage.warning("请选择 md 文件");
      return;
    }
    if (!activeId.value) {
      ElMessage.warning("请先选择知识库");
      return;
    }
    const form = new FormData();
    form.append("file", mdFile.value);
    form.append("libraryId", activeId.value);
    uploadLoading.value = true;
    fetchUploadMd(form).then((res) => {
      if (res.data.code === 200) {
        ElMessage.success("已上传，处理中");
        mdFile.value = null;
        mdName.value = "";
        showAdd.value = false;
        loadFiles();
        eventBus.emit("addComplete");
      }
    }).finally(() => {
      uploadLoading.value = false;
    });
  }

  function handleDeleteFile(f: KnowledgeFile) {
    ElMessageBox.confirm(`删除文件「${f.title}」？`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    }).then(() => {
      fetchDeleteFile({ libraryId: activeId.value, source: f.source }).then((res) => {
        if (res.data.code === 200) {
          ElMessage.success("已删除");
          loadFiles();
        }
      });
    }).catch(() => {});
  }

  return {
    libs, libsLoading, activeId,
    files, filesLoading, keyword, filteredFiles, page, pageSize, pageFiles,
    showCreate, showAdd, createForm, createLoading,
    addTab, urlForm, addUrlLoading,
    mdName, dragOver, uploadLoading,
    loadLibs, selectLib, handleCreate, confirmDelete,
    handleAddUrl, openMdPicker, handleMdInput, handleMdDrop, handleUploadMd, handleDeleteFile,
    stopPoll, resetPage,
  };
}
