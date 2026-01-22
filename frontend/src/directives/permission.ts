import type {Directive} from 'vue'
import {usePermissionStore, usePerms} from '@/stores/permission'

export const permission: Directive = {
  mounted(el, binding) {
    const { value } = binding
    const permissionStore = usePermissionStore()

    if (!value) {
      throw new Error(`need permission! Like v-permission="'system:user:add'"`)
    }

    const requiredPermissions = Array.isArray(value) ? value : [value]

    // Check if user has ANY of the required permissions (OR logic)
    // You can change to EVERY (AND logic) if needed, but usually it's OR for visibility
    const hasPermission = requiredPermissions.some(code => permissionStore.hasPermission(code))

    if (!hasPermission) {
      el.parentNode?.removeChild(el)
    }
  }
}

export const hasPerm: Directive = {
  mounted(el, binding) {
    const perms = usePerms() as any
    const action = binding.arg

    const hasPermission = perms[action]

    if (!hasPermission) {
      // el.parentNode?.removeChild(el)
      el.style.display = 'none'
    }else{
      el.style.display = 'block'
    }

  }
}