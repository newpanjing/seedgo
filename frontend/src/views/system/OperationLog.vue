<script setup lang="ts">
import { ref, computed } from 'vue'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

type Log = {
  id: number
  module: string
  action: string
  username?: string
  ip?: string
  status: 0 | 1
  createdAt: string
}

const search = ref('')
const logs = ref<Log[]>([
  { id: 1, module: '系统管理', action: '新增角色', username: 'Admin', ip: '10.0.0.1', status: 1, createdAt: '2026-01-16 10:00:00' },
  { id: 2, module: '门店管理', action: '编辑门店', username: 'Admin', ip: '10.0.0.1', status: 1, createdAt: '2026-01-16 10:10:00' },
  { id: 3, module: '账号管理', action: '停用账号', username: '店长王雷', ip: '10.0.0.2', status: 0, createdAt: '2026-01-16 10:30:00' },
])

const filtered = computed(() => {
  if (!search.value) return logs.value
  return logs.value.filter(i => i.module.includes(search.value) || i.action.includes(search.value) || (i.username || '').includes(search.value))
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">操作日志</h1>
        <p class="text-sm text-muted-foreground mt-1">记录后台重要操作行为。</p>
      </div>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border p-4">
      <div class="relative max-w-sm w-full">
        <Input v-model="search" placeholder="按模块、操作或用户名搜索" />
      </div>
    </div>

    <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>模块</TableHead>
            <TableHead>操作</TableHead>
            <TableHead>用户</TableHead>
            <TableHead>IP</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>时间</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in filtered" :key="item.id">
            <TableCell class="text-sm font-semibold">{{ item.module }}</TableCell>
            <TableCell class="text-sm">{{ item.action }}</TableCell>
            <TableCell class="text-sm">{{ item.username }}</TableCell>
            <TableCell class="text-xs text-muted-foreground">{{ item.ip }}</TableCell>
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
