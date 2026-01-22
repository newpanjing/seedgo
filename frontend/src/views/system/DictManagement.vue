<script setup lang="tsx">
import { ref, reactive, computed } from 'vue'
import { TableColumn } from '@/types/column'
import { getDicts, createDict, updateDict, deleteDict, batchDeleteDicts, getDict, type Dict, type DictItem } from '@/api/dict'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { showConfirm, showToast } from '@/lib/message'
import Form from '@/components/common/form/Form.vue'
import FormItem from '@/components/common/form/FormItem.vue'
import FormDialog from '@/components/common/form/FormDialog.vue'
import TableLayout from '@/components/common/TableLayout.vue'
import type { FormRules } from '@/lib/symbols'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useDark } from '@vueuse/core'

const isDark = useDark()

// Table columns definition
const columns = reactive<TableColumn[]>([
  { label: '编码', field: 'code', sortable: true },
  { label: '名称', field: 'name' },
  { 
    label: '字典项数量', 
    field: 'items', 
    formatter: (value: any, row: any) => {
      return <span>{row.items?.length || 0}</span>
    }
  },
  { label: '描述', field: 'description' },
  { label: '创建时间', field: 'createdAt', sortable: true }
])

const items = ref<Dict[]>([])
const tableLayoutRef = ref()
const isDialogOpen = ref(false)
const editingDict = ref<Dict | null>(null)

// Form state
const form = reactive({
  code: '',
  name: '',
  description: '',
  items: [] as DictItem[]
})

const rules = computed<FormRules>(() => ({
  code: [{ required: true, message: '请输入编码' }],
  name: [{ required: true, message: '请输入名称' }]
}))

const formRef = ref()

const resetForm = () => {
  form.code = ''
  form.name = ''
  form.description = ''
  form.items = [{ id: undefined, label: '', value: '', sort: 1, status: 1 }]
  editingDict.value = null
}

const refreshTable = () => {
  if (tableLayoutRef.value) {
    tableLayoutRef.value.fetchData()
  }
}

const fetchData = async (params: any) => {
  const res = (await getDicts({
    ...params,
    keyword: params.query
  })) as any
  items.value = res.items || []
  return {
    total: res.total,
    items: res.items
  }
}

const handleCreate = () => {
  resetForm()
  isDialogOpen.value = true
}

const handleUpdate = async (dict: Dict) => {
  // 1. 先回显基本信息并打开弹窗
  editingDict.value = dict
  form.code = dict.code
  form.name = dict.name
  form.description = dict.description || ''
  // 列表数据中可能没有 items，先置空或使用列表中的数据（如果列表已包含部分 items）
  form.items = (dict.items || []).map(i => ({ ...i }))
  isDialogOpen.value = true

  // 2. 异步拉取完整详情
  try {
    const res = (await getDict(dict.id)) as unknown as Dict
    // 再次确认当前编辑的 ID 是否一致（防止快速切换导致的数据错乱）
    if (editingDict.value?.id === res.id) {
      form.items = (res.items || []).map(i => ({ ...i }))
      // 也可以更新其他可能在列表中不完整的字段
      form.description = res.description || ''
    }
  } catch (e) {
    console.error(e)
    showToast('获取详情失败', { type: 'error' })
  }
}

const handleDelete = async (dict: Dict) => {
  const confirmed = await showConfirm('确认删除', `确定要删除字典 "${dict.name}" 吗？此操作不可恢复。`)
  if (!confirmed) return

  await deleteDict(dict.id)
  showToast('删除成功')
  refreshTable()
}

const handleBatchDelete = async (dicts: Dict[]) => {
  try {
    await batchDeleteDicts(dicts.map(d => d.id))
    showToast(`成功删除 ${dicts.length} 条记录`)
    refreshTable()
  } catch (error) {
    console.error(error)
    showToast('部分删除失败，请重试', { type: 'error' })
  }
}

const handleSave = async () => {
  try {
    if (formRef.value) {
      const valid = await formRef.value.validate()
      if (!valid) return
    }

    const data = {
      code: form.code,
      name: form.name,
      description: form.description,
      items: form.items
    }

    if (editingDict.value) {
      await updateDict(editingDict.value.id, data)
    } else {
      await createDict(data)
    }
    
    isDialogOpen.value = false
    refreshTable()
    showToast('保存成功')
  } catch (e) {
    console.error(e)
  }
}

// Inner table logic
const addItemRow = () => {
  const nextSort = form.items.length ? Math.max(...form.items.map(i => i.sort)) + 1 : 1
  form.items.push({ label: '', value: '', sort: nextSort, status: 1 })
}

const removeItemRow = (idx: number) => {
  form.items.splice(idx, 1)
}
</script>

<template>
  <div>
    <TableLayout 
      ref="tableLayoutRef" 
      title="字典" 
      :fetch-data="fetchData" 
      :columns="columns" 
      :items="items"
      :checkable="true" 
      @delete="handleDelete" 
      @batch-delete="handleBatchDelete" 
      @update="handleUpdate" 
      @create="handleCreate"
    >
    </TableLayout>

    <FormDialog 
      :visible="isDialogOpen" 
      @update:visible="isDialogOpen = $event" 
      :title="editingDict ? '编辑字典' : '新增字典'"
      width="800px"
      @confirm="handleSave"
    >
      <div class="py-2">
        <Form ref="formRef" :model="form" :rules="rules">
          <div class="grid grid-cols-2 gap-4">
            <FormItem label="编码" field="code" required>
              <Input v-model="form.code" placeholder="如 room_type" />
            </FormItem>
            <FormItem label="名称" field="name" required>
              <Input v-model="form.name" placeholder="如 房间类型" />
            </FormItem>
            <FormItem label="描述" field="description" class="col-span-2">
              <Input v-model="form.description" placeholder="字典描述" />
            </FormItem>
          </div>

          <div class="mt-6 space-y-4">
            <div class="flex items-center justify-between">
              <Label>字典项</Label>
              <Button variant="outline" size="sm" type="button" @click="addItemRow">新增字典项</Button>
            </div>
            <div class="rounded-lg border border-border overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>显示名称</TableHead>
                    <TableHead>实际值</TableHead>
                    <TableHead class="w-24">排序</TableHead>
                    <TableHead class="w-32">状态</TableHead>
                    <TableHead class="text-right w-20">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="(it, idx) in form.items" :key="idx">
                    <TableCell>
                      <Input v-model="it.label" placeholder="如：棋牌室" />
                    </TableCell>
                    <TableCell>
                      <Input v-model="it.value" placeholder="如：chess" />
                    </TableCell>
                    <TableCell>
                      <Input v-model="it.sort" type="number" />
                    </TableCell>
                    <TableCell>
                      <div class="flex items-center gap-2">
                        <input
                          type="checkbox"
                          class="toggle toggle-success toggle-sm"
                          :checked="it.status === 1"
                          @change="it.status = ($event.target as HTMLInputElement).checked ? 1 : 0"
                        />
                        <Badge :color="it.status === 1 ? 'green' : 'red'">
                          {{ it.status === 1 ? '启用' : '停用' }}
                        </Badge>
                      </div>
                    </TableCell>
                    <TableCell class="text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        class="text-xs h-8 w-8 p-0 text-destructive hover:text-destructive"
                        type="button"
                        @click="removeItemRow(idx)"
                      >
                        删除
                      </Button>
                    </TableCell>
                  </TableRow>
                  <TableRow v-if="form.items.length === 0">
                    <TableCell colspan="5" class="text-center text-xs text-muted-foreground h-24">暂无字典项</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </div>
        </Form>
      </div>
    </FormDialog>

  </div>
</template>
