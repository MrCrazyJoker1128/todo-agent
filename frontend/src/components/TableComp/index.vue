<template>
  <el-table
    v-loading="props.loading"
    stripe
    :border="true"
    :class="[hasSorted ? 'max-table has-sorted' : 'max-table']"
    ref="maxTable"
    :data="tableData"
    style="width: 100%"
    force-scroll="vertical"
    :height="tableHeight"
    tooltip-effect="light"
    :highlight-current-row="highlightCurrentRow"
    :default-sort="defaultSort"
    :lazy="props.lazy"
    :load="loadFunc"
    :row-key="props.rowKey"
    @selection-change="selectData"
    @cell-click="(row: T) => cellDblclickData(row)"
    @sort-change="sortChange"
    @expand-change="expandChange"
  >
    <!-- 多选框 -->
    <el-table-column
      :fixed="props.fixedLift"
      v-if="selection"
      :selectable="isCanSelect"
      type="selection"
      width="38"
    ></el-table-column>
    <!-- 表单嵌套 -->
    <el-table-column v-if="props.expand" type="expand">
      <template #default="scope">
        <slot name="expand" :row="scope.row" :index="scope.$index"></slot>
      </template>
    </el-table-column>
    <!-- 文本数据渲染 -->
    <template v-for="(item, index) in tableHead">
      <slot v-if="item.slot" :name="item.slot"></slot>
      <!-- show-overflow-tooltip -->
      <el-table-column
        v-else
        :show-overflow-tooltip="item.minWidth ? true : false"
        :class-name="item.className"
        :prop="item.prop"
        :label="item.label"
        :key="index + 1"
        :sortable="item.isSorted ? 'custom' : false"
        :width="item.width"
        :min-width="item.minWidth"
        :align="item.align ? item.align : align"
        :fixed="item.fixed ? item.fixed : false"
      >
      </el-table-column>
    </template>
    <template slot="empty">
      <el-empty
        :image-size="props.emptySize"
        :description="props.emptyText"
      ></el-empty>
    </template>
  </el-table>
  <div ref="pagenationRef" v-show="props.isPaging" class="row-end-center">
    <el-pagination
      :hide-on-single-page="false"
      background
      :layout="props.layout"
      :total="props.total"
      v-model:current-page="current_page"
      v-model:page-size="per_page"
      :page-sizes="[10, 20, 30]"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts" generic="T extends DefaultRow">
import type { TableInstance, TableColumnCtx, Sort } from "element-plus";
import { type PropType } from "vue";
import { useDebounceFn } from "@vueuse/core";
import type { TableColumn } from "@/types/type";
import type { DefaultRow } from "element-plus/lib/components/table/src/table/defaults.js";

defineOptions({
  name: "TableComp",
});

