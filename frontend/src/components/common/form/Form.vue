<script setup lang="ts">
import { computed, provide, reactive, toRef } from 'vue'
import { cn } from '@/lib/utils'
import { FormContextKey, type FormRules } from '@/lib/symbols'

const props = withDefaults(defineProps<{
  model?: Record<string, any>
  rules?: FormRules
  columns?: number
  gap?: number
  class?: string
}>(), {
  columns: 1,
  gap: 4
})

const errors = reactive<Record<string, string>>({})
const fields = new Set<string>()

const addField = (field: string) => fields.add(field)
const removeField = (field: string) => fields.delete(field)

// Use reactive to unwrap refs and keep them reactive
const context = reactive({
  model: toRef(props, 'model'),
  rules: toRef(props, 'rules'),
  errors,
  addField,
  removeField
})

provide(FormContextKey, context)

const validate = async () => {
  if (!props.model) return true
  
  let isValid = true
  // Clear previous errors
  Object.keys(errors).forEach(key => delete errors[key])

  // Validate all registered fields
  for (const field of fields) {
    const value = props.model[field]
    const fieldRules = props.rules?.[field]
    
    if (!fieldRules) continue
    
    const rulesArray = Array.isArray(fieldRules) ? fieldRules : [fieldRules]
    
    for (const rule of rulesArray) {
      if (rule.required) {
        if (value === undefined || value === null || value === '') {
          errors[field] = rule.message || '此项是必填项'
          isValid = false
          break 
        }
      }
      
      if (rule.pattern) {
        if (value && !rule.pattern.test(value)) {
           errors[field] = rule.message || '格式不正确'
           isValid = false
           break
        }
      }
      
      if (rule.validator) {
        try {
          const result = await rule.validator(value)
          if (result === false || typeof result === 'string') {
            errors[field] = typeof result === 'string' ? result : (rule.message || '验证失败')
            isValid = false
            break
          }
        } catch (e: any) {
          errors[field] = e.message || rule.message || '验证失败'
          isValid = false
          break
        }
      }
    }
  }
  
  return isValid
}

const clearValidate = (fieldsToClear?: string[]) => {
    if (fieldsToClear) {
        fieldsToClear.forEach(f => delete errors[f])
    } else {
        Object.keys(errors).forEach(key => delete errors[key])
    }
}

defineExpose({
  validate,
  clearValidate
})

const gridClass = computed(() => {
  return cn(
    'grid',
    `gap-${props.gap}`,
    {
      'grid-cols-1': props.columns === 1,
      'grid-cols-2': props.columns === 2,
      'grid-cols-3': props.columns === 3,
      'grid-cols-4': props.columns === 4,
      'grid-cols-6': props.columns === 6,
      'grid-cols-12': props.columns === 12,
    },
    props.class
  )
})
</script>

<template>
  <div :class="gridClass">
    <slot />
  </div>
</template>
