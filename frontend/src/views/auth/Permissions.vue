<script setup lang="tsx">
import { ref, reactive, computed } from 'vue';
import * as icons from 'lucide-vue-next';
import { TableColumn } from '@/types/column';
import { getMenus, createMenu, updateMenu, deleteMenu, batchDeleteMenus, type Menu } from '@/api/menu';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { showToast, showConfirm } from '@/lib/message';
import Form from '@/components/common/form/Form.vue';
import FormItem from '@/components/common/form/FormItem.vue';
import type { FormRules } from '@/lib/symbols';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import IconPicker from '@/components/common/IconPicker.vue';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const items = ref<Menu[]>([])
const tableLayoutRef = ref()
// 用于父级菜单选择的列表 (Flattened for Select)
const menuOptions = ref<any[]>([])

// Resolve icon string to component
const resolveIcon = (iconName?: string) => {
  if (!iconName) return undefined
  
  // Try exact match
  if ((icons as any)[iconName]) {
    return (icons as any)[iconName]
  }

  // Try PascalCase (e.g. "user" -> "User")
  const pascalName = iconName.charAt(0).toUpperCase() + iconName.slice(1)
  if ((icons as any)[pascalName]) {
    return (icons as any)[pascalName]
  }

  // Try kebab-case to PascalCase (e.g. "user-check" -> "UserCheck")
  const camelToPascal = iconName.replace(/(^\w|-\w)/g, (g) => g.replace('-', '').toUpperCase())
  if ((icons as any)[camelToPascal]) {
    return (icons as any)[camelToPascal]
  }
  
  return undefined
}

// Helper to flatten tree for select options
const flattenMenuOptions = (menus: any[], level = 0, excludeId?: number | string): any[] => {
  let result: any[] = []
  menus.forEach(menu => {
    if (excludeId && String(menu.id) === String(excludeId)) return
    // Skip buttons (type 2) from being selected as parent
    if (menu.type === 2) return

    result.push({
      ...menu,
      level, // Add level for indentation
      label: '\u00A0\u00A0\u00A0\u00A0'.repeat(level) + menu.name
    })
    if (menu.children && menu.children.length > 0) {
      result = result.concat(flattenMenuOptions(menu.children, level + 1, excludeId))
    }
  })
  return result
}

const refreshTable = () => {
  if (tableLayoutRef.value) {
    tableLayoutRef.value.fetchData()
  }
}

const isDialogOpen = ref(false)
const editingMenu = ref<Menu | null>(null)

// Form state
const form = reactive({
  parentId: undefined as string | undefined,
  type: '1',
  name: '',
  path: '',
  icon: '',
  permissionCode: '',
  permissionUrls: '',
  sort: 0,
  visible: true // boolean for switch, convert to 1/0 for API
})

const permissionUrlList = ref<string[]>([''])

const addPermissionUrl = () => {
  permissionUrlList.value.push('')
}

const removePermissionUrl = (index: number) => {
  if (permissionUrlList.value.length > 1) {
    permissionUrlList.value.splice(index, 1)
  } else {
    permissionUrlList.value[0] = ''
  }
}

const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: '请输入菜单名称' }]
}))

const formRef = ref()

const resetForm = () => {
  form.parentId = undefined
  form.type = '1'
  form.name = ''
  form.path = ''
  form.icon = ''
  form.permissionCode = ''
  form.permissionUrls = ''
  permissionUrlList.value = ['']
  form.sort = 0
  form.visible = true
  editingMenu.value = null
}

const handleDelete = async (id: number | string) => {
    console.log('Deleting menu with ID:', id, typeof id)
    if (!id || id === 'undefined' || id === 'null') {
      showToast('无法删除：菜单ID无效', { type: 'error' })
      return
    }
    const confirmed = await showConfirm('确认删除', '确定要删除该菜单及其所有子菜单吗？此操作无法恢复。')
    if (!confirmed) return

    try {
      await deleteMenu(id)
      showToast('删除成功')
      refreshTable()
    } catch (e) {
      console.error(e)
    }
  }

const handleBatchDelete = async (menus: Menu[]) => {
  try {
    const ids = menus.map(menu => menu.id).filter(id => id && id !== 'undefined' && id !== 'null')
    if (ids.length === 0) {
       showToast('未选择有效菜单', { type: 'warning' })
       return
    }
    await batchDeleteMenus(ids)
    showToast(`成功删除 ${menus.length} 条记录`)
    refreshTable()
  } catch (error) {
    console.error(error)
    showToast('部分删除失败，请重试', { type: 'error' })
  }
}

// 获取所有菜单用于选择父级
const fetchAllMenus = async (excludeId?: number | string) => {
  try {
    // 假设 pageSize 10000 能获取所有
    const res = (await getMenus({ page: 1, pageSize: 10000 })) as any
    // res.items is now a tree. Flatten it for options.
    menuOptions.value = flattenMenuOptions(res.items || [], 0, excludeId)
  } catch (error) {
    console.error('Failed to fetch menu options', error)
  }
}

