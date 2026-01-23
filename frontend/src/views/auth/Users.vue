<script setup lang="tsx">
import { TableColumn } from '@/types/column';
import { getUsers, createUser, updateUser, deleteUser, batchDeleteUsers, type User } from '@/api/user';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { showToast, showConfirm } from '@/lib/message';
import TenantSelect from '@/components/business/TenantSelect.vue';
import RoleSelect from '@/components/business/RoleSelect.vue';
import Form from '@/components/common/form/Form.vue';
import FormItem from '@/components/common/form/FormItem.vue';
import type { FormRules } from '@/lib/symbols';
import { Badge } from '@/components/ui/badge';
import { Eye, EyeOff, Building2, Edit, Trash2, Key } from 'lucide-vue-next';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore()
const isSuper = computed(() => {
  const user = authStore.currentUser
  return user?.isSuper || user?.is_super
})

const isResetPasswordOpen = ref(false)
const resetPasswordUser = ref<User | null>(null)
const resetPasswordForm = reactive({
  password: ''
})
const resetPasswordRules = computed<FormRules>(() => ({
  password: [{ required: true, message: '请输入新密码' }]
}))
const resetPasswordFormRef = ref()

const openResetPassword = (user: User) => {
  resetPasswordUser.value = user
  resetPasswordForm.password = ''
  showPassword.value = true
  isResetPasswordOpen.value = true
}

const handleResetPassword = async () => {
  try {
    if (resetPasswordFormRef.value) {
      const valid = await resetPasswordFormRef.value.validate()
      if (!valid) return
    }

    if (resetPasswordUser.value) {
      await updateUser(resetPasswordUser.value.id, {
        password: resetPasswordForm.password
      })
      showToast('密码重置成功')
      isResetPasswordOpen.value = false
    }
  } catch (error) {
    console.error(error)
  }
}

//TableColumn定义表格列信息
const columns = reactive<TableColumn[]>([
  {
    label: '用户名', field: 'username', formatter(value: any, row: any) {
      return (
        <div class="flex items-center gap-2">
          <Avatar class="h-8 w-8">
            <AvatarImage src={row.avatar} />
            <AvatarFallback>{value ? value.substring(0, 2).toUpperCase() : 'U'}</AvatarFallback>
          </Avatar>
          <span class="font-medium">{value}</span>
        </div>
      )
    },
  },
  { label: '姓名', field: 'realName' },
  { label: '手机号', field: 'phone' },
  {
    label: '角色', field: 'roles', formatter(value: any, row: any) {
      if (row.roles && row.roles.length > 0) {
        return (
          <div class="flex flex-wrap gap-1">
            {row.roles.map((r: any) => (
               <Badge color='primary' key={r.roleId || r.id}>
                 {r.roleName || r.name || r.roleId || r.id} 
               </Badge>
            ))}
          </div>
        )
      }
      // Fallback for old single role field if needed
      if (row.role) {
         return (
            <Badge color='primary'>
              {row.role}
            </Badge>
          )
      }
      return ""
    }
  },
  {
    label: '租户', field: 'tenant', formatter(value: any, row: any) {
      if (!value && !row.tenantId) return <div class="text-muted-foreground">-</div>

      return (
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <Building2 class="h-4 w-4" />
          <span class="max-w-[150px] truncate" title={value?.name || row.tenantId}>
            {value?.name || '-'}
          </span>
        </div>
      )
    }
  },
  { label: '创建时间', field: 'createdAt', sortable: true },
  {
    label: '状态', field: 'status', sortable: true, formatter(value: any, _row: any) {
      return (
        <div>
          <Badge color={value ? 'green' : 'red'}>{value ? '启用' : '禁用'}</Badge>
        </div>
      )
    }
  }
])

const items = ref<User[]>([])
const tableLayoutRef = ref()
const selectedTenantId = ref<string | number | undefined>(undefined)

const refreshTable = () => {
  if (tableLayoutRef.value) {
    tableLayoutRef.value.fetchData()
  }
}

