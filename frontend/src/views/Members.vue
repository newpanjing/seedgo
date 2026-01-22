<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useUserStore } from '../stores/users'
import { 
  Search, 
  Plus, 
  Trash2, 
  RefreshCw,
  User as UserIcon,
  Users as UsersIcon,
  Zap
} from 'lucide-vue-next'
import type { User } from '../types/user'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'

const store = useUserStore()
const searchQuery = ref('')
const levelFilter = ref<'all' | User['level']>('all')
const minBalance = ref<number | undefined>(undefined)
const selectedIds = ref<string[]>([])
const isEditModalOpen = ref(false)
const isRechargeModalOpen = ref(false)
const editingUser = ref<Partial<User>>({})
const rechargeAmount = ref(0)
const deleteDialogOpen = ref(false)
const pendingDeleteIds = ref<string[]>([])
const currentTag = ref('')

const levelOptions = [
  { label: '💎 钻石会员', value: 'diamond' },
  { label: '🥇 金牌会员', value: 'gold' },
  { label: '🥈 银牌会员', value: 'silver' },
  { label: '🥉 普通会员', value: 'normal' }
]

const openEditModal = (user?: User) => {
  if (user) {
    editingUser.value = { ...user }
    currentTag.value = user.tags?.[0] || ''
  } else {
    editingUser.value = {
      name: '',
      phone: '',
      level: 'normal',
      discount: 1,
      balance: 0,
      points: 0,
      tags: []
    }
    currentTag.value = ''
  }
  isEditModalOpen.value = true
}

const openRechargeModal = (user: User) => {
  editingUser.value = { ...user }
  rechargeAmount.value = 0
  isRechargeModalOpen.value = true
}

const closeEditModal = () => {
  isEditModalOpen.value = false
}

const closeRechargeModal = () => {
  isRechargeModalOpen.value = false
}

const saveUser = () => {
  editingUser.value.tags = currentTag.value ? [currentTag.value] : []

  if (editingUser.value.id) {
    store.updateUser(editingUser.value.id, editingUser.value)
  } else {
    store.addUser(editingUser.value as any)
  }
  isEditModalOpen.value = false
}

const handleRecharge = () => {
  if (editingUser.value.id && rechargeAmount.value > 0) {
    const newBalance = (editingUser.value.balance || 0) + rechargeAmount.value
    store.updateUser(editingUser.value.id, { balance: newBalance })
    isRechargeModalOpen.value = false
  }
}

const filteredUsers = computed(() => {
  return store.users.filter(user => {
    const keyword = searchQuery.value.trim().toLowerCase()

    const matchKeyword =
      !keyword ||
      user.name.toLowerCase().includes(keyword) ||
      user.phone.includes(keyword)

    const matchLevel =
      levelFilter.value === 'all' || user.level === levelFilter.value

    const matchBalance =
      minBalance.value == null || minBalance.value === 0
        ? true
        : user.balance >= minBalance.value

    return matchKeyword && matchLevel && matchBalance
  })
})

const page = ref(1)
const pageSize = ref(10)
const pageSizeOptions = [10, 20, 50]
const pageInput = ref('')

const totalUsers = computed(() => filteredUsers.value.length)

const totalPages = computed(() => {
  if (totalUsers.value === 0) return 0
  return Math.ceil(totalUsers.value / pageSize.value)
})

