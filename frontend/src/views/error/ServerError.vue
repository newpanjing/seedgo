<template>
  <div class="flex flex-col items-center justify-center h-full bg-background text-foreground">
    <h1 class="text-9xl font-bold text-destructive">500</h1>
    <p class="text-2xl font-medium mt-4">服务器错误</p>
    <p class="text-muted-foreground mt-2">抱歉，服务器暂时无法处理您的请求。</p>
    <Button class="mt-8" @click="handleRetry">重试</Button>
  </div>
</template>

<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useRouter } from 'vue-router'
const router = useRouter()
const handleRetry = () => {
  router.push('/')
}
</script>
