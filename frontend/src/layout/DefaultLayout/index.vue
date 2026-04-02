<template>
  <div class="layout-container">
    <BaseHeader title="AI TODO Agent" />

    <div class="layout-body">
      <TaskPanel />

      <AIAssistant>
        <template #chat-messages>
          <slot name="chat-messages"></slot>
        </template>
      </AIAssistant>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseHeader from "@/components/BaseHeader/index.vue";
import TaskPanel from "@/components/TaskPanel/index.vue";
import AIAssistant from "@/components/AIAssistant/index.vue";
import { useUserStore } from "@/store/user";

const userStore = useUserStore();
const router = useRouter();

async function checkUserInfo() {
  const success = await userStore.getUserInfo();
  if (!success) {
    router.push("/login");
  }
}
checkUserInfo();
</script>

<style scoped lang="scss">
.layout-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(135deg, #fff5f5 0%, #fff0ed 40%, #ffe8e6 100%);
  position: relative;
}

@keyframes float {
  0%,
  100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -30px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

.layout-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}
</style>
