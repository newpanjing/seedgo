<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Badge } from '@/components/ui/badge'
import { getRoles, type Role } from '@/api/role'
import { useRouter } from 'vue-router'
import { X, Check, ChevronDown } from 'lucide-vue-next'

const router = useRouter()
const props = defineProps<{
  modelValue?: string | number | (string | number)[]
  placeholder?: string
  disabled?: boolean
  clearable?: boolean
  multiple?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number | (string | number)[] | undefined): void
  (e: 'change', role: Role | Role[] | undefined): void
}>()

const roles = ref<Role[]>([])
const loading = ref(false)
const open = ref(false)

const goToCreateRole = () => {
  router.push('/system/roles')
  open.value = false
}

const fetchRoles = async () => {
  loading.value = true
  try {
    const params: any = {
      page: 1,
      pageSize: 100
    }
    
    const res = await getRoles(params) as any
    roles.value = res.items || []
  } catch (error) {
    console.error('Failed to fetch roles:', error)
    roles.value = []
  } finally {
    loading.value = false
  }
}

// Initial fetch
fetchRoles()

// Single Select Handler
const handleValueChange = (value: string) => {
  const role = roles.value.find(r => String(r.id) === value)
  if (role) {
    emit('update:modelValue', role.id)
    emit('change', role)
  } else {
    emit('update:modelValue', value)
  }
}

// Multi Select Logic
const selectedValues = computed(() => {
  if (!props.modelValue) return []
  return Array.isArray(props.modelValue) ? props.modelValue.map(String) : [String(props.modelValue)]
})

const selectedRoles = computed(() => {
  return roles.value.filter(r => selectedValues.value.includes(String(r.id)))
})

const toggleRole = (roleId: string | number) => {
  const idStr = String(roleId)
  const current = new Set(selectedValues.value)
  if (current.has(idStr)) {
    current.delete(idStr)
  } else {
    current.add(idStr)
  }
  const newValue = Array.from(current)
  // Find original roles to get correct ID types
  const selectedOriginalRoles = roles.value.filter(r => current.has(String(r.id)))
  const emittedIds = selectedOriginalRoles.map(r => r.id)
  
  emit('update:modelValue', emittedIds)
  emit('change', selectedOriginalRoles)
}

const removeRole = (e: Event, roleId: string | number) => {
  e.stopPropagation()
  toggleRole(roleId)
}

const handleClear = (e: Event) => {
  e.stopPropagation()
  emit('update:modelValue', props.multiple ? [] : undefined)
  emit('change', undefined)
}
</script>

<template>
  <div class="relative w-full">
    <!-- Single Select Mode -->
    <template v-if="!multiple">
      <Select :model-value="modelValue ? String(modelValue) : undefined" @update:model-value="handleValueChange" :disabled="disabled || loading">
        <SelectTrigger class="w-full pr-8">
          <SelectValue :placeholder="placeholder || '选择角色'" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup v-if="roles.length > 0">
            <SelectItem v-for="role in roles" :key="String(role.id)" :value="String(role.id)">
              {{ role.name }}
            </SelectItem>
          </SelectGroup>
          <div v-else class="py-6 text-center text-sm text-muted-foreground px-2">
            <p>暂无角色数据</p>
            <button class="text-primary hover:underline mt-2" @click.prevent="goToCreateRole">
              去创建角色
            </button>
          </div>
        </SelectContent>
      </Select>
    </template>

    <!-- Multi Select Mode -->
    <template v-else>
      <Popover v-model:open="open">
        <PopoverTrigger as-child>
          <div
            class="flex min-h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
            :class="{'opacity-50 cursor-not-allowed': disabled || loading}"
          >
            <div class="flex flex-wrap gap-1">
              <span v-if="selectedRoles.length === 0" class="text-muted-foreground">{{ placeholder || '选择角色' }}</span>
              <Badge
                v-for="role in selectedRoles"
                :key="role.id"
                variant="secondary"
                class="mr-1 mb-1"
              >
                {{ role.name }}
                <button
                  class="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  @click="(e) => removeRole(e, role.id)"
                  v-if="!disabled"
                >
                  <X class="h-3 w-3 text-muted-foreground hover:text-foreground" />
                </button>
              </Badge>
            </div>
            <ChevronDown class="h-4 w-4 opacity-50 shrink-0" />
          </div>
        </PopoverTrigger>
        <PopoverContent class="w-[200px] p-0" align="start">
          <div class="max-h-[300px] overflow-y-auto p-1">
             <div v-if="roles.length === 0" class="py-6 text-center text-sm text-muted-foreground px-2">
                <p>暂无角色数据</p>
                <button class="text-primary hover:underline mt-2" @click.prevent="goToCreateRole">
                  去创建角色
                </button>
              </div>
              <div
                v-for="role in roles"
                :key="role.id"
                class="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                @click="toggleRole(role.id)"
              >
                <div class="mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary"
                     :class="selectedValues.includes(String(role.id)) ? 'bg-primary text-primary-foreground' : 'opacity-50 [&_svg]:invisible'">
                  <Check class="h-4 w-4" />
                </div>
                <span>{{ role.name }}</span>
              </div>
          </div>
        </PopoverContent>
      </Popover>
    </template>
    
    <button
      v-if="clearable && ((!multiple && modelValue) || (multiple && modelValue && (modelValue as any[]).length > 0)) && !disabled"
      class="absolute right-8 top-1/2 -translate-y-1/2 text-muted-foreground/50 hover:text-foreground p-1 rounded-sm transition-colors"
      :class="{'right-2': !multiple, 'right-8': multiple}"
      @click="handleClear"
    >
      <X class="w-4 h-4" />
    </button>
  </div>
</template>
