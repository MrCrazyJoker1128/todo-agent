import { useUserStore } from "@/store/user";
import type { LoginReq } from "@/api/user/type";
import { useRequestHook } from "@/utils/hooks";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { Md5 } from "ts-md5";
import type { Router } from "vue-router";

export function useLogin(router: Router) {
  const userStore = useUserStore();
  const formRef = ref<FormInstance>();
  const form = reactive<LoginReq>({
    email: "",
    password: "",
  });

  const rules: FormRules = {
    email: [
      { required: true, message: "请输入邮箱", trigger: "blur" },
      { type: "email", message: "请输入正确的邮箱格式", trigger: "blur" },
    ],
    password: [
      { required: true, message: "请输入密码", trigger: "blur" },
      { min: 6, message: "密码长度不能少于 6 位", trigger: "blur" },
    ],
  };

  const { func: loginFunc, copyLoadingOption: loginLoadingOption } =
    useRequestHook<LoginReq>({
      requestFunc: userStore.login,
      success: () => {
        ElMessage.success("登录成功");
        router.push("/");
      },
    });

  const handleLogin = async () => {
    if (!formRef.value) return;
    await formRef.value.validate((valid) => {
      if (valid) {
        const data: LoginReq = {
          email: form.email,
          password: Md5.hashStr(form.password),
        };
        loginFunc(data);
      }
    });
  };

  return {
    formRef,
    form,
    rules,
    loginLoadingOption,
    handleLogin,
  };
}