const paginatedUsers = computed(() => {
  if (totalPages.value === 0) return []
  const currentPage = Math.min(page.value, totalPages.value)
  const start = (currentPage - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredUsers.value.slice(start, end)
})

const pageStart = computed(() => {
  if (totalUsers.value === 0) return 0
  return (page.value - 1) * pageSize.value + 1
})

const pageEnd = computed(() => {
  if (totalUsers.value === 0) return 0
  return Math.min(page.value * pageSize.value, totalUsers.value)
})

watch([filteredUsers, pageSize], () => {
  if (totalPages.value === 0) {
    page.value = 1
  } else {
    if (page.value > totalPages.value) {
      page.value = totalPages.value
    }
    if (page.value < 1) {
      page.value = 1
    }
  }
})

watch(
  () => page.value,
  () => {
    if (totalPages.value === 0) {
      pageInput.value = ''
    } else {
      pageInput.value = String(page.value)
    }
  },
  { immediate: true },
)

const handleSearch = () => {
}

const handleReset = () => {
  searchQuery.value = ''
  levelFilter.value = 'all'
  minBalance.value = undefined
  selectedIds.value = []
  page.value = 1
}

const selectedCount = computed(() => selectedIds.value.length)

const canDelete = computed(() => selectedCount.value > 0)

const isAllSelected = computed(() => {
  if (paginatedUsers.value.length === 0) return false
  return paginatedUsers.value.every(user => selectedIds.value.includes(user.id))
})

const isIndeterminate = computed(() => {
  const currentIds = paginatedUsers.value.map(user => user.id)
  const selectedOnPage = currentIds.filter(id => selectedIds.value.includes(id)).length
  return selectedOnPage > 0 && selectedOnPage < currentIds.length
})

const handleToggleSelectAll = (event: Event) => {
  const target = event.target as HTMLInputElement
  const pageIds = paginatedUsers.value.map(user => user.id)
  if (target.checked) {
    const set = new Set([...selectedIds.value, ...pageIds])
    selectedIds.value = Array.from(set)
  } else {
    const pageIdSet = new Set(pageIds)
    selectedIds.value = selectedIds.value.filter(id => !pageIdSet.has(id))
  }
}

const handleToggleRowSelection = (id: string, event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id]
    }
  } else {
    selectedIds.value = selectedIds.value.filter(existingId => existingId !== id)
  }
}

const openBulkDeleteDialog = () => {
  if (selectedIds.value.length === 0) return
  pendingDeleteIds.value = [...selectedIds.value]
  deleteDialogOpen.value = true
}

const openRowDeleteDialog = (id: string) => {
  pendingDeleteIds.value = [id]
  deleteDialogOpen.value = true
}

const handleConfirmDelete = () => {
  if (pendingDeleteIds.value.length === 0) {
    deleteDialogOpen.value = false
    return
  }

  pendingDeleteIds.value.forEach(id => {
    store.deleteUser(id)
  })

  if (selectedIds.value.length > 0) {
    const set = new Set(pendingDeleteIds.value)
    selectedIds.value = selectedIds.value.filter(id => !set.has(id))
  }

  pendingDeleteIds.value = []
  deleteDialogOpen.value = false
}

const handleJumpToPage = () => {
  if (totalPages.value === 0) {
    pageInput.value = ''
    return
  }

  const raw = pageInput.value.trim()
  if (!raw) {
    pageInput.value = String(page.value)
    return
  }

  const value = Number(raw)
  if (Number.isNaN(value)) {
    pageInput.value = String(page.value)
    return
  }

  let target = Math.floor(value)
  if (target < 1) target = 1
  if (target > totalPages.value) target = totalPages.value
  page.value = target
}

const getLevelClass = (level: User['level']) => {
  const colors: Record<User['level'], string> = {
    normal: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200 hover:bg-gray-100/80 dark:hover:bg-gray-800/80',
    silver: 'bg-slate-100 text-slate-800 dark:bg-slate-800 dark:text-slate-200 hover:bg-slate-100/80 dark:hover:bg-slate-800/80',
    gold: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200 hover:bg-yellow-100/80 dark:hover:bg-yellow-900/80',
    diamond: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200 hover:bg-purple-100/80 dark:hover:bg-purple-900/80'
  }
  return colors[level] || colors.normal
}

const getLevelText = (level: User['level']) => {
  const map: Record<string, string> = {
    diamond: '💎 钻石会员',
    gold: '🥇 金牌会员',
    silver: '🥈 银牌会员',
    normal: '🥉 普通会员'
  }
  return map[level] || level
}

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(price)
}

