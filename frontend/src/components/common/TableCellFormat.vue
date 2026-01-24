<script setup lang="ts">
import { computed } from 'vue'
import { formatDate } from '@/utils/date'

const props = defineProps({
    value: {
        //时间戳
        type: [Number, String],
        default: () => null
    },
    format: {
        type: String,
        default: 'YYYY-MM-DD HH:mm:ss'
    }
})

const diaplay = computed(() => {
    return formatDate(props.value, props.format)
})
</script>
<template>
    {{ diaplay }}
</template>