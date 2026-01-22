<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTheme } from '@/composables/useTheme'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { showConfirm } from '@/lib/message'
import ThemeSettings from '@/components/ThemeSettings.vue'
import { 
  Menu, 
  Home, 
  ChevronRight, 
  Search, 
  X, 
  Bell, 
  Sun, 
  Moon, 
  User, 
  UserPlus, 
  LogOut,
  Clock,
  Users,
  Package,
  Settings 
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'toggleSidebar'): void
}>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const isDark = useDark()
const toggleDark = useToggle(isDark)
const { initTheme } = useTheme()

const isThemeDialogOpen = ref(false)
const isSearchExpanded = ref(false)
const searchQuery = ref('')
const searchContainerRef = ref<HTMLElement | null>(null)

const routeTitleMap: Record<string, string> = {
  dashboard: '仪表盘',
  products: '商品/区域管理',
  zones: '区域管理',
  billing: '计时计费',
  analytics: '经营分析',
  members: '会员管理',
  tenantManagement: '租户管理',
  subAccounts: '账号管理',
  storeManagement: '门店管理',
  roleManagement: '角色管理',
  renewal: '续费管理',
  permissionManagement: '权限管理',
  menuManagement: '菜单管理',
  operationLog: '操作日志',
  loginLog: '登录日志',
  fileManagement: '文件管理',
  dictManagement: '字典管理',
  paramSettings: '参数设置',
  configSettings: '平台配置'
}

const currentRouteName = computed(() => {
  return routeTitleMap[route.name as string] || 'Seedgo'
})

const mockData = [
  { id: 1, title: '订单 #ORD-20240115-001', type: '订单', path: '/billing' },
  { id: 2, title: '会员 张三 (13800138000)', type: '会员', path: '/members' },
  { id: 3, title: '商品 可口可乐 330ml', type: '商品', path: '/products' },
  { id: 4, title: '系统设置 - 基础参数', type: '设置', path: '/system/params' },
  { id: 5, title: '门店 - 旗舰店', type: '门店', path: '/system/stores' },
]

