<script lang="ts" setup>
const props = defineProps({
    code: {
        type: String,
        default: ''
    }
})
const perms = usePerms()
const hasPerm = computed(() => {
    return perms[props.code] || false
})
</script>
<template>
    <slot v-if="hasPerm"></slot>
</template>