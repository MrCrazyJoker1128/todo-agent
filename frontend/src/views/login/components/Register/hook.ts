import { fetchRegister } from "@/api/user";
import type { RegisterReq, RegisterRes } from "@/api/user/type";
import type { IResult } from "@/types/type";
import { useRequestHook } from "@/utils/hooks";
import { ElMessage, type FormInstance, type FormRules } from "element-plus";
import { Md5 } from "ts-md5";
import { type ModelRef } from "vue";

export function useRegister(visible: ModelRef<boolean>) {
  const formRef = ref<FormInstance>();
  const form = reactive<RegisterReq & { confirmPassword: string }>({
    nickname: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const validateConfirmPassword = (_: any, value: string, callback: any) => {
    if (value === "") {
      callback(new Error("请再次输入密码"));
    } else if (value !== form.password) {
      callback(new Error("两次输入的密码不一致"));
    } else {
      callback();
    }
  };

  const rules: FormRules = {
    nickname: [{ required: true, message: "请输入用户名", trigger: "blur" }],
    email: [
      { required: true, message: "请输入邮箱", trigger: "blur" },
      { type: "email", message: "请输入正确的邮箱格式", trigger: "blur" },
    ],
    password: [
      { required: true, message: "请输入密码", trigger: "blur" },
      { min: 6, message: "密码长度不能少于6位", trigger: "blur" },
    ],
    confirmPassword: [
      { required: true, validator: validateConfirmPassword, trigger: "blur" },
    ],
  };
  const { func: registerFunc, copyLoadingOption: registerLoadingOption } =
    useRequestHook({
      requestFunc: fetchRegister,
      success: (res: IResult<RegisterRes>) => {
        if (res.code === 200) {
          ElMessage.success("注册成功，请登录");
          formRef.value?.resetFields();
          visible.value = false;
        }
      },
    });
  const handleRegister = async () => {
    if (!formRef.value) return;
    await formRef.value.validate((valid) => {
      if (valid) {
        const data: RegisterReq = {
          nickname: form.nickname,
          email: form.email,
          password: Md5.hashStr(form.password),
        };
        registerFunc(data);
      }
    });
  };

  return {
    formRef,
    form,
    rules,
    registerLoadingOption,
    handleRegister,
  };
}
