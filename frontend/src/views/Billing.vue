<script setup lang="ts">
import { ref, computed } from 'vue'
import { useBillingStore } from '../stores/billing'
import { useProductStore } from '../stores/products'
import { 
  Coffee, 
  Play, 
  Square, 
  RotateCcw, 
  Timer,
  Clock,
  Armchair,
  Gamepad2,
  CircleDot,
  LayoutGrid,
  Search,
  Filter,
  Utensils
} from 'lucide-vue-next'
import type { Room } from '../types/billing'
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
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

const billingStore = useBillingStore()
const productStore = useProductStore()

const selectedRoom = ref<Room | null>(null)
const showOrderModal = ref(false)
const selectedProduct = ref('')
const orderQuantity = ref(1)
const filterStatus = ref('all')
const searchQuery = ref('')

const formatTime = (isoString?: string) => {
  if (!isoString) return '--:--'
  return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const formatDuration = (minutes: number) => {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return `${h}小时 ${m}分钟`
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

const getRoomIcon = (type: string) => {
  switch (type) {
    case 'chess': return Gamepad2
    case 'billiards': return CircleDot
    case 'mahjong': return LayoutGrid
    default: return Armchair
  }
}

// Stats & Filters
const statusCounts = computed(() => {
  const counts = { all: billingStore.rooms.length, available: 0, occupied: 0, cleaning: 0, maintenance: 0 }
  billingStore.rooms.forEach(r => {
    if (counts[r.status] !== undefined) counts[r.status]++
  })
  return counts
})

const filteredRooms = computed(() => {
  return billingStore.rooms.filter(room => {
    const matchesStatus = filterStatus.value === 'all' || room.status === filterStatus.value
    const matchesSearch = room.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesStatus && matchesSearch
  })
})

const openOrderModal = (room: Room) => {
  selectedRoom.value = room
  showOrderModal.value = true
  selectedProduct.value = ''
  orderQuantity.value = 1
}

const submitOrder = () => {
  if (selectedRoom.value?.currentSessionId && selectedProduct.value) {
    billingStore.addOrderToSession(selectedRoom.value.currentSessionId, selectedProduct.value, orderQuantity.value)
    showOrderModal.value = false
  }
}

const showCheckoutModal = ref(false)
const checkoutForm = ref({
  roomCharge: 0,
  ordersTotal: 0,
  subtotal: 0,
  discountRate: 10, // 10折 = 100%
  rounding: 0,
  finalTotal: 0,
  durationMinutes: 0
})

const openCheckoutModal = (room: Room) => {
  selectedRoom.value = room
  const session = billingStore.getSessionByRoomId(room.id)
  if (!session) return

  // Calculate duration and charges
  const now = new Date().getTime()
  const start = new Date(session.startTime).getTime()
  const durationMinutes = Math.floor((now - start) / (1000 * 60))
  const hours = durationMinutes / 60
  const roomCharge = Math.ceil(hours * room.hourlyRate)
  const ordersTotal = session.orders.reduce((sum, order) => sum + order.total, 0)
  
  checkoutForm.value = {
    roomCharge,
    ordersTotal,
    subtotal: roomCharge + ordersTotal,
    discountRate: 10,
    rounding: 0,
    finalTotal: roomCharge + ordersTotal,
    durationMinutes
  }
  
  showCheckoutModal.value = true
}

const updateFinalTotal = () => {
  const { subtotal, discountRate, rounding } = checkoutForm.value
  const discounted = subtotal * (discountRate / 10)
  checkoutForm.value.finalTotal = Math.max(0, Math.floor(discounted - rounding))
}

const confirmCheckout = () => {
  if (selectedRoom.value) {
    billingStore.endSession(selectedRoom.value.id, {
      totalAmount: checkoutForm.value.finalTotal,
      discount: checkoutForm.value.subtotal - (checkoutForm.value.subtotal * (checkoutForm.value.discountRate / 10)),
      rounding: checkoutForm.value.rounding
    })
    showCheckoutModal.value = false
  }
}

// Consumable products only
const consumables = computed(() => {
  return productStore.products.filter(p => p.pricingType === 'fixed')
})

const filters = [
  { id: 'all', label: '全部' },
  { id: 'available', label: '空闲' },
  { id: 'occupied', label: '使用中' },
  { id: 'cleaning', label: '待清理' },
  { id: 'maintenance', label: '维护中' },
]
</script>

<template>
  <div class="space-y-6">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-card p-6 rounded-xl shadow-sm border border-border">
      <div>
        <h1 class="text-2xl font-bold text-foreground flex items-center">
          <Timer class="w-8 h-8 mr-3 text-primary" />
          计时计费
        </h1>
        <p class="mt-1 text-sm text-muted-foreground">
          实时管理房间状态、计费及消费订单
        </p>
      </div>
      
      <div class="flex items-center space-x-2 text-sm">
        <Badge variant="outline" color="green">
          空闲: {{ statusCounts.available }}
        </Badge>
        <Badge color="red">
          使用中: {{ statusCounts.occupied }}
        </Badge>
      </div>
    </div>

    <!-- Filters & Search -->
    <div class="bg-card p-2 rounded-xl shadow-sm border border-border flex flex-col md:flex-row justify-between items-center gap-4">
      <!-- Tabs -->
      <Tabs v-model="filterStatus" class="w-full md:w-auto">
        <TabsList>
          <TabsTrigger v-for="filter in filters" :key="filter.id" :value="filter.id">
            {{ filter.label }}
            <span class="ml-1 text-xs opacity-70">
              ({{ statusCounts[filter.id as keyof typeof statusCounts] }})
            </span>
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <!-- Search -->
      <div class="relative w-full md:w-64 mr-2">
        <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
        <Input 
          v-model="searchQuery"
          type="text" 
          placeholder="搜索房间..." 
          class="pl-9"
        />
      </div>
    </div>

    <!-- Room Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <Card 
        v-for="room in filteredRooms" 
        :key="room.id"
        class="overflow-hidden hover:shadow-md transition-all duration-300"
      >
        <!-- Card Header -->
        <div class="p-5 border-b border-border">
          <div class="flex justify-between items-start">
            <div class="flex items-center">
              <div class="p-2.5 rounded-lg bg-primary/10 text-primary mr-3">
                <component :is="getRoomIcon(room.type)" class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-bold text-foreground">{{ room.name }}</h3>
                <p class="text-xs text-muted-foreground mt-0.5">
                  ¥{{ room.hourlyRate }}/小时
                </p>
              </div>
            </div>
            <Badge :color="getStatusColor(room.status)">
              {{ getStatusText(room.status) }}
            </Badge>
          </div>
        </div>

        <!-- Card Body -->
        <CardContent class="p-5">
          <div v-if="room.status === 'occupied'" class="space-y-4">
            <!-- Timer Display -->
            <div class="text-center py-2">
              <div class="text-3xl font-mono font-bold text-gray-900 dark:text-white tracking-wider">
                {{ formatDuration(billingStore.getSessionByRoomId(room.id)?.durationMinutes || 0).split(' ')[0] }}
                <span class="text-base font-normal text-gray-400">:</span>
                {{ formatDuration(billingStore.getSessionByRoomId(room.id)?.durationMinutes || 0).split(' ')[1] }}
              </div>
              <p class="text-xs text-gray-400 mt-1">
                开始于 {{ formatTime(billingStore.getSessionByRoomId(room.id)?.startTime) }}
              </p>
            </div>

            <!-- Stats Row -->
            <div class="grid grid-cols-2 gap-3 text-center bg-muted rounded-lg p-3">
              <div>
                <p class="text-xs text-muted-foreground">当前消费</p>
                <p class="font-bold text-primary">
                  ¥{{ billingStore.getSessionByRoomId(room.id)?.orders.reduce((acc, o) => acc + o.total, 0) || 0 }}
                </p>
              </div>
              <div>
                <p class="text-xs text-muted-foreground">预计房费</p>
                <p class="font-bold text-foreground">
                  ¥{{ Math.ceil(((billingStore.getSessionByRoomId(room.id)?.durationMinutes || 0)/60) * room.hourlyRate) }}
                </p>
              </div>
            </div>
          </div>

          <div v-else class="flex flex-col items-center justify-center h-[140px] text-muted-foreground bg-muted/50 rounded-lg border border-dashed border-border">
            <component :is="room.status === 'available' ? Armchair : (room.status === 'cleaning' ? RotateCcw : Clock)" class="w-8 h-8 mb-2 opacity-50" />
            <span class="text-sm">{{ room.status === 'available' ? '等待入住' : (room.status === 'cleaning' ? '等待清理' : '暂不可用') }}</span>
          </div>
        </CardContent>

        <!-- Card Footer (Actions) -->
        <CardFooter class="p-4 pt-0">
          <div class="flex gap-2 w-full">
            <Button 
              v-if="room.status === 'available'"
              @click="billingStore.startSession(room.id)"
              class="w-full shadow-lg shadow-primary/30"
            >
              <Play class="w-4 h-4 mr-2 fill-current" />
              开始计费
            </Button>
            
            <template v-if="room.status === 'occupied'">
              <Button 
                variant="outline"
                @click="openOrderModal(room)"
                class="flex-1 border-primary text-primary bg-primary/5 hover:bg-primary/10"
              >
                <Utensils class="w-4 h-4 mr-1.5" />
                加单
              </Button>
              <Button 
                variant="outline"
                @click="openCheckoutModal(room)"
                class="flex-1"
              >
                <Square class="w-4 h-4 mr-1.5 fill-current" />
                结账
              </Button>
            </template>

            <Button 
              v-if="room.status === 'cleaning'"
              @click="billingStore.setRoomAvailable(room.id)"
              class="w-full shadow-lg shadow-primary/30"
            >
              <RotateCcw class="w-4 h-4 mr-2" />
              完成清理
            </Button>
          </div>
        </CardFooter>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-if="filteredRooms.length === 0" class="text-center py-12 bg-card rounded-xl border border-border border-dashed">
      <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
        <Filter class="w-8 h-8 text-muted-foreground" />
      </div>
      <h3 class="text-lg font-medium text-foreground">未找到相关房间</h3>
      <p class="text-muted-foreground mt-1">请尝试切换筛选状态或搜索关键词</p>
    </div>

    <!-- Order Modal -->
    <Dialog :open="showOrderModal" @update:open="showOrderModal = $event">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle class="flex items-center">
            <Coffee class="w-5 h-5 mr-2 text-primary" />
            添加消费 - {{ selectedRoom?.name }}
          </DialogTitle>
        </DialogHeader>
        
        <div class="space-y-5 py-4">
          <div class="space-y-2">
            <Label>选择商品</Label>
            <Select v-model="selectedProduct">
              <SelectTrigger>
                <SelectValue placeholder="请选择商品" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="p in consumables" :key="p.id" :value="p.id">
                  {{ p.name }} - ¥{{ p.price }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          
          <div class="space-y-2">
            <Label>数量</Label>
            <div class="flex items-center gap-2">
              <Button variant="outline" size="icon" @click="orderQuantity > 1 && orderQuantity--">
                -
              </Button>
              <Input 
                v-model.number="orderQuantity"
                type="number" 
                min="1"
                class="text-center"
              />
              <Button variant="outline" size="icon" @click="orderQuantity++">
                +
              </Button>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="showOrderModal = false">取消</Button>
          <Button @click="submitOrder" :disabled="!selectedProduct">确认添加</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Checkout Modal -->
    <Dialog :open="showCheckoutModal" @update:open="showCheckoutModal = $event">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle class="flex items-center">
            <Square class="w-5 h-5 mr-2 text-red-600" />
            结账确认 - {{ selectedRoom?.name }}
          </DialogTitle>
        </DialogHeader>
        
        <div class="space-y-4 py-4">
          <!-- Bill Details -->
          <div class="space-y-3 text-sm">
            <div class="flex justify-between items-center text-muted-foreground">
              <span>使用时长</span>
              <span class="font-medium text-foreground">{{ formatDuration(checkoutForm.durationMinutes) }}</span>
            </div>
            <div class="flex justify-between items-center text-muted-foreground">
              <span>房费</span>
              <span class="font-medium text-foreground">¥{{ checkoutForm.roomCharge }}</span>
            </div>
            <div class="flex justify-between items-center text-muted-foreground">
              <span>消费金额</span>
              <span class="font-medium text-foreground">¥{{ checkoutForm.ordersTotal }}</span>
            </div>
            <div class="flex justify-between items-center pt-2 border-t font-bold">
              <span>小计</span>
              <span>¥{{ checkoutForm.subtotal }}</span>
            </div>
          </div>

          <!-- Adjustments -->
          <div class="space-y-4 pt-2">
            <div class="space-y-2">
              <Label>折扣 (折)</Label>
              <div class="flex items-center gap-2">
                <Input 
                  v-model.number="checkoutForm.discountRate"
                  @input="updateFinalTotal"
                  type="number" 
                  min="0"
                  max="10"
                  step="0.1"
                />
                <span class="text-sm text-muted-foreground whitespace-nowrap w-24 text-right">
                  -¥{{ (checkoutForm.subtotal * (1 - checkoutForm.discountRate/10)).toFixed(1) }}
                </span>
              </div>
            </div>

            <div class="space-y-2">
              <Label>抹零</Label>
              <div class="flex items-center gap-2">
                <Input 
                  v-model.number="checkoutForm.rounding"
                  @input="updateFinalTotal"
                  type="number" 
                  min="0"
                />
                <span class="text-sm text-muted-foreground whitespace-nowrap w-24 text-right">
                  -¥{{ checkoutForm.rounding }}
                </span>
              </div>
            </div>

            <div class="flex justify-between items-center pt-4 border-t">
              <span class="text-lg font-bold">实收金额</span>
              <span class="text-2xl font-bold text-primary">¥{{ checkoutForm.finalTotal }}</span>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="showCheckoutModal = false">取消</Button>
          <Button variant="destructive" @click="confirmCheckout">确认结账</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
