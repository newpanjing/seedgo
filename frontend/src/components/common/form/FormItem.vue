<script setup lang="ts">
import { inject, computed, onMounted, onUnmounted, provide } from 'vue'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'
import { FormContextKey, FormItemContextKey, type FormContext } from '@/lib/symbols'

const props = defineProps<{
  field?: string
  label?: string
  required?: boolean
  error?: string
  description?: string
  className?: string
}>()

const formContext = inject<FormContext | undefined>(FormContextKey, undefined)

onMounted(() => {
  if (props.field && formContext) {
    formContext.addField(props.field)
  }
})

onUnmounted(() => {
  if (props.field && formContext) {
    formContext.removeField(props.field)
  }
})

const validationError = computed(() => {
  if (props.error) return props.error
  if (props.field && formContext && formContext.errors[props.field]) {
    return formContext.errors[props.field]
  }
  return undefined
})

provide(FormItemContextKey, {
  validationError
})

const isRequired = computed(() => {
  if (props.required) return true
  if (props.field && formContext && formContext.rules && formContext.rules[props.field]) {
    const rules = formContext.rules[props.field]
    const rulesArray = Array.isArray(rules) ? rules : [rules]
    return rulesArray.some(r => r.required)
  }
  return false
})
</script>

<template>
  <div :class="cn(className)">
    <div v-if="label" class="flex items-center">
      <Label :class="cn('text-xs font-bold', validationError && 'text-destructive [&_span]:!text-destructive')">
        {{ label }}
      </Label>
      <span v-if="isRequired" class="text-destructive text-xs">*</span>
    </div>
    <slot />
    <p v-if="description" class="text-xs text-muted-foreground">
      {{ description }}
    </p>
    <p v-if="validationError" class="text-xs font-medium text-destructive">
      {{ validationError }}
    </p>
  </div>
</template>
