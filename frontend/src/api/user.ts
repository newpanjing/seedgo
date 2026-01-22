import request from '@/utils/request'

export interface User {
  id: number
  username: string
  realName?: string
  phone?: string
  email?: string
  status: number
  isMain?: number
  tenantId?: number
  avatar?: string
  role?: string
  roles?: { roleId: string }[]
  createdAt: string
  updatedAt: string
}

export interface GetUsersParams {
  page?: number
  pageSize?: number
  keyword?: string
  tenantId?: number
}

export interface UserListResult {
  items: User[]
  total: number
}

export function getUsers(params: GetUsersParams) {
  return request<UserListResult>({
    url: '/system/users',
    method: 'get',
    params,
  })
}

export function getUser(id: number) {
  return request<User>({
    url: `/system/users/${id}`,
    method: 'get',
  })
}

export function createUser(data: Partial<User> & { password?: string; roleId?: string; roleIds?: number[]; tenantId?: number }) {
  return request<User>({
    url: '/system/users',
    method: 'post',
    data,
  })
}

export function updateUser(id: number, data: Partial<User> & { password?: string; roleId?: string; roleIds?: number[] }) {
  return request<User>({
    url: `/system/users/${id}`,
    method: 'put',
    data,
  })
}

export function deleteUser(id: number) {
  return request({
    url: `/system/users/${id}`,
    method: 'delete',
  })
}

export function batchDeleteUsers(ids: number[]) {
  return request({
    url: '/system/users/batch-delete',
    method: 'post',
    data: { ids }
  })
}
