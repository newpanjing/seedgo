<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getDictByCode } from '@/api/dict'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

const props = defineProps<{
  dictCode: string
  modelValue?: string | number
  placeholder?: string
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const items = ref<any[]>([])

const fetchDict = async () => {
  if (!props.dictCode) return
  const res = await getDictByCode(props.dictCode) as any
  if (res && res.items) {
    items.value = res.items.filter((i: any) => i.status === 1)
  }
}

watch(() => props.dictCode, fetchDict)

onMounted(fetchDict)

const handleChange = (val: string) => {
  emit('update:modelValue', val)
  emit('change', val)
}
</script>

<template>
  <Select :model-value="String(modelValue || '')" @update:model-value="handleChange">
    <SelectTrigger>
      <SelectValue :placeholder="placeholder || '请选择'" />
    </SelectTrigger>
    <SelectContent>
      <SelectGroup>
        <SelectItem v-for="item in items" :key="item.value" :value="item.value">
          {{ item.label }}
        </SelectItem>
      </SelectGroup>
    </SelectContent>
  </Select>
</template>
