<script lang="tsx" setup>
import {getStores} from '@/api/store'
import {createTenant, deleteTenant, getTenants, type Tenant, updateTenant} from '@/api/tenant'
import FormDialog from '@/components/common/form/FormDialog.vue'
import FormSheet from '@/components/common/form/FormSheet.vue'
import TableLayout from '@/components/common/TableLayout.vue'
import {Badge} from '@/components/ui/badge'
import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {Label} from '@/components/ui/label'
import {Switch} from '@/components/ui/switch'
import {showToast} from '@/lib/message'
import {TableColumn} from '@/types/column'
import {formatDate} from '@/utils/date'
import {ChevronDown, ChevronUp, Loader2, RefreshCcw, Store, User as UserIcon} from 'lucide-vue-next'
import {ref} from 'vue'

const tableLayoutRef = ref()
const isTenantDrawerOpen = ref(false)
const selectedTenant = ref<Tenant | null>(null)
const isStoreDrawerOpen = ref(false)
const storeList = ref<any[]>([]) 
const storeLoading = ref(false)

// Drawer Logic
const openTenantDrawer = (tenant: Tenant) => {
  selectedTenant.value = tenant
  isTenantDrawerOpen.value = true
}

const openStoreDrawer = async (tenant: Tenant) => {
  if (!tenant || !tenant.id) {
    console.error('Tenant ID is missing')
    return
  }
  selectedTenant.value = tenant
  isStoreDrawerOpen.value = true
  storeLoading.value = true
  storeList.value = [] // Reset list
  try {
    const res = await getStores({ tenantId: tenant.id, pageSize: 100 }) as any
    storeList.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch stores:', error)
    showToast('获取门店列表失败', { type: 'error' })
  } finally {
    storeLoading.value = false
  }
}

// Edit/Create State
const isEditModalOpen = ref(false)
const isCreate = ref(true)
const editingTenant = ref<Partial<Tenant> & { username?: string; password?: string; realName?: string; phone?: string; mainUserId?: string }>({})
const errors = ref<Record<string, string>>({})
const nameInputRef = ref<HTMLInputElement | null>(null)
const usernameInputRef = ref<HTMLInputElement | null>(null)
const passwordInputRef = ref<HTMLInputElement | null>(null)
const phoneInputRef = ref<HTMLInputElement | null>(null)

const showOptionalTenantInfo = ref(false)

// Fetch Data
const fetchData = async (params: any) => {
  try {
    const res: any = await getTenants({
      page: params.page,
      pageSize: params.pageSize,
      keyword: params.keyword
    })
    return {
      total: res.total,
      items: res.items || []
    }
  } catch (error) {
    console.error(error)
    showToast('获取租户列表失败', { type: 'error' })
    return { total: 0, items: [] }
  }
}

const refreshTable = () => {
  if (tableLayoutRef.value) {
    tableLayoutRef.value.fetchData()
  }
}

// Modal Logic
const generateCredentials = () => {
  const randomStr = Math.random().toString(36).slice(-8)
  editingTenant.value.username = `user_${randomStr}`
  editingTenant.value.password = Math.random().toString(36).slice(-8) + 'Aa1!'
}

const openCreateModal = () => {
  isCreate.value = true
  editingTenant.value = {
    status: 1
  }
  errors.value = {}
  showOptionalTenantInfo.value = false
  generateCredentials()
  isEditModalOpen.value = true
}

const openEditModal = (tenant: Tenant) => {
  isCreate.value = false
  editingTenant.value = { ...tenant }
  errors.value = {}
  showOptionalTenantInfo.value = false
  
  if (tenant.users && tenant.users.length > 0) {
      const admin = tenant.users[0]
      editingTenant.value.mainUserId = admin.id
      editingTenant.value.username = admin.username
      editingTenant.value.realName = admin.realName
      editingTenant.value.phone = admin.phone
  }
  isEditModalOpen.value = true
}

const handleDelete = async (tenant: Tenant) => {
  try {
    await deleteTenant(tenant.id)
    showToast('租户删除成功')
    refreshTable()
  } catch (error) {
    console.error(error)
    showToast('删除租户失败', { type: 'error' })
  }
}

