<script setup lang="ts">
import { getStores } from '@/api/store'
import { createTenant, deleteTenant, getTenants, updateTenant, type Tenant } from '@/api/tenant'
import { getUsers, type User } from '@/api/user'
import FormDialog from '@/components/common/form/FormDialog.vue'
import FormSheet from '@/components/common/form/FormSheet.vue'
import Pagination from '@/components/common/Pagination.vue'
import TableActions from '@/components/common/TableActions.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { showToast } from '@/lib/message'
import { formatDate } from '@/utils/date'
import { getAdminUsername } from '@/utils/tenant'
import {
  Building2,
  ChevronDown,
  ChevronUp,
  Plus,
  RefreshCcw,
  RefreshCw,
  Search,
  Store,
  User as UserIcon
} from 'lucide-vue-next'
import { onMounted, ref, watch } from 'vue'

const tenants = ref<Tenant[]>([])
const total = ref(0)
const loading = ref(false)
const searchQuery = ref('')
const page = ref(1)
const pageSize = ref(10)

// Drawer State
const isTenantDrawerOpen = ref(false)
const selectedTenant = ref<Tenant | null>(null)
const isStoreDrawerOpen = ref(false)
const storeList = ref<any[]>([]) // Should type this properly if possible
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
    console.log('Fetching stores for tenant:', tenant.id)
    const res = await getStores({ tenantId: tenant.id, pageSize: 100 }) as any
    console.log('Stores fetched:', res)
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
const tenantUsers = ref<User[]>([])

