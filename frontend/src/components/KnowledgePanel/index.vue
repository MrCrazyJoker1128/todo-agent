<template>
  <div class="w-full h-full col-start-start">
    <SearchArea @search="initTableList" />
    <div class="flex-1 w-full">
      <TableComp
        :tableHeadConfig="tableHeadConfig"
        :tableLoadData="tableList"
        :loading="pageOption.loading"
        v-model:currentPage="pageOption.current_page"
        v-model:pageSize="pageOption.per_page"
        :total="pageOption.total"
        :getTableData="getKnowledgeUrlsList"
      >
        <template #status>
          <el-table-column prop="status" label="状态" show-overflow-tooltip>
            <template #default="scope">
              {{ getStatusOptionsLabel(scope.row.status) }}
            </template>
          </el-table-column>
        </template>
        <template #url>
          <el-table-column prop="url" label="地址" show-overflow-tooltip>
            <template #default="scope">
              <a :href="scope.row.url" target="_blank">{{ scope.row.url }}</a>
            </template>
          </el-table-column>
        </template>
        <template #operation>
          <el-table-column prop="operation" label="操作" width="80">
            <template #default="scope">
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(scope.row.id)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </template>
      </TableComp>
    </div>
  </div>
</template>

<script setup lang="ts">
import { fetchDeleteKnowledgeUrl } from "@/api/knowledge";
import SearchArea from "./components/SearchArea/index.vue";
import { useGetKnowledgeListHooks, useStatusOptions } from "./hooks";
import { ElMessage, ElMessageBox } from "element-plus";
defineOptions({
  name: "KnowledgePanel",
  components: {
    SearchArea,
  },
});

const { getStatusOptionsLabel } = useStatusOptions();

const {
  initTableList,
  getKnowledgeUrlsList,
  tableHeadConfig,
  pageOption,
  tableList,
} = useGetKnowledgeListHooks();
onMounted(() => {
  initTableList();
});
function handleDelete(id: string) {
  ElMessageBox.confirm("确认删除吗？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    fetchDeleteKnowledgeUrl({ id }).then((res) => {
      if (res.data.code === 200) {
        ElMessage.success(res.data.message);
        initTableList();
      }
    });
  });
}
</script>

<style lang="scss" scoped></style>
