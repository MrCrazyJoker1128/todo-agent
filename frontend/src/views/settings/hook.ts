import {
  fetchGetModelConfig,
  fetchGetSystemKEK,
  fetchGetUserInfo,
  fetchSaveModelConfig,
  fetchSaveSystemKEK,
  fetchUpdatePassword,
  fetchUpdateProfile,
} from "@/api/user";
import type {
  ModelConfigInfo,
  SaveModelConfigReq,
  UpdatePasswordReq,
  UpdateProfileReq,
  UserInfo,
} from "@/api/user/type";
import type { IResult } from "@/types/type";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { useRequestHook } from "@/utils/hooks";

export function useModelConfig() {
  const formRef = ref<FormInstance>();
  const form = reactive<SaveModelConfigReq>({
    apiKey: "",
    baseUrl: "",
    model: "",
    temperature: undefined,
    maxTokens: undefined,
    systemPrompt: "",
    embedBaseUrl: "",
    embedModel: "",
  });
  const configInfo = ref<ModelConfigInfo | null>(null);
  const loadingConfig = ref(false);

  const rules: FormRules = {
    // 已配置过则可留空（= 不更新 Key），只有首次配置必须填。
    apiKey: [
      {
        validator: (_rule, value, callback) => {
          if (configInfo.value?.configured) return callback();
          if (!String(value ?? "").trim())
            return callback(new Error("请输入模型 API Key"));
          return callback();
        },
        trigger: "blur",
      },
    ],
    temperature: [
      {
        validator: (_rule, value, callback) => {
          if (value === undefined || value === null) return callback();
          if (typeof value !== "number" || value < 0 || value > 2)
            return callback(new Error("temperature 必须在 0~2 之间"));
          return callback();
        },
        trigger: "blur",
      },
    ],
    maxTokens: [
      {
        validator: (_rule, value, callback) => {
          if (value === undefined || value === null) return callback();
          if (!Number.isInteger(value) || value <= 0)
            return callback(new Error("maxTokens 必须大于 0"));
          return callback();
        },
        trigger: "blur",
      },
    ],
  };

  async function loadConfig() {
    loadingConfig.value = true;
    try {
      const res = await fetchGetModelConfig();
      if (res.data.code === 200) {
        configInfo.value = res.data.data;
        // 脱敏 key 不回填输入框；已配值回填方便修改，null 留空（placeholder 展示默认）
        form.baseUrl = res.data.data.baseUrl ?? "";
        form.model = res.data.data.model ?? "";
        // null = 未设置 → undefined（输入框留空、placeholder 展示默认）
        form.temperature = res.data.data.temperature ?? undefined;
        form.maxTokens = res.data.data.maxTokens ?? undefined;
        form.systemPrompt = res.data.data.systemPrompt ?? "";
        form.embedBaseUrl = res.data.data.embedBaseUrl ?? "";
        form.embedModel = res.data.data.embedModel ?? "";
      }
    } finally {
      loadingConfig.value = false;
    }
  }

  const { func: saveFunc, copyLoadingOption: saveLoadingOption } =
    useRequestHook<SaveModelConfigReq>({
      requestFunc: fetchSaveModelConfig,
      success: () => {
        ElMessage.success("模型配置已保存（Key 已加密存储）");
        form.apiKey = "";
        loadConfig();
      },
    });

  function handleSave() {
    if (!formRef.value) return;
    formRef.value.validate((valid) => {
      if (valid) {
        // undefined → null：明确告知后端“未设置”（JSON 丢 undefined 会让后端无法区分缺字段与清参数）
        saveFunc({
          ...form,
          temperature: form.temperature ?? null,
          maxTokens: form.maxTokens ?? null,
        });
      }
    });
  }

  loadConfig();

  return {
    formRef,
    form,
    rules,
    configInfo,
    loadingConfig,
    saveLoadingOption,
    handleSave,
    loadConfig,
  };
}

