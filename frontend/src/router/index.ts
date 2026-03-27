import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/login/index.vue"),
    },
    {
      path: "/",
      component: () => import("@/layout/DefaultLayout/index.vue"),
    },
  ],
});

export default router;
