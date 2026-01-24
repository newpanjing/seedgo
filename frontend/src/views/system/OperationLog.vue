<script setup lang="tsx">
import { ref, computed, reactive, onMounted, onBeforeUnmount } from 'vue'
import { TableColumn } from '@/types/column'
import { getOperationLogs, type OperationLog } from '@/api/operation-log'
import { Badge } from '@/components/ui/badge'
import { Eye } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import FormSheet from '@/components/common/form/FormSheet.vue'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import DateRangeSelect from '@/components/common/DateRangeSelect.vue'

import TableCellFormat from '@/components/common/TableCellFormat.vue'

const isMounted = ref(false)
onMounted(() => {
  isMounted.value = true
})

// Filters
const status = ref<string>('')
const dateRange = ref<[string | undefined, string | undefined]>([undefined, undefined])

const columns: TableColumn[] = [
  {
    label: '模块/操作',
    field: 'module',
    formatter: (_: any, row: OperationLog) => (
      <div class="flex flex-col gap-1">
        <div class="font-medium">{row.method}</div>
        <div class="text-xs text-muted-foreground truncate max-w-[200px]" title={row.path}>{row.path}</div>
      </div>
    )
  },
  { label: '用户', field: 'username' },
  { label: 'IP', field: 'ip' },
  {
    label: '状态',
    field: 'status',
    formatter: (val: number) => (
      <Badge variant="outline" class={val >= 200 && val < 300 ? 'text-green-600 border-green-200 bg-green-50' : 'text-red-600 border-red-200 bg-red-50'}>
        {val}
      </Badge>
    )
  },
  {
    label: '耗时',
    field: 'latency',
    formatter: (val: number) => <span class={val > 1000 ? 'text-orange-500' : ''}>{val}ms</span>
  },
  { 
    label: '操作时间', 
    field: 'operationTime', 
    sortable: true,
    formatter: (val: any) => <TableCellFormat value={val} />
  }
]

const tableLayoutRef = ref()
const refreshTable = () => {
  tableLayoutRef.value?.fetchData()
}

const fetchData = async (params: any) => {
  const res = (await getOperationLogs({
    ...params,
    status: status.value ? Number(status.value) : undefined,
    startDate: dateRange.value?.[0],
    endDate: dateRange.value?.[1],
  })) as any
  return {
    total: res.total,
    items: res.items
  }
}

// Detail View
const isViewOpen = ref(false)
const viewedLog = ref<OperationLog | null>(null)

const handleView = (log: OperationLog) => {
  viewedLog.value = log
  isViewOpen.value = true
}

// Helper to format JSON or large text
const formatJson = (str: string) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}
</script>

<template>
  <div>
    <TableLayout
      ref="tableLayoutRef"
      title="操作日志"
      :columns="columns"
      :fetch-data="fetchData"
      :show-create="false"
      :show-update="false"
      :show-delete="false"
      :checkable="false"
    >
      <template #filters>
        <select 
            v-model="status" 
            @change="refreshTable"
            class="h-8 w-[120px] appearance-none rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
        >
            <option value="" disabled selected>状态</option>
            <option value="">全部</option>
            <option value="200">成功</option>
            <option value="500">失败</option>
        </select>
        <DateRangeSelect v-model="dateRange" @update:model-value="refreshTable" class="h-8" />
      </template>

      <template #actions="{ Row }">
        <Button variant="ghost" size="icon" class="h-8 w-8" @click="handleView(Row)" title="查看详情">
          <Eye class="w-4 h-4" />
        </Button>
      </template>
    </TableLayout>

    <FormSheet :visible="isViewOpen" @update:visible="isViewOpen = $event" title="日志详情">
      <div v-if="viewedLog" class="space-y-6">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <Label class="text-muted-foreground">操作用户</Label>
            <div class="mt-1 font-medium">{{ viewedLog.username || '-' }} (ID: {{ viewedLog.userId }})</div>
          </div>
          <div>
            <Label class="text-muted-foreground">客户端IP</Label>
            <div class="mt-1 font-medium">{{ viewedLog.ip }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">请求方式</Label>
            <div class="mt-1 font-medium">{{ viewedLog.method }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">状态码</Label>
            <div class="mt-1">
              <Badge variant="outline" :class="viewedLog.status >= 200 && viewedLog.status < 300 ? 'text-green-600 border-green-200 bg-green-50' : 'text-red-600 border-red-200 bg-red-50'">
                {{ viewedLog.status }}
              </Badge>
            </div>
          </div>
          <div class="col-span-2">
            <Label class="text-muted-foreground">请求路径</Label>
            <div class="mt-1 font-mono text-sm break-all bg-muted p-2 rounded">{{ viewedLog.path }}</div>
          </div>
           <div class="col-span-2">
            <Label class="text-muted-foreground">User Agent</Label>
            <div class="mt-1 text-sm text-muted-foreground break-all">{{ viewedLog.userAgent }}</div>
          </div>
          <div class="col-span-2" v-if="viewedLog.query">
            <Label class="text-muted-foreground">Query 参数</Label>
            <pre class="mt-1 font-mono text-xs bg-muted p-2 rounded overflow-x-auto">{{ viewedLog.query }}</pre>
          </div>
          <div class="col-span-2" v-if="viewedLog.body">
            <Label class="text-muted-foreground">请求体</Label>
            <pre class="mt-1 font-mono text-xs bg-muted p-2 rounded overflow-x-auto whitespace-pre-wrap">{{ formatJson(viewedLog.body) }}</pre>
          </div>
          <div class="col-span-2" v-if="viewedLog.errorMessage">
            <Label class="text-muted-foreground text-red-500">错误信息</Label>
            <pre class="mt-1 font-mono text-xs bg-red-50 text-red-600 p-2 rounded overflow-x-auto whitespace-pre-wrap">{{ viewedLog.errorMessage }}</pre>
          </div>
        </div>
      </div>
    </FormSheet>
  </div>
</template>

