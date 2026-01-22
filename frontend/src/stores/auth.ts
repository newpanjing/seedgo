import {defineStore} from 'pinia'
import {computed, ref} from 'vue'
import {login as apiLogin, type LoginParams} from '@/api/auth'
// import type { User } from '../types/user' // This might be Member user, we might need Admin user type.

export const useAuthStore = defineStore('auth', () => {
  // We'll use 'any' for now or define a proper User type for the logged-in user (likely Admin/Staff)
  const currentUser = ref<any | null>(null)
  const token = ref<string | null>(localStorage.getItem('token') || sessionStorage.getItem('token'))

  // Initialize from localStorage or sessionStorage if available
  const storedUser = localStorage.getItem('currentUser') || sessionStorage.getItem('currentUser')
  if (storedUser) {
    try {
      currentUser.value = JSON.parse(storedUser)
    } catch (e) {
      localStorage.removeItem('currentUser')
      sessionStorage.removeItem('currentUser')
    }
  }

  const isAuthenticated = computed(() => !!token.value)

  const login = async (params: LoginParams, remember: boolean = false) => {
    try {
      const res = await apiLogin(params)
      // Assuming res contains { token, user }
      if (res.token) {
        token.value = res.token
        if (remember) {
          localStorage.setItem('token', res.token)
        } else {
          sessionStorage.setItem('token', res.token)
        }
      }
      
      if (res.user) {
        currentUser.value = res.user
        if (remember) {
          localStorage.setItem('currentUser', JSON.stringify(res.user))
        } else {
          sessionStorage.setItem('currentUser', JSON.stringify(res.user))
        }
      }
      
      return true
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    currentUser.value = null
    token.value = null
    localStorage.removeItem('currentUser')
    localStorage.removeItem('token')
    sessionStorage.removeItem('currentUser')
    sessionStorage.removeItem('token')
  }

  return {
    currentUser,
    isAuthenticated,
    login,
    logout
  }
})
