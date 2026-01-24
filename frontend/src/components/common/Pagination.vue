<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-vue-next'

interface Props {
  page: number
  pageSize: number
  total: number
  loading?: boolean
  pageSizeOptions?: number[]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  pageSizeOptions: () => [10, 20, 50, 100]
})

const emit = defineEmits<{
  (e: 'update:page', value: number): void
  (e: 'update:pageSize', value: number): void
  (e: 'change'): void
}>()

const pageStart = computed(() => {
  if (props.total === 0) return 0
  return (props.page - 1) * props.pageSize + 1
})

const pageEnd = computed(() => {
  if (props.total === 0) return 0
  return Math.min(props.page * props.pageSize, props.total)
})

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const handlePageChange = (newPage: number) => {
  if (newPage < 1 || newPage > totalPages.value) return
  emit('update:page', newPage)
  emit('change')
}

const handlePageSizeChange = (val: string) => {
  const newSize = Number(val)
  emit('update:pageSize', newSize)
  // Reset to page 1 when page size changes to avoid out of bounds
  emit('update:page', 1)
  emit('change')
}

// Page Numbers Logic
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = props.page
  const delta = 1
  const range: (number | string)[] = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      range.push(i)
    }
  } else {
    range.push(1)

    if (current > delta + 3) {
      range.push('...')
    } else {
        // If no ellipsis, fill the gap
        for(let i = 2; i < Math.max(2, current - delta); i++) {
            range.push(i)
        }
    }

    const start = Math.max(2, current - delta)
    const end = Math.min(total - 1, current + delta)

    // Adjust start/end to keep constant number of items if possible? 
    // Simplified: just show range
    for (let i = start; i <= end; i++) {
      range.push(i)
    }

    if (current < total - (delta + 2)) {
      range.push('...')
    } else {
        for(let i = Math.min(total - 1, current + delta) + 1; i < total; i++) {
            range.push(i)
        }
    }

    range.push(total)
  }
  
  // Deduplicate just in case logic overlaps
  return [...new Set(range)]
})

// Jump to Page Logic
const jumpPage = ref(props.page)

watch(() => props.page, (val) => {
  jumpPage.value = val
})

const handleJump = () => {
  let p = Number(jumpPage.value)
  if (isNaN(p)) {
    p = 1
  }
  if (p < 1) p = 1
  if (p > totalPages.value) p = totalPages.value
  
  jumpPage.value = p
  
  if (p !== props.page) {
    handlePageChange(p)
  }
}
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-4 px-4 py-3 border-t border-border bg-muted/30 text-xs text-muted-foreground">
    <div>
      共 {{ total }} 条
    </div>
    <div class="flex items-center gap-4">
      <!-- Page Size Select -->
      <div class="relative">
          <select 
            :value="pageSize" 
            @change="e => handlePageSizeChange((e.target as HTMLSelectElement).value)"
            class="h-8 w-[100px] appearance-none rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
          >
            <option v-for="size in pageSizeOptions" :key="size" :value="size">
                {{ size }} 条/页
            </option>
          </select>
      </div>

      <div class="flex items-center gap-2" v-if="totalPages > 1">
        <Button
          variant="outline"
          size="icon"
          class="h-8 w-8"
          :disabled="page === 1 || total === 0 || loading"
          @click="handlePageChange(page - 1)"
        >
          <ChevronLeft class="h-4 w-4" />
        </Button>

        <template v-for="(p, index) in visiblePages" :key="index">
            <span v-if="p === '...'" class="flex items-center justify-center w-8 h-8">
                <MoreHorizontal class="h-4 w-4" />
            </span>
            <Button
                v-else
                :variant="p === page ? 'default' : 'outline'"
                size="icon"
                class="h-8 w-8"
                :disabled="loading"
                @click="handlePageChange(p as number)"
            >
                {{ p }}
            </Button>
        </template>

        <Button
          variant="outline"
          size="icon"
          class="h-8 w-8"
          :disabled="page === totalPages || total === 0 || loading"
          @click="handlePageChange(page + 1)"
        >
          <ChevronRight class="h-4 w-4" />
        </Button>
      </div>

      <!-- Jump to Page -->
      <div class="flex items-center gap-2" v-if="totalPages > 1">
        <span>跳至</span>
        <Input
            v-model="jumpPage"
            type="number"
            :min="1"
            :max="totalPages"
            class="h-8 w-12 px-1 text-center"
            @keydown.enter="handleJump"
            @blur="handleJump"
        />
      </div>
    </div>
  </div>
</template>
