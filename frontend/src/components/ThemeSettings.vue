<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Sun, Moon } from 'lucide-vue-next'
import { useDark } from '@vueuse/core'
import { useTheme } from '@/composables/useTheme'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const isDark = useDark()
const { currentTheme, applyThemeColor, themeColors } = useTheme()
</script>

<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>主题设置</DialogTitle>
      </DialogHeader>

      <div class="space-y-6 py-2">
        <div class="space-y-3">
          <p class="text-sm font-medium">颜色模式</p>
          <div class="grid grid-cols-3 gap-3">
            <button type="button"
              class="flex flex-col items-center justify-center gap-2 rounded-lg border p-3 hover:bg-muted/60 transition-all"
              :class="!isDark ? 'border-primary bg-primary/5' : 'border-border'" @click="isDark = false">
              <Sun class="w-6 h-6" />
              <span class="text-xs text-muted-foreground">浅色模式</span>
            </button>
            <button type="button"
              class="flex flex-col items-center justify-center gap-2 rounded-lg border p-3 hover:bg-muted/60 transition-all"
              :class="isDark ? 'border-primary bg-primary/5' : 'border-border'" @click="isDark = true">
              <Moon class="w-6 h-6" />
              <span class="text-xs text-muted-foreground">深色模式</span>
            </button>
            <!-- System mode is handled by vueuse/core useDark default behavior if we clear storage, but for now simple toggle is enough. Or we can just have 2 options. -->
          </div>
        </div>

        <div class="space-y-3">
          <p class="text-sm font-medium">主题颜色</p>
          <div class="flex flex-wrap gap-4">
            <button v-for="(item, key) in themeColors" :key="key" type="button"
              class="relative flex h-8 w-8 items-center justify-center rounded-full transition-all hover:scale-110 focus:outline-none"
              @click="applyThemeColor(key)" :title="item.label">
              <span v-if="currentTheme === key" class="absolute -inset-1.5 rounded-full border-2 border-primary" />
              <span class="h-full w-full rounded-full shadow-sm ring-1 ring-inset ring-black/10"
                :style="{ backgroundColor: item.color }" />
            </button>
          </div>
        </div>
      </div>

      <DialogFooter class="justify-end">
        <Button variant="outline" @click="$emit('update:open', false)">
          关闭
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
