<script setup lang="ts">
import { computed, ref, watch, provide, inject } from 'vue'
import { cn } from '@/lib/utils'
import { Check, ChevronRight, ChevronDown, Minus, Circle } from 'lucide-vue-next'
import * as icons from 'lucide-vue-next'

interface TreeItem {
  id: string | number
  name: string
  icon?: string
  children?: TreeItem[]
  [key: string]: any
}

const props = defineProps<{
  items: TreeItem[]
  modelValue?: (string | number)[]
  level?: number
  defaultExpand?: (item: TreeItem) => boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: (string | number)[]): void
}>()

// --- Expansion Logic ---
const expandedIds = ref<Set<string | number>>(new Set())

const initExpanded = () => {
  const newExpanded = new Set<string | number>()
  if (props.defaultExpand) {
    const traverse = (nodes: TreeItem[]) => {
      nodes.forEach(item => {
        if (props.defaultExpand!(item)) {
          newExpanded.add(item.id)
        }
        if (item.children) {
          traverse(item.children)
        }
      })
    }
    traverse(props.items)
  }
  expandedIds.value = newExpanded
}

watch(() => props.items, initExpanded, { immediate: true })

const toggleExpand = (id: string | number) => {
  const newExpanded = new Set(expandedIds.value)
  if (newExpanded.has(id)) {
    newExpanded.delete(id)
  } else {
    newExpanded.add(id)
  }
  expandedIds.value = newExpanded
}

const isExpanded = (id: string | number) => expandedIds.value.has(id)

// --- Selection Logic ---

const isRoot = computed(() => !props.level || props.level === 0)
const parentMap = ref(new Map<string | number, string | number>())

const getAllDescendantIds = (item: TreeItem): (string | number)[] => {
  let ids: (string | number)[] = []
  if (item.children) {
    item.children.forEach(child => {
      ids.push(child.id)
      ids = ids.concat(getAllDescendantIds(child))
    })
  }
  return ids
}

// Build parent map (Only at Root)
const buildParentMap = (nodes: TreeItem[], parentId?: string | number) => {
  nodes.forEach(node => {
    if (parentId) {
      parentMap.value.set(node.id, parentId)
    }
    if (node.children) {
      buildParentMap(node.children, node.id)
    }
  })
}

if (isRoot.value) {
  watch(() => props.items, (newItems) => {
    parentMap.value.clear()
    buildParentMap(newItems)
  }, { immediate: true, deep: true })
}

// Root Toggle Handler
const rootHandleToggle = (item: TreeItem, checked: boolean) => {
  const currentSelected = new Set(props.modelValue || [])
  const descendants = getAllDescendantIds(item)
  
  if (checked) {
    // 1. Add self
    currentSelected.add(item.id)
    // 2. Add descendants
    descendants.forEach(id => currentSelected.add(id))
    
    // 3. Add ancestors
    let currentId = item.id
    while (true) {
      const pId = parentMap.value.get(currentId)
      if (!pId) break
      currentSelected.add(pId)
      currentId = pId
    }
  } else {
    // 1. Remove self
    currentSelected.delete(item.id)
    // 2. Remove descendants
    descendants.forEach(id => currentSelected.delete(id))
    
    // 3. (Optional) Check if ancestors should be unchecked? 
    // Requirement: "Child selected -> Parent selected". 
    // It doesn't imply "Child unchecked -> Parent unchecked".
    // So we leave ancestors alone.
  }
  
  emit('update:modelValue', Array.from(currentSelected))
}

// Provide/Inject Context
const treeContext = inject<{ handleToggle: (item: TreeItem, checked: boolean) => void }>('treeContext', null as any)

const context = isRoot.value ? {
  handleToggle: rootHandleToggle
} : null

if (isRoot.value) {
  provide('treeContext', context)
}

// Local Toggle Wrapper
const toggleSelection = (item: TreeItem, checked: boolean) => {
  if (treeContext) {
    treeContext.handleToggle(item, checked)
  } else if (context) {
    context.handleToggle(item, checked)
  } else {
    // Fallback for standalone usage without root context (shouldn't happen in recursion)
    const currentSelected = new Set(props.modelValue || [])
    const descendants = getAllDescendantIds(item)
    if (checked) {
      currentSelected.add(item.id)
      descendants.forEach(id => currentSelected.add(id))
    } else {
      currentSelected.delete(item.id)
      descendants.forEach(id => currentSelected.delete(id))
    }
    emit('update:modelValue', Array.from(currentSelected))
  }
}

