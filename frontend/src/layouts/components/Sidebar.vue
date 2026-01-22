<script setup lang="ts">
import { ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { usePermissionStore } from '@/stores/permission'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { showConfirm } from '@/lib/message'
import { 
  LayoutDashboard, 
  LogOut, 
  PanelLeftClose, 
  PanelLeftOpen, 
  X,
  ChevronRight,
  Circle
} from 'lucide-vue-next'

const props = defineProps<{
  isSidebarOpen: boolean
  isSidebarCollapsed: boolean
}>()

const emit = defineEmits<{
  (e: 'update:isSidebarOpen', value: boolean): void
  (e: 'update:isSidebarCollapsed', value: boolean): void
}>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const openMenus = ref<Record<string, boolean>>({})

const toggleSubMenu = (name: string) => {
  openMenus.value[name] = !openMenus.value[name]
}

const isMenuOpen = (name: string) => !!openMenus.value[name]

// Auto open menus based on current route
watch(
  () => [route.path, permissionStore.menuTree],
  ([path, menuTree]) => {
    if (!menuTree || (menuTree as any[]).length === 0) return
    
    // Find which menu contains the current path
    const findParent = (items: any[]): string | undefined => {
      for (const item of items) {
        if (item.children) {
          if (item.children.some((child: any) => child.href === path)) {
            return item.name
          }
          const found = findParent(item.children)
          if (found) return found
        }
      }
      return undefined
    }

    const parentName = findParent(permissionStore.menuTree)
    if (parentName) {
      openMenus.value[parentName] = true
    }
  },
  { immediate: true, deep: true }
)

const handleLogout = async () => {
  const confirmed = await showConfirm('确定要退出登录吗？', '退出后需要重新登录才能访问系统。')
  if (confirmed) {
    authStore.logout()
    router.push('/login')
  }
}

const toggleCollapse = () => {
  emit('update:isSidebarCollapsed', !props.isSidebarCollapsed)
}

const closeSidebar = () => {
  emit('update:isSidebarOpen', false)
}
</script>

<template>
  <div>
    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isSidebarOpen" 
      class="fixed inset-0 z-40 bg-background/80 backdrop-blur-sm lg:hidden"
      @click="closeSidebar"
    ></div>

    <!-- Sidebar -->
    <aside 
      class="fixed inset-y-0 left-0 z-50 bg-card transform transition-all duration-300 ease-in-out shadow-xl lg:shadow-none border-r border-border flex flex-col"
      :class="[
        isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0',
        isSidebarCollapsed ? 'w-20' : 'w-64'
      ]"
    >
      <!-- Sidebar Header -->
      <div class="flex items-center h-16 px-4 border-b border-border shrink-0 transition-all duration-300"
        :class="[isSidebarCollapsed ? 'justify-center' : 'justify-between']"
      >
        <!-- Expanded State: Logo + Text + Close Button -->
        <template v-if="!isSidebarCollapsed">
          <div class="flex items-center space-x-3 overflow-hidden">
            <div class="h-8 w-8 bg-primary rounded-lg flex items-center justify-center text-primary-foreground shrink-0">
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <span class="text-xl font-bold tracking-tight whitespace-nowrap">Seedgo</span>
          </div>
          
          <!-- Collapse Button (Desktop) -->
          <Button 
            variant="ghost"
            size="icon"
            @click="toggleCollapse"
            class="hidden lg:flex"
            title="折叠菜单"
          >
            <PanelLeftClose class="w-5 h-5" />
          </Button>
        </template>
        
        <!-- Collapsed State: Expand Button Only -->
        <template v-else>
          <Button 
            variant="ghost"
            size="icon"
            @click="toggleCollapse"
            class="hidden lg:flex text-primary bg-primary/10 hover:bg-primary/20"
            title="展开菜单"
          >
            <PanelLeftOpen class="w-6 h-6" />
          </Button>
        </template>

        <!-- Mobile Close Button (Always visible on mobile if open) -->
        <Button 
          variant="ghost" 
          size="icon"
          @click="closeSidebar" 
          class="lg:hidden absolute right-4"
        >
          <X class="w-6 h-6" />
        </Button>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto overflow-x-hidden">
        <!-- Dashboard -->
        <RouterLink
          to="/"
          class="group flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200"
          :class="[
            route.path === '/' || route.name === 'dashboard'
              ? 'bg-primary text-primary-foreground shadow-md'
              : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground',
            isSidebarCollapsed ? 'justify-center' : ''
          ]"
          @click="closeSidebar"
          :title="isSidebarCollapsed ? '仪表盘' : ''"
        >
          <LayoutDashboard 
            class="w-5 h-5 transition-colors shrink-0"
            :class="[
              route.path === '/' || route.name === 'dashboard' ? 'text-primary-foreground' : 'text-muted-foreground group-hover:text-foreground',
              isSidebarCollapsed ? '' : 'mr-3'
            ]"
          />
          <span v-if="!isSidebarCollapsed" class="whitespace-nowrap transition-opacity duration-200">仪表盘</span>
          <div v-if="!isSidebarCollapsed && (route.path === '/' || route.name === 'dashboard')" class="ml-auto w-1.5 h-1.5 rounded-full bg-primary-foreground/90 shadow-[0_0_4px_currentColor]"></div>
        </RouterLink>

        <!-- Divider -->
        <!-- <div class="my-2 mx-1 border-b border-border/50"></div> -->
        
        <template v-for="item in permissionStore.menuTree" :key="item.name">
          <RouterLink
            v-if="!item.children || item.children.length === 0"
            :to="item.href"
            class="group flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200"
            :class="[
              route.path === item.href
                ? 'bg-primary text-primary-foreground shadow-md'
                : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground',
              isSidebarCollapsed ? 'justify-center' : ''
            ]"
            @click="closeSidebar"
            :title="isSidebarCollapsed ? item.name : ''"
          >
            <component 
              :is="item.icon" 
              class="w-5 h-5 transition-colors shrink-0"
              :class="[
                route.path === item.href ? 'text-primary-foreground' : 'text-muted-foreground group-hover:text-foreground',
                isSidebarCollapsed ? '' : 'mr-3'
              ]"
            />
            <span v-if="!isSidebarCollapsed" class="whitespace-nowrap transition-opacity duration-200">{{ item.name }}</span>
            <div v-if="!isSidebarCollapsed && route.path === item.href" class="ml-auto w-1.5 h-1.5 rounded-full bg-primary-foreground/90 shadow-[0_0_4px_currentColor]"></div>
          </RouterLink>

          <div v-else class="space-y-1">
            <!-- Collapsed State: Dropdown Menu -->
            <DropdownMenu v-if="isSidebarCollapsed">
              <DropdownMenuTrigger as-child>
                <button
                  type="button"
                  class="w-full group flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200 justify-center"
                  :class="[
                    item.children.some((child: any) => route.path === child.href)
                      ? 'bg-primary/10 text-primary'
                      : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground'
                  ]"
                  :title="item.name"
                >
                  <component 
                    :is="item.icon" 
                    class="w-5 h-5 transition-colors shrink-0"
                    :class="[
                      item.children.some((child: any) => route.path === child.href) ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'
                    ]"
                  />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent side="right" align="start" class="w-48 ml-2">
                <DropdownMenuLabel>{{ item.name }}</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem v-for="child in item.children" :key="child.name" as-child>
                  <RouterLink 
                    :to="child.href"
                    class="w-full cursor-pointer flex items-center"
                    :class="[
                      route.path === child.href ? 'text-primary font-medium' : ''
                    ]"
                  >
                    <component 
                      :is="child.icon || Circle" 
                      class="w-4 h-4 mr-2 shrink-0"
                      :class="[
                        route.path === child.href ? 'text-primary' : 'text-muted-foreground'
                      ]"
                    />
                    <span class="flex-1">{{ child.name }}</span>
                    <div v-if="route.path === child.href" class="w-1.5 h-1.5 rounded-full bg-primary shadow-[0_0_4px_currentColor]"></div>
                  </RouterLink>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <!-- Expanded State: Accordion -->
            <template v-else>
              <button
                type="button"
                class="w-full group flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200"
                :class="[
                  item.children.some((child: any) => route.path === child.href)
                    ? 'bg-primary/10 text-primary'
                    : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground'
                ]"
                @click="toggleSubMenu(item.name)"
              >
                <component 
                  :is="item.icon" 
                  class="w-5 h-5 transition-colors shrink-0 mr-3"
                  :class="[
                    item.children.some((child: any) => route.path === child.href) ? 'text-primary' : 'text-muted-foreground group-hover:text-foreground'
                  ]"
                />
                <span class="whitespace-nowrap transition-opacity duration-200">{{ item.name }}</span>
                <ChevronRight
                  class="w-4 h-4 ml-auto transition-transform"
                  :class="isMenuOpen(item.name) ? 'transform rotate-90' : ''"
                />
              </button>

              <div
                v-if="isMenuOpen(item.name)"
                class="pl-9 space-y-1"
              >
                <RouterLink
                  v-for="child in item.children"
                  :key="child.name"
                  :to="child.href"
                  class="group flex items-center px-3 py-2 text-xs font-medium rounded-lg transition-all duration-200"
                  :class="[
                    route.path === child.href
                      ? 'bg-primary text-primary-foreground shadow-md'
                      : 'text-muted-foreground hover:bg-primary/10 hover:text-foreground'
                  ]"
                  @click="closeSidebar"
                >
                  <component 
                    :is="child.icon || Circle" 
                    class="w-4 h-4 mr-2 shrink-0"
                    :class="[
                      route.path === child.href ? 'text-primary-foreground' : 'text-muted-foreground group-hover:text-foreground'
                    ]"
                  />
                  <span class="whitespace-nowrap transition-opacity duration-200">{{ child.name }}</span>
                  <div v-if="route.path === child.href" class="ml-auto w-1.5 h-1.5 rounded-full bg-primary-foreground/90 shadow-[0_0_4px_currentColor]"></div>
                </RouterLink>
              </div>
            </template>
          </div>
        </template>
      </nav>
      
      <!-- Footer Actions -->
      <div class="p-4 border-t border-border shrink-0 space-y-2">
        <!-- Logout -->
        <Button 
          variant="ghost"
          @click="handleLogout"
          class="w-full justify-start text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
          :class="[isSidebarCollapsed ? 'justify-center px-2' : 'px-3']"
          :title="isSidebarCollapsed ? '退出登录' : ''"
        >
          <LogOut class="w-5 h-5 shrink-0" :class="[isSidebarCollapsed ? '' : 'mr-3']" />
          <span v-if="!isSidebarCollapsed" class="whitespace-nowrap">退出登录</span>
        </Button>
      </div>
    </aside>
  </div>
</template>