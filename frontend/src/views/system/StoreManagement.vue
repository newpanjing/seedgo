<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Search, Building2, Phone, MapPin, Users } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
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
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'

type Store = {
  id: number
  name: string
  phone: string
  address: string
  subAccountCount: number
  status: 'open' | 'closed'
}

const search = ref('')
const stores = ref<Store[]>([
  {
    id: 1,
    name: '主门店（总店）',
    phone: '021-66668888',
    address: '上海市徐汇区漕河泾科技园 A 座 8 楼',
    subAccountCount: 8,
    status: 'open',
  },
  {
    id: 2,
    name: '人民广场门店',
    phone: '021-66669999',
    address: '上海市黄浦区南京东路 399 号',
    subAccountCount: 5,
    status: 'open',
  },
  {
    id: 3,
    name: '测试门店（待开业）',
    phone: '021-00000000',
    address: '待定',
    subAccountCount: 0,
    status: 'closed',
  },
])

const isDialogOpen = ref(false)
const editing = ref<Store | null>(null)
const formName = ref('')
const formPhone = ref('')
const formAddress = ref('')

const filteredStores = computed(() => {
  if (!search.value) return stores.value
  return stores.value.filter((item) => {
    return (
      item.name.includes(search.value) ||
      item.phone.includes(search.value) ||
      item.address.includes(search.value)
    )
  })
})

const openCreate = () => {
  editing.value = null
  formName.value = ''
  formPhone.value = ''
  formAddress.value = ''
  isDialogOpen.value = true
}

const openEdit = (store: Store) => {
  editing.value = store
  formName.value = store.name
  formPhone.value = store.phone
  formAddress.value = store.address
  isDialogOpen.value = true
}

const handleSave = () => {
  if (!formName.value) {
    return
  }

  if (editing.value) {
    stores.value = stores.value.map((item) =>
      item.id === editing.value?.id
        ? {
            ...item,
            name: formName.value,
            phone: formPhone.value,
            address: formAddress.value,
          }
        : item,
    )
  } else {
    const nextId =
      stores.value.length > 0
        ? Math.max(...stores.value.map((i) => i.id)) + 1
        : 1

    stores.value.push({
      id: nextId,
      name: formName.value,
      phone: formPhone.value,
      address: formAddress.value,
      subAccountCount: 0,
      status: 'open',
    })
  }

  isDialogOpen.value = false
}

const toggleStatus = (store: Store) => {
  stores.value = stores.value.map((item) =>
    item.id === store.id
      ? {
          ...item,
          status: item.status === 'open' ? 'closed' : 'open',
        }
      : item,
  )
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">门店管理</h1>
        <p class="text-sm text-muted-foreground mt-1">
          管理多门店基础信息，并与子账号进行绑定关系维护。
        </p>
      </div>
      <Button @click="openCreate">
        <Plus class="w-4 h-4 mr-2" />
        新增门店
      </Button>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border p-4">
      <div class="flex flex-col md:flex-row gap-4 md:items-center md:justify-between">
        <div class="relative max-w-sm w-full">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            v-model="search"
            placeholder="按门店名称、电话或地址搜索..."
            class="pl-9"
          />
        </div>

        <div class="flex items-center gap-3 text-xs text-muted-foreground">
          <div class="flex items-center gap-1">
            <span class="inline-block h-2 w-2 rounded-full bg-emerald-500" />
            营业中
          </div>
          <div class="flex items-center gap-1">
            <span class="inline-block h-2 w-2 rounded-full bg-rose-500" />
            已停用/未开业
          </div>
        </div>
      </div>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
      <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/30">
        <div class="flex items-center gap-2">
          <Button size="sm" @click="openCreate">
            新增
          </Button>
          <Button
            size="sm"
            variant="default"
            disabled
            class="bg-emerald-600 text-emerald-50 hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            编辑
          </Button>
          <Button
            size="sm"
            variant="destructive"
            disabled
            class="disabled:opacity-50 disabled:cursor-not-allowed"
          >
            删除
          </Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>门店信息</TableHead>
            <TableHead>联系方式</TableHead>
            <TableHead>已绑定子账号数</TableHead>
            <TableHead>状态</TableHead>
            <TableHead class="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="item in filteredStores"
            :key="item.id"
            class="group"
          >
            <TableCell>
              <div class="flex items-start gap-3">
                <div class="h-9 w-9 flex-shrink-0 rounded-lg bg-primary/10 text-primary flex items-center justify-center group-hover:scale-110 transition-transform">
                  <Building2 class="w-4 h-4" />
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-semibold">{{ item.name }}</span>
                    <Badge
                      v-if="item.id === 1"
                      variant="outline"
                    >
                      主门店
                    </Badge>
                  </div>
                  <p class="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                    <MapPin class="w-3 h-3" />
                    <span class="truncate max-w-xs">{{ item.address }}</span>
                  </p>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <div class="flex items-center gap-1 text-sm">
                <Phone class="w-4 h-4 text-muted-foreground" />
                <span>{{ item.phone }}</span>
              </div>
            </TableCell>
            <TableCell>
              <div class="flex items-center gap-1 text-sm">
                <Users class="w-4 h-4 text-muted-foreground" />
                <span>{{ item.subAccountCount }}</span>
              </div>
            </TableCell>
            <TableCell>
              <Badge
                variant="outline"
                :color="item.status === 'open' ? 'green' : 'red'"
              >
                {{ item.status === 'open' ? '营业中' : '已停用' }}
              </Badge>
            </TableCell>
            <TableCell class="text-right space-x-1">
              <Button
                variant="ghost"
                size="sm"
                class="text-xs"
                @click="openEdit(item)"
              >
                编辑
              </Button>
              <Button
                variant="ghost"
                size="sm"
                class="text-xs text-muted-foreground"
                @click="toggleStatus(item)"
              >
                {{ item.status === 'open' ? '停用' : '启用' }}
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <Dialog v-model:open="isDialogOpen">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{{ editing ? '编辑门店' : '新增门店' }}</DialogTitle>
        </DialogHeader>

        <div class="space-y-4 py-2">
          <div class="space-y-2">
            <Label for="store-name">门店名称</Label>
            <Input
              id="store-name"
              v-model="formName"
              placeholder="例如：人民广场门店"
            />
          </div>

          <div class="space-y-2">
            <Label for="store-phone">联系电话</Label>
            <Input
              id="store-phone"
              v-model="formPhone"
              placeholder="用于对外展示与内部联系"
            />
          </div>

          <div class="space-y-2">
            <Label for="store-address">地址</Label>
            <Input
              id="store-address"
              v-model="formAddress"
              placeholder="方便顾客导航与内部管理"
            />
          </div>
        </div>

        <DialogFooter class="mt-2">
          <Button variant="outline" @click="isDialogOpen = false">
            取消
          </Button>
          <Button @click="handleSave">
            {{ editing ? '保存修改' : '创建门店' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
