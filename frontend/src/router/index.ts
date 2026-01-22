import {createRouter, createWebHistory} from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'
import {useAuthStore} from '../stores/auth'
import {usePermissionStore} from '../stores/permission'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/Login.vue')
    },
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('../views/Dashboard.vue')
        },
        {
          path: 'products',
          name: 'products',
          component: () => import('../views/Products.vue')
        },
        {
          path: 'billing',
          name: 'billing',
          component: () => import('../views/Billing.vue')
        },
        {
          path: 'analytics',
          name: 'analytics',
          component: () => import('../views/Analytics.vue')
        },
        {
          path: 'members',
          name: 'members',
          component: () => import('../views/Members.vue')
        },
        {
          path: 'notifications',
          name: 'notifications',
          component: () => import('../views/Notifications.vue')
        },
        {
          path: 'system/tenants',
          name: 'tenantManagement',
          component: () => import('../views/system/Tenants.vue')
        },

        {
          path: 'system/users',
          name: 'users',
          component: () => import('../views/auth/Users.vue')
        },
        {
          path: 'system/roles',
          name: 'roleManagement',
          component: () => import('../views/auth/Role.vue')
        },
        {
          path: 'system/stores',
          name: 'storeManagement',
          component: () => import('../views/system/StoreManagement.vue')
        },
        {
          path: 'system/renewal',
          name: 'renewal',
          component: () => import('../views/renewal/Renewal.vue')
        },
        {
          path: 'system/permissions',
          name: 'permissions',
          component: () => import('../views/auth/Permissions.vue')
        },
        {
          path: 'system/operation-logs',
          name: 'operationLog',
          component: () => import('../views/system/OperationLog.vue')
        },
        {
          path: 'system/login-logs',
          name: 'loginLog',
          component: () => import('../views/login/LoginLog.vue')
        },
        {
          path: 'system/dicts',
          name: 'dictManagement',
          component: () => import('../views/system/DictManagement.vue')
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('../views/Profile.vue')
        },
        {
          path: '/403',
          name: 'forbidden',
          component: () => import('../views/error/Forbidden.vue')
        },
        {
          path: '/500',
          name: 'server-error',
          component: () => import('../views/error/ServerError.vue')
        },
        {
          path: '/:pathMatch(.*)*',
          name: 'not-found',
          component: () => import('../views/error/NotFound.vue')
        }
      ]
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()
  
  if (to.path === '/500' || to.name === 'server-error') {
    return next()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (authStore.isAuthenticated) {
    // 1. Fetch permissions if needed
    if (permissionStore.menuTree.length === 0) {
      try {
        await permissionStore.fetchPermissions()
      } catch (error) {
        console.error('Failed to fetch permissions', error)
        // return next('/500')
      }
    }

    // 2. White List Check
    const whiteList = ['/profile', '/notifications', '/403', '/404', '/500', '/', '/dashboard'];
    if (whiteList.includes(to.path) || to.name === 'dashboard' || to.name === 'not-found' || to.name === 'forbidden' || to.name === 'server-error') {
      return next()
    }

    // 3. Permission Check
    // Check if the target route path is in the allowed menu paths
    // We check to.path (resolved path) against menuPaths.
    // NOTE: This assumes menuPaths contains the full path like /system/users
    
    if (permissionStore.menuPaths.has(to.path)) {
      next()
    } else {
      // If path is valid route but not in permissions -> 403
      next('/403')
    }
  } else {
    next()
  }
})

export default router
