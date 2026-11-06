<template>
  <nav class="wb-rail">
    <button
      v-for="item in navs"
      :key="item.path"
      :class="['wb-rail-btn', { active: isActive(item.path) }]"
      :title="item.label"
      @click="goNav(item.path)"
    >
      <el-badge
        :value="item.count"
        :hidden="!item.count"
        :max="99"
        class="wb-rail-badge"
      >
        <el-icon :size="19"><component :is="item.icon" /></el-icon>
      </el-badge>
      <span class="wb-rail-label">{{ item.label }}</span>
    </button>
    <div class="wb-rail-foot">
      <div class="wb-sys-dot" :class="sysOk ? 'ok' : 'bad'"></div>
      <span class="wb-mono wb-sys-text">{{ sysOk ? "SYS OK" : "SYS ERR" }}</span>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { House, Files, List, Setting } from "@element-plus/icons-vue";
import type { Component } from "vue";
import { useNavStore } from "@/store/nav";
import eventBus from "@/utils/eventbus";

const route = useRoute();
const router = useRouter();
const navStore = useNavStore();

const props = defineProps<{
  sysOk?: boolean;
}>();

interface NavItem {
  label: string;
  path: string;
  icon: Component;
  count?: number;
}
const navs = computed<NavItem[]>(() => [
  { label: "首页", path: "/home", icon: House },
  { label: "任务", path: "/tasks", icon: List, count: navStore.todayCount },
  { label: "知识库", path: "/knowledge", icon: Files },
  { label: "设置", path: "/settings", icon: Setting },
]);

function isActive(path: string) {
  return route.path.startsWith(path);
}

function goNav(path: string) {
  router.push(path);
}

onMounted(() => {
  navStore.refreshTodayCount();
});
eventBus.on("task-changed", () => {
  navStore.refreshTodayCount();
});
</script>

<style scoped lang="scss">
.wb-rail {
  width: 76px;
  flex-shrink: 0;
  background: var(--wb-bg-deep);
  border-right: 1px solid var(--wb-line-soft);
  display: flex;
  flex-direction: column;
  align-items: stretch;
  padding: 12px 8px;
  gap: 4px;
  z-index: 5;
}

.wb-rail-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 4px 8px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 8px;
  color: var(--wb-text-3);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  i {
    font-size: 19px;
  }

  &:hover {
    color: var(--wb-text-2);
    background: var(--wb-panel-2);
  }

  &.active {
    color: var(--wb-accent);
    background: var(--wb-accent-dim);
    border-color: var(--wb-accent-line);

    &::before {
      content: "";
      position: absolute;
      left: -8px;
      top: 20%;
      bottom: 20%;
      width: 3px;
      border-radius: 2px;
      background: var(--wb-accent);
      box-shadow: 0 0 8px rgba(94, 234, 212, 0.8);
    }
  }
}

.wb-rail-label {
  font-size: 11px;
  line-height: 1;
}

 .wb-rail-badge {
   display: inline-flex;
   :deep(.el-badge__content) {
     top: 2px;
     right: 2px;
   }
 }

.wb-rail-foot {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding-top: 12px;
  border-top: 1px solid var(--wb-line-soft);
}

.wb-sys-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;

  &.ok {
    background: var(--wb-green);
    box-shadow: 0 0 8px rgba(52, 211, 153, 0.9);
    animation: wb-pulse 2s infinite;
  }
  &.bad {
    background: var(--wb-red);
    box-shadow: 0 0 8px rgba(248, 113, 113, 0.9);
  }
}

.wb-sys-text {
  font-size: 9px;
  letter-spacing: 0.1em;
  color: var(--wb-text-3);
}

@keyframes wb-pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}
</style>