const saveTenant = async () => {
  // Validation
  errors.value = {}
  let hasError = false
  
  if (!editingTenant.value.name) {
    errors.value.name = '请输入租户名称'
    hasError = true
  }

  if (isCreate.value) {
     if (!editingTenant.value.username) {
        errors.value.username = '请输入管理员账号'
        hasError = true
     }
     if (!editingTenant.value.password) {
        errors.value.password = '请输入初始密码'
        hasError = true
     }
     if (!editingTenant.value.phone) {
        errors.value.phone = '请输入管理员手机号'
        hasError = true
     }
  }

  if (hasError) {
      setTimeout(() => {
          if (errors.value.name && nameInputRef.value) nameInputRef.value.focus()
          else if (errors.value.username && usernameInputRef.value) usernameInputRef.value.focus()
          else if (errors.value.password && passwordInputRef.value) passwordInputRef.value.focus()
          else if (errors.value.phone && phoneInputRef.value) phoneInputRef.value.focus()
      }, 0)
      return
  }

  try {
    if (isCreate.value) {
      await createTenant(editingTenant.value)
      showToast('租户创建成功')
    } else {
      if (editingTenant.value.id) {
         const updateData = { ...editingTenant.value }
         delete (updateData as any).users
         delete (updateData as any)._count
         delete (updateData as any).createdAt
         delete (updateData as any).updatedAt
         
         await updateTenant(editingTenant.value.id, updateData)
         showToast('租户更新成功')
      }
    }
    isEditModalOpen.value = false
    refreshTable()
  } catch (error: any) {
    console.error(error)
    if (error.response?.data?.message) {
        if (error.response.data.message.includes('username')) {
            errors.value.username = '用户名已存在'
        } else if (error.response.data.message.includes('phone')) {
            errors.value.phone = '手机号已存在'
        } else {
            showToast(error.response.data.message, { type: 'error' })
        }
    } else {
        showToast(isCreate.value ? '创建失败' : '更新失败', { type: 'error' })
    }
  }
}

const columns: TableColumn[] = [
  { label: '租户信息', field: 'name', minWidth: '150px' },
  { 
    label: '管理员账号', 
    field: 'users', 
    formatter: (val: any, row: any) => {
        const username = row.users?.[0]?.username || '-'
        return (
            <div class="flex items-center gap-2">
                <Badge variant="outline" class="gap-1">
                    <UserIcon class="w-3 h-3" />
                    {username}
                </Badge>
            </div>
        )
    }
  },
  { label: '联系人', field: 'contactName' },
  { label: '联系电话', field: 'contactPhone' },
  { 
    label: '门店数量', 
    field: '_count', 
    formatter: (val: any, row: any) => {
        const count = row._count?.stores || 0
        return (
             <div class="flex items-center gap-2 cursor-pointer hover:text-primary" onClick={(e) => { e.stopPropagation(); openStoreDrawer(row); }}>
                <Store class="w-4 h-4" />
                <span class="font-medium">{count}</span>
            </div>
        )
    }
  },
  { 
    label: '状态', 
    field: 'status', 
    formatter: (val: number) => (
      <Badge color={val === 1 ? 'green' : 'gray'}>
        {val === 1 ? '启用' : '禁用'}
      </Badge>
    )
  },
  { label: '创建时间', field: 'createdAt', formatter: (val: string) => formatDate(val) },
]
</script>

