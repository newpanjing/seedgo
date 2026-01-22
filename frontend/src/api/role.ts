import request from '@/utils/request'

export interface Role {
  id: string
  tenantId: string
  name: string
  code?: string
  description?: string
  isSystem: number
  createdAt: string
  updatedAt: string
}

export const getRoles = (params?: any) => {
  return request({
    url: '/system/roles',
    method: 'get',
    params
  })
}

export const getRole = (id: string) => {
  return request({
    url: `/system/roles/${id}`,
    method: 'get'
  })
}

export const createRole = (data: any) => {
  return request({
    url: '/system/roles',
    method: 'post',
    data
  })
}

export const updateRole = (id: string, data: any) => {
  return request({
    url: `/system/roles/${id}`,
    method: 'put',
    data
  })
}

export const deleteRole = (id: string) => {
  return request({
    url: `/system/roles/${id}`,
    method: 'delete'
  })
}
