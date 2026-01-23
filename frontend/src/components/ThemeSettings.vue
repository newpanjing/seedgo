<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Sun, Moon, Monitor, Check, Plus } from 'lucide-vue-next'
import { useColorMode } from '@vueuse/core'
import { useTheme } from '@/composables/useTheme'
import { ref } from 'vue'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const customColorInput = ref<HTMLInputElement | null>(null)

const triggerCustomColor = () => {
  customColorInput.value?.click()
}

const handleCustomColor = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.value) {
    applyThemeColor(target.value)
  }
}

const mode = useColorMode({
  emitAuto: true,
  selector: 'html',
  attribute: 'class',
  modes: {
    dark: 'dark',
    light: '',
  },
})

const { currentTheme, applyThemeColor, themeColors } = useTheme()

// Filter colors to match requirements
const colorList = [
  "#4F46E5", "#334155", "#1D4ED8", "#B91C1C", "#0F766E", "#71717A", // 企业级
  "#F59E0B", "#F43F5E", "#F97316", "#FACC15", "#D946EF", "#EF4444", // 暖色系
  "#0EA5E9", "#10B981", "#8B5CF6", "#06B6D4", "#84CC16",   // 冷色系
];

const displayColors = computed(() => {
  const result: Record<string, { color: string }> = {}
  colorList.forEach(color => {
    result[color] = { color }
  })
  return result
})

const handleClose = () => {
  emit('update:open', false)
}
</script>

