<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { useMessageHost } from '@/lib/message'

const {
  confirmState,
  alertState,
  toasts,
  notices,
  resolveConfirm,
  closeConfirm,
  resolveAlert,
  closeAlert,
  removeToast,
  removeNotice,
} = useMessageHost()

const handleConfirmClose = (open: boolean) => {
  if (!open && confirmState.open) {
    closeConfirm()
  } else {
    confirmState.open = open
  }
}

const handleAlertClose = (open: boolean) => {
  if (!open && alertState.open) {
    closeAlert()
  } else {
    alertState.open = open
  }
}
</script>

<template>
  <RouterView />

  <ConfirmDialog
    :open="confirmState.open"
    :title="confirmState.title"
    :description="confirmState.message"
    :confirm-text="confirmState.confirmText"
    :cancel-text="confirmState.cancelText"
    @update:open="handleConfirmClose"
    @confirm="resolveConfirm(true)"
    @cancel="resolveConfirm(false)"
  />

  <Dialog :open="alertState.open" @update:open="handleAlertClose">
    <DialogContent class="sm:max-w-[380px]">
      <DialogHeader>
        <DialogTitle>{{ alertState.title }}</DialogTitle>
        <DialogDescription>
          {{ alertState.message }}
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button @click="resolveAlert">
          {{ alertState.confirmText }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <div class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none w-full max-w-sm px-4 sm:px-0">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-center gap-3 rounded-lg border px-4 py-3 text-sm shadow-lg transition-all"
        :class="[
          toast.type === 'success' && 'bg-white border-emerald-200 text-emerald-600 dark:bg-emerald-950/30 dark:border-emerald-900 dark:text-emerald-400',
          toast.type === 'error' && 'bg-white border-red-200 text-red-600 dark:bg-red-950/30 dark:border-red-900 dark:text-red-400',
          toast.type === 'warning' && 'bg-white border-amber-200 text-amber-600 dark:bg-amber-950/30 dark:border-amber-900 dark:text-amber-400',
          toast.type === 'info' && 'bg-white border-slate-200 text-slate-600 dark:bg-slate-900 dark:border-slate-800 dark:text-slate-400',
        ]"
      >
        <!-- Icons -->
        <svg v-if="toast.type === 'success'" class="w-5 h-5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><path d="m9 11 3 3L22 4"/></svg>
        <svg v-if="toast.type === 'error'" class="w-5 h-5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>
        <svg v-if="toast.type === 'warning'" class="w-5 h-5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>
        <svg v-if="toast.type === 'info'" class="w-5 h-5 shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>

        <span class="flex-1 font-medium">{{ toast.message }}</span>
        <button
          type="button"
          class="ml-2 text-current opacity-60 hover:opacity-100 transition-opacity"
          @click="removeToast(toast.id)"
        >
          <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
        </button>
      </div>
    </TransitionGroup>
  </div>

  <div class="fixed bottom-4 right-4 z-40 space-y-2 max-w-sm">
    <div
      v-for="notice in notices"
      :key="notice.id"
      class="rounded-md border px-4 py-3 text-xs shadow-sm bg-card"
      :class="[
        notice.type === 'success' && 'border-emerald-500',
        notice.type === 'error' && 'border-red-500',
        notice.type === 'warning' && 'border-amber-500',
        notice.type === 'info' && 'border-primary/40',
      ]"
    >
      <div class="flex items-start gap-2">
        <span class="flex-1 leading-relaxed">{{ notice.message }}</span>
        <button
          type="button"
          class="ml-2 text-xs opacity-60 hover:opacity-100"
          @click="removeNotice(notice.id)"
        >
          ×
        </button>
      </div>
    </div>
  </div>
</template>

<style>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
