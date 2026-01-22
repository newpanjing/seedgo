<script setup lang="ts">
import { computed } from 'vue'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

const props = defineProps({
    title: {
        type: String,
        default: '表单'
    },
    description: {
        type: String,
        default: ''
    },
    visible: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['update:visible', 'confirm'])

const isOpen = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
})

const close = () => {
    isOpen.value = false
}
</script>
<template>
    <Sheet v-model:open="isOpen">
        <SheetContent class="w-[400px] sm:w-[540px] flex flex-col h-full p-0 gap-0">
            <div class="px-6 py-4 flex-none">
                <SheetHeader>
                    <SheetTitle>{{ title }}</SheetTitle>
                    <SheetDescription v-if="description">
                        {{ description }}
                    </SheetDescription>
                </SheetHeader>
            </div>
            <Separator />
            <div class="flex-1 overflow-y-auto px-6 py-4">
                <slot />
            </div>
            <div class="px-6 py-4 flex-none bg-background">
                <div class="flex justify-end gap-2">
                    <Button class="w-full" variant="outline" @click="close">关闭</Button>
                </div>
            </div>
        </SheetContent>
    </Sheet>
</template>