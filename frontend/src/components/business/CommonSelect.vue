<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue'
import { Check, ChevronDown, X, Search, Loader2, RefreshCw } from 'lucide-vue-next'
import request from '@/utils/request'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { PopoverRoot, PopoverTrigger, PopoverContent, PopoverPortal } from 'radix-vue'
import { useDebounceFn, useInfiniteScroll } from '@vueuse/core'

const props = defineProps<{
  modelValue?: string | number
  placeholder?: string
  url: string
  valueField?: string
  labelField?: string
  searchParam?: string
  params?: Record<string, any>
  width?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number | undefined): void
  (e: 'change', value: any): void
}>()

// Default values for optional props
const kField = props.valueField || 'value'
const vField = props.labelField || 'label'
const sParam = props.searchParam || 'keyword'
const width = props.width || '150px'

const open = ref(false)
const loading = ref(false)
const items = ref<any[]>([])
const searchQuery = ref('')
const selectedItem = ref<any>(undefined)
const scrollEl = ref<HTMLElement | null>(null)

// Pagination state
const page = ref(1)
const pageSize = 10
const total = ref(0)
const hasMore = ref(true)

const fetchData = async (keyword = '', isLoadMore = false) => {
  if (!props.url) return
  if (!isLoadMore) {
    page.value = 1
    items.value = []
    hasMore.value = true
  }
  
  if (!hasMore.value && isLoadMore) return

  loading.value = true
  try {
    const res = await request<any>({
      url: props.url,
      method: 'get',
      params: {
        page: page.value,
        pageSize: pageSize,
        sortBy: 'id',
        sortOrder: 'desc',
        [sParam]: keyword,
        ...props.params
      }
    }) as any
    
    let newItems: any[] = []
    
    // Handle different response structures
    if (res && Array.isArray(res.items)) {
      newItems = res.items
      total.value = res.total || 0
    } else if (Array.isArray(res)) {
      newItems = res
      total.value = res.length
    } else {
      newItems = []
      total.value = 0
    }

    if (isLoadMore) {
      items.value = [...items.value, ...newItems]
    } else {
      items.value = newItems
    }

    // Check if we have more data
    hasMore.value = items.value.length < total.value
    if (newItems.length < pageSize) {
      hasMore.value = false
    }

    if (isLoadMore && newItems.length > 0) {
       page.value++
    } else if (!isLoadMore && newItems.length > 0) {
       // Prepared for next page
       page.value++
    }

  } catch (error) {
    console.error('Failed to fetch data:', error)
    if (!isLoadMore) items.value = []
  } finally {
    loading.value = false
  }
}

useInfiniteScroll(
  scrollEl,
  () => {
    if (!loading.value && hasMore.value) {
      fetchData(searchQuery.value, true)
    }
  },
  { distance: 10 }
)

const debouncedSearch = useDebounceFn((value: string) => {
  fetchData(value, false)
}, 300)

watch(searchQuery, (val) => {
  debouncedSearch(val)
})

watch(() => props.modelValue, (val) => {
  if (!val) {
    selectedItem.value = undefined
    return
  }
  
  // Try to find in current items first
  if (items.value.length > 0) {
    const found = items.value.find(item => String(item[kField]) === String(val))
    if (found) {
      selectedItem.value = found
      return
    }
  }

  // If not found and we have a value, we might need to fetch the specific item
  // Or just display the ID? For now, if not in list, we can't show label easily without another API call.
  // Assuming the initial list might contain it or we rely on pre-loaded list.
  // If we want to support showing label for selected item not in current page, 
  // we would need an API to get item by ID.
}, { immediate: true })

watch(items, () => {
  if (props.modelValue) {
    const found = items.value.find(item => String(item[kField]) === String(props.modelValue))
    if (found) selectedItem.value = found
  }
})

onMounted(() => {
  fetchData()
})

const handleSelect = (item: any) => {
  selectedItem.value = item
  emit('update:modelValue', item[kField])
  emit('change', item)
  open.value = false
}

const handleClear = (e: Event) => {
  e.stopPropagation()
  selectedItem.value = undefined
  emit('update:modelValue', undefined)
  emit('change', undefined)
}

const handleRefresh = () => {
    fetchData(searchQuery.value, false)
}
</script>

<template>
  <PopoverRoot v-model:open="open">
    <PopoverTrigger as-child>
      <Button variant="outline" role="combobox" :aria-expanded="open" class="justify-between px-3 font-normal"
        :style="{ width }"
        :class="!selectedItem && 'text-muted-foreground'">
        <span class="truncate">
          {{ selectedItem ? selectedItem[vField] : (placeholder || '请选择') }}
        </span>
        <div class="flex items-center gap-1">
          <div v-if="selectedItem" @click="handleClear" class="rounded-sm hover:bg-muted p-0.5 cursor-pointer z-10">
            <X class="h-4 w-4 opacity-50 hover:opacity-100" />
          </div>
          <ChevronDown class="ml-1 h-4 w-4 shrink-0 opacity-50 transition-transform duration-200" :class="{ 'rotate-180': open }" />
        </div>
      </Button>
    </PopoverTrigger>
    <PopoverPortal>
      <PopoverContent
        :side-offset="4"
        class="p-0 bg-popover text-popover-foreground rounded-md border shadow-md z-50 outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
        :style="{ width: 'var(--radix-popover-trigger-width)' }"
      >
        <div class="flex items-center border-b px-3 gap-2">
          <Search class="h-4 w-4 shrink-0 opacity-50" />
          <Input v-model="searchQuery"
            class="flex h-9 w-full rounded-md !bg-transparent py-2 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-none focus:ring-0 focus-visible:ring-0 px-0 shadow-none text-popover-foreground"
            :placeholder="'搜索...'" />
          <Button variant="ghost" size="icon" class="h-6 w-6 shrink-0" @click="handleRefresh" :disabled="loading">
            <RefreshCw class="h-3 w-3" :class="{ 'animate-spin': loading }" />
          </Button>
        </div>
        <div ref="scrollEl" class="max-h-[300px] overflow-y-auto p-1">
          <Transition name="fade" mode="out-in">
            <div v-if="loading && items.length === 0" class="flex items-center justify-center py-6 text-sm text-muted-foreground">
              <Loader2 class="h-4 w-4 animate-spin mr-2" />
              加载中...
            </div>
            <div v-else>
              <div v-if="items.length === 0" class="py-6 text-center text-sm text-muted-foreground">
                未找到数据
              </div>
              <div v-for="item in items" :key="item[kField]" @click="handleSelect(item)"
                class="relative flex cursor-default select-none items-center rounded-sm px-2 py-2.5 my-1 text-sm outline-none hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 data-[disabled]:pointer-events-none data-[disabled]:opacity-50 cursor-pointer"
                :class="selectedItem && String(selectedItem[kField]) === String(item[kField]) ? 'bg-accent text-accent-foreground dark:bg-accent/20 dark:text-foreground' : ''">
                <Check class="mr-2 h-4 w-4 shrink-0"
                  :class="selectedItem && String(selectedItem[kField]) === String(item[kField]) ? 'opacity-100' : 'opacity-0'" />
                <span class="truncate">{{ item[vField] }}</span>
              </div>
              <div v-if="loading && items.length > 0" class="py-2 text-center text-xs text-muted-foreground flex items-center justify-center">
                <Loader2 class="h-3 w-3 animate-spin mr-1" /> 加载更多...
              </div>
            </div>
          </Transition>
        </div>
      </PopoverContent>
    </PopoverPortal>
  </PopoverRoot>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