const isDialogOpen = ref(false)
const isViewOpen = ref(false)
const editingUser = ref<User | null>(null)
const viewedUser = ref<User | null>(null)
const showPassword = ref(true)

// Form state
const form = reactive({
  username: '',
  realName: '',
  phone: '',
  password: '',
  email: '',
  roleIds: [] as number[]
})

const rules = computed<FormRules>(() => ({
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: !editingUser.value, message: '请输入密码' }],
  roleIds: [{ required: true, message: '请选择角色' }]
}))

const formRef = ref()

const formTenantId = computed(() => {
  if (editingUser.value && editingUser.value.tenantId) {
    return editingUser.value.tenantId
  }
  return selectedTenantId.value
})

const resetForm = () => {
  form.username = ''
  form.realName = ''
  form.phone = ''
  form.password = ''
  form.email = ''
  editingUser.value = null
  form.roleIds = []
  showPassword.value = false
}

const handleDelete = async (user: User) => {
  const confirmed = await showConfirm('确认删除', `确定要删除用户 "${user.username}" 吗？`)
  if (!confirmed) return
  await deleteUser(user.id)
  showToast('删除成功')
  refreshTable()
}

const handleBatchDelete = async (users: User[]) => {
  try {
    await batchDeleteUsers(users.map(user => user.id))
    showToast(`成功删除 ${users.length} 条记录`)
    refreshTable()
  } catch (error) {
    console.error(error)
    showToast('部分删除失败，请重试', { type: 'error' })
  }
}

const handleCreate = () => {
  resetForm()
  isDialogOpen.value = true
}

const handleUpdate = (user: User) => {
  editingUser.value = user
  form.username = user.username
  form.realName = user.realName || ''
  form.phone = user.phone || ''
  form.email = user.email || ''
  form.password = '' // Don't fill password
  showPassword.value = false
  if (user.roles && user.roles.length > 0) {
    // Map roleId or id, filtering out undefined values
    form.roleIds = user.roles.map((r: any) => Number(r.roleId || r.id)).filter(id => id && !isNaN(id))
  } else {
    form.roleIds = []
  }
  isDialogOpen.value = true
}

const handleView = (user: User) => {
  viewedUser.value = user
  isViewOpen.value = true
}

const fetchData = async (params: any) => {
  const res = (await getUsers({
    ...params,
    tenantId: selectedTenantId.value,
    roleId: selectRoleId.value
  })) as any
  items.value = res.items || []
  return {
    total: res.total,
    items: res.items
  }
}

const handleSave = async () => {
  try {
    if (formRef.value) {
      const valid = await formRef.value.validate()
      if (!valid) return
    }

    if (editingUser.value) {
      await updateUser(editingUser.value.id, {
        username: form.username,
        realName: form.realName,
        phone: form.phone,
        email: form.email,
        ...(form.password ? { password: form.password } : {}),
        roleIds: form.roleIds
      })
      showToast('更新成功')
    } else {
      await createUser({
        username: form.username,
        realName: form.realName,
        phone: form.phone,
        email: form.email,
        password: form.password,
        roleIds: form.roleIds,
        tenantId: formTenantId.value ? Number(formTenantId.value) : undefined
      })
      showToast('创建成功')
    }
    isDialogOpen.value = false
    // Refreshed by saveAndRefresh
  } catch (e) {
    console.error(e)
  }
}

// Override handleSave to call refreshTable
const saveAndRefresh = async () => {
  await handleSave()
  refreshTable()
}
const selectRoleId = ref<string | number | undefined>(undefined)


</script>

