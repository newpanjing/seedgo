<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { inject } from 'vue'
import { useVModel } from '@vueuse/core'
import { cn } from '@/lib/utils'
import { FormItemContextKey } from '@/lib/symbols'

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  (e: 'update:modelValue', payload: string | number): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})

const formItem = inject(FormItemContextKey, undefined)
</script>

<template>
  <input
    v-model="modelValue"
    :class="cn(
      'w-full px-3 py-1 h-9 rounded-xl border text-sm outline-none transition-all duration-200 bg-white dark:bg-slate-950 border-slate-200 dark:border-slate-800 text-slate-900 dark:text-slate-50 placeholder:text-slate-400 focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/5 dark:focus:ring-indigo-500/20',
      formItem?.validationError?.value && 'border-destructive focus:border-destructive focus:ring-destructive/20',
      props.class
    )"
  >
</template>
