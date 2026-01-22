<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { showToast } from '@/lib/message'
import { Eye, EyeOff } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)
const showPassword = ref(false)
const rememberMe = ref(false)

onMounted(() => {
  const savedUsername = localStorage.getItem('savedUsername')
  if (savedUsername) {
    username.value = savedUsername
    rememberMe.value = true
  }
})

const handleLogin = async () => {
  if (!username.value || !password.value) {
    showToast('请输入用户名和密码', { type: 'warning' })
    return
  }

  loading.value = true
  try {
    await authStore.login({
      username: username.value,
      password: password.value
    }, rememberMe.value)

    if (rememberMe.value) {
      localStorage.setItem('savedUsername', username.value)
    } else {
      localStorage.removeItem('savedUsername')
    }

    showToast('登录成功', { type: 'success' })
    router.push('/')
  } catch (e: any) {
    // Error is handled in request interceptor, but we can log it
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex bg-background text-foreground">
    <!-- Left Side: Banner & Intro -->
    <div class="hidden lg:flex lg:w-1/2 bg-muted relative overflow-hidden items-center justify-center">
      <!-- Abstract Background Pattern -->
      <div class="absolute inset-0 opacity-40">
        <div class="absolute top-0 -left-4 w-72 h-72 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
        <div class="absolute top-0 -right-4 w-72 h-72 bg-blue-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
        <div class="absolute -bottom-8 left-20 w-72 h-72 bg-indigo-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
      </div>
      
      <!-- Content Container -->
      <div class="relative z-10 w-full max-w-md px-8">
        <div class="mb-10">
           <!-- Logo & Title -->
           <div class="flex items-center space-x-3 mb-6">
             <div class="h-10 w-10 flex items-center justify-center rounded-lg bg-primary text-primary-foreground shadow-lg">
              <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
             </div>
             <span class="text-2xl font-bold tracking-tight">Seedgo</span>
           </div>
           
           <h1 class="text-4xl font-extrabold mb-4 leading-tight">
             让管理更智慧<br/>让服务更温情
           </h1>
           <p class="text-muted-foreground text-lg leading-relaxed">
             一站式智能管家解决方案，为您提供全场景、全流程的高效管理体验。
           </p>
        </div>
        
        <!-- Feature Pills -->
        <div class="flex flex-wrap gap-3">
          <div class="px-4 py-2 bg-background/80 backdrop-blur rounded-full shadow-sm border border-border flex items-center space-x-2">
             <div class="w-2 h-2 rounded-full bg-blue-500"></div>
             <span class="text-sm font-medium">数据可视化</span>
          </div>
          <div class="px-4 py-2 bg-background/80 backdrop-blur rounded-full shadow-sm border border-border flex items-center space-x-2">
             <div class="w-2 h-2 rounded-full bg-indigo-500"></div>
             <span class="text-sm font-medium">智能排班</span>
          </div>
          <div class="px-4 py-2 bg-background/80 backdrop-blur rounded-full shadow-sm border border-border flex items-center space-x-2">
             <div class="w-2 h-2 rounded-full bg-purple-500"></div>
             <span class="text-sm font-medium">多端协同</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="flex-1 flex flex-col justify-center py-12 px-4 sm:px-6 lg:px-20 xl:px-24 bg-background">
       <div class="mx-auto w-full max-w-sm lg:w-96">
          <!-- Mobile Logo -->
          <div class="lg:hidden mb-10">
             <div class="flex items-center space-x-2">
               <div class="h-8 w-8 flex items-center justify-center rounded-lg bg-primary text-primary-foreground">
                 <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                 </svg>
               </div>
               <span class="text-xl font-bold">Seedgo</span>
             </div>
          </div>

          <div class="mb-8">
             <h2 class="text-2xl font-bold">欢迎回来</h2>
             <p class="mt-2 text-sm text-muted-foreground">
               请输入用户名和密码登录
             </p>
          </div>

          <form class="space-y-6" @submit.prevent="handleLogin">
            <div class="space-y-2">
              <Label for="username">账号</Label>
              <Input 
                id="username" 
                v-model="username" 
                type="text" 
                placeholder="手机号 / 用户名" 
                required 
              />
            </div>

            <div class="space-y-2">
              <Label for="password">密码</Label>
              <div class="relative">
                <Input 
                  id="password" 
                  v-model="password" 
                  :type="showPassword ? 'text' : 'password'" 
                  placeholder="请输入密码" 
                  class="pr-10"
                  required 
                />
                <button 
                  type="button" 
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground focus:outline-none"
                  @click="showPassword = !showPassword"
                >
                  <Eye v-if="!showPassword" class="h-4 w-4" />
                  <EyeOff v-else class="h-4 w-4" />
                </button>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <input 
                  id="remember" 
                  v-model="rememberMe" 
                  type="checkbox" 
                  class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                >
                <Label for="remember" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  记住我
                </Label>
              </div>
              <a href="#" class="text-sm font-medium text-primary hover:underline">
                忘记密码?
              </a>
            </div>

            <div>
              <Button type="submit" class="w-full" :disabled="loading">
                <span v-if="loading" class="mr-2">
                  <svg class="animate-spin h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </span>
                {{ loading ? '登录中...' : '登 录' }}
              </Button>
            </div>
          </form>

          <p class="mt-6 text-center text-xs text-muted-foreground">
            登录即代表您同意 <a href="#" class="text-primary hover:underline">服务条款</a> 和 <a href="#" class="text-primary hover:underline">隐私协议</a>
          </p>
       </div>
    </div>
  </div>
</template>