<template>
  <div>
    <TableLayout ref="tableLayoutRef" title="账号" :fetch-data="fetchData" :columns="columns" :items="items"
      :checkable="true" @delete="handleDelete" @batch-delete="handleBatchDelete" @update="handleUpdate"
      @view="handleView" @create="handleCreate">
      <template #filters>
        <div v-if="isSuper">
          <TenantSelect v-model="selectedTenantId" placeholder="按租户筛选" @change="refreshTable" />
        </div>
        <CommonSelect v-model="selectRoleId" url="/common/options/roles" placeholder="搜索角色"></CommonSelect>
      </template>
      <template #actions="{ Row, View, Update, Delete }">
        <div class="flex items-center justify-end gap-1">
          <component :is="View" />
          <component :is="Update" />

          <div v-has:reset class="tooltip" data-tip="重置密码">
            <Button variant="ghost" size="icon" class="h-8 w-8" @click="openResetPassword(Row)" title="重置密码">
              <Key class="w-4 h-4" />
            </Button>
          </div>

          <component :is="Delete" />
        </div>
      </template>
    </TableLayout>

    <FormDialog :visible="isDialogOpen" @update:visible="isDialogOpen = $event" :title="editingUser ? '编辑账号' : '新增账号'"
      @confirm="saveAndRefresh">
      <div class="py-2">
        <Form ref="formRef" :model="form" :rules="rules">
          <div class="grid grid-cols-2 gap-4">
            <FormItem label="用户名" field="username" required>
              <Input v-model="form.username" placeholder="登录账号" />
            </FormItem>
            <FormItem v-if="!editingUser" label="密码" field="password">
              <div class="relative">
                <Input v-model="form.password" :type="showPassword ? 'text' : 'password'" placeholder="登录密码"
                  class="pr-10" />
                <button type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
                  @click="showPassword = !showPassword">
                  <Eye v-if="!showPassword" class="w-4 h-4" />
                  <EyeOff v-else class="w-4 h-4" />
                </button>
              </div>
            </FormItem>
            <FormItem label="手机号" field="phone" required>
              <Input v-model="form.phone" placeholder="联系电话" />
            </FormItem>
            <FormItem label="真实姓名" field="realName">
              <Input v-model="form.realName" placeholder="员工姓名" />
            </FormItem>
            <FormItem label="角色" field="roleIds">
              <div v-if="editingUser && editingUser.isMain === 1"
                class="flex h-10 w-full rounded-md border border-input bg-muted px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 items-center text-muted-foreground">
                主账号
              </div>
              <RoleSelect v-else v-model="form.roleIds" placeholder="选择角色" multiple />
            </FormItem>
          </div>

        </Form>
      </div>
    </FormDialog>

    <FormSheet :visible="isViewOpen" @update:visible="isViewOpen = $event" title="账号详情">
      <div v-if="viewedUser" class="space-y-6">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <Label class="text-muted-foreground">用户名</Label>
            <div class="mt-1 font-medium">{{ viewedUser.username }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">真实姓名</Label>
            <div class="mt-1 font-medium">{{ viewedUser.realName || '-' }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">手机号</Label>
            <div class="mt-1 font-medium">{{ viewedUser.phone || '-' }}</div>
          </div>
          <div>
            <Label class="text-muted-foreground">创建时间</Label>
            <div class="mt-1 font-medium">{{ new Date(viewedUser.createdAt).toLocaleString() }}</div>
          </div>
        </div>
      </div>
    </FormSheet>

    <FormDialog :visible="isResetPasswordOpen" @update:visible="isResetPasswordOpen = $event" title="重置密码" width="400px"
      @confirm="handleResetPassword">
      <div class="py-2">
        <Form ref="resetPasswordFormRef" :model="resetPasswordForm" :rules="resetPasswordRules">
          <FormItem label="新密码" field="password" required>
            <div class="relative">
              <Input v-model="resetPasswordForm.password" :type="showPassword ? 'text' : 'password'"
                placeholder="请输入新密码" class="pr-10" />
              <button type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
                @click="showPassword = !showPassword">
                <Eye v-if="!showPassword" class="w-4 h-4" />
                <EyeOff v-else class="w-4 h-4" />
              </button>
            </div>
          </FormItem>
        </Form>
      </div>
    </FormDialog>
  </div>
</template>
