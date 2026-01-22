<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProductStore } from '../stores/products'
import { useBillingStore } from '../stores/billing'
import type { Product } from '../types/product'
import type { Room } from '../types/billing'
import { 
  Plus, 
  Search, 
  Filter, 
  Edit, 
  Trash2, 
  Package, 
  Tag, 
  Clock, 
  Box,
  Armchair,
  LayoutGrid,
  Gamepad2,
  CircleDot
} from 'lucide-vue-next'
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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'

const store = useProductStore()
const billingStore = useBillingStore()
const activeTab = ref<string>('products')
const searchQuery = ref('')
const selectedCategory = ref<string>('all')

// Room Management State
const isRoomModalOpen = ref(false)
const isDeleteRoomModalOpen = ref(false)
const editingRoom = ref<Partial<Room>>({})
const deletingRoom = ref<Room | null>(null)

// Room Actions
const openRoomModal = (room?: Room) => {
  editingRoom.value = room ? { ...room } : { type: 'chess', status: 'available', hourlyRate: 0 }
  isRoomModalOpen.value = true
}

const closeRoomModal = () => {
  isRoomModalOpen.value = false
  editingRoom.value = {}
}

const saveRoom = () => {
  if (editingRoom.value.id) {
    billingStore.updateRoom(editingRoom.value.id, editingRoom.value)
  } else if (editingRoom.value.name && editingRoom.value.type && editingRoom.value.hourlyRate) {
    billingStore.addRoom({
      name: editingRoom.value.name,
      type: editingRoom.value.type,
      hourlyRate: editingRoom.value.hourlyRate
    })
  }
  closeRoomModal()
}

const openDeleteRoomModal = (room: Room) => {
  deletingRoom.value = room
  isDeleteRoomModalOpen.value = true
}

const closeDeleteRoomModal = () => {
  isDeleteRoomModalOpen.value = false
  deletingRoom.value = null
}

const confirmDeleteRoom = () => {
  if (deletingRoom.value) {
    billingStore.deleteRoom(deletingRoom.value.id)
    closeDeleteRoomModal()
  }
}

const getRoomIcon = (type: string) => {
  switch (type) {
    case 'chess': return Gamepad2
    case 'billiards': return CircleDot
    case 'mahjong': return LayoutGrid
    default: return Armchair
  }
}

const getRoomTypeText = (type: string) => {
  const map: Record<string, string> = {
    chess: '棋牌室',
    mahjong: '麻将房',
    billiards: '台球桌'
  }
  return map[type] || type
}

// Modal States
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const editingProduct = ref<Partial<Product>>({})
const deletingProduct = ref<Product | null>(null)

// Actions
const openEditModal = (product: Product) => {
  editingProduct.value = { ...product }
  isEditModalOpen.value = true
}

const closeEditModal = () => {
  isEditModalOpen.value = false
  editingProduct.value = {}
}

const saveEdit = () => {
  if (editingProduct.value.id) {
    store.updateProduct(editingProduct.value.id, editingProduct.value)
    closeEditModal()
  }
}

const openDeleteModal = (product: Product) => {
  deletingProduct.value = product
  isDeleteModalOpen.value = true
}

const closeDeleteModal = () => {
  isDeleteModalOpen.value = false
  deletingProduct.value = null
}

const confirmDelete = () => {
  if (deletingProduct.value) {
    store.deleteProduct(deletingProduct.value.id)
    closeDeleteModal()
  }
}

const filteredProducts = computed(() => {
  return store.products.filter(product => {
    const matchesSearch = product.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = selectedCategory.value === 'all' || product.categoryId === selectedCategory.value
    return matchesSearch && matchesCategory
  })
})

const getCategoryName = (id: string) => store.getCategoryName(id)

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency: 'CNY' }).format(price)
}

const getStatusColor = (status: Room['status']) => {
  switch (status) {
    case 'available': return 'green'
    case 'occupied': return 'red'
    case 'cleaning': return 'yellow'
    case 'maintenance': return 'gray'
    default: return 'gray'
  }
}

