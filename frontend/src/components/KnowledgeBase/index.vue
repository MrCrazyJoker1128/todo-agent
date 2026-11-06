<template>
  <div class="wb-scope wb-kb">
    <div class="wb-kb-head">
      <span class="wb-mono wb-eyebrow">KNOWLEDGE</span>
      <h2 class="wb-font-display">知识库</h2>
    </div>
    <div class="wb-kb-body">
      <el-card class="wb-kb-libs">
        <template #header>
          <div class="wb-kb-libs-head">
            <span>库列表</span>
            <el-button size="small" type="primary" @click="showCreate = true"
              >新建</el-button
            >
          </div>
        </template>
        <div v-loading="libsLoading" class="wb-kb-libs-list">
          <div
            v-for="lib in libs"
            :key="lib.id"
            :class="['wb-kb-lib', { active: lib.id === activeId }]"
            @click="selectLib(lib.id)"
          >
            <span class="wb-kb-lib-name">{{ lib.name }}</span>
            <el-button
              size="small"
              type="danger"
              text
              @click.stop="confirmDelete(lib)"
              >删除</el-button
            >
          </div>
          <el-empty v-if="!libs.length" description="暂无知识库" />
        </div>
      </el-card>
      <el-card v-if="activeId" v-loading="filesLoading" class="wb-kb-files">
        <template #header>
          <div class="wb-kb-libs-head">
            <span>文件</span>
            <el-button size="small" type="primary" @click="showAdd = true"
              >添加文件</el-button
            >
          </div>
        </template>
        <el-form :inline="true" class="wb-kb-searchrow" @submit.prevent>
          <el-form-item label="标题">
            <el-input
              v-model="keyword"
              placeholder="标题关键字"
              clearable
              style="width: 200px"
              @input="resetPage"
            />
          </el-form-item>
        </el-form>
        <el-table :data="pageFiles" style="width: 100%">
          <el-table-column prop="title" label="标题" show-overflow-tooltip />
          <el-table-column prop="source" label="来源" show-overflow-tooltip />
          <el-table-column prop="chunkCount" label="分片" width="70" />
          <el-table-column label="状态" width="90">
            <template #default="scope">
              <el-tag
                v-if="!scope.row.status || scope.row.status === 'completed'"
                size="small"
                type="success"
                >已完成</el-tag
              >
              <el-tag
                v-else-if="scope.row.status === 'pending'"
                size="small"
                type="warning"
                >处理中</el-tag
              >
              <el-tag
                v-else-if="scope.row.status === 'failed'"
                size="small"
                type="danger"
                :title="scope.row.error"
                >失败</el-tag
              >
              <el-tag v-else size="small" type="info">{{
                scope.row.status
              }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="scope">
              <el-button
                size="small"
                type="danger"
                @click="handleDeleteFile(scope.row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="filteredFiles.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          class="wb-kb-pager"
        />
      </el-card>
      <el-empty v-else description="请选择或新建知识库" />
    </div>
    <el-dialog v-model="showCreate" title="新建知识库" width="400">
      <el-form label-width="70px">
        <el-form-item label="库名">
          <el-input v-model="createForm.name" placeholder="库名" clearable />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="createForm.description"
            placeholder="描述（可选）"
            clearable
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="handleCreate"
          >确定</el-button
        >
      </template>
    </el-dialog>
    <el-dialog v-model="showAdd" title="添加文件" width="440">
      <el-radio-group v-model="addTab" class="wb-kb-addtype">
        <el-radio-button value="url">输入地址</el-radio-button>
        <el-radio-button value="md">md 文档</el-radio-button>
      </el-radio-group>
      <div v-if="addTab === 'url'" class="wb-kb-addbody">
        <el-input
          v-model="urlForm.url"
          placeholder="输入 html 页面地址"
          clearable
        />
      </div>
      <div v-else class="wb-kb-addbody">
        <div
          class="wb-kb-drop"
          :class="{ 'is-dragover': dragOver }"
          @click="openMdPicker"
          @dragover.prevent="dragOver = true"
          @dragleave="dragOver = false"
          @drop.stop.prevent="handleMdDrop"
        >
          <input
            id="mdFileInput"
            type="file"
            accept=".md,.markdown"
            hidden
            @change="handleMdInput"
          />
          <el-icon class="wb-kb-drop-icon"><Plus /></el-icon>
          <div>拖入或<em>点击上传</em> md 文件（≤2MB）</div>
        </div>
        <div v-if="mdName" class="wb-kb-mdname">{{ mdName }}</div>
      </div>
      <template #footer>
        <el-button @click="showAdd = false">取消</el-button>
        <el-button
          v-if="addTab === 'url'"
          type="primary"
          :loading="addUrlLoading"
          @click="handleAddUrl"
          >添加</el-button
        >
        <el-button
          v-else
          type="primary"
          :loading="uploadLoading"
          :disabled="!mdName"
          @click="handleUploadMd"
          >上传</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus } from "@element-plus/icons-vue";