export function useSystemKEK() {
  const kekForm = reactive({ kek: "" });
  const kekConfigured = ref<boolean | null>(null);
  const saveKekLoading = ref(false);

  async function loadKEK() {
    try {
      const res = await fetchGetSystemKEK();
      if (res.data.code === 200) {
        kekConfigured.value = res.data.data.configured;
      }
    } catch {
      kekConfigured.value = null;
    }
  }

  const { func: saveKekFunc } = useRequestHook<{ kek: string }>({
    requestFunc: fetchSaveSystemKEK,
    success: () => {
      ElMessage.success("系统 KEK 已入库");
      kekForm.kek = "";
      loadKEK();
    },
  });

  function handleSaveKek() {
    if (!kekForm.kek) {
      ElMessage.warning("请输入 KEK 口令");
      return;
    }
    saveKekLoading.value = true;
    Promise.resolve(saveKekFunc({ kek: kekForm.kek })).finally(() => {
      saveKekLoading.value = false;
    });
  }

  loadKEK();

  return { kekForm, kekConfigured, saveKekLoading, handleSaveKek, loadKEK };
}

export function useProfile() {
  const profileRef = ref<FormInstance>();
  const passwordRef = ref<FormInstance>();
  const userInfo = ref<UserInfo | null>(null);
  const loadingProfile = ref(false);
  const profileForm = reactive<UpdateProfileReq>({ nickname: "", email: "" });
  const passwordForm = reactive<UpdatePasswordReq & { confirmPassword: string }>({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  const profileRules: FormRules = {
    email: [{ type: "email", message: "邮箱格式不正确", trigger: "blur" }],
  };
  const passwordRules: FormRules = {
    oldPassword: [{ required: true, message: "请输入旧密码", trigger: "blur" }],
    newPassword: [
      { required: true, message: "请输入新密码", trigger: "blur" },
      { min: 6, message: "新密码至少 6 位", trigger: "blur" },
    ],
    confirmPassword: [
      { required: true, message: "请再次输入新密码", trigger: "blur" },
      {
        validator: (_rule, value, callback) => {
          if (value !== passwordForm.newPassword) return callback(new Error("两次输入的新密码不一致"));
          return callback();
        },
        trigger: "blur",
      },
    ],
  };

  async function loadProfile() {
    loadingProfile.value = true;
    try {
      const res = await fetchGetUserInfo();
      if (res.data.code === 200) {
        userInfo.value = res.data.data;
        profileForm.nickname = res.data.data.nickname ?? "";
        profileForm.email = res.data.data.email ?? "";
      }
    } finally {
      loadingProfile.value = false;
    }
  }

  const { func: saveProfileFunc, copyLoadingOption: saveProfileLoading } =
    useRequestHook<UpdateProfileReq, UserInfo>({
      requestFunc: fetchUpdateProfile,
      success: (res: IResult<UserInfo>) => {
        ElMessage.success("个人资料已更新");
        userInfo.value = res.data;
        profileForm.nickname = res.data.nickname ?? "";
        profileForm.email = res.data.email ?? "";
      },
    });

  const { func: savePasswordFunc, copyLoadingOption: savePasswordLoading } =
    useRequestHook<UpdatePasswordReq>({
      requestFunc: fetchUpdatePassword,
      success: () => {
        ElMessage.success("密码修改成功");
        passwordForm.oldPassword = "";
        passwordForm.newPassword = "";
        passwordForm.confirmPassword = "";
      },
    });

  function handleSaveProfile() {
    if (!profileRef.value) return;
    profileRef.value.validate((valid) => {
      if (valid) {
        if (!profileForm.nickname && !profileForm.email) {
          ElMessage.warning("昵称和邮箱至少填一个");
          return;
        }
        saveProfileFunc({ ...profileForm });
      }
    });
  }

  function handleSavePassword() {
    if (!passwordRef.value) return;
    passwordRef.value.validate((valid) => {
      if (valid) {
        savePasswordFunc({
          oldPassword: passwordForm.oldPassword,
          newPassword: passwordForm.newPassword,
        });
      }
    });
  }

  loadProfile();

  return {
    profileRef,
    passwordRef,
    userInfo,
    loadingProfile,
    profileForm,
    passwordForm,
    profileRules,
    passwordRules,
    saveProfileLoading,
    savePasswordLoading,
    handleSaveProfile,
    handleSavePassword,
  };
}
