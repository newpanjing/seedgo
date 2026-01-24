<script setup lang="ts">
import {cn} from '@/lib/utils';
import {TableColumn} from '@/types/column';
import {ArrowDown, ArrowUp, ArrowUpDown, ChevronDown, ChevronRight} from 'lucide-vue-next';
import type {HTMLAttributes} from 'vue';
import {computed, defineComponent, h, type PropType, ref, watch} from 'vue';
import TableCellFormat from './TableCellFormat.vue';

const CellRenderer = defineComponent({
    name: 'CellRenderer',
    props: {
        row: {
            type: Object,
            required: true
        },
        column: {
            type: Object as PropType<TableColumn>,
            required: true
        },
        value: {
            type: [String, Number, Boolean, Object],
            default: null
        }
    },
    setup(props) {
        return () => {
            const { row, column, value } = props;
            
            // Check for date fields
            if (['createdAt', 'updatedAt', 'deletedAt'].includes(column.field)) {
                return h(TableCellFormat, { value: value as any });
            }

            if (column.formatter) {
                const result = column.formatter(value, row);
                // Simple check for HTML string to support "pure html"
                if (typeof result === 'string' && /<\/?[a-z][\s\S]*>/i.test(result)) {
                    return h('span', { innerHTML: result })
                }
                return result;
            }
            return value;
        };
    }
});

const props = withDefaults(defineProps<{
    class?: HTMLAttributes['class'],
    columns: TableColumn[],
    items: any[],
    sortBy?: string,
    sortOrder?: 'asc' | 'desc' | null,
    checkable?: boolean,
    defaultExpandLevel?: number,
    shouldExpandNode?: (row: any, level: number) => boolean,
    //是否启用操作
    operation?: boolean
}>(), {
    operation: true
})

const visibleColumns = computed(() => {
    return props.columns.filter(col => {
        if (col.field === 'operation') {
            return props.operation
        }
        return true
    })
})

const emit = defineEmits(['sort', 'selection-change'])

// Selection state
const selectedIds = ref<Set<number | string>>(new Set())

// Tree state
const expandedRowIds = ref<Set<number | string>>(new Set())

// Flatten items for tree display
const visibleItems = computed(() => {
    const result: any[] = []

    const flatten = (items: any[], level = 0) => {
        items.forEach(item => {
            result.push({ ...item, level, hasChildren: item.children && item.children.length > 0 })
            if (item.children && item.children.length > 0 && expandedRowIds.value.has(item.id)) {
                flatten(item.children, level + 1)
            }
        })
    }

    flatten(props.items)
    return result
})

const toggleRow = (id: number | string) => {
    if (expandedRowIds.value.has(id)) {
        expandedRowIds.value.delete(id)
    } else {
        expandedRowIds.value.add(id)
    }
}

// Computed properties for selection
const isAllSelected = computed(() => {
    return props.items.length > 0 && selectedIds.value.size === props.items.length
})

const isIndeterminate = computed(() => {
    return selectedIds.value.size > 0 && selectedIds.value.size < props.items.length
})

// Toggle all selection
const toggleSelectAll = (e: Event) => {
    const checked = (e.target as HTMLInputElement).checked
    if (checked) {
        props.items.forEach(item => selectedIds.value.add(item.id))
    } else {
        selectedIds.value.clear()
    }
    emit('selection-change', Array.from(selectedIds.value))
}

// Toggle single selection
const toggleSelect = (id: number | string) => {
    if (selectedIds.value.has(id)) {
        selectedIds.value.delete(id)
    } else {
        selectedIds.value.add(id)
    }
    emit('selection-change', Array.from(selectedIds.value))
}

// Watch items change to clean up selection
watch(() => props.items, (newItems) => {
    selectedIds.value.clear()
    expandedRowIds.value.clear()

    // Handle default expansion
    if (props.defaultExpandLevel && props.defaultExpandLevel > 0) {
        const expandRecursive = (items: any[], level: number) => {
            items.forEach(item => {
                if (level < props.defaultExpandLevel!) {
                    if (item.children && item.children.length > 0) {
                        // Check custom predicate if provided
                        if (props.shouldExpandNode && !props.shouldExpandNode(item, level)) {
                            return
                        }
                        expandedRowIds.value.add(item.id)
                        expandRecursive(item.children, level + 1)
                    }
                }
            })
        }
        expandRecursive(newItems, 0)
    }

    emit('selection-change', [])
}, { immediate: true })

const sortState = ref<{ field: string, order: 'asc' | 'desc' | null }>({
    field: props.sortBy || '',
    order: props.sortOrder || null
})

watch(() => [props.sortBy, props.sortOrder], ([newSortBy, newSortOrder]) => {
    sortState.value.field = newSortBy || ''
    sortState.value.order = (newSortOrder as 'asc' | 'desc' | null) ?? null
})

