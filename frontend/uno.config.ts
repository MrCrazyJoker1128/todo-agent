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
      primary: "#ff8a80",
    },
  },
});
