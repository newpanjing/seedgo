<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Check, ChevronsUpDown, X, Search, Loader2 } from 'lucide-vue-next'
import { getTenants, type Tenant } from '@/api/tenant'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { PopoverRoot, PopoverTrigger, PopoverContent, PopoverPortal } from 'radix-vue'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  modelValue?: string | number
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number | undefined): void
  (e: 'change', value: Tenant | undefined): void
}>()

const open = ref(false)
const loading = ref(false)
const tenants = ref<Tenant[]>([])
const searchQuery = ref('')
const selectedTenant = ref<Tenant | undefined>(undefined)

const fetchTenants = async (keyword = '') => {
  loading.value = true
  try {
    const res = await getTenants({
      page: 1,
      pageSize: 10,
      name: keyword, // Changed from keyword to name based on backend TenantService
      sortBy: 'createdAt',
      sortOrder: 'desc'
    }) as any
    console.log('Tenants loaded:',  res.items)
    tenants.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch tenants:', error)
  } finally {
    loading.value = false
  }
}

const debouncedSearch = useDebounceFn((value: string) => {
  fetchTenants(value)
}, 300)

watch(searchQuery, (val) => {
  debouncedSearch(val)
})

watch(() => props.modelValue, (val) => {
  if (!val) {
    selectedTenant.value = undefined
    return
  }
  // Ensure tenants are loaded or wait for them?
  // Since fetchTenants is called onMounted, we might need to wait if immediate is true.
  // But tenants is empty initially.
  if (tenants.value.length > 0) {
    const found = tenants.value.find(t => String(t.id) === String(val))
    if (found) {
      selectedTenant.value = found
    }
  }
}, { immediate: true })

watch(tenants, () => {
  if (props.modelValue) {
    const found = tenants.value.find(t => String(t.id) === String(props.modelValue))
    if (found) selectedTenant.value = found
  }
})

onMounted(() => {
  fetchTenants()
})

const handleSelect = (tenant: Tenant) => {
  selectedTenant.value = tenant
  emit('update:modelValue', tenant.id)
  emit('change', tenant)
  open.value = false
}

const handleClear = (e: Event) => {
  e.stopPropagation()
  selectedTenant.value = undefined
  emit('update:modelValue', undefined)
  emit('change', undefined)
}
</script>

<template>
  <PopoverRoot v-model:open="open">
    <PopoverTrigger as-child>
      <Button variant="outline" role="combobox" :aria-expanded="open" class="w-[200px] justify-between px-3 font-normal"
        :class="!selectedTenant && 'text-muted-foreground'">
        <span class="truncate">
          {{ selectedTenant ? selectedTenant.name : (placeholder || '选择租户') }}
        </span>
        <div class="flex items-center gap-1">
          <div v-if="selectedTenant" @click="handleClear" class="rounded-sm hover:bg-muted p-0.5 cursor-pointer z-10">
            <X class="h-4 w-4 opacity-50 hover:opacity-100" />
          </div>
          <ChevronsUpDown class="ml-1 h-4 w-4 shrink-0 opacity-50" />
        </div>
      </Button>
    </PopoverTrigger>
    <PopoverPortal>
      <PopoverContent
        class="w-[200px] p-0 bg-popover text-popover-foreground rounded-md border shadow-md z-50 outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2">
        <div class="flex items-center border-b px-3">
          <Search class="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input v-model="searchQuery"
            class="flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-none focus-visible:ring-0 px-0"
            placeholder="搜索租户..." />
        </div>
        <div class="max-h-[300px] overflow-y-auto p-1">
          <div v-if="loading" class="flex items-center justify-center py-6 text-sm text-muted-foreground">
            <Loader2 class="h-4 w-4 animate-spin mr-2" />
            加载中...
          </div>
          <template v-else>
            <div v-if="tenants.length === 0" class="py-6 text-center text-sm text-muted-foreground">
              未找到租户
            </div>
            <div v-for="tenant in tenants" :key="tenant.id" @click="handleSelect(tenant)"
              class="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 cursor-pointer"
              :class="String(selectedTenant?.id) === String(tenant.id) ? 'bg-accent text-accent-foreground' : ''">
              <Check class="mr-2 h-4 w-4"
                :class="String(selectedTenant?.id) === String(tenant.id) ? 'opacity-100' : 'opacity-0'" />
                {{ tenant.name }}
            </div>
          </template>
        </div>
      </PopoverContent>
    </PopoverPortal>
  </PopoverRoot>
</template>
