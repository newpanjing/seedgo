<script setup lang="ts">
import { computed } from 'vue'
import DateSelect from './DateSelect.vue'
import { ArrowRight } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  modelValue?: [string | undefined, string | undefined] | []
  type?: 'date' | 'time' | 'datetime'
  placeholder?: [string, string] | string[]
  disabled?: boolean
  clearable?: boolean
}>(), {
  type: 'date',
  modelValue: () => [],
  placeholder: () => ['开始时间', '结束时间'],
  clearable: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: [string | undefined, string | undefined]): void
  (e: 'change', value: [string | undefined, string | undefined]): void
}>()

const startValue = computed({
  get: () => props.modelValue?.[0],
  set: (val) => {
    const end = props.modelValue?.[1]
    const newVal: [string | undefined, string | undefined] = [val, end]
    emit('update:modelValue', newVal)
    emit('change', newVal)
  }
})

const endValue = computed({
  get: () => props.modelValue?.[1],
  set: (val) => {
    const start = props.modelValue?.[0]
    const newVal: [string | undefined, string | undefined] = [start, val]
    emit('update:modelValue', newVal)
    emit('change', newVal)
  }
})
</script>

<template>
  <div class="flex items-center gap-2">
    <DateSelect
      v-model="startValue"
      :type="type"
      :placeholder="placeholder[0]"
      :disabled="disabled"
      :clearable="clearable"
      :max-date="endValue"
      class="flex-1"
    />
    <span class="text-muted-foreground flex-shrink-0">
        <ArrowRight class="h-4 w-4" />
    </span>
    <DateSelect
      v-model="endValue"
      :type="type"
      :placeholder="placeholder[1]"
      :disabled="disabled"
      :clearable="clearable"
      :min-date="startValue"
      class="flex-1"
    />
  </div>
</template>
