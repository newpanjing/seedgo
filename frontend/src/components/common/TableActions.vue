<script setup lang="tsx">
import { Button } from '@/components/ui/button'
import { Eye, Edit, Trash2 } from 'lucide-vue-next'
import { showConfirm } from '@/lib/message'

interface Props {
  showView?: boolean
  showEdit?: boolean
  showDelete?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showView: true,
  showEdit: true,
  showDelete: true
})

const emit = defineEmits<{
  (e: 'view'): void
  (e: 'edit'): void
  (e: 'delete'): void
}>()

const handleDelete = async () => {
  const confirmed = await showConfirm('确认删除', '确定要删除这条记录吗？此操作无法撤销。', {
    confirmText: '确认删除',
    cancelText: '取消'
  })
  if (confirmed) {
    emit('delete')
  }
}

// 核心逻辑：定义一个生成器，把事件预先绑定
const ViewBtn = (props: any) => (
  <Button 
    variant="ghost" 
    size="icon" 
    class="h-8 w-8" 
    title="查看" 
    onClick={() => emit('view')} // 内部直接处理掉 click
    {...props} // 仍然允许父组件传入额外的 class 或 style
  >
    <Eye class="w-4 h-4" />
  </Button>
);

const UpdateBtn = (props: any) => (
  <Button 
    variant="ghost" 
    size="icon" 
    class="h-8 w-8" 
    title="编辑" 
    onClick={() => emit('edit')} // 内部直接处理掉 click
    {...props} // 仍然允许父组件传入额外的 class 或 style
  >
    <Edit class="w-4 h-4" />
  </Button>
);

const DeleteBtn = (props: any) => (
  <Button 
    variant="ghost" 
    size="icon" 
    class="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10" 
    title="删除" 
    onClick={() => handleDelete()} // 内部直接处理掉 click
    {...props} // 仍然允许父组件传入额外的 class 或 style
  >
    <Trash2 class="w-4 h-4" />
  </Button>
);
</script>

<template>
  <div class="flex items-center justify-end gap-1">
    <slot :view="ViewBtn" :update="UpdateBtn" :delete="DeleteBtn">
      <component :is="ViewBtn" />
      <component :is="UpdateBtn" />
      <component :is="DeleteBtn" />
    </slot>
  </div>
</template>
