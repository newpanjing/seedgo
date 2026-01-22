<script setup lang="ts">
import { ref } from 'vue'
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const systemName = ref('Seedgo')
const logoUrl = ref('')

const handleSave = () => {
  console.log('保存系统配置', {
    systemName: systemName.value,
    logoUrl: logoUrl.value,
  })
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold tracking-tight">系统配置</h1>
        <p class="text-sm text-muted-foreground mt-1">
          配置品牌 Logo 和系统显示名称，统一前后台展示风格。
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <Card class="lg:col-span-2">
        <CardHeader>
          <CardTitle>基础信息</CardTitle>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-2">
            <Label for="system-name">系统显示名称</Label>
            <Input
              id="system-name"
              v-model="systemName"
              placeholder="例如：智享门店管家"
            />
            <p class="text-xs text-muted-foreground">
              将显示在左上角 Logo 区域以及登录页标题位置。
            </p>
          </div>

          <div class="space-y-2">
            <Label>品牌 Logo</Label>
            <div class="flex items-center gap-4">
              <div class="h-16 w-16 rounded-xl border border-dashed border-border flex items-center justify-center bg-muted/40 text-muted-foreground text-xs">
                Logo
              </div>
              <div class="flex-1 space-y-2">
                <Input
                  v-model="logoUrl"
                  placeholder="粘贴 Logo 图片地址，或后续接入上传能力"
                />
                <p class="text-xs text-muted-foreground">
                  推荐使用正方形或圆角图标，尺寸不小于 128×128 像素。
                </p>
              </div>
            </div>
          </div>

          <div class="flex justify-end">
            <Button @click="handleSave">
              保存配置
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>展示预览</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div>
              <p class="text-xs text-muted-foreground mb-2">侧边栏效果</p>
              <div class="flex items-center space-x-3 px-3 py-2 rounded-lg border border-border bg-muted/40">
                <div class="h-8 w-8 rounded-lg flex items-center justify-center bg-primary text-primary-foreground text-sm">
                  {{ systemName.charAt(0) || 'S' }}
                </div>
                <div class="truncate text-sm font-semibold">
                  {{ systemName || 'Seedgo' }}
                </div>
              </div>
            </div>

            <div>
              <p class="text-xs text-muted-foreground mb-2">登录页标题</p>
              <div class="px-3 py-2 rounded-lg border border-dashed border-border">
                <p class="text-lg font-bold">
                  欢迎使用 {{ systemName || 'Seedgo' }}
                </p>
                <p class="text-xs text-muted-foreground mt-1">
                  统一门店、会员与计费管理的一体化工作台。
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