// Fetch Data
const fetchData = async () => {
  loading.value = true
  try {
    const res: any = await getTenants({
      page: page.value,
      pageSize: pageSize.value,
      name: searchQuery.value
    })
    if (res.items) {
        tenants.value = res.items
        total.value = res.total
    } else {
        tenants.value = []
        total.value = 0
    }
  } catch (error) {
    console.error(error)
    showToast('获取租户列表失败', { type: 'error' })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

// Watchers
watch([page, pageSize], () => {
  fetchData()
})

const handleSearch = () => {
  page.value = 1
  fetchData()
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
  // Exclude sensitive fields if any, though tenant usually doesn't have password directly
  editingTenant.value = { ...tenant }
  errors.value = {}
  showOptionalTenantInfo.value = false
  isEditModalOpen.value = true
}

const saveTenant = async () => {
  errors.value = {}
  let hasError = false
  
  // Validation
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
          errors.value.password = '请输入管理员密码'
          hasError = true
      }
      if (!editingTenant.value.phone) {
          errors.value.phone = '请输入管理员手机号'
          hasError = true
      }
  }

  if (hasError) {
      // Focus first error
      /* nextTick not strictly needed with standard ref updates but safer for DOM */
      setTimeout(() => {
          if (errors.value.name && nameInputRef.value) {
              nameInputRef.value.focus() // Note: ref binding needed in template
          } else if (errors.value.username && usernameInputRef.value) {
              usernameInputRef.value.focus()
          } else if (errors.value.password && passwordInputRef.value) {
              passwordInputRef.value.focus()
          } else if (errors.value.phone && phoneInputRef.value) {
              phoneInputRef.value.focus()
          }
      }, 0)
      return
  }

  try {
    if (isCreate.value) {
      await createTenant(editingTenant.value)
      showToast('租户创建成功')
    } else {
      if (editingTenant.value.id) {
         // Create a copy and remove non-updatable fields
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
    fetchData()
  } catch (error: any) {
    // Toast is already handled in request.ts for non-200 responses.
    // However, if it's a network error or something else, we might want to log it.
    // To avoid double toast, we can rely on the interceptor or check if the error was already handled.
    // The request interceptor returns Promise.reject(error), so we catch it here.
    // We should not showToast here if the interceptor already did.
    console.error('Operation failed:', error)
  }
}

const handleDelete = async (id: string) => {
    try {
        await deleteTenant(id)
        showToast('删除成功')
        fetchData()
    } catch (error: any) {
        console.error('Delete failed:', error)
    }
}

// Helpers

</script>

<template>
  <div class="space-y-2">
    <!-- Filters -->
    <div class="py-2 rounded-xl shadow-sm flex justify-between items-center">
      <div class="flex flex-col md:flex-row gap-4 items-center">
        <div class="relative flex-1 md:max-w-sm">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-4 h-4"/>
          <Input 
            v-model="searchQuery" 
            placeholder="搜索租户名称..." 
            class="pl-9"
            @keyup.enter="handleSearch"
          />
        </div>
        <div class="flex gap-2">
          <Button @click="handleSearch" variant="outline">
            <Search class="w-4 h-4 mr-2" />
            查询
          </Button>
        </div>
      </div>
      <div class="flex gap-2">
         <div>
           <Button @click="openCreateModal" variant="outline">
            <Plus class="w-4 h-4 mr-2" />
            新增租户
          </Button>
        </div>
          <Button variant="outline" @click="fetchData">
          <RefreshCw class="w-4 h-4 mr-1" />
        </Button>
      </div>
    </div>

    <!-- Tenant List -->
    <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>租户信息</TableHead>
            <TableHead>管理员账号</TableHead>
            <TableHead>联系人</TableHead>
            <TableHead>联系电话</TableHead>
            <TableHead>门店数量</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead class="text-center">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="tenant in tenants" :key="tenant.id" class="group">
            <TableCell>
              <div class="flex items-center cursor-pointer hover:underline" @click="openTenantDrawer(tenant)">
                <div class="h-10 w-10 flex-shrink-0 bg-primary/10 rounded-full flex items-center justify-center text-primary">
                  <Building2 class="w-5 h-5" />
                </div>
                <div class="ml-4">
                  <div class="text-sm font-bold text-foreground">{{ tenant.name }}</div>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <div class="flex items-center gap-2">
                 <Badge variant="outline" class="gap-1">
                    <UserIcon class="w-3 h-3" />
                    {{ getAdminUsername(tenant) }}
                 </Badge>
              </div>
            </TableCell>
            <TableCell>
              <div class="text-sm">
                <div>{{ tenant.contactName || '-' }}</div>
              </div>
            </TableCell>
            <TableCell>
              <div class="text-sm">
                 {{ tenant.contactPhone || '-' }}
              </div>
            </TableCell>
            <TableCell>
              <div class="flex items-center gap-2 cursor-pointer hover:text-primary" @click="openStoreDrawer(tenant)">
                <Store class="w-4 h-4" />
                <span class="font-medium">{{ tenant._count?.stores || 0 }}</span>
              </div>
            </TableCell>
            <TableCell>
              <Badge :color="tenant.status === 1 ? 'green' : 'gray'">
                {{ tenant.status === 1 ? '启用' : '禁用' }}
              </Badge>
            </TableCell>
            <TableCell>
              <div class="text-sm text-muted-foreground">
                {{ formatDate(tenant.createdAt) }}
              </div>
            </TableCell>
            <TableCell class="text-center">
              <TableActions
                @view="openTenantDrawer(tenant)"
                @edit="openEditModal(tenant)"
                @delete="handleDelete(tenant.id)"
              />
            </TableCell>
          </TableRow>
          <TableRow v-if="tenants.length === 0 && !loading">
            <TableCell colspan="7" class="text-center py-10 text-sm text-muted-foreground">
              暂无租户数据。
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination -->
      <Pagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        :loading="loading"
        @update:page="(val: number) => page = val"
      />
    </div>

    <!-- Edit/Create Modal -->
    <FormDialog 
    :title="isCreate ? '新增租户' : '编辑租户'"
    :description="isCreate ? '填写租户信息及初始管理员账号。默认门店将自动创建。' : '修改租户基本信息。'"
    @confirm="saveTenant"
    v-model:visible="isEditModalOpen">
        <div class="grid gap-6 py-4">
          <!-- Tenant Info -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-medium text-muted-foreground">租户信息</h3>
              <Button 
                variant="ghost" 
                size="sm" 
                class="h-6 text-xs"
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
                      placeholder="请输入租户名称" 
                      :class="{'border-red-500 focus-visible:ring-red-500': errors.name}"
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
                     <Label for="status-mode" class="cursor-pointer whitespace-nowrap">{{ editingTenant.status === 1 ? '启用' : '禁用' }}</Label>
                  </div>
                </div>
              </div>
              
              <div v-if="showOptionalTenantInfo" class="grid grid-cols-2 gap-4 animate-in fade-in slide-in-from-top-2 duration-200">
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

          <!-- Admin Info (Only for Create) -->
          <div v-if="isCreate" class="space-y-4 border-t pt-4">
            <div class="flex items-center justify-between">
              <h3 class="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <UserIcon class="w-4 h-4" />
                  管理员账号
              </h3>
              <Button variant="outline" size="sm" class="h-7 text-xs" @click="generateCredentials">
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
                        placeholder="登录账号" 
                        :class="{'border-red-500 focus-visible:ring-red-500': errors.username}"
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
                        type="text" 
                        placeholder="登录密码" 
                        :class="{'border-red-500 focus-visible:ring-red-500': errors.password}"
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
                        placeholder="管理员手机号" 
                        :class="{'border-red-500 focus-visible:ring-red-500': errors.phone}"
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
      title="租户详情">
        <div class="mt-6 space-y-6" v-if="selectedTenant">
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
             <div class="grid grid-cols-2 gap-4" v-if="selectedTenant.users && selectedTenant.users.length">
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
    <FormSheet v-model:visible="isStoreDrawerOpen" title="门店列表" :description="selectedTenant ? `租户 ${selectedTenant.name} 下的所有门店` : ''">
        <div class="mt-6">
           <div v-if="storeLoading" class="flex justify-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
           </div>
           <div v-else-if="storeList.length === 0" class="text-center py-8 text-muted-foreground">
              暂无门店数据
           </div>
           <div v-else class="space-y-4">
              <div v-for="store in storeList" :key="store.id" class="flex items-start justify-between p-4 rounded-lg border bg-card">
                 <div>
                    <div class="font-medium flex items-center gap-2">
                      {{ store.name }}
                      <Badge variant="outline" class="text-xs" v-if="store.code">{{ store.code }}</Badge>
                    </div>
                    <div class="text-sm text-muted-foreground mt-1" v-if="store.address">{{ store.address }}</div>
                    <div class="text-xs text-muted-foreground mt-2">
                       创建于 {{ formatDate(store.createdAt) }}
                    </div>
                 </div>
                 <Badge :variant="store.status === 1 ? 'default' : 'secondary'">
                    {{ store.status === 1 ? '营业中' : '休息中' }}
                 </Badge>
              </div>
           </div>
        </div>
    </FormSheet>
  </div>
</template>
