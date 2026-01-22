<script setup lang="ts">
import { computed } from 'vue'
import BaseChart from '../components/BaseChart.vue'
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { 
  DollarSign, 
  Users, 
  ShoppingBag, 
  Activity, 
  TrendingUp, 
  TrendingDown,
  Plus,
  UserPlus,
  FileText,
  Clock,
  CheckCircle2,
  MoreHorizontal
} from 'lucide-vue-next'

const stats = [
  {
    title: '今日营收',
    value: '¥ 2,850',
    trend: '+12.5%',
    trendUp: true,
    icon: DollarSign,
    color: 'text-primary',
    bg: 'bg-primary/10'
  },
  {
    title: '活跃会员',
    value: '128',
    trend: '+4.2%',
    trendUp: true,
    icon: Users,
    color: 'text-primary',
    bg: 'bg-primary/10'
  },
  {
    title: '今日订单',
    value: '45',
    trend: '-2.1%',
    trendUp: false,
    icon: ShoppingBag,
    color: 'text-primary',
    bg: 'bg-primary/10'
  },
  {
    title: '房间使用率',
    value: '85%',
    trend: '+5.4%',
    trendUp: true,
    icon: Activity,
    color: 'text-primary',
    bg: 'bg-primary/10'
  }
]

// --- Charts Configuration ---
const revenueOption = computed(() => ({
  grid: { left: '3%', right: '4%', bottom: '3%', top: '10%', containLabel: true },
  tooltip: { trigger: 'axis' },
  xAxis: { 
    type: 'category', 
    data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
    axisLine: { show: false },
    axisTick: { show: false }
  },
  yAxis: { 
    type: 'value',
    splitLine: { lineStyle: { type: 'dashed', opacity: 0.5 } }
  },
  series: [
    {
      name: '营收',
      type: 'line',
      smooth: true,
      symbol: 'none',
      areaStyle: {
        opacity: 0.2,
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: '#ea580c' },
            { offset: 1, color: 'rgba(234, 88, 12, 0.05)' }
          ]
        }
      },
      lineStyle: { width: 3, color: '#ea580c' },
      data: [1500, 2300, 2240, 2180, 1350, 1470, 2600]
    }
  ]
}))

const occupancyOption = computed(() => ({
  grid: { left: '3%', right: '4%', bottom: '3%', top: '10%', containLabel: true },
  tooltip: { trigger: 'axis' },
  xAxis: { 
    type: 'category', 
    data: ['V01', 'V02', 'V03', 'V05', 'V06', 'V08'],
    axisLine: { show: false },
    axisTick: { show: false }
  },
  yAxis: { 
    type: 'value',
    splitLine: { lineStyle: { type: 'dashed', opacity: 0.5 } }
  },
  series: [
    {
      name: '时长(小时)',
      type: 'bar',
      barWidth: '40%',
      itemStyle: { borderRadius: [4, 4, 0, 0], color: '#4b5563' },
      data: [8, 12, 5, 9, 14, 6]
    }
  ]
}))

import { useRouter } from 'vue-router'

const router = useRouter()

// --- Quick Actions ---
const quickActions = [
  { name: '开单结账', icon: Plus, color: 'bg-primary/10 text-primary border border-primary/20 hover:bg-primary/15', action: () => router.push('/billing') },
  { name: '新增会员', icon: UserPlus, color: 'bg-muted text-foreground border border-border hover:bg-muted/80', action: () => router.push('/members') },
  { name: '快速核销', icon: CheckCircle2, color: 'bg-primary/10 text-primary border border-primary/20 hover:bg-primary/15', action: () => {} },
  { name: '查看报表', icon: FileText, color: 'bg-muted text-foreground border border-border hover:bg-muted/80', action: () => router.push('/analytics') },
]

// --- Recent Activity ---
const activities = [
  { id: 1, content: '会员 张三 完成了充值', time: '10分钟前', icon: DollarSign, color: 'bg-primary/10 text-primary' },
  { id: 2, content: 'V05 房间开始计费', time: '25分钟前', icon: Clock, color: 'bg-muted text-foreground' },
  { id: 3, content: '新增商品 \"可乐\" 库存 50', time: '1小时前', icon: ShoppingBag, color: 'bg-muted text-foreground' },
  { id: 4, content: '会员 李四 预约了 V08', time: '2小时前', icon: Users, color: 'bg-primary/10 text-primary' },
]
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold tracking-tight">仪表盘</h1>
      <div class="text-sm text-muted-foreground">
        {{ new Date().toLocaleDateString() }}
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card v-for="stat in stats" :key="stat.title">
        <CardContent class="p-6">
          <div class="flex justify-between items-start">
            <div>
              <p class="text-sm font-medium text-muted-foreground">{{ stat.title }}</p>
              <h3 class="mt-2 text-3xl font-bold">{{ stat.value }}</h3>
            </div>
            <div :class="`p-3 rounded-lg ${stat.bg}`">
              <component :is="stat.icon" :class="`w-6 h-6 ${stat.color}`" />
            </div>
          </div>
          <div class="mt-4 flex items-center text-sm">
            <span :class="stat.trendUp ? 'text-amber-600' : 'text-rose-500'" class="flex items-center font-medium">
              <component :is="stat.trendUp ? TrendingUp : TrendingDown" class="w-4 h-4 mr-1" />
              {{ stat.trend }}
            </span>
            <span class="text-muted-foreground ml-2">较昨日</span>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Charts & Quick Actions -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Chart (Revenue) -->
      <Card class="lg:col-span-2">
          <CardHeader class="flex flex-row items-center justify-between">
          <CardTitle>本周营收趋势</CardTitle>
          <Button variant="ghost" size="sm" class="text-primary hover:text-primary/80">查看详情</Button>
        </CardHeader>
        <CardContent>
          <BaseChart :option="revenueOption" height="300px" />
        </CardContent>
      </Card>

      <!-- Quick Actions & Secondary Chart -->
      <div class="space-y-6">
        <!-- Quick Actions -->
        <Card>
          <CardHeader>
            <CardTitle>快捷操作</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-2 gap-4">
              <Button
                v-for="action in quickActions"
                :key="action.name"
                :class="`${action.color} h-24 flex flex-col items-center justify-center gap-2 hover:opacity-90 transition-opacity`"
                @click="action.action"
              >
                <component :is="action.icon" class="w-6 h-6 mb-1" />
                <span>{{ action.name }}</span>
              </Button>
            </div>
          </CardContent>
        </Card>

        <!-- Room Usage Chart -->
        <Card>
          <CardHeader>
            <CardTitle>今日房间时长</CardTitle>
          </CardHeader>
          <CardContent>
            <BaseChart :option="occupancyOption" height="200px" />
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- Recent Activity -->
    <Card>
      <CardHeader>
        <CardTitle>近期活动</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="space-y-6">
          <div v-for="activity in activities" :key="activity.id" class="flex items-start">
            <div :class="`flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center ${activity.color}`">
              <component :is="activity.icon" class="w-5 h-5" />
            </div>
            <div class="ml-4 flex-1">
              <p class="text-sm font-medium">{{ activity.content }}</p>
              <p class="text-xs text-muted-foreground mt-1">{{ activity.time }}</p>
            </div>
            <Button variant="ghost" size="icon" class="text-muted-foreground">
              <MoreHorizontal class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
