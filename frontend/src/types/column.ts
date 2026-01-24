export interface TableColumn {
    label: string;
    field: string;
    width?: string;
    minWidth?: string;
    align?: 'left' | 'center' | 'right';
    formatter?: (value: any, row: any) => any;
    sortable?: boolean;
    //默认为升序
    sortOrder?: 'asc' | 'desc';
}