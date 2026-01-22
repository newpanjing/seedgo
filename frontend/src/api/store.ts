import request from '@/utils/request'

export interface Store {
  id: string
  tenantId: string
  name: string
  code?: string
  phone?: string
  address?: string
  status: number
  openingHours?: string
  createdAt: string
  updatedAt: string
}

export interface GetStoresParams {
  page?: number
  pageSize?: number
  name?: string
  tenantId?: string
}

export interface StoreListResult {
  items: Store[]
  total: number
}

export function getStores(params: GetStoresParams) {
  return request<StoreListResult>({
    url: '/business/stores',
    method: 'get',
    params,
  })
}

export function getStore(id: string) {
  return request<Store>({
    url: `/business/stores/${id}`,
    method: 'get',
  })
}

export function createStore(data: Partial<Store>) {
  return request<Store>({
    url: '/business/stores',
    method: 'post',
    data,
  })
}

export function updateStore(id: string, data: Partial<Store>) {
  return request<Store>({
    url: `/business/stores/${id}`,
    method: 'put',
    data,
  })
}

export function deleteStore(id: string) {
  return request({
    url: `/business/stores/${id}`,
    method: 'delete',
  })
}
