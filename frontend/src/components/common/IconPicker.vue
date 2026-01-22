<script setup lang="ts">
import { ref, computed } from 'vue'
import * as icons from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Search } from 'lucide-vue-next'

const props = defineProps<{
  modelValue?: string
}>()

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const searchQuery = ref('')
const displayLimit = ref(100)

// Filter out non-component exports and create a list of icon names
const iconList = Object.keys(icons)
  .filter(key => key !== 'default' && key !== 'createLucideIcon' && /^[A-Z]/.test(key))
  .sort()

const filteredIcons = computed(() => {
  const query = searchQuery.value.toLowerCase()
  if (!query) return iconList
  return iconList.filter(name => name.toLowerCase().includes(query))
})

const visibleIcons = computed(() => {
  return filteredIcons.value.slice(0, displayLimit.value)
})

const handleSelect = (iconName: string) => {
  emit('update:modelValue', iconName)
  open.value = false
}

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (target.scrollTop + target.clientHeight >= target.scrollHeight - 50) {
    if (displayLimit.value < filteredIcons.value.length) {
      displayLimit.value += 50
    }
  }
}

// Reset limit when search changes or popover opens
const onOpenChange = (val: boolean) => {
  if (val) {
    displayLimit.value = 100
    searchQuery.value = ''
  }
  open.value = val
}

// Dynamically get the icon component
const getIconComponent = (name: string) => {
  return (icons as any)[name]
}
</script>

<template>
  <Popover :open="open" @update:open="onOpenChange">
    <PopoverTrigger as-child>
      <Button
        variant="outline"
        role="combobox"
        class="w-full justify-between"
      >
        <div class="flex items-center gap-2">
          <template v-if="modelValue">
            <component
              :is="getIconComponent(modelValue)"
              v-if="getIconComponent(modelValue)"
              class="w-4 h-4"
            />
            <span>{{ modelValue }}</span>
          </template>
          <span v-else class="text-muted-foreground">选择图标...</span>
        </div>
        <Search class="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-[340px] p-0" align="start">
      <div class="p-2 border-b">
        <div class="relative">
          <Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            v-model="searchQuery"
            placeholder="搜索图标..."
            class="pl-8 h-9"
          />
        </div>
      </div>
      <div 
        class="h-[300px] overflow-y-auto p-2"
        @scroll="handleScroll"
      >
        <div class="grid grid-cols-4 gap-2">
          <button
            v-for="name in visibleIcons"
            :key="name"
            class="flex flex-col items-center justify-center p-2 rounded-md hover:bg-accent hover:text-accent-foreground transition-colors gap-1 h-16"
            :class="{ 'bg-accent text-accent-foreground': modelValue === name }"
            @click="handleSelect(name)"
            :title="name"
          >
            <component :is="getIconComponent(name)" class="w-6 h-6" />
            <span class="text-[9px] text-muted-foreground truncate w-full text-center">{{ name }}</span>
          </button>
        </div>
        <div v-if="filteredIcons.length === 0" class="p-4 text-center text-sm text-muted-foreground">
          未找到图标
        </div>
        <div v-else-if="visibleIcons.length < filteredIcons.length" class="p-2 text-center text-xs text-muted-foreground">
          滚动加载更多...
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