function handleSort(column: TableColumn) {
    if (!column.sortable) return

    if (sortState.value.field === column.field) {
        if (sortState.value.order === 'asc') {
            sortState.value.order = 'desc'
        } else if (sortState.value.order === 'desc') {
            sortState.value.order = null
            sortState.value.field = ''
        } else {
            sortState.value.order = 'asc'
        }
    } else {
        sortState.value.field = column.field
        sortState.value.order = 'asc'
    }

    emit('sort', { field: sortState.value.field, order: sortState.value.order })
}

const align = (column: TableColumn) => column.align === 'center' ? 'text-center' : column.align === 'right' ? 'text-right' : 'text-left'
const alignFlex = (column: TableColumn) => column.align === 'center' ? 'justify-center' : column.align === 'right' ? 'justify-end' : 'justify-start'

</script>
<template>
    <div class="relative w-full overflow-auto">
        <table :class="cn('w-full caption-bottom text-sm', props.class)">
            <thead class="[&_tr]:border-b">
                <tr class="border-b bg-muted">
                    <th v-if="checkable" class="h-12 px-4 align-middle font-medium text-muted-foreground w-[50px]">
                        <input type="checkbox" class="checkbox checkbox-sm" :checked="isAllSelected"
                            :indeterminate="isIndeterminate" @change="toggleSelectAll" />
                    </th>
                    <th v-for="column in visibleColumns" :key="column.field" @click="handleSort(column)"
                    :class="cn('h-12 px-4 align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0', align(column), column.sortable && 'cursor-pointer select-none hover:text-foreground')"
                    :style="{ width: column.width, minWidth: column.minWidth }">
                    <div :class="cn('flex items-center gap-2', alignFlex(column))">
                        {{ column.label }}
                        <template v-if="column.sortable">
                            <ArrowUp v-if="sortState.field === column.field && sortState.order === 'asc'"
                                class="h-4 w-4" />
                            <ArrowDown v-else-if="sortState.field === column.field && sortState.order === 'desc'"
                                class="h-4 w-4" />
                            <ArrowUpDown v-else class="h-4 w-4 opacity-50" />
                        </template>
                    </div>
                </th>
            </tr>
        </thead>
        <tbody v-if="visibleItems.length > 0" class="[&_tr:last-child]:border-0">
            <tr v-for="(item, index) in visibleItems" :key="item.id"
                :class="cn('group border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted', index % 2 === 1 ? 'bg-muted/30' : '')">
                <td v-if="checkable" class="group-hover:bg-muted/50 px-4 py-2 align-middle">
                    <input type="checkbox" class="checkbox checkbox-sm" :checked="selectedIds.has(item.id)"
                        @change="toggleSelect(item.id)" />
                </td>
                <td v-for="(column, colIndex) in visibleColumns" :key="column.field"
                    :class="cn('group-hover:bg-muted/50 px-4 py-2 align-middle [&:has([role=checkbox])]:pr-0', align(column))">
                        
                        <!-- First column (not selection) gets indentation and toggle -->
                        <div v-if="colIndex === 0" class="flex items-center">
                            <div :style="{ width: `${item.level * 20}px` }" class="shrink-0"></div>
                            <button v-if="item.hasChildren" @click.stop="toggleRow(item.id)"
                                class="mr-1 p-0.5 hover:bg-muted rounded cursor-pointer">
                                <ChevronDown v-if="expandedRowIds.has(item.id)" class="w-4 h-4" />
                                <ChevronRight v-else class="w-4 h-4" />
                            </button>
                            <span v-else class="w-5 mr-1 inline-block"></span> <!-- Placeholder for alignment -->
                            
                            <template v-if="column.formatter">
                                <CellRenderer :value="item[column.field]" :row="item" :column="column" />
                            </template>
                            <template v-else-if="column.field === 'operation'">
                                <slot name="actions" :row="item"></slot>
                            </template>
                            <template v-else>
                                <CellRenderer :value="item[column.field]" :row="item" :column="column" />
                            </template>
                        </div>

                        <!-- Other columns -->
                        <template v-else>
                            <template v-if="column.formatter">
                                <CellRenderer :value="item[column.field]" :row="item" :column="column" />
                            </template>
                            <template v-else-if="column.field === 'operation'">
                                <div class="flex items-center justify-end gap-1">
                                    <slot name="actions" :row="item"></slot>
                                </div>
                            </template>
                            <template v-else>
                                <CellRenderer :value="item[column.field]" :row="item" :column="column" />
                            </template>
                        </template>
                    </td>
                </tr>
            </tbody>
            <tbody v-else>
                <tr>
                    <td :colspan="checkable ? columns.length + 1 : columns.length" class="h-24 text-center">
                        No results.
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>