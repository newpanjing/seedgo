<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseChart from '../components/BaseChart.vue'
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const timeRange = ref('week')

// --- Revenue Charts ---
const revenueOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [
    {
      name: '总营收',
      type: 'bar',
      data: [1200, 2000, 1500, 800, 2800, 3500, 3200],
      itemStyle: { color: '#ea580c' }
    },
    {
      name: '会员充值',
      type: 'bar',
      data: [500, 1000, 800, 200, 1500, 2000, 1000],
      itemStyle: { color: '#f97316' }
    }
  ]
}))

const categoryOption = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { top: '5%', left: 'center' },
  series: [
    {
      name: '销售额',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: { show: false, position: 'center' },
      emphasis: {
        label: { show: true, fontSize: 20, fontWeight: 'bold' }
      },
      data: [
        { value: 1048, name: '棋牌室' },
        { value: 735, name: '台球' },
        { value: 580, name: '麻将' },
        { value: 484, name: '饮料零食' }
      ]
    }
  ]
}))

// --- Member Charts ---
const memberGrowthOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', boundaryGap: false, data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'] },
  yAxis: { type: 'value' },
  series: [
    {
      name: '新增会员',
      type: 'line',
      smooth: true,
      areaStyle: {
        opacity: 0.5,
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: '#ea580c' },
            { offset: 1, color: 'rgba(234, 88, 12, 0.01)' }
          ]
        }
      },
      data: [2, 5, 3, 4, 8, 12, 10],
      itemStyle: { color: '#ea580c' }
    }
  ]
}))

const memberLevelOption = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { top: 'bottom' },
  series: [
    {
      name: '会员等级',
      type: 'pie',
      radius: '50%',
  data: [
    { value: 50, name: '普通会员', itemStyle: { color: '#e5e7eb' } },
    { value: 30, name: '银牌会员', itemStyle: { color: '#d1d5db' } },
    { value: 15, name: '金牌会员', itemStyle: { color: '#f59e0b' } },
    { value: 5, name: '钻石会员', itemStyle: { color: '#ea580c' } }
      ]
    }
  ]
}))

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', boundaryGap: false, data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'] },
  yAxis: { type: 'value' },
  series: [
    {
      name: '客流量',
      type: 'line',
      smooth: true,
      areaStyle: { color: '#ea580c', opacity: 0.2 },
      data: [5, 2, 10, 25, 45, 80, 60],
      itemStyle: { color: '#ea580c' }
    }
  ]
}))
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold tracking-tight">数据分析</h1>
      <div class="w-[180px]">
        <Select v-model="timeRange">
          <SelectTrigger>
            <SelectValue placeholder="选择时间范围" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="week">本周</SelectItem>
            <SelectItem value="month">本月</SelectItem>
            <SelectItem value="year">全年</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>

    <Tabs default-value="revenue" class="space-y-4">
      <TabsList>
        <TabsTrigger value="revenue">营收分析</TabsTrigger>
        <TabsTrigger value="members">会员分析</TabsTrigger>
      </TabsList>
      
      <TabsContent value="revenue" class="space-y-4">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>营收统计</CardTitle>
            </CardHeader>
            <CardContent>
              <BaseChart :option="revenueOption" height="350px" />
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>商品销售占比</CardTitle>
            </CardHeader>
            <CardContent>
              <BaseChart :option="categoryOption" height="350px" />
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>客流量趋势</CardTitle>
          </CardHeader>
          <CardContent>
            <BaseChart :option="trendOption" height="350px" />
          </CardContent>
        </Card>
      </TabsContent>
      
      <TabsContent value="members" class="space-y-4">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>会员增长趋势</CardTitle>
            </CardHeader>
            <CardContent>
              <BaseChart :option="memberGrowthOption" height="350px" />
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle>会员等级分布</CardTitle>
            </CardHeader>
            <CardContent>
              <BaseChart :option="memberLevelOption" height="350px" />
            </CardContent>
          </Card>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>