<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)">
    <DialogContent class="sm:max-w-xl rounded-3xl backdrop-blur-xl border border-border/50 dark:border-white/10 shadow-2xl duration-300 ease-[cubic-bezier(0.34,1.56,0.64,1)]">
      <DialogHeader>
        <DialogTitle>主题设置</DialogTitle>
      </DialogHeader>

      <!-- Main Content -->
      <div class="space-y-8 py-4">
        
        <!-- Appearance Mode -->
        <div class="space-y-3">
          <h3 class="text-sm font-medium text-muted-foreground">外观模式</h3>
          <div class="grid grid-cols-3 gap-3">
            <button 
              type="button"
              @click="mode = 'light'"
              class="relative flex items-center justify-center gap-2 rounded-lg border p-2.5 transition-all duration-300 hover:bg-muted/50 active:scale-95"
              :class="mode === 'light' ? 'border-primary bg-primary/5' : 'border-muted hover:border-border'"
            >
              <Sun class="w-4 h-4 transition-colors duration-300" :class="mode === 'light' ? 'text-primary' : 'text-muted-foreground'" />
              <span class="text-sm font-medium transition-colors duration-300" :class="mode === 'light' ? 'text-primary' : 'text-muted-foreground'">浅色</span>
            </button>

            <button 
              type="button"
              @click="mode = 'dark'"
              class="relative flex items-center justify-center gap-2 rounded-lg border p-2.5 transition-all duration-300 hover:bg-muted/50 active:scale-95"
              :class="mode === 'dark' ? 'border-primary bg-primary/5' : 'border-muted hover:border-border'"
            >
              <Moon class="w-4 h-4 transition-colors duration-300" :class="mode === 'dark' ? 'text-primary' : 'text-muted-foreground'" />
              <span class="text-sm font-medium transition-colors duration-300" :class="mode === 'dark' ? 'text-primary' : 'text-muted-foreground'">深色</span>
            </button>

            <button 
              type="button"
              @click="mode = 'auto'"
              class="relative flex items-center justify-center gap-2 rounded-lg border p-2.5 transition-all duration-300 hover:bg-muted/50 active:scale-95"
              :class="mode === 'auto' ? 'border-primary bg-primary/5' : 'border-muted hover:border-border'"
            >
              <Monitor class="w-4 h-4 transition-colors duration-300" :class="mode === 'auto' ? 'text-primary' : 'text-muted-foreground'" />
              <span class="text-sm font-medium transition-colors duration-300" :class="mode === 'auto' ? 'text-primary' : 'text-muted-foreground'">跟随系统</span>
            </button>
          </div>
        </div>

        <!-- Accent Color -->
        <div class="space-y-4">
          <h3 class="text-sm font-medium text-muted-foreground">强调色</h3>
          <div class="flex items-center gap-4 flex-wrap">
            <button
              v-for="(item, key) in displayColors"
              :key="key"
              type="button"
              @click="applyThemeColor(key as any)"
              class="group relative flex h-10 w-10 items-center justify-center rounded-full transition-all duration-300 hover:scale-110 focus:outline-none active:scale-90"
            >
              <span 
                class="absolute inset-0 rounded-full transition-all duration-300"
                :class="currentTheme === key ? 'scale-110 ring-2 ring-primary ring-offset-2 ring-offset-background' : 'scale-100'"
                :style="{ backgroundColor: item.color }"
              ></span>
              <Transition
                enter-active-class="transition-all duration-300 ease-[cubic-bezier(0.34,1.56,0.64,1)]"
                enter-from-class="opacity-0 scale-0 -rotate-90"
                enter-to-class="opacity-100 scale-100 rotate-0"
                leave-active-class="transition-all duration-200 ease-in"
                leave-from-class="opacity-100 scale-100 rotate-0"
                leave-to-class="opacity-0 scale-0 -rotate-90"
              >
                <Check v-if="currentTheme === key" class="z-10 w-5 h-5 text-white drop-shadow-sm" />
              </Transition>
            </button>

            <!-- Custom Color Button -->
            <div class="relative">
              <input
                ref="customColorInput"
                type="color"
                class="absolute inset-0 opacity-0 w-full h-full cursor-pointer"
                @input="handleCustomColor"
              />
              <button
                type="button"
                @click="triggerCustomColor"
                class="group relative flex h-10 w-10 items-center justify-center rounded-full transition-all duration-300 hover:scale-110 focus:outline-none active:scale-90 border border-dashed border-muted-foreground/50 hover:border-primary hover:bg-primary/5"
              >
                <Plus class="w-5 h-5 text-muted-foreground transition-colors group-hover:text-primary" />
              </button>
            </div>
          </div>
        </div>

        <!-- Live Preview Card -->
        <div class="space-y-3">
          <h3 class="text-sm font-medium text-muted-foreground">预览</h3>
          <div class="rounded-xl border border-border bg-muted/30 p-2 transition-colors duration-500">
            <div class="relative overflow-hidden rounded-lg border border-border bg-background shadow-sm transition-all duration-500 hover:shadow-md">
              <div class="flex items-center border-b border-border p-1.5 bg-muted/20 transition-colors duration-300">
                <div class="flex space-x-1.5">
                  <div class="h-1.5 w-1.5 rounded-full bg-red-500/80 transition-transform duration-300 hover:scale-125"></div>
                  <div class="h-1.5 w-1.5 rounded-full bg-yellow-500/80 transition-transform duration-300 hover:scale-125"></div>
                  <div class="h-1.5 w-1.5 rounded-full bg-green-500/80 transition-transform duration-300 hover:scale-125"></div>
                </div>
              </div>
              <div class="p-2 space-y-2">
                <div class="h-1 w-1/3 rounded-full bg-muted animate-pulse"></div>
                <div class="rounded-lg border border-border bg-card p-2 shadow-sm transition-colors duration-300">
                   <div class="flex items-center gap-2">
                      <div class="h-6 w-6 rounded-full bg-primary/10 flex items-center justify-center transition-colors duration-300">
                          <span class="text-primary text-xs font-bold transition-colors duration-300">A</span>
                      </div>
                      <div class="space-y-1 flex-1">
                          <div class="h-1.5 w-12 rounded-full bg-primary/20 animate-pulse"></div>
                          <div class="h-1 w-8 rounded-full bg-muted-foreground/20 animate-pulse"></div>
                      </div>
                      <Button size="sm" class="h-6 text-[10px] rounded-full px-2 shadow-sm transition-all duration-300 active:scale-95 hover:scale-105 hover:shadow-md">
                        Button
                      </Button>
                   </div>
                </div>
                <div class="space-y-1">
                  <div class="h-1 w-full rounded-full bg-muted/50 animate-pulse delay-75"></div>
                  <div class="h-1 w-5/6 rounded-full bg-muted/50 animate-pulse delay-150"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

      </div>

      <DialogFooter class="gap-2 sm:gap-0">
        <Button variant="outline" @click="handleClose">取消</Button>
        <Button @click="handleClose">保存</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