const handleCreate = async () => {
  resetForm()
  await fetchAllMenus()
  isDialogOpen.value = true
}

const handleAddSub = async (row: Menu) => {
  resetForm()
  await fetchAllMenus()
  form.parentId = String(row.id)
  isDialogOpen.value = true
}

const handleEdit = async (menu: Menu) => {
  try {
    resetForm()
    await fetchAllMenus(menu.id)
    editingMenu.value = menu
    form.parentId = menu.parentId ? String(menu.parentId) : undefined
    form.type = String(menu.type || 1)
    form.name = menu.name
    form.path = menu.path || ''
    form.icon = menu.icon || ''
    form.permissionCode = menu.permissionCode || ''
    // Ensure permissionUrls is a string before splitting
    form.permissionUrls = menu.permissionUrls ? String(menu.permissionUrls) : ''
    
    if (form.permissionUrls && typeof form.permissionUrls === 'string') {
        permissionUrlList.value = form.permissionUrls.split(/[\n,]/).map(u => u.trim())
    } else {
        permissionUrlList.value = ['']
    }
    form.sort = menu.sort
    form.visible = menu.visible
    isDialogOpen.value = true
  } catch (error) {
    console.error('Failed to open edit dialog:', error)
    showToast('打开编辑框失败', { type: 'error' })
  }
}

const fetchData = async (params: any) => {
  // Override pageSize to huge number to disable pagination effect effectively
  params.pageSize = 10000
  const res = (await getMenus(params)) as any
  items.value = res.items || []
  return {
    total: res.total,
    items: res.items // Return tree directly, DataTable handles it now
  }
}

const handleSave = async () => {
  try {
    form.permissionUrls = permissionUrlList.value.map(u => u.trim()).filter(Boolean).join(',')
    if (formRef.value) {
      const valid = await formRef.value.validate()
      if (!valid) return
    }

    const data = {
      parentId: form.parentId && form.parentId !== '0' ? parseInt(form.parentId) : undefined,
      type: Number(form.type),
      name: form.name,
      path: form.path,
      icon: form.icon,
      permissionCode: form.permissionCode,
      permissionUrls: form.permissionUrls,
      sort: Number(form.sort),
      visible: form.type === '2' ? false : form.visible
    }

    // Convert parentId to string/number based on API requirement.
    // My API definition says number | string.
    // If backend expects BigInt (serialized as string/number), let's keep it consistent.
    // BaseController usually handles basic types.

    if (editingMenu.value) {
      await updateMenu(editingMenu.value.id, data)
    } else {
      await createMenu(data)
    }
    isDialogOpen.value = false
    refreshTable()
  } catch (e) {
    console.error(e)
    // showToast('保存失败', { type: 'error' })
  }
}

// 定义表格列
const columns = reactive<TableColumn[]>([
  {
    label: '菜单名称', field: 'name', align: 'left', formatter(value: any, row: any) {
      const Icon = row.icon ? (icons as any)[row.icon] : null
      return (
        <div class="flex items-center gap-2">
          {Icon ? <Icon class="w-4 h-4" /> : <span class="w-4 h-4"></span>}
          <span>{value}</span>
        </div>
      )
    }
  },
  { label: '路由路径', field: 'path' },
  { label: '权限标识', field: 'permissionCode' },
  { label: '排序', field: 'sort', sortable: true },
  {
    label: '可见', field: 'visible', sortable: true, formatter(value: any, row: any) {
      if (row.type !== 1) return null
      return (
        <div>
          <Badge color={value ? 'green' : 'gray'}>{value ? '显示' : '隐藏'}</Badge>
        </div>
      )
    }
  },
  {
     label: '操作', field: 'operation', width: '180px', formatter(value: any, row: any) {
       if (!row || !row.id || row.id === 'undefined') {
         console.warn('Row missing valid ID in formatter:', row)
       }
       const btnClass = "inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-8 w-8"
       return (
         <div class="flex items-center gap-2">
            <button class={btnClass} onClick={(e) => { e.stopPropagation(); handleAddSub(row); }} title="新增子菜单">
              <icons.Plus class="w-4 h-4" />
            </button>
            <button class={btnClass} onClick={(e) => { e.stopPropagation(); handleEdit(row); }} title="编辑">
              <icons.Edit2 class="w-4 h-4" />
            </button>
            <button class={btnClass + " text-destructive hover:text-destructive"} onClick={(e) => {
              e.stopPropagation();
              const id = row.id;
              if (!id || id === 'undefined' || id === 'null') {
                showToast(`无法删除：无效的ID (${id})`, { type: 'error' });
                return;
              }
              handleDelete(id);
            }} title="删除">
              <icons.Trash2 class="w-4 h-4" />
            </button>
         </div>
       )
     },
   }
 ])