const filteredSearchResults = computed(() => {
  if (!searchQuery.value) return []
  return mockData.filter(item => 
    item.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    item.type.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const handleSearchFocus = () => {
  isSearchExpanded.value = true
}

const handleClickOutside = (event: MouseEvent) => {
  if (searchContainerRef.value && !searchContainerRef.value.contains(event.target as Node)) {
    if (!searchQuery.value) {
      isSearchExpanded.value = false
    }
  }
}

const notifications = [
  {
    id: 1,
    title: '系统维护通知',
    content: '系统将于今晚凌晨 02:00 进行例行维护...',
    time: '10分钟前',
    read: false
  },
  {
    id: 2,
    title: '新订单提醒',
    content: '收到一个新的订单 #ORD-20240115-001...',
    time: '30分钟前',
    read: false
  },
  {
    id: 3,
    title: '会员注册成功',
    content: '新会员 张三 (13800138000) 已完成注册。',
    time: '1小时前',
    read: true
  }
]

const handleLogout = async () => {
  const confirmed = await showConfirm('确定要退出登录吗？', '退出后需要重新登录才能访问系统。')
  if (confirmed) {
    authStore.logout()
    router.push('/login')
  }
}

const handleSwitchAccount = async () => {
  const confirmed = await showConfirm('确定要切换账号吗？', '切换账号将退出当前登录状态。')
  if (confirmed) {
    authStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  initTheme()
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <header class="flex items-center justify-between h-16 px-4 sm:px-6 bg-card shadow-sm border-b border-border shrink-0 z-10">
    <!-- Left: Mobile Toggle & Breadcrumbs -->
    <div class="flex items-center flex-1">
      <Button 
        variant="ghost" 
        size="icon" 
        @click="emit('toggleSidebar')" 
        class="lg:hidden mr-4 text-muted-foreground hover:text-foreground"
      >
        <Menu class="w-6 h-6" />
      </Button>
      
      <nav class="hidden sm:flex items-center text-sm font-medium text-muted-foreground">
        <RouterLink to="/" class="hover:text-primary transition-colors">
          <Home class="w-4 h-4" />
        </RouterLink>
        <ChevronRight class="w-4 h-4 mx-2 text-muted" />
        <span class="text-foreground font-semibold">{{ currentRouteName }}</span>
      </nav>
    </div>

    <!-- Center: Empty -->
    <div class="hidden md:flex flex-1"></div>

    <!-- Right: Actions -->
    <div class="flex items-center justify-end flex-1 space-x-2 sm:space-x-3">
      <!-- Expandable Search Bar -->
      <div 
        ref="searchContainerRef"
        class="relative flex items-center h-9 transition-all duration-300 ease-in-out"
        :class="[isSearchExpanded ? 'w-64 sm:w-80' : 'w-9']"
      >
        <Button 
          v-if="!isSearchExpanded"
          variant="ghost" 
          size="icon" 
          class="rounded-full w-9 h-9 shrink-0 hover:bg-muted"
          @click.stop="isSearchExpanded = true"
        >
          <Search class="h-5 w-5 text-muted-foreground" />
        </Button>

        <div 
          class="absolute right-0 flex items-center w-full transition-all duration-300"
          :class="[isSearchExpanded ? 'opacity-100 visible translate-x-0' : 'opacity-0 invisible translate-x-4 pointer-events-none']"
        >
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search class="h-4 w-4 text-muted-foreground" />
          </div>
          <Input 
            v-model="searchQuery"
            type="text" 
            class="block w-full pl-9 pr-8 h-9 rounded-full bg-muted/50 focus-visible:ring-primary/20 focus-visible:border-primary border-none shadow-none"
            placeholder="搜索内容..." 
            @focus="handleSearchFocus"
          />
          <button 
            v-if="isSearchExpanded"
            @click="isSearchExpanded = false; searchQuery = ''"
            class="absolute right-2.5 text-muted-foreground hover:text-foreground p-0.5 rounded-full hover:bg-muted transition-colors"
          >
            <X class="h-4 w-4" />
          </button>

          <!-- Search Results Dropdown -->
          <div 
            v-if="isSearchExpanded && searchQuery && filteredSearchResults.length > 0"
            class="absolute top-full mt-2 right-0 w-80 bg-card rounded-xl shadow-2xl border border-border z-50 overflow-hidden"
          >
            <div class="py-2">
              <div class="px-4 py-1.5 text-[11px] font-bold text-muted-foreground/60 uppercase tracking-widest bg-muted/30">
                搜索结果
              </div>
              <router-link
                v-for="result in filteredSearchResults"
                :key="result.id"
                :to="result.path"
                class="flex items-center px-4 py-2.5 text-sm text-foreground hover:bg-primary/5 transition-all cursor-pointer group"
                @click="isSearchExpanded = false; searchQuery = ''"
              >
                <div class="mr-3 p-2 rounded-lg bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all">
                  <component 
                    :is="result.type === '订单' ? Clock : (result.type === '会员' ? Users : (result.type === '商品' ? Package : Settings))" 
                    class="w-4 h-4"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="font-medium truncate group-hover:text-primary transition-colors">{{ result.title }}</div>
                  <div class="text-[10px] text-muted-foreground uppercase tracking-tight">{{ result.type }}</div>
                </div>
                <ChevronRight class="w-4 h-4 text-muted-foreground opacity-0 group-hover:opacity-100 transform translate-x-[-4px] group-hover:translate-x-0 transition-all" />
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Notifications -->
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon" class="rounded-full w-9 h-9 relative">
            <Bell class="w-5 h-5" />
            <span class="absolute top-2 right-2 block h-2 w-2 rounded-full bg-destructive ring-2 ring-background"></span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-80 p-0">
          <div class="flex items-center justify-between px-4 py-3 border-b">
            <span class="font-semibold text-sm">通知</span>
            <span class="text-xs text-muted-foreground">3 条未读</span>
          </div>
          <div class="py-2">
            <DropdownMenuItem 
              v-for="notification in notifications" 
              :key="notification.id"
              class="px-4 py-3 cursor-pointer items-start gap-3 focus:bg-muted/50"
            >
              <div class="relative mt-1">
                <div class="w-2 h-2 rounded-full bg-primary" v-if="!notification.read"></div>
                <div class="w-2 h-2 rounded-full bg-muted" v-else></div>
              </div>
              <div class="flex-1 space-y-1">
                <p class="text-sm font-medium leading-none">{{ notification.title }}</p>
                <p class="text-xs text-muted-foreground line-clamp-2">
                  {{ notification.content }}
                </p>
                <p class="text-[10px] text-muted-foreground pt-1">{{ notification.time }}</p>
              </div>
            </DropdownMenuItem>
          </div>
          <div class="p-2 border-t bg-muted/20">
            <Button 
              variant="ghost" 
              size="sm" 
              class="w-full text-xs h-8 text-primary"
              @click="$router.push('/notifications')"
            >
              查看全部通知
              <ChevronRight class="w-3 h-3 ml-1" />
            </Button>
          </div>
        </DropdownMenuContent>
      </DropdownMenu>

      <!-- Theme Toggle -->
      <Button 
        variant="ghost" 
        size="icon" 
        @click="toggleDark()" 
        class="rounded-full w-9 h-9"
      >
        <Sun v-if="isDark" class="w-5 h-5 text-yellow-400" />
        <Moon v-else class="w-5 h-5" />
      </Button>
      
      <div class="h-6 w-px bg-border mx-2"></div>

      <!-- User Profile -->
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" class="relative h-9 w-9 rounded-full p-0 overflow-hidden">
            <Avatar class="h-9 w-9">
              <AvatarImage :src="authStore.currentUser?.avatar" alt="User avatar" />
              <AvatarFallback class="bg-primary/10 text-primary font-medium">
                {{ (authStore.currentUser?.realName || authStore.currentUser?.name || authStore.currentUser?.username || 'A').charAt(0).toUpperCase() }}
              </AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" side="bottom" :side-offset="8" class="w-56">
          <DropdownMenuLabel class="font-normal">
            <div class="flex flex-col space-y-1">
              <div class="flex items-center gap-2 px-1 py-1.5">
                <Avatar class="h-8 w-8">
                  <AvatarImage :src="authStore.currentUser?.avatar" alt="User avatar" />
                  <AvatarFallback class="bg-primary/10 text-primary font-medium">
                    {{ (authStore.currentUser?.realName || authStore.currentUser?.name || authStore.currentUser?.username || 'A').charAt(0).toUpperCase() }}
                  </AvatarFallback>
                </Avatar>
                <div class="grid flex-1 text-left text-sm leading-tight">
                  <span class="truncate font-semibold">{{ authStore.currentUser?.realName || authStore.currentUser?.name || authStore.currentUser?.username }}</span>
                  <span class="truncate text-xs text-muted-foreground">{{ authStore.currentUser?.email || authStore.currentUser?.phone || 'No contact info' }}</span>
                </div>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="isThemeDialogOpen = true">
            <Sun class="mr-2 h-4 w-4" />
            <span>主题设置</span>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="$router.push('/profile')">
            <User class="mr-2 h-4 w-4" />
            <span>个人中心</span>
          </DropdownMenuItem>
          <DropdownMenuItem @click="handleSwitchAccount">
            <UserPlus class="mr-2 h-4 w-4" />
            <span>切换账号</span>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleLogout" class="text-destructive focus:text-destructive">
            <LogOut class="mr-2 h-4 w-4" />
            <span>退出登录</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
    
    <ThemeSettings v-model:open="isThemeDialogOpen" />
  </header>
</template>