import { useKnowledgeBase } from "./hook";

defineOptions({ name: "KnowledgeBase" });

const {
  libs,
  libsLoading,
  activeId,
  files,
  filesLoading,
  keyword,
  filteredFiles,
  page,
  pageSize,
  pageFiles,
  showCreate,
  showAdd,
  createForm,
  createLoading,
  addTab,
  urlForm,
  addUrlLoading,
  mdName,
  dragOver,
  uploadLoading,
  loadLibs,
  selectLib,
  handleCreate,
  confirmDelete,
  handleAddUrl,
  openMdPicker,
  handleMdInput,
  handleMdDrop,
  handleUploadMd,
  handleDeleteFile,
  stopPoll,
  resetPage,
} = useKnowledgeBase();

loadLibs();
void files;

onUnmounted(() => {
  stopPoll();
});
</script>
<style scoped lang="scss">
.wb-kb {
  padding: 16px 20px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}
.wb-kb-head {
  margin-bottom: 16px;
}
.wb-eyebrow {
  font-size: 10px;
  letter-spacing: 0.22em;
  color: var(--wb-text-3);
}
.wb-kb-head h2 {
  margin: 4px 0 0;
  font-size: 24px;
  font-weight: 700;
}
.wb-kb-body {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}
.wb-kb-libs {
  width: 260px;
  flex-shrink: 0;
}
.wb-kb-libs-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.wb-kb-libs-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.wb-kb-lib {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  &:hover {
    background: var(--wb-panel-2);
  }
  &.active {
    background: var(--wb-accent-dim);
  }
}
.wb-kb-lib-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.wb-kb-files {
  flex: 1;
  min-width: 0;
}
.wb-kb-searchrow {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}
 .wb-kb-mdname {
   margin-top: 8px;
   display: flex;
   gap: 8px;
   align-items: center;
 }
 .wb-kb-addtype {
   margin-bottom: 12px;
 }
 .wb-kb-addbody {
   margin-bottom: 4px;
 }
 .wb-kb-drop {
   display: flex;
   flex-direction: column;
   align-items: center;
   justify-content: center;
   gap: 8px;
   min-height: 140px;
   padding: 20px;
   border: 1px dashed var(--wb-border);
   border-radius: 8px;
   background: var(--wb-panel-2);
   color: var(--wb-text-3);
   cursor: pointer;
   transition: border-color 0.15s, background 0.15s;
   &:hover, &.is-dragover {
     border-color: var(--wb-accent);
     background: var(--wb-accent-dim);
   }
   em {
     font-style: normal;
     color: var(--wb-accent);
   }
 }
 .wb-kb-drop-icon {
   font-size: 24px;
 }
 .wb-kb-pager {
   margin-top: 12px;
   justify-content: flex-end;
 }
 </style>