const props = defineProps({
  //表头
  tableHeadConfig: {
    type: Array as PropType<TableColumn[]>,
    required: true,
    default: () => [],
  },
  //数据
  tableLoadData: {
    type: Array as PropType<T[]>,
    required: true,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  //表格高度
  tableHeight: {
    type: [Number, String],
    default: () => {
      return "100%";
    },
  },
  //对齐
  align: {
    type: String,
    default: "left",
  },
  //多选
  selection: {
    type: Boolean,
    default: false,
  },
  //固定列
  fixedLift: {
    type: Boolean,
    default: false,
  },
  //点击行索引
  rowIndex: {
    type: Number,
    default: 0,
  },
  //选择单行数据时使用色块表示
  highlightCurrentRow: {
    type: Boolean,
    default: true,
  },
  defaultSort: {
    type: Object as PropType<Sort>,
    defalut() {
      return {};
    },
  },
  //单元格文字样式
  cellStyle: {
    type: Array,
    default: () => {
      return [
        {
          name: "",
          color: "",
        },
      ];
    },
  },
  //单元格回调方法
  cellClickFn: {
    type: Array,
    default: () => {
      return [
        {
          name: "",
          columnCall: "",
        },
      ];
    },
  },
  //表单嵌套
  expand: {
    type: Boolean,
    default: false,
  },
  emptySize: {
    type: Number,
    default: 100,
  },
  emptyText: {
    type: String,
    default: "暂无数据",
  },
  //行样式
  rowClassName: {
    type: Array,
    default: () => {
      return [
        {
          index: "",
          rowClass: "",
        },
      ];
    },
  },
  isCanSelect: {
    type: Function as PropType<(row: T) => boolean>,
    default: () => {
      return true;
    },
  },
  layout: {
    type: String,
    default: "total,sizes,prev, pager, next,jumper",
  },
  total: {
    type: Number,
    default: 0,
  },
  getTableData: {
    type: Function,
    default: () => {
      return () => {};
    },
  },
  //是否需要分页
  isPaging: {
    type: Boolean,
    default: true,
  },
  lazy: {
    type: Boolean,
    default: false,
  },
  rowKey: {
    type: String,
    default: "id",
  },
});
/**
 * 定义一个名为 current_page 的模型
 * 该模型用于表示页面索引，其类型为数字，默认值为 1
 */
const current_page = defineModel("currentPage", {
  type: Number,
  default: 1,
});
const per_page = defineModel("perPage", {
  type: Number,
  default: 10,
});
// 使用useDebounceFn自定义钩子来创建一个防抖函数
// 这里的作用是防止getTableData方法在短时间内被频繁调用
// 参数是一个箭头函数，它将getTableData作为实际执行的操作
// 500是延迟执行的时间，单位为毫秒
const debounced = useDebounceFn(() => {
  props.getTableData();
}, 500);
/**
 * 处理每页条数变化的函数
 * 当用户改变每页显示的数据条数时调用此函数
 * @param value {number} - 用户选择的每页显示的条数
 */
function handleSizeChange(value: number): void {
  current_page.value = 1;
  per_page.value = value;
  debounced();
}
/**
 * 处理当前页码变化的函数
 * 当用户点击分页导航时调用此函数
 * @param value {number} - 用户选择的新页码
 */
function handleCurrentChange(value: number): void {
  current_page.value = value;
  debounced();
}

const emit = defineEmits([
  "rowClick",
  "selectData",
  "sortChange",
  "load",
  "expandChange",
]);
const maxTable = ref<TableInstance>();

const tableData = computed(() => {
  return props.tableLoadData;
});
const tableHead = computed((): TableColumn[] => {
  return props.tableHeadConfig;
});
const tableIndex = computed((): number => {
  return props.rowIndex;
});
const hasSorted = computed(() => {
  return !props.tableHeadConfig.some((item) => {
    return item?.isSorted;
  });
});

const pageHeight = ref<number>(32);
onUpdated(() => {
  pageHeight.value = pagenationRef?.value?.offsetHeight ?? 0;
});
const pagenationRef = ref<HTMLDivElement>();
const tableHeight = computed((): string | undefined => {
  if (props.isPaging) {
    return typeof props.tableHeight === "number"
      ? `${props.tableHeight - pageHeight.value}px`
      : `calc(${props.tableHeight} - ${pageHeight.value}px)`;
  } else {
    return typeof props.tableHeight === "number"
      ? `${props.tableHeight}px`
      : props.tableHeight;
  }
});
//单元格点击方法
function cellDblclickData<T>(row: T): void {
  emit("rowClick", row);
}
//获取选中数据
function selectData<T>(selection: T[]) {
  emit("selectData", selection);
}
function sortChange<T>(column: TableColumnCtx<T & DefaultRow>) {
  emit("sortChange", column);
}
type ExpandChange = ((row: T, expandedRows: T[]) => void) &
  ((row: T, expanded: boolean) => void);
const expandChange: ExpandChange = (row: T, expanded: T[] | boolean) => {
  emit("expandChange", row, expanded);
};
//单选行数据显示色块跳转不同行的方法
function setCurrent(index: number) {
  maxTable.value?.setCurrentRow(props.tableLoadData[index]);
}
watch(
  () => tableIndex,
  (newValue) => {
    setCurrent(newValue.value);
  },
);
watch(
  () => tableData,
  () => {
    nextTick(() => {
      maxTable.value?.doLayout();
    });
  },
);

function loadFunc(row: T, _treeNode: unknown, resolve: (data: T[]) => void) {
  void _treeNode;
  emit("load", row, resolve);
}
</script>

<style scoped></style>
