<script setup lang="ts">
import { ref, computed } from 'vue'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

type Login = {
  id: number
  username?: string
  ip?: string
  status: 0 | 1
  userAgent?: string
  createdAt: string
}

const search = ref('')
const logs = ref<Login[]>([
  { id: 1, username: 'Admin', ip: '10.0.0.1', status: 1, userAgent: 'Chrome', createdAt: '2026-01-16 09:00:00' },
  { id: 2, username: '店长王雷', ip: '10.0.0.2', status: 0, userAgent: 'Safari', createdAt: '2026-01-16 09:30:00' },
])

const filtered = computed(() => {
  if (!search.value) return logs.value
  return logs.value.filter(i => (i.username || '').includes(search.value) || (i.ip || '').includes(search.value))
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-bold tracking-tight">登录日志</h1>
      <p class="text-sm text-muted-foreground mt-1">记录用户登录行为。</p>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border p-4">
      <div class="relative max-w-sm w-full">
        <Input v-model="search" placeholder="按账号或IP搜索" />
      </div>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>账号</TableHead>
            <TableHead>IP</TableHead>
            <TableHead>UA</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>时间</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in filtered" :key="item.id">
            <TableCell class="text-sm font-semibold">{{ item.username }}</TableCell>
            <TableCell class="text-xs text-muted-foreground">{{ item.ip }}</TableCell>
            <TableCell class="text-xs text-muted-foreground">{{ item.userAgent }}</TableCell>
            <TableCell>
              <Badge
                variant="outline"
                :color="item.status===1 ? 'green' : 'red'"
              >
                {{ item.status===1 ? '成功' : '失败' }}
              </Badge>
            </TableCell>
            <TableCell class="text-xs text-muted-foreground">{{ item.createdAt }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
