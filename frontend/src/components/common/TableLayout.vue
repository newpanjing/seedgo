<script lang="tsx" setup>
import { showConfirm } from '@/lib/message'
import { TableColumn } from '@/types/column'
//导入图标
import { Search, Plus, RefreshCw, Loader2, Trash2 } from 'lucide-vue-next'



const props = defineProps({
    title: {
        type: String,
        default: ''
    },
    fetchData: {
        type: Function as PropType<(params: any) => Promise<{ total: number, items: any[] }>>,
        default: () => { }
    },
    //定义表格列信息TableColumn
    columns: {
        type: Array as PropType<TableColumn[]>,
        default: () => []
    },
    checkable: {
        type: Boolean,
        default: false
    },
    defaultExpandLevel: {
        type: Number,
        default: 0
    },
    noPagination: {
        type: Boolean,
        default: false
    },
    shouldExpandNode: {
        type: Function as PropType<(row: any, level: number) => boolean>,
        default: null
    },
    showView: {
        type: Boolean,
        default: false
    },
    showUpdate: {
        type: Boolean,
        default: true
    },
    showDelete: {
        type: Boolean,
        default: true
    },
})

const columns = computed(() => {
    //如果没有名为operation的列，添加一个
    if (!props.columns.some(column => column.field === 'operation')) {
        props.columns.push({
            label: '操作',
            field: 'operation',
            width: '100px',
            align: 'center'
        })
    }
    return props.columns
})

const emit = defineEmits(['save', 'delete', 'update', 'view', 'create', 'selection-change', 'batch-delete'])

const searchQuery = ref('')
const page = ref(1)
const pageSize = ref(10)
const sortBy = ref('')
const sortOrder = ref<'asc' | 'desc' | null>(null)

const params = computed(() => ({
    keyword: searchQuery.value,
    page: page.value,
    pageSize: pageSize.value,
    sortBy: sortBy.value,
    sortOrder: sortOrder.value
}))

const handleSearch = () => {
    fetchData()
}

const handleSort = ({ field, order }: { field: string, order: 'asc' | 'desc' | null }) => {
    sortBy.value = field
    sortOrder.value = order
    fetchData()
}

//openCreateModal
const openCreateModal = () => {
    emit('create')
}
const total = ref(0)
const items = ref<any[]>([])

const loading = ref(false)
//判断是否有错误的变量
const error = reactive({
    hasError: false,
    message: ''
})
const fetchData = async () => {
    loading.value = true
    error.hasError = false
    error.message = ''
    try {
        let res = await props.fetchData(params.value);
        total.value = res.total
        items.value = res.items || []
    } catch (e) {
        error.hasError = true
        error.message = e instanceof Error ? e.message : '未知错误'
    } finally {
        loading.value = false
    }

}

onMounted(() => {
    fetchData()
})


const selectedItems = ref<any[]>([])

const handleSelectionChange = (selection: any[]) => {
    selectedItems.value = selection
    emit('selection-change', selection)
}

const handleBatchDelete = async () => {
    //显示提示
    if (await showConfirm('确认删除', `确认删除选中的 ${selectedItems.value.length} 条吗？`)) {
        emit('batch-delete', selectedItems.value)
    }
}

defineExpose({
    fetchData
})

watch([page, pageSize], () => {
    fetchData()
})



// 核心逻辑：定义一个生成器，把事件预先绑定
const ViewBtn = (row: any) => (
    <BaseAction code="view">
        <Button
            variant="ghost"
            size="icon"
            class="h-8 w-8"
            title="查看"
            onClick={() => emit('view', row)} // 内部直接处理掉 click
        >
            <Eye class="w-4 h-4" />
        </Button>
    </BaseAction>
);

const UpdateBtn = (row: any) => (
    <BaseAction code="update">
        <Button
            variant="ghost"
            size="icon"
            class="h-8 w-8"
            title="编辑"
            onClick={() => emit('update', row)} // 内部直接处理掉 click
        >
            <Edit class="w-4 h-4" />
        </Button>
    </BaseAction>
);

const DeleteBtn = (row: any) => (
    <BaseAction code="delete">
        <Button
            variant="ghost"
            size="icon"
            class="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            title="删除"
            onClick={() => emit('delete', row)} // 内部直接处理掉 click
        >
            <Trash2 class="w-4 h-4" />
        </Button>
    </BaseAction>
);

