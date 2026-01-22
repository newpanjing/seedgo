import type {AxiosError, AxiosInstance, AxiosResponse, InternalAxiosRequestConfig} from 'axios'
import axios from 'axios'
import {showToast} from '@/lib/message'

// Define the response structure
interface ApiResponse<T = any> {
  code: number
  msg?: string
  message?: string
  data: T
}

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
})

// Request interceptor
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add token if available
    const token = localStorage.getItem('token') || sessionStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data
    // Assuming 0 or 200 is success. Adjust based on actual backend agreement.
    // Common pattern: code 200/0 is success.
    // Also accept 201 (Created)
    if (res.code !== 200 && res.code !== 0 && res.code !== 201) {
      const errorMessage = res.msg || res.message || 'Error'
      showToast(errorMessage, { type: 'error' })
      return Promise.reject(new Error(errorMessage))
    }
    return res.data
  },
  (error: AxiosError<ApiResponse>) => {
    let message = '请求失败'
    const status = error.response?.status
    
    if (status === 401) {
      message = '登录已过期，请重新登录'
      // Clear all auth data
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('currentUser')
      
      showToast(message, { type: 'error' })
      
      // Use window.location to ensure full reset and avoid circular dependencies
      // Delay slightly to let toast show? No, immediate is better for security/UX flow.
      setTimeout(() => {
         window.location.href = '/login'
      }, 500)
      
      return Promise.reject(error)
    } 
    
    if (status === 403) {
      message = '您没有权限执行此操作'
      showToast(message, { type: 'error' })
      return Promise.reject(error)
    }

    if (error.response && error.response.data) {
      message = error.response.data.msg || error.response.data.message || message
    } else if (error.message) {
      message = error.message
    }
    
    showToast(message, { type: 'error' })
    return Promise.reject(error)
  }
)

export default service