const isSelected = (id: string | number) => props.modelValue?.includes(id) || false

const isIndeterminate = (item: TreeItem) => {
  // If not selected, it's not indeterminate (it's unchecked)
  // Wait, standard checkbox: 
  // - Checked: true/false
  // - Indeterminate: true/false (visual override)
  // Here, if parent is NOT in modelValue, it is strictly unchecked.
  // But if we enforce "Child Selected -> Parent Selected", then parent MUST be in modelValue if indeterminate.
  if (!isSelected(item.id)) return false
  
  if (!item.children || item.children.length === 0) return false
  
  const descendants = getAllDescendantIds(item)
  if (descendants.length === 0) return false
  
  const selectedDescendants = descendants.filter(id => props.modelValue?.includes(id))
  // Indeterminate if some but not all descendants are selected
  return selectedDescendants.length > 0 && selectedDescendants.length < descendants.length
}

const getCheckboxState = (item: TreeItem) => {
  const selected = isSelected(item.id)
  const indeterminate = isIndeterminate(item)
  return { selected, indeterminate }
}

const resolveIcon = (iconName?: string) => {
  if (!iconName) return undefined
  
  // Try exact match
  if ((icons as any)[iconName]) {
    return (icons as any)[iconName]
  }

  // Try PascalCase (e.g. "user" -> "User")
  const pascalName = iconName.charAt(0).toUpperCase() + iconName.slice(1)
  if ((icons as any)[pascalName]) {
    return (icons as any)[pascalName]
  }

  // Try kebab-case to PascalCase (e.g. "user-check" -> "UserCheck")
  const camelToPascal = iconName.replace(/(^\w|-\w)/g, (g) => g.replace('-', '').toUpperCase())
  if ((icons as any)[camelToPascal]) {
    return (icons as any)[camelToPascal]
  }
  
  return undefined
}
</script>

<template>
  <div :class="cn('flex flex-col gap-1', level && level > 0 ? 'pl-6 border-l ml-2 border-border/50' : '')">
    <div v-for="item in items" :key="item.id" class="flex flex-col">
      <div class="flex items-center gap-2 py-1.5 group hover:bg-muted/50 rounded px-1 -ml-1">
        <button 
          v-if="item.children && item.children.length > 0"
          @click.stop="toggleExpand(item.id)"
          class="p-0.5 hover:bg-muted rounded-sm text-muted-foreground"
        >
          <component 
            :is="isExpanded(item.id) ? ChevronDown : ChevronRight" 
            class="h-4 w-4"
          />
        </button>
        <div v-else class="w-5"></div>

        <div class="relative flex items-center justify-center mr-2">
          <!-- Hidden real checkbox for accessibility/form -->
          <input 
            type="checkbox" 
            :checked="isSelected(item.id)"
            @change="(e) => toggleSelection(item, (e.target as HTMLInputElement).checked)"
            class="peer h-4 w-4 shrink-0 rounded border border-primary ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 appearance-none checked:bg-primary checked:text-primary-foreground"
          />
          <!-- Visual Indicators -->
          <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
            <Minus 
              v-if="isIndeterminate(item)" 
              class="h-3 w-3 text-primary-foreground" 
            />
            <Check 
              v-else-if="isSelected(item.id)" 
              class="h-3 w-3 text-primary-foreground" 
            />
          </div>
        </div>
        
        <component 
          v-if="item.icon && resolveIcon(item.icon)"
          :is="resolveIcon(item.icon)" 
          class="h-4 w-4 mr-2 text-muted-foreground"
        />
        
        <span class="text-sm cursor-pointer select-none flex-1" @click="toggleSelection(item, !isSelected(item.id))">
          {{ item.name }}
        </span>
      </div>
      
      <div v-if="item.children && item.children.length > 0 && isExpanded(item.id)">
        <Tree 
          :items="item.children"
          :model-value="modelValue"
          @update:model-value="(val) => emit('update:modelValue', val)"
          :level="(level || 0) + 1"
          :default-expand="defaultExpand"
        />
      </div>
    </div>
  </div>
</template>
