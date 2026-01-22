<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div
    :class="
      cn(
        'flex flex-row justify-between items-center border-t -mx-6 px-6 pt-4 mt-4',
        props.class,
      )
    "
  >
    <slot />
  </div>
</template>
