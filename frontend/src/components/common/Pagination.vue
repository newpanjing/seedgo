<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'

interface Props {
  page: number
  pageSize: number
  total: number
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  (e: 'update:page', value: number): void
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
  emit('update:page', newPage)
  emit('change')
}
</script>

<template>
  <div class="flex items-center justify-between px-4 py-3 border-t border-border bg-muted/30 text-xs text-muted-foreground">
    <div>
      共 {{ total }} 条记录，当前显示
      <span class="font-medium text-foreground">{{ pageStart }}-{{ pageEnd }}</span>
    </div>
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2" v-if="totalPages > 1">
        <Button
          variant="outline"
          size="sm"
          :disabled="page === 1 || total === 0 || loading"
          @click="handlePageChange(page - 1)"
        >
          上一页
        </Button>
        <Button
          variant="outline"
          size="sm"
          :disabled="page === totalPages || total === 0 || loading"
          @click="handlePageChange(page + 1)"
        >
          下一页
        </Button>
      </div>
    </div>
  </div>
</template>