const getStatusText = (status: Room['status']) => {
  const map: Record<string, string> = {
    available: '空闲',
    occupied: '使用中',
    cleaning: '待清理',
    maintenance: '维护中'
  }
  return map[status] || status
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-card p-6 rounded-xl shadow-sm border border-border">
      <div>
        <h1 class="text-2xl font-bold text-foreground flex items-center">
          <Package class="w-8 h-8 mr-3 text-primary" />
          {{ activeTab === 'products' ? '商品管理' : '计费区管理' }}
        </h1>
        <p class="mt-1 text-sm text-muted-foreground">
          {{ activeTab === 'products' ? '管理您的商品库存、价格及分类信息' : '管理棋牌室、包厢及台球桌等计费资源' }}
        </p>
      </div>
      
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full sm:w-auto">
        <Button 
          v-if="activeTab === 'products'"
          @click="() => {} /* Implement Add Product */"
        >
          <Plus class="w-4 h-4 mr-2" />
          新增商品
        </Button>
        <Button 
          v-else
          @click="openRoomModal()"
        >
          <Plus class="w-4 h-4 mr-2" />
          新增房间
        </Button>
      </div>
    </div>

    <Tabs v-model="activeTab" class="w-full">
      <TabsList>
        <TabsTrigger value="products">商品管理</TabsTrigger>
        <TabsTrigger value="zones">计费区管理</TabsTrigger>
      </TabsList>

      <TabsContent value="products" class="space-y-6 mt-6">
        <!-- Filters & Toolbar -->
        <div class="bg-card p-4 rounded-xl shadow-sm border border-border flex flex-col md:flex-row gap-4 items-center justify-between">
          <div class="flex flex-1 w-full gap-4">
            <!-- Search -->
            <div class="relative flex-1 max-w-md">
              <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input 
                v-model="searchQuery"
                type="text" 
                placeholder="搜索商品名称、编号..." 
                class="pl-9"
              />
            </div>
            
            <!-- Category Filter -->
            <div class="w-[180px]">
              <Select v-model="selectedCategory">
                <SelectTrigger>
                  <div class="flex items-center">
                    <Filter class="w-4 h-4 mr-2 text-muted-foreground" />
                    <SelectValue placeholder="选择分类" />
                  </div>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">所有分类</SelectItem>
                  <SelectItem v-for="cat in store.categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          
          <!-- Stats Summary -->
          <div class="hidden md:flex items-center space-x-4 text-sm text-muted-foreground">
            <span class="flex items-center">
              <Box class="w-4 h-4 mr-1" />
              共 {{ store.products.length }} 件商品
            </span>
          </div>
        </div>

        <!-- Product List -->
        <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
          <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/30">
            <div class="flex items-center gap-2">
              <Button
                size="sm"
                @click="() => {} /* Implement Add Product */"
              >
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
                <TableHead>商品信息</TableHead>
                <TableHead>分类</TableHead>
                <TableHead>价格/方式</TableHead>
                <TableHead>库存状态</TableHead>
                <TableHead class="text-right">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="product in filteredProducts" :key="product.id" class="group">
                <TableCell>
                  <div class="flex items-center">
                    <div class="h-12 w-12 flex-shrink-0 bg-primary/10 rounded-lg flex items-center justify-center text-2xl shadow-sm group-hover:scale-105 transition-transform">
                      {{ product.imageUrl ? '' : '📦' }}
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-bold text-foreground">{{ product.name }}</div>
                      <div class="text-xs text-muted-foreground mt-0.5">ID: {{ product.id.slice(0, 8) }}</div>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">
                    <Tag class="w-3 h-3 mr-1" />
                    {{ getCategoryName(product.categoryId) }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div class="flex flex-col">
                    <span class="text-sm font-bold text-foreground">{{ formatPrice(product.price) }}</span>
                    <span class="text-xs text-muted-foreground flex items-center mt-0.5">
                      <Clock class="w-3 h-3 mr-1" v-if="product.pricingType === 'hourly'" />
                      {{ product.pricingType === 'hourly' ? '/小时' : (product.pricingType === 'package' ? '套餐' : '固定价格') }}
                    </span>
                  </div>
                </TableCell>
                <TableCell>
                  <div class="flex flex-col">
                     <div class="flex items-center">
                        <span :class="`text-sm font-medium ${product.stock < 10 ? 'text-orange-600' : 'text-foreground'}`">
                          {{ product.stock }}
                        </span>
                        <span class="text-xs text-muted-foreground ml-1">件</span>
                     </div>
                     <div class="mt-1 flex items-center">
                       <span class="relative flex h-2 w-2 mr-2">
                          <span v-if="product.status === 'active'" class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                          <span :class="`relative inline-flex rounded-full h-2 w-2 ${product.status === 'active' ? 'bg-green-500' : 'bg-red-500'}`"></span>
                        </span>
                        <span :class="`text-xs ${product.status === 'active' ? 'text-green-600' : 'text-red-600'}`">
                          {{ product.status === 'active' ? '销售中' : '已停用' }}
                        </span>
                     </div>
                  </div>
                </TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button variant="ghost" size="icon" @click="openEditModal(product)" title="编辑">
                      <Edit class="w-4 h-4 text-primary" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="openDeleteModal(product)" title="删除">
                      <Trash2 class="w-4 h-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
              <TableRow v-if="filteredProducts.length === 0">
                <TableCell colspan="5" class="h-32 text-center">
                  <div class="flex flex-col items-center justify-center text-muted-foreground">
                    <Search class="w-8 h-8 mb-2 opacity-20" />
                    <p>未找到相关商品</p>
                    <p class="text-xs mt-1">请尝试调整搜索关键词或筛选条件</p>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </TabsContent>

      <TabsContent value="zones" class="space-y-6 mt-6">
         <!-- Zone List -->
         <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
           <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/30">
             <div class="flex items-center gap-2">
               <Button size="sm" @click="openRoomModal()">
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
                 <TableHead>房间信息</TableHead>
                 <TableHead>类型</TableHead>
                 <TableHead>费率 (元/小时)</TableHead>
                 <TableHead>状态</TableHead>
                 <TableHead class="text-right">操作</TableHead>
               </TableRow>
             </TableHeader>
             <TableBody>
               <TableRow v-for="room in billingStore.rooms" :key="room.id" class="group">
                 <TableCell>
                   <div class="flex items-center">
                     <div class="h-10 w-10 flex-shrink-0 bg-primary/10 rounded-lg flex items-center justify-center text-primary">
                       <component :is="getRoomIcon(room.type)" class="w-5 h-5" />
                     </div>
                     <div class="ml-4">
                       <div class="text-sm font-bold text-foreground">{{ room.name }}</div>
                       <div class="text-xs text-muted-foreground mt-0.5">ID: {{ room.id.slice(0, 8) }}</div>
                     </div>
                   </div>
                 </TableCell>
                 <TableCell>
                   <Badge variant="outline">
                     {{ getRoomTypeText(room.type) }}
                   </Badge>
                 </TableCell>
                 <TableCell>
                   <span class="text-sm font-bold text-foreground">¥{{ room.hourlyRate }}</span>
                 </TableCell>
                 <TableCell>
                   <Badge :color="getStatusColor(room.status)">
                     {{ getStatusText(room.status) }}
                   </Badge>
                 </TableCell>
                 <TableCell class="text-right">
                   <div class="flex items-center justify-end space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                     <Button variant="ghost" size="icon" @click="openRoomModal(room)" title="编辑">
                       <Edit class="w-4 h-4 text-primary" />
                     </Button>
                     <Button variant="ghost" size="icon" @click="openDeleteRoomModal(room)" title="删除">
                       <Trash2 class="w-4 h-4 text-destructive" />
                     </Button>
                   </div>
                 </TableCell>
               </TableRow>
             </TableBody>
           </Table>
         </div>
      </TabsContent>
    </Tabs>

    <!-- Edit Product Dialog -->
    <Dialog :open="isEditModalOpen" @update:open="isEditModalOpen = $event">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>编辑商品</DialogTitle>
          <DialogDescription>
            修改商品信息，点击保存提交更改。
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="name" class="text-right">商品名称</Label>
            <Input id="name" v-model="editingProduct.name" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="price" class="text-right">价格</Label>
            <Input id="price" v-model.number="editingProduct.price" type="number" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="stock" class="text-right">库存</Label>
            <Input id="stock" v-model.number="editingProduct.stock" type="number" class="col-span-3" />
          </div>
          <!-- Add more fields as needed -->
        </div>
        <DialogFooter>
          <Button variant="outline" @click="closeEditModal">取消</Button>
          <Button @click="saveEdit">保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Product Dialog -->
    <Dialog :open="isDeleteModalOpen" @update:open="isDeleteModalOpen = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>确认删除</DialogTitle>
          <DialogDescription>
            您确定要删除商品 "{{ deletingProduct?.name }}" 吗？此操作无法撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="closeDeleteModal">取消</Button>
          <Button variant="destructive" @click="confirmDelete">删除</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Room Dialog -->
    <Dialog :open="isRoomModalOpen" @update:open="isRoomModalOpen = $event">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{{ editingRoom.id ? '编辑房间' : '新增房间' }}</DialogTitle>
          <DialogDescription>
            {{ editingRoom.id ? '修改房间信息。' : '填写新房间的信息。' }}
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="room-name" class="text-right">房间名称</Label>
            <Input id="room-name" v-model="editingRoom.name" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="room-type" class="text-right">类型</Label>
            <div class="col-span-3">
              <Select v-model="editingRoom.type">
                <SelectTrigger>
                  <SelectValue placeholder="选择类型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="chess">棋牌室</SelectItem>
                  <SelectItem value="mahjong">麻将房</SelectItem>
                  <SelectItem value="billiards">台球桌</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="hourly-rate" class="text-right">费率 (元/小时)</Label>
            <Input id="hourly-rate" v-model.number="editingRoom.hourlyRate" type="number" class="col-span-3" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="closeRoomModal">取消</Button>
          <Button @click="saveRoom">保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Room Dialog -->
    <Dialog :open="isDeleteRoomModalOpen" @update:open="isDeleteRoomModalOpen = $event">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>确认删除房间</DialogTitle>
          <DialogDescription>
            您确定要删除房间 "{{ deletingRoom?.name }}" 吗？此操作无法撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="closeDeleteRoomModal">取消</Button>
          <Button variant="destructive" @click="confirmDeleteRoom">删除</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
