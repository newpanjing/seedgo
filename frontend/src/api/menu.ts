import request from '@/utils/request'

export interface Menu {
  id: number | string
  parentId?: number | string
  name: string
  type: number // 1: Menu, 2: Button
  path?: string
  icon?: string
  permissionCode?: string
  permissionUrls?: string
  sort: number
  visible: boolean
  children?: Menu[]
  createdAt: string
  updatedAt: string
}

export interface GetMenusParams {
  page?: number
  pageSize?: number
  keyword?: string
}

export interface MenuListResult {
  items: Menu[]
  total: number
}

export function getMenus(params: GetMenusParams) {
  return request<MenuListResult>({
    url: '/system/permissions',
    method: 'get',
    params,
  })
}

export function getMenu(id: number | string) {
  return request<Menu>({
    url: `/system/permissions/${id}`,
    method: 'get',
  })
}

export function createMenu(data: Partial<Menu>) {
  return request<Menu>({
    url: '/system/permissions',
    method: 'post',
    data,
  })
}

export function updateMenu(id: number | string, data: Partial<Menu>) {
  return request<Menu>({
    url: `/system/permissions/${id}`,
    method: 'put',
    data,
  })
}

export function deleteMenu(id: number | string) {
  return request({
    url: `/system/permissions/${id}`,
    method: 'delete',
  })
}

export function batchDeleteMenus(ids: (number | string)[]) {
  return request({
    url: '/system/permissions/batch-delete',
    method: 'post',
    data: { ids }
  })
}
