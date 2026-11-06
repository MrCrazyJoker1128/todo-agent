import {
  defineConfig,
  presetUno,
  presetAttributify,
  presetIcons,
} from "unocss";

export default defineConfig({
  presets: [presetUno(), presetAttributify(), presetIcons()],
  shortcuts: [
    [
      /^row-(start|center|end|between|around|evenly)-(start|center|end|stretch)$/,
      ([, justify, items]) => `flex justify-${justify} items-${items}`,
    ],
    [
      /^col-(start|center|end|between|around|evenly)-(start|center|end|stretch)$/,
      ([, justify, items]) => `flex flex-col justify-${justify} items-${items}`,
    ],
  ],
  theme: {
    colors: {
      // 品牌主色：引用 theme.css 全局 var（--brand 系列）
      brand: {
        DEFAULT: "var(--brand)",
        strong: "var(--brand-strong)",
        100: "var(--brand-100)",
        200: "var(--brand-200)",
        300: "var(--brand-300)",
        deep: "var(--brand-deep)",
        ink: "var(--brand-ink)",
      },
      // 兼容旧 primary 用法
      primary: "var(--brand)",
      // 状态色
      amber: "var(--status-amber)",
      danger: "var(--status-red)",
      success: "var(--status-green)",
      info: "var(--status-blue)",
      // 界面层级
      surface: {
        0: "var(--surface-0)",
        1: "var(--surface-1)",
        2: "var(--surface-2)",
        3: "var(--surface-3)",
        4: "var(--surface-4)",
      },
    },
    fontFamily: {
      display: "var(--font-display)",
      body: "var(--font-body)",
      mono: "var(--font-mono)",
    },
  },
});
