import {reactive, ref} from 'vue'

type ToastType = 'info' | 'success' | 'warning' | 'error'

type ConfirmState = {
  open: boolean
  title: string
  message: string
  confirmText: string
  cancelText: string
  resolve: ((value: boolean) => void) | null
}

type AlertState = {
  open: boolean
  title: string
  message: string
  confirmText: string
  resolve: (() => void) | null
}

type ToastItem = {
  id: number
  type: ToastType
  message: string
}

const confirmState = reactive<ConfirmState>({
  open: false,
  title: '',
  message: '',
  confirmText: '确认',
  cancelText: '取消',
  resolve: null,
})

const alertState = reactive<AlertState>({
  open: false,
  title: '',
  message: '',
  confirmText: '我知道了',
  resolve: null,
})

const toasts = ref<ToastItem[]>([])
const notices = ref<ToastItem[]>([])

let seed = 1

export function showConfirm(title: string, message: string, options?: Partial<Pick<ConfirmState, 'confirmText' | 'cancelText'>>): Promise<boolean> {
  confirmState.title = title
  confirmState.message = message
  confirmState.confirmText = options?.confirmText ?? '确认'
  confirmState.cancelText = options?.cancelText ?? '取消'
  confirmState.open = true

  return new Promise<boolean>(resolve => {
    confirmState.resolve = resolve
  })
}

export function showAlert(title: string, message: string, options?: Partial<Pick<AlertState, 'confirmText'>>): Promise<void> {
  alertState.title = title
  alertState.message = message
  alertState.confirmText = options?.confirmText ?? '我知道了'
  alertState.open = true

  return new Promise<void>(resolve => {
    alertState.resolve = resolve
  })
}

export function showToast(message: string, options?: { type?: ToastType; duration?: number }) {
  const id = seed++
  const type: ToastType = options?.type ?? 'success'
  const duration = options?.duration ?? 3000

  toasts.value.push({
    id,
    type,
    message,
  })

  if (duration > 0) {
    window.setTimeout(() => {
      removeToast(id)
    }, duration)
  }
}

export function showNotice(message: string, options?: { type?: ToastType; duration?: number }) {
  const id = seed++
  const type: ToastType = options?.type ?? 'info'
  const duration = options?.duration ?? 0

  notices.value.push({
    id,
    type,
    message,
  })

  if (duration > 0) {
    window.setTimeout(() => {
      removeNotice(id)
    }, duration)
  }
}

function removeToast(id: number) {
  toasts.value = toasts.value.filter(item => item.id !== id)
}

function removeNotice(id: number) {
  notices.value = notices.value.filter(item => item.id !== id)
}

function resolveConfirm(value: boolean) {
  if (confirmState.resolve) {
    confirmState.resolve(value)
  }
  confirmState.resolve = null
  confirmState.open = false
}

function resolveAlert() {
  if (alertState.resolve) {
    alertState.resolve()
  }
  alertState.resolve = null
  alertState.open = false
}

function closeConfirm() {
  resolveConfirm(false)
}

function closeAlert() {
  resolveAlert()
}

export function useMessageHost() {
  return {
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
  }
}

