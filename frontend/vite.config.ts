import vue from "@vitejs/plugin-vue";
import { loadEnv, defineConfig } from "vite";
import { resolve } from "path";

import UnoCSS from "unocss/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import Icons from "unplugin-icons/vite";
import IconsResolver from "unplugin-icons/resolver";
import vueDevTools from "vite-plugin-vue-devtools";

const pathSrc = resolve(__dirname, "src");

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [
      vue(),
      UnoCSS(),
      AutoImport({
        imports: ["vue", "vue-router", "@vueuse/core"],
        resolvers: [
          // 自动导入 Element Plus 相关函数，如：ElMessage, ElMessageBox... (带样式)
          ElementPlusResolver(),
          // 自动导入图标组件
          IconsResolver({}),
        ],
        eslintrc: {
          enabled: true, // 是否自动生成 eslint 规则，建议生成之后设置 false
          filepath: "./.eslintrc-auto-import.json", // 指定自动导入函数 eslint 规则的文件
        },
        vueTemplate: true, // 是否在 vue 模板中自动导入
        dts: resolve(pathSrc, "types", "auto-imports.d.ts"),
      }),
      Components({
        resolvers: [
          // 自动导入 Element Plus 组件
          ElementPlusResolver(),
          // 自动注册图标组件
          IconsResolver({
            enabledCollections: ["ep"], // element-plus图标库，其他图标库 https://icon-sets.iconify.design/
          }),
        ],
        dts: resolve(pathSrc, "types", "components.d.ts"), //  自动导入组件类型声明文件位置，默认根目录
      }),
      Icons({
        // 自动安装图标库
        autoInstall: true,
      }),
      vueDevTools(),
    ],
    resolve: {
      alias: {
        "@": resolve(__dirname, "src"),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@use "@/style/variables.scss" as *;`,
          api: "modern-compiler",
        },
      },
    },
    server: {
      host: "0.0.0.0",
      port: 3000,
      open: true,
      proxy: {
        // 代理 /dev-api 的请求
        [env.VITE_APP_BASE_API]: {
          changeOrigin: true,
          target: "http://127.0.0.1:8080/api",
          rewrite: (path: string) =>
            path.replace(new RegExp("^" + env.VITE_APP_BASE_API), ""),
        },
      },
    },
  };
});
