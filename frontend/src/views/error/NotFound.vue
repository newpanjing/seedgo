<template>
  <div class="flex flex-col items-center justify-center h-full bg-background text-foreground">
    <h1 class="text-9xl font-bold text-primary">404</h1>
    <p class="text-2xl font-medium mt-4">页面未找到</p>
    <p class="text-muted-foreground mt-2">抱歉，您访问的页面不存在。</p>
    <Button class="mt-8" @click="$router.push('/')">返回首页</Button>
  </div>
</template>

<script setup lang="ts">
import { Button } from '@/components/ui/button'
</script>
