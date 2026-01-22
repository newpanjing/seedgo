<script lang="ts" setup>
defineProps({
    title: {
        type: String,
        default: '数据加载失败'
    },
    message: {
        type: String,
        default: '抱歉，由于网络波动或服务器响应异常，我们无法加载最新的项目数据。请检查您的网络连接或稍后再试。'
    }
})
const emit = defineEmits(['retry'])

</script>
<template>
    <div class="mt-10 mb-10 flex flex-col items-center justify-center">
        <div class="max-w-md w-full text-center">
            <div class="mb-6 relative inline-block">
                <div class="absolute inset-0 bg-rose-100 rounded-full blur-2xl opacity-50 animate-pulse">
                </div>
                <div
                    class="relative w-20 h-20 bg-rose-50 rounded-2xl flex items-center justify-center mx-auto border border-rose-100 shadow-inner">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                        class="lucide lucide-circle-alert w-10 h-10 text-rose-500" aria-hidden="true">
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" x2="12" y1="8" y2="12"></line>
                        <line x1="12" x2="12.01" y1="16" y2="16"></line>
                    </svg>
                </div>
            </div>
            <h3 class="text-xl font-bold text-slate-800 mb-2 dark:text-white" v-html="title"></h3>
            <p class="text-slate-500 mb-8 leading-relaxed" v-html="message">
            </p>
            <div class="flex flex-col sm:flex-row gap-3 justify-center">
                <button @click="emit('retry')"
                    class="dark:bg-white dark:text-black inline-flex items-center justify-center gap-2 px-6 py-3 bg-slate-900 hover:bg-slate-800 text-white rounded-xl font-semibold transition-all hover:shadow-lg active:scale-95">
                    <RefreshCwIcon class="w-4 h-4" />
                    重新尝试加载
                </button>
            </div>
            <div class="mt-8 pt-6 border-t border-slate-100 dark:border-slate-700"><code
                    class="text-[10px] text-slate-400 uppercase tracking-widest font-mono dark:text-slate-500">Error Code: 503 SERVICE_UNAVAILABLE</code>
            </div>
        </div>
    </div>
</template>