<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{ class?: HTMLAttributes['class'] }>()
</script>

<template>
  <label :class="cn('label text-slate-500 dark:[&_span]:text-white', props.class)">
    <span class="label-text">
      <slot />
    </span>
  </label>
</template>
