<template>
  <div class="row-between-center w-full">
    <el-form inline :model="searchReq" class="demo-form-inline">
      <el-form-item label="标题" class="w-[200px]">
        <el-input v-model="searchReq.title" placeholder="标题" clearable />
      </el-form-item>
      <el-form-item label="作者" class="w-[200px]">
        <el-input v-model="searchReq.author" placeholder="作者" clearable />
      </el-form-item>
      <el-form-item label="域名" class="w-[200px]">
        <el-input v-model="searchReq.domain" placeholder="域名" clearable />
      </el-form-item>
      <el-form-item label="状态" class="w-[200px]">
        <el-select v-model="searchReq.status" placeholder="状态" clearable>
          <el-option
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
    <el-form :inline="true">
      <el-form-item>
        <el-input
          v-model="knowledgeUrlObj.url"
          placeholder="请输入知识库URL"
          :disabled="loadingOption.loading"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          @click="handleAdd"
          :loading="loadingOption.loading"
          >添加</el-button
        >
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import type { GetKnowledgeUrlsSearchOptions } from "@/api/knowledge/type";
import { useCreateKnowledgeUrlHooks } from "./hooks";
import { useReactiveReset } from "@/utils/hooks";
import eventBus from "@/utils/eventbus";
import { useStatusOptions } from "../../hooks";
defineOptions({
  name: "SearchArea",
});
const { loadingOption, knowledgeUrlObj, handleAdd } =
  useCreateKnowledgeUrlHooks();

const emit = defineEmits(["search"]);
const { state: searchReq, resetState: resetSearchReq } =
  useReactiveReset<GetKnowledgeUrlsSearchOptions>({
    title: "",
    author: "",
    domain: "",
    status: "",
  });
const { statusOptions } = useStatusOptions();
function handleReset() {
  resetSearchReq();
  emit("search", searchReq);
}
function handleSearch() {
  emit("search", searchReq);
}

eventBus.on("addComplete", () => {
  console.log("addComplete");
  handleSearch();
});
defineExpose({
  searchReq,
});
</script>

<style lang="scss" scoped></style>
