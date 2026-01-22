<script setup lang="ts">
import { computed } from 'vue'
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'outline'
type BadgeColor = 'primary' | 'red' | 'yellow' | 'green' | 'gray'

interface Props {
  variant?: BadgeVariant
  color?: BadgeColor
  class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  color: 'gray',
})

const colorClasses: Record<BadgeColor, { solid: string; outline: string }> = {
  primary: {
    solid: 'bg-primary/10 text-primary border border-primary/20',
    outline: 'border border-primary/40 text-primary bg-transparent',
  },
  red: {
    solid: 'bg-red-50 text-red-700 border border-red-200 dark:bg-red-950/40 dark:text-red-300 dark:border-red-800',
    outline: 'border border-red-300 text-red-600 bg-transparent dark:border-red-700 dark:text-red-300',
  },
  yellow: {
    solid: 'bg-amber-50 text-amber-700 border border-amber-200 dark:bg-amber-950/40 dark:text-amber-200 dark:border-amber-700',
    outline: 'border border-amber-300 text-amber-600 bg-transparent dark:border-amber-700 dark:text-amber-200',
  },
  green: {
    solid: 'bg-emerald-50 text-emerald-700 border border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-200 dark:border-emerald-700',
    outline: 'border border-emerald-300 text-emerald-600 bg-transparent dark:border-emerald-700 dark:text-emerald-200',
  },
  gray: {
    solid: 'bg-zinc-100 text-zinc-800 border border-zinc-200 dark:bg-zinc-800/70 dark:text-zinc-100 dark:border-zinc-700',
    outline: 'border border-zinc-300 text-zinc-700 bg-transparent dark:border-zinc-600 dark:text-zinc-200',
  },
}

const badgeColorClass = computed(() => {
  const tone: BadgeColor =
    props.color ||
    (props.variant === 'destructive'
      ? 'red'
      : props.variant === 'secondary'
        ? 'gray'
        : 'primary')

  const style = props.variant === 'outline' ? 'outline' : 'solid'

  return colorClasses[tone][style]
})
</script>

<template>
  <div
    :class="cn(
      'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium whitespace-nowrap',
      badgeColorClass,
      props.class,
    )"
  >
    <slot />
  </div>
</template>
