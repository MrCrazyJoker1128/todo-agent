<template>
  <header class="wb-topbar">
    <div class="wb-topbar-left">
      <div class="wb-logo">
        <span class="wb-logo-mark">◈</span>
        <span class="wb-logo-text wb-font-display">{{ title }}</span>
      </div>
      <span class="wb-env-tag wb-mono">OPS // LIVE</span>
    </div>

    <div class="wb-topbar-center wb-mono">
      <span class="wb-clock">{{ clock }}</span>
      <span class="wb-date">{{ dateStr }}</span>
    </div>

    <div class="wb-topbar-right">
      <div class="wb-user">
        <span class="wb-user-avatar">{{ userInitial }}</span>
        <span class="wb-user-name">{{ userStore.userInfo?.nickname || "—" }}</span>
      </div>
      <button class="wb-icon-btn" title="模型配置" @click="goSettings">
        <i class="i-ep-setting"></i>
      </button>
      <button class="wb-icon-btn" title="退出登录" @click="handleLogout">
        <i class="i-ep-switch-button"></i>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUserStore } from "@/store/user";

defineProps<{
  title?: string;
}>();

const router = useRouter();
const userStore = useUserStore();

const userInitial = computed(() =>
  (userStore.userInfo?.nickname || "U").slice(0, 1).toUpperCase(),
);

const now = ref(new Date());
let timer: number | undefined;
onMounted(() => {
  timer = window.setInterval(() => {
    now.value = new Date();
  }, 1000);
});
onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});

const pad = (n: number) => String(n).padStart(2, "0");
const clock = computed(() => {
  const d = now.value;
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
});
const week = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
const dateStr = computed(() => {
  const d = now.value;
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${week[d.getDay()]}`;
});

function goSettings() {
  router.push("/settings");
}
function handleLogout() {
  userStore.logout();
  router.push("/login");
}
</script>

<style scoped lang="scss">
.wb-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 52px;
  padding: 0 16px;
  background: var(--wb-bg-deep);
  border-bottom: 1px solid var(--wb-line-soft);
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.wb-topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.wb-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.wb-logo-mark {
  color: var(--wb-accent);
  font-size: 18px;
  line-height: 1;
  text-shadow: 0 0 12px rgba(94, 234, 212, 0.6);
}

.wb-logo-text {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: var(--wb-text);
  white-space: nowrap;
}

.wb-env-tag {
  font-size: 10px;
  color: var(--wb-accent);
  border: 1px solid var(--wb-accent-line);
  background: var(--wb-accent-dim);
  border-radius: 4px;
  padding: 2px 6px;
  letter-spacing: 0.08em;
  white-space: nowrap;
}

.wb-topbar-center {
  display: flex;
  align-items: baseline;
  gap: 10px;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.wb-clock {
  font-size: 15px;
  font-weight: 600;
  color: var(--wb-text);
}

.wb-date {
  font-size: 11px;
  color: var(--wb-text-3);
  letter-spacing: 0.08em;
}

.wb-topbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.wb-user {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-right: 8px;
  border-right: 1px solid var(--wb-line-soft);
  margin-right: 4px;
}

.wb-user-avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: linear-gradient(135deg, #134e4a, #0f766e);
  border: 1px solid var(--wb-accent-line);
  color: var(--wb-accent);
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--wb-font-display);
}

.wb-user-name {
  font-size: 13px;
  color: var(--wb-text-2);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wb-icon-btn {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 6px;
  color: var(--wb-text-3);
  cursor: pointer;
  font-size: 15px;
  transition: all 0.2s;

  &:hover {
    color: var(--wb-accent);
    background: var(--wb-accent-dim);
    border-color: var(--wb-accent-line);
  }
}

@media (max-width: 720px) {
  .wb-topbar-center {
    display: none;
  }
  .wb-user-name {
    display: none;
  }
}
</style>
