import { createApp } from "vue";
import "virtual:uno.css";
import "element-plus/theme-chalk/src/index.scss";
import "./style/element-plus.scss";
import "./style.css";
import App from "./App.vue";
import router from "./router";
import pinia from "./store";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

const app = createApp(App);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}
app.use(pinia);
app.use(router);
app.mount("#app");
