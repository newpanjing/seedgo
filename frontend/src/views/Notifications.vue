<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { AlertTriangle, Bell, CheckCircle, Clock, Info } from 'lucide-vue-next'
import { ref } from 'vue'

const notifications = ref([
  {
    id: 1,
    title: '系统维护通知',
    content: '系统将于今晚凌晨 02:00 进行例行维护，预计耗时 2 小时，期间将无法访问系统。',
    time: '10分钟前',
    type: 'warning',
    read: false
  },
  {
    id: 2,
    title: '新订单提醒',
    content: '收到一个新的订单 #ORD-20240115-001，请及时处理。',
    time: '30分钟前',
    type: 'info',
    read: false
  },
  {
    id: 3,
    title: '会员注册成功',
    content: '新会员 张三 (13800138000) 已完成注册。',
    time: '1小时前',
    type: 'success',
    read: true
  },
  {
    id: 4,
    title: '库存预警',
    content: '商品 "可口可乐 330ml" 库存不足 10 件，请及时补货。',
    time: '2小时前',
    type: 'error',
    read: true
  },
  {
    id: 5,
    title: '续费提醒',
    content: '您的租户服务将于 7 天后到期，请及时续费以免影响使用。',
    time: '1天前',
    type: 'warning',
    read: true
  }
])

const markAsRead = (id: number) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    notification.read = true
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
}

const getIcon = (type: string) => {
  switch (type) {
    case 'warning': return AlertTriangle
    case 'error': return Info
    case 'success': return CheckCircle
    default: return Bell
  }
}

const getIconColor = (type: string) => {
  switch (type) {
    case 'warning': return 'text-orange-500'
    case 'error': return 'text-red-500'
    case 'success': return 'text-green-500'
    default: return 'text-blue-500'
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold tracking-tight">消息通知</h2>
        <p class="text-muted-foreground mt-2">查看和管理您的所有系统通知消息。</p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="markAllAsRead">
          全部已读
        </Button>
      </div>
    </div>

    <div class="grid gap-4">
      <Card v-for="notification in notifications" :key="notification.id" class="transition-all hover:shadow-md" :class="{ 'bg-muted/30': notification.read }">
        <div class="p-4 flex items-start gap-4">
          <div class="p-2 rounded-full bg-background border shrink-0">
            <component :is="getIcon(notification.type)" class="w-5 h-5" :class="getIconColor(notification.type)" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-2 mb-1">
              <h4 class="font-semibold text-sm" :class="{ 'text-muted-foreground': notification.read }">
                {{ notification.title }}
              </h4>
              <span class="text-xs text-muted-foreground whitespace-nowrap flex items-center">
                <Clock class="w-3 h-3 mr-1" />
                {{ notification.time }}
              </span>
            </div>
            <p class="text-sm text-muted-foreground mb-2 line-clamp-2">
              {{ notification.content }}
            </p>
            <div class="flex items-center justify-between">
              <Badge variant="secondary" class="text-xs font-normal">
                {{ notification.type === 'warning' ? '系统通知' : (notification.type === 'error' ? '重要提醒' : '业务消息') }}
              </Badge>
              <Button 
                v-if="!notification.read" 
                variant="ghost" 
                size="sm" 
                class="h-6 text-xs"
                @click="markAsRead(notification.id)"
              >
                标为已读
              </Button>
            </div>
          </div>
          <div v-if="!notification.read" class="w-2 h-2 rounded-full bg-red-500 shrink-0 mt-2"></div>
        </div>
      </Card>
    </div>
  </div>
</template>