const perms = usePerms()
const slots = useSlots();
const hasOperation = computed(() => {
    // 1. 检查是否有默认按钮的权限
    const hasBasePermission = perms.update || perms.delete;

    // 2. 检查父组件是否实现了 actions 插槽
    const hasCustomSlot = !!slots.actions;

    // 只要满足其一，就需要显示“操作”列
    return hasBasePermission || hasCustomSlot;
})

</script>

<template>
    <div class="space-y-2">
        <!-- Filters -->
        <div class="bg-card py-3 px-4 rounded-xl shadow-sm flex justify-between items-center">
            <div class="flex flex-col md:flex-row gap-4 items-center">
                <div class="relative flex-1 md:max-w-sm">
                    <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-4 h-4" />
                    <Input v-model="searchQuery" placeholder="搜索..." class="pl-9 h-8" @keyup.enter="handleSearch" />
                </div>
                <!-- 额外的查询条件 slot -->
                <slot name="filters"></slot>
                <div class="flex gap-2">
                    <Button @click="handleSearch" variant="outline">
                        <Search class="w-4 h-4 mr-2" />
                        查询
                    </Button>
                </div>
            </div>

        </div>


        <div class="overflow-hidden py-2 flex gap-2 items-center justify-between">
            <slot class="toolbar">
                <div class="flex-1 flex flex-wrap gap-2 items-center">
                    <div v-has:create>
                        <Button @click="openCreateModal" variant="default">
                            <Plus class="w-4 h-4 mr-2" />
                            新增{{ title }}
                        </Button>
                    </div>
                    <slot name="toolbar-extra" :selectedItems="selectedItems" :items="items">

                    </slot>
                    <div v-has:delete v-if="selectedItems.length > 0" class="flex items-center gap-2 mr-2">
                        <Button variant="outline" size="sm"
                            class="border-red-500 bg-red-50 text-red-500 hover:bg-red-100 hover:text-red-600 hover:border-red-600"
                            @click="handleBatchDelete">
                            <Trash2 class="w-4 h-4 mr-2" />
                            删除 {{ selectedItems.length }} 条
                        </Button>
                    </div>
                </div>
            </slot>
            <slot name="refresh">
                <Button variant="outline" @click="fetchData" :disabled="loading">
                    <Loader2 v-if="loading" class="w-4 h-4 mr-1 animate-spin" />
                    <RefreshCw v-else class="w-4 h-4 mr-1" />
                </Button>
            </slot>
        </div>

        <!--  List -->
        <div class="bg-card rounded-xl shadow-sm border border-border overflow-hidden">
            <!-- table -->
            <div class="relative min-h-[200px]">
                <Transition name="fade">
                    <div v-if="loading"
                        class="absolute inset-0 z-10 flex items-center justify-center bg-background/50 backdrop-blur-[1px]">
                        <Loader2 class="w-8 h-8 animate-spin text-primary" />
                    </div>
                </Transition>
                <Transition name="fade">
                    <FetchDataError v-if="error.hasError" :message="error.message" @retry="fetchData" />
                    <DataTable v-else
                        :class="{ 'transition-opacity duration-300': true, 'opacity-50 pointer-events-none': loading }"
                        :columns="columns" :items="items" :sort-by="sortBy" :sort-order="sortOrder"
                        :checkable="checkable" :operation="hasOperation" :default-expand-level="defaultExpandLevel"
                        :should-expand-node="shouldExpandNode" @sort="handleSort"
                        @selection-change="handleSelectionChange">
                        <template #actions="{ row }">
                            <div class="flex items-center justify-end gap-1">
                                <slot name="actions" :View="() => ViewBtn(row)" :Update="() => UpdateBtn(row)"
                                    :Delete="() => DeleteBtn(row)" :Row="row">
                                    <component v-if="showView" :is="ViewBtn(row)" />
                                    <component v-if="showUpdate" :is="UpdateBtn(row)" />
                                    <component v-if="showDelete" :is="DeleteBtn(row)" />
                                </slot>
                            </div>
                        </template>
                    </DataTable>
                </Transition>
            </div>
            <!-- Pagination -->
            <Pagination v-if="!noPagination" :page="page" :page-size="pageSize" :total="total" :loading="loading"
                @update:page="(val: number) => page = val" />
        </div>

    </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
