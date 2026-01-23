import request from '@/utils/request'

export interface Permission {
  id: string | number
  parentId?: string | number
  name: string
  permissionCode: string // This maps to permissionCode from backend
  type: number // 1: Menu, 2: Button
  path?: string
  icon?: string
  visible?: number | boolean
  sort?: number
  children?: Permission[]
}

export function getPermissionTree(params?: any) {
  return request<Permission[]>({
    url: '/common/user/permissions',
    method: 'get',
    params
  })
}
