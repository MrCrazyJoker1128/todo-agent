<template>
  <div class="wb-scope wb-workbench">
    <BaseHeader title="TODO OPS" />

    <div class="wb-body">
      <NavRail :sys-ok="sysOk" />

      <section class="wb-main">
        <router-view />
      </section>

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
import NavRail from "@/components/NavRail/index.vue";
import AIAssistant from "@/components/AIAssistant/index.vue";
import { useUserStore } from "@/store/user";

const userStore = useUserStore();
const router = useRouter();

const sysOk = ref(true);

async function checkUserInfo() {
  const success = await userStore.getUserInfo();
  if (!success) {
    router.push("/login");
  }
}
checkUserInfo();
</script>

<style scoped lang="scss">
.wb-workbench {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background:
    radial-gradient(1200px 500px at 70% -10%, rgba(94, 234, 212, 0.06), transparent 60%),
    radial-gradient(900px 500px at 0% 110%, rgba(96, 165, 250, 0.05), transparent 60%),
    var(--wb-bg);
  color: var(--wb-text);
  font-family: var(--wb-font-body);
  overflow: hidden;
}

.wb-body {
  display: flex;
  flex: 1;
  min-height: 0;
}

.wb-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: var(--wb-bg);
  background-image:
    linear-gradient(rgba(94, 234, 212, 0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(94, 234, 212, 0.025) 1px, transparent 1px);
  background-size: 28px 28px;
}
</style>
