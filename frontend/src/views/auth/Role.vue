<script setup lang="tsx">
import { getPermissionTree, type Permission } from '@/api/permission';
import { createRole, deleteRole, getRole, getRoles, updateRole, type Role } from '@/api/role';
import Form from '@/components/common/form/Form.vue';
import FormDialog from '@/components/common/form/FormDialog.vue';
import FormItem from '@/components/common/form/FormItem.vue';
import TableLayout from '@/components/common/TableLayout.vue';
import Tree from '@/components/common/Tree.vue';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { showConfirm, showToast } from '@/lib/message';
import type { FormRules } from '@/lib/symbols';
import { TableColumn } from '@/types/column';
import { onMounted, reactive, ref, computed } from 'vue';

const items = ref<Role[]>([])
const tableLayoutRef = ref()
const permissionTree = ref<Permission[]>([])

const fetchData = async (params: any) => {
  return (await getRoles(params)) as any
}

onMounted(async () => {
  try {
    permissionTree.value = (await getPermissionTree()) as any
  } catch (error) {
    console.error('Failed to load permission tree:', error)
  }
})

const refreshTable = () => {
  if (tableLayoutRef.value) {
    tableLayoutRef.value.fetchData()
  }
}

const isDialogOpen = ref(false)
const editingRole = ref<Role | null>(null)

const form = reactive({
  name: '',
  description: '',
  permissionIds: [] as (string | number)[]
})

const getAllPermissionIds = (nodes: Permission[]): (string | number)[] => {
  let ids: (string | number)[] = []
  nodes.forEach(node => {
    ids.push(node.id)
    if (node.children && node.children.length > 0) {
      ids.push(...getAllPermissionIds(node.children))
    }
  })
  return ids
}

const allPermissionIds = computed(() => getAllPermissionIds(permissionTree.value))

const isAllSelected = computed(() => {
  if (allPermissionIds.value.length === 0) return false
  return allPermissionIds.value.every(id => form.permissionIds.includes(id))
})

const handleToggleAll = () => {
  if (isAllSelected.value) {
    form.permissionIds = []
  } else {
    form.permissionIds = [...allPermissionIds.value]
  }
}

const rules: FormRules = {
  name: [{ required: true, message: '请输入角色名称' }]
}

const formRef = ref()

const resetForm = () => {
  form.name = ''
  form.description = ''
  form.permissionIds = []
  editingRole.value = null
}

const handleCreate = () => {
  resetForm()
  isDialogOpen.value = true
}

const handleUpdate = async (row: Role) => {
  editingRole.value = row
  form.name = row.name
  form.description = row.description || ''
  form.permissionIds = []

  try {
    const res = await getRole(row.id) as any
    if (res && res.permissionIds) {
      form.permissionIds = res.permissionIds
    }
  } catch (error) {
    console.error('Failed to fetch role details:', error)
  }

  isDialogOpen.value = true
}

const handleDelete = async (row: any) => {
  if (await showConfirm('确认删除', '确定要删除该角色吗？')) {
    await deleteRole(row.id)
    showToast('删除成功')
    refreshTable()
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate()
  if (!valid) return

  try {
    if (editingRole.value) {
      await updateRole(editingRole.value.id, form)
      showToast('更新成功')
    } else {
      await createRole(form)
      showToast('创建成功')
    }
    isDialogOpen.value = false
    refreshTable()
  } catch (error: any) {
    console.error(error)
  }
}

const shouldExpandNode = (node: any) => {
  if (node.children && node.children.length > 0) {
    const hasOnlyButtons = node.children.every((c: any) => c.type === 2)
    return !hasOnlyButtons
  }
  return false
}

const columns: TableColumn[] = [
  { label: '角色名称', field: 'name', sortable: true },
  { label: '描述', field: 'description' },
  { label: '创建时间', field: 'createdAt', sortable: true, width: '180px' },
  {
    label: '操作',
    field: 'operation',
    width: '180px',
  }
]
</script>

<template>
  <div>
    <TableLayout ref="tableLayoutRef" title="角色" :fetch-data="fetchData" :columns="columns" @delete="handleDelete"
      @update="handleUpdate" @create="handleCreate">
    </TableLayout>

    <FormDialog :visible="isDialogOpen" @update:visible="isDialogOpen = $event" :title="editingRole ? '编辑角色' : '新增角色'"
      @confirm="handleSave">
      <div class="py-2">
        <Form ref="formRef" :model="form" :rules="rules">
          <div class="grid grid-cols-2 gap-4">
            <FormItem label="角色名称" field="name" required>
              <Input v-model="form.name" placeholder="请输入角色名称" />
            </FormItem>
            <FormItem label="描述" field="description">
              <Input v-model="form.description" placeholder="请输入描述" />
            </FormItem>
          </div>
          <div>
            <FormItem label="权限配置" field="permissionIds">
              <div class="space-y-2">
                <div class="border rounded-md p-2 max-h-[300px] overflow-y-auto">
                  <div v-if="permissionTree.length === 0" class="text-sm text-muted-foreground text-center py-4">
                    暂无权限数据
                  </div>
                  <Tree v-else :items="permissionTree" v-model="form.permissionIds" :default-expand="shouldExpandNode" />
                </div>
                <div class="flex items-center justify-between bg-muted/50 p-3 rounded-md border text-sm">
                  <div class="text-muted-foreground">
                    已选择 <span class="font-medium text-primary mx-1">{{ form.permissionIds.length }}</span> 项权限
                  </div>
                  <Button type="button" variant="outline" size="sm" class="h-8" @click="handleToggleAll">
                    {{ isAllSelected ? '取消全选' : '全选' }}
                  </Button>
                </div>
              </div>
            </FormItem>
          </div>
        </Form>
      </div>
    </FormDialog>
  </div>
</template>
