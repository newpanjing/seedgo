<script setup lang="ts">
import { computed } from 'vue'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Save } from 'lucide-vue-next'

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
const confirm = () => {
    emit('confirm')
}
</script>
<template>
    <Dialog v-model:open="isOpen">
        <DialogContent class="sm:max-w-[600px]">
            <DialogHeader>
                <DialogTitle>{{ title }}</DialogTitle>
                <DialogDescription v-if="description">
                    {{ description }}
                </DialogDescription>
            </DialogHeader>
            <slot />
            <DialogFooter>
                <Button variant="ghost" @click="close">取消</Button>
                <Button @click="confirm">
                    <Save class="w-4 h-4 mr-2" />
                    保存
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>