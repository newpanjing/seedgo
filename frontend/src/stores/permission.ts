import {defineStore} from 'pinia'
import {computed, ref} from 'vue'
import {getPermissionTree, type Permission} from '@/api/permission'
import * as icons from 'lucide-vue-next'
import {useAuthStore} from './auth'
import {useRoute} from 'vue-router'
import router from '@/router'

export const usePermissionStore = defineStore('permission', () => {
  const menuTree = ref<any[]>([])
  const rawPermissions = ref<Permission[]>([])
  const permissions = ref<Set<string>>(new Set())
  const menuPaths = ref<Set<string>>(new Set())
  const isLoading = ref(false)

  // Resolve icon string to component
  const resolveIcon = (iconName?: string) => {
    if (!iconName) return (icons as any)['Circle']
    
    // Try exact match
    if ((icons as any)[iconName]) {
      return (icons as any)[iconName]
    }

    // Try PascalCase (e.g. "user" -> "User")
    const pascalName = iconName.charAt(0).toUpperCase() + iconName.slice(1)
    if ((icons as any)[pascalName]) {
      return (icons as any)[pascalName]
    }

    // Try kebab-case to PascalCase (e.g. "user-check" -> "UserCheck")
    const camelToPascal = iconName.replace(/(^\w|-\w)/g, (g) => g.replace('-', '').toUpperCase())
    if ((icons as any)[camelToPascal]) {
      return (icons as any)[camelToPascal]
    }
    
    // Default fallback
    return (icons as any)['Circle']
  }

  // Recursive function to process tree for sidebar
  const processMenuTree = (items: Permission[], parentId: string | number | null = null): any[] => {
    return items
      .filter(item => {
        // Only show visible menus
        // Handle both number (1) and boolean (true)
        const isVisible = item.visible === 1 || item.visible === true
        if (item.type !== 1 || !isVisible) return false

        // Filter by parentId
        if (parentId === null) {
          // Top level: parentId is 0, -1, null, undefined, or empty string
          return !item.parentId || item.parentId === 0 || item.parentId === -1 || item.parentId === '0' || item.parentId === '-1'
        }
        
        // Child level: parentId matches
        return item.parentId == parentId
      })
      .map(item => {
        const menu: any = {
          name: item.name,
          href: item.path || '',
          icon: resolveIcon(item.icon),
        }

        // Try to find children in the current list (Flat List pattern)
        let children = processMenuTree(items, item.id)

        // If no children found in flat list, check item.children (Tree pattern)
        if (children.length === 0 && item.children && item.children.length > 0) {
          children = processMenuTree(item.children, item.id)
        }

        if (children.length > 0) {
          menu.children = children
        }

        return menu
      })
  }

  // Recursive function to extract all permission codes and paths
  const extractPermissions = (items: Permission[]) => {
    items.forEach(item => {
      if (item.permissionCode) {
        permissions.value.add(item.permissionCode)
      }
      if (item.path) {
        // Ensure path starts with /
        const path = item.path.startsWith('/') ? item.path : '/' + item.path
        menuPaths.value.add(path)
      }
      if (item.children && item.children.length > 0) {
        extractPermissions(item.children)
      }
    })
  }

  const fetchPermissions = async () => {
    if (isLoading.value) return
    isLoading.value = true
    try {
      const res = await getPermissionTree() as any
      // res is the array because request interceptor returns res.data

      const items = Array.isArray(res) ? res : []
      console.log('Raw permissions:', items)
      permissions.value.clear()
      menuPaths.value.clear()
      rawPermissions.value = items // Store raw tree
      extractPermissions(items)
      menuTree.value = processMenuTree(items)

    } catch (error) {
      console.error('Failed to fetch permissions', error)
      menuTree.value = []
      rawPermissions.value = []
      permissions.value.clear()
      menuPaths.value.clear()
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const hasPermission = (code: string) => {
    return permissions.value.has(code)
  }

  const getModulePermissions = (path: string): string[] => {
    const codes: string[] = []

    const findNode = (nodes: Permission[], targetPath: string): Permission | undefined => {
      for (const node of nodes) {
        const nodePath = node.path?.startsWith('/') ? node.path : '/' + (node.path || '');
        const target = targetPath.startsWith('/') ? targetPath : '/' + targetPath;

        // Exact match or match without trailing slash
        if (node.type === 1 && (nodePath === target || nodePath === target + '/')) {
          return node;
        }
        if (node.children) {
          const found = findNode(node.children, targetPath);
          if (found) return found;
        }
      }
      return undefined;
    }

    const node = findNode(rawPermissions.value, path)

    if (node) {
      // If the menu node itself has a code, include it
      if (node.permissionCode) {
        codes.push(node.permissionCode)
      }

      if (node.children) {
        node.children.forEach(child => {
          if (child.type === 2 && child.permissionCode) {
            codes.push(child.permissionCode)
          }
        })
      }
    }
    return codes
  }

  const clearPermissions = () => {
    menuTree.value = []
    rawPermissions.value = []
    permissions.value.clear()
    menuPaths.value.clear()
  }

  return {
    menuTree,
    permissions,
    menuPaths,
    isLoading,
    fetchPermissions,
    hasPermission,
    getModulePermissions,
    clearPermissions
  }
})


export interface ModulePermissions {
  create: boolean
  update: boolean
  delete: boolean
  view: boolean
  import: boolean
  export: boolean
  [key: string]: boolean
}


export const usePerms = ():any => {
  const store = usePermissionStore()
  const authStore = useAuthStore()
  const route = useRoute() || router.currentRoute.value

  const perms = computed<ModulePermissions>(() => {
    // If super user, return all permissions true
    if (authStore.currentUser?.isSuper) {
      return {
        create: true,
        update: true,
        delete: true,
        view: true,
        import: true,
        export: true
      }
    }

    const codes = store.getModulePermissions(route.path)
    const p: ModulePermissions = {
      create: false,
      update: false,
      delete: false,
      view: false,
      import: false,
      export: false
    }

    codes.forEach(code => {
      const parts = code.split(':')
      const action = parts[parts.length - 1]

      if (['create', 'add'].includes(action)) p.create = true
      if (['update', 'edit'].includes(action)) p.update = true
      if (['delete', 'remove'].includes(action)) p.delete = true
      if (['view', 'list', 'query', 'get'].includes(action)) p.view = true
      if (['import'].includes(action)) p.import = true
      if (['export'].includes(action)) p.export = true

      p[action] = true
    })
    return p
  })

  return {
    ...perms.value,
    hasPermission: (code: string) => {
      if (authStore.currentUser?.isSuper) return true
      return store.hasPermission(code)
    },
    //当前模块的所有权限
    permissions: computed(() => store.getModulePermissions(route.path))
  }
}
