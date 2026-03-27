import type { AxiosResponse } from "axios";
import { defineStore } from "pinia";
import { useLocalStorage } from "@vueuse/core";
import { fetchLogin, fetchGetUserInfo } from "@/api/user";
import type { LoginReq, LoginRes, UserInfo } from "@/api/user/type";
import type { IResult } from "@/types/type";

export const useUserStore = defineStore("user", () => {
  const token = useLocalStorage<string>("token", "");
  const userInfo = ref<UserInfo | null>(null);

  const isLoggedIn = computed(() => !!token.value);

  function login(req: LoginReq): Promise<IResult<LoginRes>> {
    return new Promise<IResult<LoginRes>>((resolve, reject) => {
      fetchLogin(req).then((res) => {
        if (res.data.code === 200) {
          console.log(res.data);
          token.value = res.data.data.token;
          resolve(res.data);
        } else if (res.data.code === 401 || res.data.code === 403) {
          logout();
        } else {
          reject(res.data);
        }
      });
    });
  }

  async function getUserInfo() {
    if (!token.value) {
      return false;
    }
    try {
      const res = (await fetchGetUserInfo()) as unknown as AxiosResponse<
        IResult<UserInfo>
      >;
      if (res.data.code === 200) {
        userInfo.value = res.data.data;
        return true;
      }
      if (res.data.code === 401 || res.data.code === 403) {
        logout();
      }
      return false;
    } catch {
      // logout();
      return false;
    }
  }

  function logout() {
    token.value = "";
    userInfo.value = null;
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    getUserInfo,
    logout,
  };
});
