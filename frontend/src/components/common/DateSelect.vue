<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Calendar as CalendarIcon, Clock, ChevronLeft, ChevronRight, X, ChevronsLeft, ChevronsRight } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

const props = withDefaults(defineProps<{
  modelValue?: string | Date
  type?: 'date' | 'time' | 'datetime'
  placeholder?: string
  format?: string
  disabled?: boolean
  clearable?: boolean
  minDate?: string | Date
  maxDate?: string | Date
}>(), {
  type: 'date',
  clearable: true
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | undefined): void
  (e: 'change', value: string | undefined): void
}>()

const open = ref(false)
const now = new Date()
const currentYear = ref(now.getFullYear())
const currentMonth = ref(now.getMonth()) // 0-11
const selectedDate = ref<Date | null>(null)

// Time state
const time = ref({
  h: 0,
  m: 0,
  s: 0
})

// Formatting helpers
const pad = (n: number) => n.toString().padStart(2, '0')

const formatDate = (date: Date, type: string) => {
  const y = date.getFullYear()
  const m = pad(date.getMonth() + 1)
  const d = pad(date.getDate())
  const h = pad(date.getHours())
  const min = pad(date.getMinutes())
  const s = pad(date.getSeconds())

  if (type === 'date') return `${y}-${m}-${d}`
  if (type === 'time') return `${h}:${min}:${s}`
  return `${y}-${m}-${d} ${h}:${min}:${s}`
}

const parseDate = (val: string | Date | undefined): Date | null => {
  if (!val) return null
  if (val instanceof Date) return val
  const d = new Date(val)
  return isNaN(d.getTime()) ? null : d
}

// Watch modelValue
watch(() => props.modelValue, (val) => {
  const date = parseDate(val)
  if (date) {
    selectedDate.value = date
    currentYear.value = date.getFullYear()
    currentMonth.value = date.getMonth()
    time.value = {
      h: date.getHours(),
      m: date.getMinutes(),
      s: date.getSeconds()
    }
  } else {
    selectedDate.value = null
    // Don't reset current view on clear, keep it where it was or reset to now?
    // resetting to now is safer
    // currentYear.value = now.getFullYear()
    // currentMonth.value = now.getMonth()
  }
}, { immediate: true })

// Calendar Logic
const days = ['日', '一', '二', '三', '四', '五', '六']

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  
  const daysInMonth = lastDay.getDate()
  const startDayOfWeek = firstDay.getDay() // 0-6
  
  const res = []
  
  // Previous month days
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = startDayOfWeek - 1; i >= 0; i--) {
    let m = month - 1
    let y = year
    if (m < 0) {
      m = 11
      y--
    }
    res.push({
      day: prevMonthLastDay - i,
      month: m,
      year: y,
      currentMonth: false
    })
  }
  
  // Current month days
  for (let i = 1; i <= daysInMonth; i++) {
    res.push({
      day: i,
      month: month,
      year: year,
      currentMonth: true
    })
  }
  
  // Next month days
  const remaining = 42 - res.length
  for (let i = 1; i <= remaining; i++) {
    let m = month + 1
    let y = year
    if (m > 11) {
      m = 0
      y++
    }
    res.push({
      day: i,
      month: m,
      year: y,
      currentMonth: false
    })
  }
  
  return res
})

