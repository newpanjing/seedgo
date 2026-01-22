import request from '@/utils/request'

export interface Tenant {
  id: string
  name: string
  contactName?: string
  contactPhone?: string
  contactEmail?: string
  status: number
  createdAt: string
  users?: any[] // Admin users
  _count?: {
    stores: number
  }
}

export const getTenants = (params?: any) => {
  return request({
    url: '/tenant/tenants',
    method: 'get',
    params
  })
}

export const createTenant = (data: any) => {
  return request({
    url: '/tenant/tenants',
    method: 'post',
    data
  })
}

export const updateTenant = (id: string, data: any) => {
  return request({
    url: `/tenant/tenants/${id}`,
    method: 'put',
    data
  })
}

export const deleteTenant = (id: string) => {
  return request({
    url: `/tenant/tenants/${id}`,
    method: 'delete'
  })
}