const formatDate = (isoString: string) => {
  return new Date(isoString).toLocaleDateString()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-card p-6 rounded-xl shadow-sm border border-border">
      <div>
        <h1 class="text-2xl font-bold text-foreground flex items-center">
          <UsersIcon class="w-8 h-8 mr-3 text-primary" />
          会员管理
        </h1>
        <p class="mt-1 text-sm text-muted-foreground">
          管理会员信息、等级权益及消费记录
        </p>
      </div>
      <Button @click="openEditModal()">
        <Plus class="w-4 h-4 mr-2" />
        新增会员
      </Button>
    </div>

    <!-- Filters -->
    <div class="bg-card p-4 rounded-xl shadow-sm border border-border">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-3 md:flex-1">
          <div>
            <Label class="mb-1 block text-xs text-muted-foreground">关键字</Label>
            <div class="relative">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 transform text-muted-foreground w-4 h-4" />
              <Input 
                v-model="searchQuery"
                type="text" 
                placeholder="搜索姓名或手机号..." 
                class="pl-9"
              />
            </div>
          </div>
          <div>
            <Label class="mb-1 block text-xs text-muted-foreground">会员等级</Label>
            <Select v-model="levelFilter">
              <SelectTrigger>
                <SelectValue placeholder="全部等级" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部等级</SelectItem>
                <SelectItem
                  v-for="option in levelOptions"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ option.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label class="mb-1 block text-xs text-muted-foreground">最低余额（元）</Label>
            <Input 
              v-model.number="minBalance"
              type="number"
              min="0"
              placeholder="不限"
            />
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <Button variant="outline" @click="handleReset">
            重置
          </Button>
          <Button @click="handleSearch">
            <Search class="w-4 h-4 mr-2" />
            查询
          </Button>
        </div>
      </div>
    </div>

    <!-- User List -->
    <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
      <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/30">
        <div class="flex items-center gap-2">
          <Button size="sm" @click="openEditModal()">
            <Plus class="w-4 h-4 mr-1" />
            新增
          </Button>
          <Button
            v-if="canDelete"
            size="sm"
            variant="destructive"
            @click="openBulkDeleteDialog"
          >
            <Trash2 class="w-4 h-4 mr-1" />
            删除 {{ selectedCount }} 行
          </Button>
        </div>
        <div class="flex items-center gap-2">
          <Button size="sm" variant="ghost" @click="handleReset">
            <RefreshCw class="w-4 h-4 mr-1" />
            刷新
          </Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-10">
              <input
                type="checkbox"
                class="checkbox checkbox-sm checkbox-primary"
                :checked="isAllSelected"
                :indeterminate="isIndeterminate"
                @change="handleToggleSelectAll"
              />
            </TableHead>
            <TableHead>会员信息</TableHead>
            <TableHead>等级</TableHead>
            <TableHead>资产信息</TableHead>
            <TableHead>累计消费</TableHead>
            <TableHead>最近访问</TableHead>
            <TableHead class="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="user in paginatedUsers" :key="user.id" class="group">
            <TableCell class="w-10">
              <input
                type="checkbox"
                class="checkbox checkbox-sm checkbox-primary"
                :checked="selectedIds.includes(user.id)"
                @change="handleToggleRowSelection(user.id, $event)"
              />
            </TableCell>
            <TableCell>
              <div class="flex items-center">
                <div class="h-10 w-10 flex-shrink-0 bg-primary/10 rounded-full flex items-center justify-center text-primary group-hover:scale-110 transition-transform">
                  <UserIcon class="w-5 h-5" />
                </div>
                <div class="ml-4">
                  <div class="text-sm font-bold text-foreground">{{ user.name }}</div>
                  <div class="text-xs text-muted-foreground mt-0.5">{{ user.phone }}</div>
                  <div class="flex flex-wrap gap-1 mt-1">
                    <Badge 
                      v-for="tag in user.tags" 
                      :key="tag"
                      variant="outline"
                      class="px-1.5 py-0 text-[10px]"
                    >
                      {{ tag }}
                    </Badge>
                  </div>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <div class="flex flex-col items-start gap-1">
                <Badge :class="getLevelClass(user.level)">
                  {{ getLevelText(user.level) }}
                </Badge>
                <span v-if="user.discount && user.discount < 1" class="text-xs text-destructive font-medium px-2">
                  {{ (user.discount * 10).toFixed(1) }}折
                </span>
              </div>
            </TableCell>
            <TableCell>
              <div class="space-y-1 text-xs">
                <div class="flex items-center gap-2">
                  <span class="inline-flex items-center rounded-full bg-emerald-50 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300 px-2 py-0.5 text-[11px]">
                    余额 {{ formatPrice(user.balance) }}
                  </span>
                  <Button variant="ghost" size="sm" class="h-6 px-2 text-[11px]" @click="openRechargeModal(user)">
                    <Zap class="w-3 h-3 mr-1" />
                    充值
                  </Button>
                </div>
                <div class="flex items-center gap-2 text-muted-foreground">
                  <span>积分 {{ user.points }}</span>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <div class="text-sm font-semibold">
                {{ formatPrice(user.totalConsumed) }}
              </div>
              <div class="text-xs text-muted-foreground">
                累计消费
              </div>
            </TableCell>
            <TableCell>
              <div class="text-sm">
                {{ formatDate(user.lastVisit) }}
              </div>
              <div class="text-xs text-muted-foreground">
                最近到店
              </div>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  variant="outline"
                  size="sm"
                  class="h-7 px-2 text-xs"
                  @click="openEditModal(user)"
                >
                  编辑
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  class="h-7 px-2 text-xs text-destructive"
                  @click="openRowDeleteDialog(user.id)"
                >
                  删除
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableRow v-if="paginatedUsers.length === 0">
            <TableCell colspan="7" class="text-center py-10 text-sm text-muted-foreground">
              暂无会员数据，请先新增会员。
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-4 py-3 border-t border-border bg-muted/30 text-xs text-muted-foreground">
        <div>
          共 {{ totalUsers }} 条会员记录，当前显示
          <span class="font-medium text-foreground">{{ pageStart }}-{{ pageEnd }}</span>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2 text-xs text-muted-foreground">
            <span>每页</span>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  class="h-8 px-2 text-xs font-normal"
                >
                  {{ pageSize }} 条
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" side="top">
                <DropdownMenuItem
                  v-for="size in pageSizeOptions"
                  :key="size"
                  @click="pageSize = size"
                  :class="size === pageSize ? 'bg-accent text-accent-foreground' : ''"
                >
                  {{ size }} 条 / 页
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="page === 1 || totalUsers === 0"
              @click="page = page - 1"
            >
              上一页
            </Button>
            <div class="flex items-center gap-1">
              <button
                v-for="p in totalPages"
                :key="p"
                type="button"
                class="px-2 py-1 rounded-md text-xs"
                :class="p === page ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'"
                @click="page = p"
              >
                {{ p }}
              </button>
            </div>
            <Button
              variant="outline"
              size="sm"
              :disabled="page === totalPages || totalUsers === 0"
              @click="page = page + 1"
            >
              下一页
            </Button>
            <div class="flex items-center gap-1">
              <span>跳至</span>
              <Input
                v-model="pageInput"
                class="w-12 h-8 px-1 text-xs"
                @keydown.enter.prevent="handleJumpToPage"
                @blur="handleJumpToPage"
              />
              <span>页</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Member Modal -->
    <Dialog v-model:open="isEditModalOpen">
      <DialogContent class="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>{{ editingUser.id ? '编辑会员' : '新增会员' }}</DialogTitle>
          <DialogDescription>
            {{ editingUser.id ? '修改会员基础信息。' : '填写新会员的基础信息。' }}
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">姓名</Label>
            <Input v-model="editingUser.name" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">手机号</Label>
            <Input v-model="editingUser.phone" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">等级</Label>
            <div class="col-span-3">
              <Select v-model="editingUser.level">
                <SelectTrigger>
                  <SelectValue placeholder="选择等级" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="diamond">💎 钻石会员</SelectItem>
                  <SelectItem value="gold">🥇 金牌会员</SelectItem>
                  <SelectItem value="silver">🥈 银牌会员</SelectItem>
                  <SelectItem value="normal">🥉 普通会员</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">折扣</Label>
            <Input
              v-model.number="editingUser.discount"
              type="number"
              step="0.1"
              min="0.1"
              max="1"
              class="col-span-3"
            />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">备注标签</Label>
            <Input
              v-model="currentTag"
              placeholder="如：常客、团长、大客户等"
              class="col-span-3"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="closeEditModal">
            取消
          </Button>
          <Button @click="saveUser">
            保存
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Recharge Modal -->
    <Dialog v-model:open="isRechargeModalOpen">
      <DialogContent class="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle>会员充值</DialogTitle>
          <DialogDescription>
            为会员余额账户充值，可用于消费抵扣。
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">会员</Label>
            <div class="col-span-3 text-sm font-medium">
              {{ editingUser.name }}
            </div>
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">当前余额</Label>
            <div class="col-span-3 text-sm text-muted-foreground">
              {{ formatPrice(editingUser.balance || 0) }}
            </div>
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label class="text-right text-xs">充值金额</Label>
            <Input
              v-model.number="rechargeAmount"
              type="number"
              min="1"
              class="col-span-3"
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="closeRechargeModal">
            取消
          </Button>
          <Button @click="handleRecharge">
            确认充值
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirm Dialog -->
    <ConfirmDialog
      v-model:open="deleteDialogOpen"
      title="删除会员"
      :description="pendingDeleteIds.length > 1 ? `确认删除选中的 ${pendingDeleteIds.length} 个会员吗？` : '确认删除该会员吗？'"
      confirm-text="确认删除"
      cancel-text="取消"
      @confirm="handleConfirmDelete"
    />
  </div>
</template>
