<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import { useResizeObserver } from '@vueuse/core'
import { useDark } from '@vueuse/core'

const props = defineProps<{
  option: any
  height?: string
}>()

const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null
const isDark = useDark()

const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value, isDark.value ? 'dark' : undefined)
    chart.setOption({
      ...props.option,
      backgroundColor: 'transparent'
    })
  }
}

watch(() => props.option, (newOption) => {
  chart?.setOption(newOption)
}, { deep: true })

watch(isDark, () => {
  if (chart) {
    chart.dispose()
    initChart()
  }
})

useResizeObserver(chartRef, () => {
  chart?.resize()
})

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  chart?.dispose()
})
</script>

<template>
  <div ref="chartRef" :style="{ height: height || '400px', width: '100%' }"></div>
</template>
