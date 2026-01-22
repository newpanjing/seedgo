import type {Directive} from 'vue'
import {useAuthStore} from '@/stores/auth'

export const superDirective: Directive = {
  mounted(el) {
    const authStore = useAuthStore()
    const user = authStore.currentUser
    // Check both camelCase (Prisma default) and snake_case (DB column/legacy)
    // Also handle boolean or number (1/0)
    const isSuper = user?.isSuper || user?.is_super

    if (!isSuper) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el) {
    const authStore = useAuthStore()
    const user = authStore.currentUser
    const isSuper = user?.isSuper || user?.is_super

    if (!isSuper) {
      el.parentNode?.removeChild(el)
    }
  }
}