<template>
  <TableLayout
    ref="tableLayoutRef"
    :columns="columns"
    :fetch-data="fetchData"
    title="租户"
    @create="openCreateModal"
    @delete="handleDelete"
    @update="openEditModal"
    @view="openTenantDrawer"
  >
    <template #actions="{ Row, View, Update, Delete }">
      <Button 
        class="h-8 w-8"
        size="icon" 
        title="查看门店"
        variant="ghost"
        @click="openStoreDrawer(Row)"
      >
        <Store class="w-4 h-4" />
      </Button>
      <component :is="View" />
      <component :is="Update" />
      <component :is="Delete" />
    </template>
  </TableLayout>

  <!-- Edit/Create Modal -->
  <FormDialog 
    :title="isCreate ? '新增租户' : '编辑租户'"
    :description="isCreate ? '填写租户信息及初始管理员账号。默认门店将自动创建。' : '修改租户基本信息。'"
    @confirm="saveTenant"
    v-model:visible="isEditModalOpen"
  >
    <div class="grid gap-6 py-4">
      <!-- Tenant Info -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-medium text-muted-foreground">租户信息</h3>
          <Button 
            class="h-6 text-xs"
            size="sm" 
            variant="ghost"
            @click="showOptionalTenantInfo = !showOptionalTenantInfo"
          >
            {{ showOptionalTenantInfo ? '隐藏可选信息' : '显示更多信息' }}
            <ChevronUp v-if="showOptionalTenantInfo" class="ml-1 w-3 h-3" />
            <ChevronDown v-else class="ml-1 w-3 h-3" />
          </Button>
        </div>
        
        <div class="space-y-4">
          <div class="grid grid-cols-12 gap-4">
            <div class="space-y-2 col-span-9">
              <Label>租户名称 <span class="text-red-500">*</span></Label>
              <div class="relative">
                <Input 
                  ref="nameInputRef"
                  v-model="editingTenant.name" 
                  :class="{'border-red-500 focus-visible:ring-red-500': errors.name}"
                  placeholder="请输入租户名称"
                />
                <span v-if="errors.name" class="text-xs text-red-500 absolute -bottom-5 left-0">{{ errors.name }}</span>
              </div>
            </div>
            <div class="space-y-2 col-span-3">
              <Label>状态</Label>
              <div class="flex items-center space-x-2 h-10">
                 <Switch 
                    id="status-mode" 
                    :checked="editingTenant.status === 1"
                    @update:checked="(val: boolean) => editingTenant.status = val ? 1 : 0"
                 />
                 <Label class="cursor-pointer whitespace-nowrap" for="status-mode">{{ editingTenant.status === 1 ? '启用' : '禁用' }}</Label>
              </div>
            </div>
          </div>

          <div v-show="showOptionalTenantInfo" class="grid grid-cols-2 gap-4 animate-in fade-in slide-in-from-top-2 duration-300">
             <div class="space-y-2">
              <Label>联系人</Label>
              <Input v-model="editingTenant.contactName" placeholder="联系人姓名" />
            </div>
             <div class="space-y-2">
              <Label>联系电话</Label>
              <Input v-model="editingTenant.contactPhone" placeholder="联系电话" />
            </div>
             <div class="space-y-2 col-span-2">
              <Label>联系邮箱</Label>
              <Input v-model="editingTenant.contactEmail" placeholder="联系邮箱" />
            </div>
          </div>
        </div>
      </div>

      <div class="h-px bg-border"></div>

      <!-- Admin Info -->
      <div v-if="isCreate" class="space-y-4">
        <div class="flex items-center justify-between">
            <h3 class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <UserIcon class="w-4 h-4" />
                管理员账号
            </h3>
            <Button class="h-7 text-xs" size="sm" variant="outline" @click="generateCredentials">
            <RefreshCcw class="w-3 h-3 mr-1" />
            刷新账号密码
            </Button>
        </div>
        
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label>用户名 <span class="text-red-500">*</span></Label>
                <div class="relative">
                    <Input 
                    ref="usernameInputRef"
                    v-model="editingTenant.username" 
                    :class="{'border-red-500 focus-visible:ring-red-500': errors.username}"
                    placeholder="登录账号"
                    />
                    <span v-if="errors.username" class="text-xs text-red-500 absolute -bottom-5 left-0">{{ errors.username }}</span>
                </div>
            </div>
            <div class="space-y-2">
                <Label>密码 <span class="text-red-500">*</span></Label>
                <div class="relative">
                    <Input 
                    ref="passwordInputRef"
                    v-model="editingTenant.password" 
                    :class="{'border-red-500 focus-visible:ring-red-500': errors.password}"
                    placeholder="登录密码" 
                    type="text"
                    />
                    <span v-if="errors.password" class="text-xs text-red-500 absolute -bottom-5 left-0">{{ errors.password }}</span>
                </div>
            </div>
            <div class="space-y-2">
                <Label>手机号 <span class="text-red-500">*</span></Label>
                <div class="relative">
                    <Input 
                    ref="phoneInputRef"
                    v-model="editingTenant.phone" 
                    :class="{'border-red-500 focus-visible:ring-red-500': errors.phone}"
                    placeholder="管理员手机号"
                    />
                    <span v-if="errors.phone" class="text-xs text-red-500 absolute -bottom-5 left-0">{{ errors.phone }}</span>
                </div>
            </div>
            <div class="space-y-2">
                <Label>真实姓名</Label>
                <Input v-model="editingTenant.realName" placeholder="管理员姓名" />
            </div>
        </div>
      </div>
    </div>
  </FormDialog>

  <!-- Tenant Detail Drawer -->
  <FormSheet 
    v-model:visible="isTenantDrawerOpen" 
    :description="selectedTenant ? '租户详情' : ''"
    title="租户详情"
  >
    <div v-if="selectedTenant" class="mt-6 space-y-6">
      <div class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <Label class="text-muted-foreground">租户名称</Label>
              <div class="font-medium mt-1">{{ selectedTenant.name }}</div>
            </div>
            <div>
              <Label class="text-muted-foreground">状态</Label>
              <div class="mt-1">
                <Badge :color="selectedTenant.status === 1 ? 'green' : 'gray'">
                {{ selectedTenant.status === 1 ? '启用' : '禁用' }}
              </Badge>
              </div>
            </div>
            <div>
              <Label class="text-muted-foreground">联系人</Label>
              <div class="font-medium mt-1">{{ selectedTenant.contactName || '-' }}</div>
            </div>
            <div>
              <Label class="text-muted-foreground">联系电话</Label>
              <div class="font-medium mt-1">{{ selectedTenant.contactPhone || '-' }}</div>
            </div>
            <div class="col-span-2">
              <Label class="text-muted-foreground">联系邮箱</Label>
              <div class="font-medium mt-1">{{ selectedTenant.contactEmail || '-' }}</div>
            </div>

            <div class="col-span-2">
              <Label class="text-muted-foreground">创建时间</Label>
              <div class="font-medium mt-1">{{ formatDate(selectedTenant.createdAt) }}</div>
            </div>
          </div>
      </div>
      
      <div class="border-t pt-4">
          <h4 class="font-medium mb-3">管理员信息</h4>
          <div v-if="selectedTenant.users && selectedTenant.users.length" class="grid grid-cols-2 gap-4">
            <div>
                <Label class="text-muted-foreground">账号</Label>
                <div class="font-medium mt-1">{{ selectedTenant.users[0].username }}</div>
            </div>
              <div>
                <Label class="text-muted-foreground">手机号</Label>
                <div class="font-medium mt-1">{{ selectedTenant.users[0].phone || '-' }}</div>
            </div>
              <div>
                <Label class="text-muted-foreground">邮箱</Label>
                <div class="font-medium mt-1">{{ selectedTenant.users[0].email || '-' }}</div>
            </div>
          </div>
          <div v-else class="text-sm text-muted-foreground">暂无管理员信息</div>
      </div>
    </div>
  </FormSheet>

  <!-- Store List Drawer -->
  <FormSheet
    v-model:visible="isStoreDrawerOpen"
    :title="selectedTenant?.name ? `${selectedTenant.name} - 门店列表` : '门店列表'"
    description="该租户下的所有门店"
  >
    <div v-if="storeLoading" class="flex justify-center py-10">
       <Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
    </div>
    <div v-else-if="storeList.length === 0" class="text-center py-10 text-muted-foreground">
      暂无门店数据
    </div>
    <div v-else class="space-y-4 mt-4">
      <div 
        v-for="store in storeList" 
        :key="store.id"
        class="flex items-center p-3 border rounded-lg hover:bg-muted/50 transition-colors"
      >
        <div class="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary mr-3">
          <Store class="w-5 h-5" />
        </div>
        <div class="flex-1">
          <div class="font-medium">{{ store.name }}</div>
          <div class="text-xs text-muted-foreground mt-0.5">
             {{ store.address || '暂无地址' }}
          </div>
        </div>
        <Badge variant="outline">{{ store.code || '无编码' }}</Badge>
      </div>
    </div>
  </FormSheet>
</template>