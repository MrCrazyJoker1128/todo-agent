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
      children: [
        {
          path: "",
          redirect: "/home",
        },
        {
          path: "home",
          name: "home",
          component: () => import("@/views/home/index.vue"),
        },
        {
          path: "tasks",
          name: "tasks",
          component: () => import("@/views/tasks/index.vue"),
        },
        {
          path: "settings",
          name: "settings",
          component: () => import("@/views/settings/index.vue"),
        },
        {
          path: "knowledge",
          name: "knowledge",
          component: () => import("@/views/knowledge/index.vue"),
        },
      ],
    },
  ],
});

export default router;
