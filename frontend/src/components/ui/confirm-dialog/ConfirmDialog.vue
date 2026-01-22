<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

const props = withDefaults(defineProps<{
  open: boolean
  title?: string
  description?: string
  confirmText?: string
  cancelText?: string
  class?: HTMLAttributes['class']
}>(), {
  open: false,
  title: '确认操作',
  description: '',
  confirmText: '确认',
  cancelText: '取消',
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const handleOpenChange = (value: boolean) => {
  emit('update:open', value)
}

const handleCancel = () => {
  emit('cancel')
  emit('update:open', false)
}

const handleConfirm = () => {
  emit('confirm')
}
</script>

<template>
  <Dialog :open="open" @update:open="handleOpenChange">
    <DialogContent :class="props.class">
      <DialogHeader>
        <DialogTitle>{{ title }}</DialogTitle>
      </DialogHeader>
      
      <div class="py-2">
        <DialogDescription v-if="description" class="text-base text-foreground">
          {{ description }}
        </DialogDescription>
        <slot name="description" v-else />
      </div>

      <slot />
      <DialogFooter>
        <Button variant="outline" @click="handleCancel">
          {{ cancelText }}
        </Button>
        <Button variant="destructive" @click="handleConfirm">
          {{ confirmText }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

