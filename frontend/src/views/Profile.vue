<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { 
  User as UserIcon, 
  Mail, 
  Shield, 
  Key, 
  Camera,
  Save,
  Eye,
  EyeOff,
  Check,
  Info
} from 'lucide-vue-next'
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter
} from '@/components/ui/dialog'
import { getUserInfo, updateUserInfo, changePassword } from '@/api/auth'
import { showToast } from '@/lib/message'

const authStore = useAuthStore()
const user = ref(authStore.currentUser)

const formData = ref({
  name: '',
  email: '',
  phone: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const errors = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const isSaving = ref(false)
const isPasswordDialogOpen = ref(false)
const isShaking = ref(false)

const triggerShake = () => {
  isShaking.value = true
  setTimeout(() => {
    isShaking.value = false
  }, 500)
}

// Password Visibility State
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// Password Strength Logic
const passwordStrength = computed(() => {
  const pwd = formData.value.newPassword
  if (!pwd) return { score: 0, label: '', color: 'bg-slate-200' }

  let score = 0
  if (pwd.length >= 8) score++
  if (/[A-Z]/.test(pwd)) score++
  if (/[a-z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd)) score++
  if (/[^A-Za-z0-9]/.test(pwd)) score++

  // Cap at 4 for 4 bars
  if (score > 4) score = 4

  const labels = ['弱', '一般', '强', '非常强']
  const colors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-green-500']
  
  // Adjust score index for labels/colors (0 score means empty or very weak)
  const index = Math.max(0, score - 1)

  return {
    score, // 0-5
    label: labels[index] || '弱',
    color: colors[index] || 'bg-slate-200'
  }
})

const passwordSuggestions = computed(() => {
  const pwd = formData.value.newPassword
  const suggestions = []
  
  if (pwd.length < 8) suggestions.push('至少 8 个字符')
  if (!/[A-Z]/.test(pwd)) suggestions.push('包含大写字母')
  if (!/[a-z]/.test(pwd)) suggestions.push('包含小写字母')
  if (!/[0-9]/.test(pwd)) suggestions.push('包含数字')
  if (!/[^A-Za-z0-9]/.test(pwd)) suggestions.push('包含特殊字符')
  
  return suggestions
})

const loadUserProfile = async () => {
  try {
    const res: any = await getUserInfo()
    if (res) {
      user.value = res
      // Update store and storage to sync data across components (e.g. MainLayout)
      authStore.currentUser = res
      const storageKey = 'currentUser'
      if (localStorage.getItem(storageKey)) {
        localStorage.setItem(storageKey, JSON.stringify(res))
      } else if (sessionStorage.getItem(storageKey)) {
        sessionStorage.setItem(storageKey, JSON.stringify(res))
      }
      
      formData.value.name = res.realName || res.name || ''
      formData.value.email = res.email || ''
      formData.value.phone = res.phone || ''
    }
  } catch (error) {
    console.error('Failed to load profile', error)
  }
}

onMounted(() => {
  loadUserProfile()
})

const saveProfile = async () => {
  isSaving.value = true
  try {
    await updateUserInfo({
      name: formData.value.name,
      phone: formData.value.phone,
      email: formData.value.email
    })
    showToast('个人资料已更新', { type: 'success' })
    await loadUserProfile() // Reload to update UI
  } catch (error) {
    console.error(error)
  } finally {
    isSaving.value = false
  }
}

const validateForm = () => {
  let isValid = true
  errors.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  }

  if (!formData.value.currentPassword) {
    errors.value.currentPassword = '请输入当前密码'
    isValid = false
  }

  if (!formData.value.newPassword) {
    errors.value.newPassword = '请输入新密码'
    isValid = false
  } else if (formData.value.newPassword.length < 8) {
    errors.value.newPassword = '密码长度至少为 8 位'
    isValid = false
  }

  if (!formData.value.confirmPassword) {
    errors.value.confirmPassword = '请确认新密码'
    isValid = false
  } else if (formData.value.newPassword !== formData.value.confirmPassword) {
    errors.value.confirmPassword = '两次输入的密码不一致'
    isValid = false
  }

  return isValid
}

const handlePasswordChange = async () => {
  if (!validateForm()) {
    triggerShake()
    return
  }
  
  isSaving.value = true
  try {
    await changePassword({
      oldPassword: formData.value.currentPassword,
      newPassword: formData.value.newPassword
    })
    showToast('密码已修改', { type: 'success' })
    isPasswordDialogOpen.value = false
    formData.value.currentPassword = ''
    formData.value.newPassword = ''
    formData.value.confirmPassword = ''
    errors.value = { currentPassword: '', newPassword: '', confirmPassword: '' }
  } catch (error) {
    console.error(error)
    triggerShake()
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <Card>
      <CardContent class="p-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 class="text-2xl font-bold flex items-center">
            <UserIcon class="w-8 h-8 mr-3 text-primary" />
            个人中心
          </h1>
          <p class="mt-1 text-sm text-muted-foreground">
            管理您的个人信息及账户安全设置
          </p>
        </div>
      </CardContent>
    </Card>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Profile Card -->
      <div class="lg:col-span-1 space-y-6">
        <Card>
          <CardContent class="p-6 text-center">
            <div class="relative inline-block">
              <div class="h-24 w-24 mx-auto rounded-full bg-gradient-to-tr from-blue-500 to-indigo-600 p-[3px] shadow-lg">
                <Avatar class="h-full w-full">
                  <AvatarFallback class="bg-background text-4xl font-bold text-primary">
                    {{ formData.name?.charAt(0) || user?.username?.charAt(0) || 'A' }}
                  </AvatarFallback>
                </Avatar>
              </div>
              <Button size="icon" variant="secondary" class="absolute bottom-0 right-0 h-8 w-8 rounded-full shadow-md" title="更换头像">
                <Camera class="w-4 h-4" />
              </Button>
            </div>
            
            <h2 class="mt-4 text-xl font-bold">{{ formData.name || user?.username || 'User' }}</h2>
            <p class="text-sm text-muted-foreground">{{ user?.isSuper ? '超级管理员' : '普通用户' }}</p>
            
            <div class="mt-6 flex justify-center space-x-4">
              <div class="text-center px-4 py-2 bg-muted/50 rounded-lg">
                <div class="text-sm font-bold">--</div>
                <div class="text-xs text-muted-foreground">登录次数</div>
              </div>
              <div class="text-center px-4 py-2 bg-muted/50 rounded-lg">
                <div class="text-sm font-bold">--</div>
                <div class="text-xs text-muted-foreground">服务时长</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Contact Info (Read Only) -->
        <Card>
          <CardHeader>
             <CardTitle class="text-lg font-medium flex items-center">
                <Shield class="w-5 h-5 mr-2 text-muted-foreground" />
                账户信息
             </CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center text-sm">
              <UserIcon class="w-4 h-4 mr-3 text-muted-foreground" />
              <span class="text-muted-foreground w-20">用户ID</span>
              <span class="font-medium">{{ user?.id?.toString().slice(0, 8) || 'unknown' }}</span>
            </div>
            <div class="flex items-center text-sm">
              <Mail class="w-4 h-4 mr-3 text-muted-foreground" />
              <span class="text-muted-foreground w-20">邮箱</span>
              <span class="font-medium">{{ formData.email || '未设置' }}</span>
            </div>
            <div class="flex items-center text-sm">
              <Shield class="w-4 h-4 mr-3 text-muted-foreground" />
              <span class="text-muted-foreground w-20">角色</span>
              <div class="flex flex-wrap gap-1">
                <span v-if="user?.isSuper" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">
                  超级管理员
                </span>
                <template v-else-if="user?.roles && user.roles.length > 0">
                  <span v-for="role in user.roles" :key="role.roleId || role.id" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">
                    {{ role.roleName || role.name || '未知角色' }}
                  </span>
                </template>
                <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">
                  {{ user?.role || '普通用户' }}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Right Column: Settings Forms -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Basic Info Form -->
        <Card>
          <CardHeader>
             <CardTitle>基本资料</CardTitle>
          </CardHeader>
          <CardContent>
            <form @submit.prevent="saveProfile" class="space-y-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <Label>显示名称</Label>
                  <Input v-model="formData.name" />
                </div>
                <div class="space-y-2">
                  <Label>联系电话</Label>
                  <Input v-model="formData.phone" type="tel" />
                </div>
              </div>
              
              <div class="space-y-2">
                <Label>邮箱</Label>
                <Input v-model="formData.email" type="email" />
              </div>

              <div class="flex justify-end">
                <Button type="submit" :disabled="isSaving">
                  <Save class="w-4 h-4 mr-2" />
                  {{ isSaving ? '保存中...' : '保存更改' }}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <!-- Security Settings -->
        <Card>
          <CardHeader>
            <CardTitle class="flex items-center justify-between">
              <div class="flex items-center">
                <Key class="w-5 h-5 mr-2 text-muted-foreground" />
                安全设置
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent>
             <div class="flex items-center justify-between">
                <div>
                   <div class="text-sm font-medium">登录密码</div>
                   <div class="text-xs text-muted-foreground">定期修改密码可以保护账号安全</div>
                </div>
                <Dialog v-model:open="isPasswordDialogOpen">
                  <DialogTrigger as-child>
                    <Button variant="outline">修改密码</Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>修改密码</DialogTitle>
                    </DialogHeader>
                    
                    <!-- Password Rules Alert -->
                    <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 flex gap-3 text-sm text-blue-700 dark:text-blue-300">
                      <Info class="w-5 h-5 shrink-0 mt-0.5" />
                      <div>
                        <p class="font-semibold mb-1">密码安全规则</p>
                        <ul class="list-disc list-inside space-y-0.5 text-xs opacity-90">
                          <li>长度至少 8 位</li>
                          <li>需包含数字、大小写字母</li>
                          <li>建议包含特殊字符</li>
                        </ul>
                      </div>
                    </div>

                    <form 
                      @submit.prevent="handlePasswordChange" 
                      class="space-y-4 py-2"
                      :class="{ 'animate-shake': isShaking }"
                    >
                      <div class="space-y-2">
                        <Label :class="{ 'text-red-500': errors.currentPassword }">当前密码</Label>
                        <div class="relative">
                          <Input 
                            v-model="formData.currentPassword" 
                            :type="showCurrentPassword ? 'text' : 'password'" 
                            class="pr-10 transition-colors"
                            :class="{ 'border-red-500 focus-visible:ring-red-500': errors.currentPassword }"
                          />
                          <button 
                            type="button" 
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                            @click="showCurrentPassword = !showCurrentPassword"
                          >
                            <Eye v-if="!showCurrentPassword" class="w-4 h-4" />
                            <EyeOff v-else class="w-4 h-4" />
                          </button>
                        </div>
                        <span v-if="errors.currentPassword" class="text-xs text-red-500">{{ errors.currentPassword }}</span>
                      </div>
                      <div class="space-y-2">
                        <Label :class="{ 'text-red-500': errors.newPassword }">新密码</Label>
                        <div class="relative">
                          <Input 
                            v-model="formData.newPassword" 
                            :type="showNewPassword ? 'text' : 'password'" 
                            class="pr-10 transition-colors"
                            :class="{ 'border-red-500 focus-visible:ring-red-500': errors.newPassword }"
                          />
                          <button 
                            type="button" 
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                            @click="showNewPassword = !showNewPassword"
                          >
                            <Eye v-if="!showNewPassword" class="w-4 h-4" />
                            <EyeOff v-else class="w-4 h-4" />
                          </button>
                        </div>
                        <span v-if="errors.newPassword" class="text-xs text-red-500">{{ errors.newPassword }}</span>
                        
                        <!-- Password Strength Indicator -->
                        <div v-if="formData.newPassword" class="space-y-2 mt-2">
                          <div class="flex gap-1 h-1">
                            <div 
                              v-for="i in 4" 
                              :key="i"
                              class="flex-1 rounded-full transition-colors duration-300"
                              :class="i <= passwordStrength.score ? passwordStrength.color : 'bg-slate-100 dark:bg-slate-800'"
                            ></div>
                          </div>
                          <p class="text-xs text-right" :class="passwordStrength.score > 0 ? 'text-foreground' : 'text-muted-foreground'">
                            强度: {{ passwordStrength.label }}
                          </p>
                          
                          <!-- Suggestions -->
                          <div class="grid grid-cols-2 gap-1 mt-2">
                            <div 
                              v-for="(suggestion, idx) in passwordSuggestions" 
                              :key="idx"
                              class="text-xs flex items-center text-muted-foreground"
                            >
                              <div class="w-1 h-1 rounded-full bg-slate-300 mr-2"></div>
                              {{ suggestion }}
                            </div>
                            <div 
                              v-if="passwordSuggestions.length === 0" 
                              class="col-span-2 text-xs text-green-600 flex items-center"
                            >
                              <Check class="w-3 h-3 mr-1" /> 密码强度符合要求
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="space-y-2">
                        <Label :class="{ 'text-red-500': errors.confirmPassword }">确认新密码</Label>
                        <div class="relative">
                          <Input 
                            v-model="formData.confirmPassword" 
                            :type="showConfirmPassword ? 'text' : 'password'" 
                            class="pr-10 transition-colors"
                            :class="{ 'border-red-500 focus-visible:ring-red-500': errors.confirmPassword }"
                          />
                          <button 
                            type="button" 
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                            @click="showConfirmPassword = !showConfirmPassword"
                          >
                            <Eye v-if="!showConfirmPassword" class="w-4 h-4" />
                            <EyeOff v-else class="w-4 h-4" />
                          </button>
                        </div>
                        <span v-if="errors.confirmPassword" class="text-xs text-red-500">{{ errors.confirmPassword }}</span>
                      </div>
                      <DialogFooter>
                         <Button type="submit" :disabled="isSaving" class="ml-auto">
                            <Save class="w-4 h-4 mr-2" />
                            确认修改
                         </Button>
                      </DialogFooter>
                    </form>
                  </DialogContent>
                </Dialog>
             </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-4px); }
  20%, 40%, 60%, 80% { transform: translateX(4px); }
}

.animate-shake {
  animation: shake 0.5s cubic-bezier(.36,.07,.19,.97) both;
}
</style>