const prevMonth = () => {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

const nextMonth = () => {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

const prevYear = () => currentYear.value--
const nextYear = () => currentYear.value++

const isSelected = (day: any) => {
  if (!selectedDate.value) return false
  return selectedDate.value.getDate() === day.day &&
         selectedDate.value.getMonth() === day.month && // Note: day.month handles overflow correctly? 
         // logic above: month-1 for prev can be -1? No, new Date handles it but logic above sets manual values.
         // Let's fix logic:
         // In calendarDays: 
         // if month is 0, prev month is -1? No.
         // Let's rely on full Date comparison or careful index.
         // Simplified: day.month stored as 0-11 or actual index.
         // For prev month of Jan (0): year-1, month 11.
         // For next month of Dec (11): year+1, month 0.
         // My logic in calendarDays handles year rollover.
         // So I should match year/month/day.
         
         // Wait, day.month in calendarDays:
         // prev: month - 1. If month=0, it's -1. 
         // next: month + 1. If month=11, it's 12.
         // I should normalize this in calendarDays or here.
         // Let's normalize in calendarDays for easier comparison.
         
         // Actually, let's just construct a string or timestamp for comparison.
         // Or normalize day.month/year in generation.
         // Re-checking generation:
         // year: month===0 ? year-1 : year. Correct.
         // month: month-1. If month=0, becomes -1. This is technically wrong for display if used directly, but for comparison we need to handle it.
         // Let's fix generation to store correct 0-11 month.
         
         // Fix in generation:
         // prev: let m = month - 1; let y = year; if(m<0) { m=11; y--; }
         // next: let m = month + 1; let y = year; if(m>11) { m=0; y++; }
         
         // I will implement this fix inside computed.
         
         selectedDate.value.getFullYear() === day.year &&
         selectedDate.value.getMonth() === day.month
}

const isToday = (day: any) => {
  const d = new Date()
  return d.getDate() === day.day &&
         d.getMonth() === day.month &&
         d.getFullYear() === day.year
}

const isDisabled = (day: any) => {
  const current = new Date(day.year, day.month, day.day)
  
  if (props.minDate) {
    const min = parseDate(props.minDate)
    if (min) {
      min.setHours(0, 0, 0, 0)
      if (current < min) return true
    }
  }
  
  if (props.maxDate) {
    const max = parseDate(props.maxDate)
    if (max) {
      max.setHours(0, 0, 0, 0)
      if (current > max) return true
    }
  }
  
  return false
}

const handleDateSelect = (day: any) => {
  if (isDisabled(day)) return
  const newDate = new Date(day.year, day.month, day.day)
  
  // Preserve time if datetime
  if (props.type === 'datetime' || props.type === 'time') {
    newDate.setHours(time.value.h)
    newDate.setMinutes(time.value.m)
    newDate.setSeconds(time.value.s)
  } else {
    // For 'date', maybe reset time to 00:00:00?
    // Usually yes for pure date.
    newDate.setHours(0, 0, 0, 0)
  }
  
  selectedDate.value = newDate
  
  if (props.type === 'date') {
    emitValue()
    open.value = false
  } else {
    // For datetime, just select date, user confirms or changes time?
    // Common UX: Select date, then time.
    // Or just update internal value.
    emitValue()
  }
}

// Time Logic
const timeOptions = computed(() => {
  const hours = Array.from({ length: 24 }, (_, i) => i)
  const minutes = Array.from({ length: 60 }, (_, i) => i)
  const seconds = Array.from({ length: 60 }, (_, i) => i)
  return { hours, minutes, seconds }
})

const updateTime = (type: 'h' | 'm' | 's', val: number) => {
  time.value[type] = val
  if (selectedDate.value) {
    const newDate = new Date(selectedDate.value)
    if (type === 'h') newDate.setHours(val)
    if (type === 'm') newDate.setMinutes(val)
    if (type === 's') newDate.setSeconds(val)
    selectedDate.value = newDate
    emitValue()
  } else {
    // If no date selected yet, maybe select today?
    // Or just keep time state waiting for date.
    // Usually for datetime, we need a date.
    // If type is 'time', we might use a dummy date.
    if (props.type === 'time') {
      const d = new Date()
      d.setHours(time.value.h, time.value.m, time.value.s)
      selectedDate.value = d
      emitValue()
    }
  }
}

const emitValue = () => {
  if (!selectedDate.value) {
    emit('update:modelValue', undefined)
    emit('change', undefined)
    return
  }
  const val = formatDate(selectedDate.value, props.type)
  emit('update:modelValue', val)
  emit('change', val)
}

const handleClear = (e: Event) => {
  e.stopPropagation()
  selectedDate.value = null
  emit('update:modelValue', undefined)
  emit('change', undefined)
}

const handleConfirm = () => {
  open.value = false
}

const handleNow = () => {
  const d = new Date()
  selectedDate.value = d
  currentYear.value = d.getFullYear()
  currentMonth.value = d.getMonth()
  time.value = { h: d.getHours(), m: d.getMinutes(), s: d.getSeconds() }
  emitValue()
  open.value = false
}

const displayText = computed(() => {
  if (!selectedDate.value) return ''
  return formatDate(selectedDate.value, props.type)
})

</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button variant="outline" :class="cn('justify-start text-left font-normal hover:scale-100 active:scale-100', !selectedDate && 'text-muted-foreground')" :disabled="disabled">
        <CalendarIcon v-if="type !== 'time'" class="mr-2 h-4 w-4" />
        <Clock v-else class="mr-2 h-4 w-4" />
        <span class="truncate flex-1">
          {{ displayText || placeholder || (type === 'time' ? '选择时间' : '选择日期') }}
        </span>
        <div v-if="selectedDate && clearable" @click="handleClear" class="ml-2 rounded-sm hover:bg-muted p-0.5 cursor-pointer">
          <X class="h-4 w-4 opacity-50 hover:opacity-100" />
        </div>
      </Button>
    </PopoverTrigger>
    <PopoverContent :class="cn('p-0', type === 'time' ? 'w-auto' : 'w-[280px]')" align="start">
      <!-- Date View -->
      <div v-if="type !== 'time'" class="p-3">
        <!-- Header -->
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-1">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="prevYear"><ChevronsLeft class="h-4 w-4" /></Button>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="prevMonth"><ChevronLeft class="h-4 w-4" /></Button>
            </div>
            <div class="text-sm font-medium">
              {{ currentYear }}年 {{ currentMonth + 1 }}月
            </div>
            <div class="flex items-center gap-1">
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="nextMonth"><ChevronRight class="h-4 w-4" /></Button>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="nextYear"><ChevronsRight class="h-4 w-4" /></Button>
            </div>
          </div>
          
          <!-- Grid -->
          <div class="grid grid-cols-7 gap-1 text-center text-xs mb-1">
            <div v-for="d in days" :key="d" class="text-muted-foreground w-8 h-8 leading-8">{{ d }}</div>
          </div>
          <div class="grid grid-cols-7 gap-1 text-center text-sm">
            <div v-for="item in calendarDays" :key="`${item.year}-${item.month}-${item.day}`"
              class="w-8 h-8 leading-8 rounded-md cursor-pointer"
              :class="[
                !item.currentMonth ? 'text-muted-foreground/50' : '',
                isDisabled(item) ? 'opacity-30 cursor-not-allowed' : '',
                isSelected(item) ? 'bg-primary text-primary-foreground hover:bg-primary hover:text-primary-foreground' : (isDisabled(item) ? '' : 'hover:bg-accent hover:text-accent-foreground'),
                isToday(item) && !isSelected(item) ? 'bg-accent text-accent-foreground' : ''
              ]"
              @click="handleDateSelect(item)">
              {{ item.day }}
            </div>
          </div>
        </div>
        
        <!-- Time View -->
        <div v-if="type === 'time' || type === 'datetime'" class="border-t p-3 bg-muted/20">
          <div class="flex items-center justify-center gap-2">
            <div class="flex flex-col items-center">
              <span class="text-xs text-muted-foreground mb-1">时</span>
              <div class="h-32 w-14 overflow-y-auto border rounded-md bg-background no-scrollbar">
                <div v-for="h in timeOptions.hours" :key="h"
                  class="text-center text-sm py-1 cursor-pointer hover:bg-accent"
                  :class="time.h === h ? 'bg-primary/20 text-primary font-medium' : ''"
                  @click="updateTime('h', h)">
                  {{ pad(h) }}
                </div>
              </div>
            </div>
            <span class="font-bold mt-4">:</span>
            <div class="flex flex-col items-center">
              <span class="text-xs text-muted-foreground mb-1">分</span>
              <div class="h-32 w-14 overflow-y-auto border rounded-md bg-background no-scrollbar">
                <div v-for="m in timeOptions.minutes" :key="m"
                  class="text-center text-sm py-1 cursor-pointer hover:bg-accent"
                  :class="time.m === m ? 'bg-primary/20 text-primary font-medium' : ''"
                  @click="updateTime('m', m)">
                  {{ pad(m) }}
                </div>
              </div>
            </div>
            <span class="font-bold mt-4">:</span>
            <div class="flex flex-col items-center">
              <span class="text-xs text-muted-foreground mb-1">秒</span>
              <div class="h-32 w-14 overflow-y-auto border rounded-md bg-background no-scrollbar">
                <div v-for="s in timeOptions.seconds" :key="s"
                  class="text-center text-sm py-1 cursor-pointer hover:bg-accent"
                  :class="time.s === s ? 'bg-primary/20 text-primary font-medium' : ''"
                  @click="updateTime('s', s)">
                  {{ pad(s) }}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Footer -->
        <div class="border-t p-2 flex justify-between items-center bg-muted/20">
           <Button variant="ghost" size="sm" @click="handleNow">此刻</Button>
           <Button size="sm" @click="handleConfirm">确定</Button>
        </div>
      </PopoverContent>
  </Popover>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