// Expand logic: collapse nodes that have button children (type 2)
const shouldExpandNode = (row: any) => {
  if (row.children && row.children.length > 0) {
    // If any child is a button (type 2), do not expand this node by default
    if (row.children.some((child: any) => child.type === 2)) {
      return false
    }
  }
  return true
}

</script>

<template>
  <div>
    <TableLayout ref="tableLayoutRef" title="菜单管理" :fetch-data="fetchData" :columns="columns" :items="items"
      :default-expand-level="99" :no-pagination="true" :should-expand-node="shouldExpandNode"
      @delete="handleDelete" @batch-delete="handleBatchDelete" @edit="handleEdit" @create="handleCreate">
    </TableLayout>

    <FormDialog :visible="isDialogOpen" @update:visible="isDialogOpen = $event" :title="editingMenu ? '编辑菜单' : '新增菜单'"
      @confirm="handleSave">
      <div class="py-2">
        <Form ref="formRef" :model="form" :rules="rules">
          <div class="grid gap-4">

            <!-- Parent Menu -->
            <FormItem label="上级菜单" field="parentId">
              <Select v-model="form.parentId">
                <SelectTrigger>
                  <SelectValue placeholder="选择上级菜单" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">顶级菜单</SelectItem>
                  <SelectItem v-for="menu in menuOptions" :key="menu.id" :value="String(menu.id)">
                    <div class="flex items-center">
                      <span :style="{ width: (menu.level * 16) + 'px', display: 'inline-block' }" class="shrink-0"></span>
                      <component 
                         v-if="menu.icon && resolveIcon(menu.icon)"
                         :is="resolveIcon(menu.icon)" 
                         class="w-4 h-4 mr-2 shrink-0 text-muted-foreground"
                       />
                      <span>{{ menu.name }}</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </FormItem>

            <!-- Type Switcher -->
            <FormItem label="菜单类型" field="type" class="w-full">
              <RadioGroup v-model="form.type" class="flex gap-4" default-value="1">
                <div class="flex items-center space-x-2">
                  <RadioGroupItem id="type-menu" value="1" />
                  <Label for="type-menu">页面</Label>
                </div>
                <div class="flex items-center space-x-2">
                  <RadioGroupItem id="type-button" value="2" />
                  <Label for="type-button">权限/按钮</Label>
                </div>
              </RadioGroup>
            </FormItem>

            <!-- Main Info -->
            <div class="grid grid-cols-2 gap-4">
              <FormItem label="菜单名称" field="name" required class="col-span-1">
                <Input v-model="form.name" placeholder="请输入菜单名称" />
              </FormItem>

              <FormItem label="权限标识" field="permissionCode" class="col-span-1">
                <Input v-model="form.permissionCode" placeholder="如: system:menu:list" />
              </FormItem>

              <FormItem label="排序" field="sort" class="col-span-1">
                <Input v-model="form.sort" type="number" placeholder="数字越小越靠前" />
              </FormItem>

              <FormItem v-if="form.type === '1'" label="是否显示" field="visible" class="col-span-1 flex items-end">
                <div class="flex items-center space-x-2 h-9">
                  <Switch :checked="form.visible" @update:checked="form.visible = $event" />
                  <Label class="text-sm text-muted-foreground font-normal">
                    {{ form.visible ? '在菜单中显示' : '在菜单中隐藏' }}
                  </Label>
                </div>
              </FormItem>
            </div>

            <FormItem label="权限URL" field="permissionUrls" class="w-full">
               <div class="flex flex-col gap-2 w-full">
                 <div v-for="(url, index) in permissionUrlList" :key="index" class="flex items-center gap-2">
                   <Input v-model="permissionUrlList[index]" placeholder="[METHOD:]URL (e.g. GET:/api/users)" />
                   <Button variant="ghost" size="icon" @click="removePermissionUrl(index)" :disabled="permissionUrlList.length === 1 && !permissionUrlList[0]" type="button">
                     <icons.Trash2 class="w-4 h-4 text-destructive" />
                   </Button>
                 </div>
                 <Button variant="outline" size="sm" class="w-full mt-1 border-dashed" @click="addPermissionUrl" type="button">
                   <icons.Plus class="w-4 h-4 mr-2" /> 添加权限URL
                 </Button>
                 <p class="text-xs text-muted-foreground mt-1">
                   后端API权限控制，支持多个URL。格式: [METHOD:]URL，默认为 ALL。
                 </p>
               </div>
            </FormItem>

            <!-- Page Specific Settings -->
            <div v-if="form.type === '1'" class="grid grid-cols-2 gap-4 border-t pt-4">
              <FormItem label="路由路径" field="path" class="col-span-1">
                <Input v-model="form.path" placeholder="请输入路由路径" />
              </FormItem>

              <FormItem label="图标" field="icon" class="col-span-1">
                <IconPicker v-model="form.icon" />
              </FormItem>
            </div>

          </div>
        </Form>
      </div>
    </FormDialog>
  </div>
</template>
