<script setup lang="ts">
import { computed, ref } from 'vue'
import { useDark } from '@vueuse/core'
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Clock, Crown, ShieldCheck, CheckCircle2 } from 'lucide-vue-next'

const isDark = useDark()
const currentPlan = ref<'basic' | 'pro' | 'enterprise'>('pro')
const expireAt = ref('2026-12-31')

const plans = [
  {
    id: 'basic',
    name: '基础版',
    price: '¥ 199/月',
    desc: '适合单店基础经营管理',
    features: ['支持 1 家门店', '基础会员与计费', '基础数据报表'],
    highlight: false,
  },
  {
    id: 'pro',
    name: '专业版',
    price: '¥ 399/月',
    desc: '适合多场景门店的日常经营',
    features: ['支持 3 家门店', '高级计费与套餐', '多角色子账号管理', '经营分析大屏'],
    highlight: true,
  },
  {
    id: 'enterprise',
    name: '旗舰版',
    price: '¥ 799/月',
    desc: '适合连锁品牌与高客流门店',
    features: ['不限门店数量', '专属客户成功服务', '自定义报表导出', '优先技术支持'],
    highlight: false,
  },
]

const daysLeft = computed(() => {
  const end = new Date(expireAt.value).getTime()
  const now = Date.now()
  const diff = Math.max(end - now, 0)
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

const getStatusText = computed(() => {
  if (daysLeft.value === 0) return '已到期，请尽快续费以避免功能受限'
  if (daysLeft.value <= 15) return '即将到期，建议尽快续费保障正常使用'
  if (daysLeft.value <= 30) return '30 天内到期，可提前续费锁定优惠价格'
  return '账户状态正常，可根据需要随时调整套餐'
})

const selectPlan = (id: 'basic' | 'pro' | 'enterprise') => {
  currentPlan.value = id
}

const handleRenew = (planId: string) => {
  console.log('选择续费套餐', planId)
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">续费管理</h1>
        <p class="text-sm text-muted-foreground mt-1">
          查看当前账户的到期时间和使用套餐，选择合适的续费方案。
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card class="lg:col-span-1">
        <CardHeader>
          <CardTitle>当前账户状态</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-muted-foreground">当前套餐</p>
              <p class="text-lg font-semibold flex items-center gap-2">
                <Crown class="w-4 h-4 text-amber-500" />
                {{ currentPlan === 'basic' ? '基础版' : currentPlan === 'pro' ? '专业版' : '旗舰版' }}
              </p>
            </div>
            <Badge variant="outline" color="green" class="flex items-center gap-1">
              <ShieldCheck class="w-3 h-3" />
              正常使用中
            </Badge>
          </div>

          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs text-muted-foreground flex items-center gap-1">
                <Clock class="w-3 h-3" />
                到期时间
              </p>
              <p class="text-sm font-medium">{{ expireAt }}</p>
            </div>
            <div class="text-right">
              <p class="text-xs text-muted-foreground">剩余天数</p>
              <p class="text-xl font-bold text-amber-600">
                {{ daysLeft }}
              </p>
            </div>
          </div>

          <div
            class="rounded-lg px-3 py-2 text-xs leading-relaxed border"
            :class="isDark ? 'bg-amber-900/40 text-amber-200 border-amber-800' : 'bg-amber-50 text-amber-700 border-amber-100'"
          >
            {{ getStatusText }}
          </div>

          <div class="text-xs text-muted-foreground space-y-1">
            <p class="flex items-center gap-1">
              <CheckCircle2 class="w-3 h-3" />
              如需要开具发票，可续费后联系客户成功人员处理。
            </p>
            <p class="flex items-center gap-1">
              <CheckCircle2 class="w-3 h-3" />
              套餐升级将在下一结算周期自动生效，不影响当前使用。
            </p>
          </div>
        </CardContent>
      </Card>

      <Card class="lg:col-span-2">
        <CardHeader>
          <CardTitle>选择续费套餐</CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs default-value="month" class="space-y-4">
            <TabsList>
              <TabsTrigger value="month">按月续费</TabsTrigger>
              <TabsTrigger value="year">按年续费（更优惠）</TabsTrigger>
            </TabsList>
            <TabsContent value="month" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div
                  v-for="plan in plans"
                  :key="plan.id"
                  class="rounded-xl border border-border bg-card p-4 flex flex-col justify-between"
                  :class="plan.highlight ? 'border-amber-400 shadow-sm' : ''"
                >
                  <div class="space-y-3">
                    <div class="flex items-center justify-between">
                      <div class="font-semibold">
                        {{ plan.name }}
                      </div>
                      <Badge
                      v-if="plan.highlight"
                      variant="outline"
                      color="yellow"
                    >
                      推荐
                    </Badge>
                    </div>
                    <div class="text-xl font-bold text-amber-600">
                      {{ plan.price }}
                    </div>
                    <p class="text-xs text-muted-foreground">
                      {{ plan.desc }}
                    </p>
                    <ul class="mt-2 space-y-1 text-xs text-muted-foreground">
                      <li
                        v-for="feature in plan.features"
                        :key="feature"
                        class="flex items-center gap-1"
                      >
                        <CheckCircle2 class="w-3 h-3 text-emerald-500" />
                        <span>{{ feature }}</span>
                      </li>
                    </ul>
                  </div>

                  <div class="mt-4 flex items-center justify-between">
                    <Badge
                      v-if="currentPlan === plan.id"
                      variant="outline"
                      color="green"
                    >
                      当前正在使用
                    </Badge>
                    <div v-else />
                    <Button
                      size="sm"
                      :variant="currentPlan === plan.id ? 'outline' : 'default'"
                      @click="() => { selectPlan(plan.id as 'basic' | 'pro' | 'enterprise'); handleRenew(plan.id); }"
                    >
                      {{ currentPlan === plan.id ? '续费当前套餐' : '切换并续费' }}
                    </Button>
                  </div>
                </div>
              </div>
            </TabsContent>

            <TabsContent value="year">
              <div
                class="rounded-lg border border-dashed p-4 text-xs"
                :class="isDark ? 'border-amber-700 bg-amber-900/40 text-amber-200' : 'border-amber-300 bg-amber-50 text-amber-800'"
              >
                年付套餐与月付价格按比例计算，实际收费以商务沟通为准。此处用于产品展示与交互体验模拟。
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
