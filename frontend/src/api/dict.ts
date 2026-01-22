import request from '@/utils/request'

export interface DictItem {
  id?: number | string
  label: string
  value: string
  sort: number
  status: number
}

export interface Dict {
  id: number | string
  code: string
  name: string
  description?: string
  items?: DictItem[]
  createdAt?: string
  updatedAt?: string
}

export function getDictByCode(code: string) {
  return request({
    url: `/system/dicts/code/${code}`,
    method: 'get'
  })
}

export function getDicts(params?: any) {
  return request({
    url: '/system/dicts',
    method: 'get',
    params
  })
}

export function getDict(id: number | string) {
  return request({
    url: `/system/dicts/${id}`,
    method: 'get'
  })
}

export function createDict(data: any) {
  return request({
    url: '/system/dicts',
    method: 'post',
    data
  })
}

export function updateDict(id: number | string, data: any) {
  return request({
    url: `/system/dicts/${id}`,
    method: 'put',
    data
  })
}

export function deleteDict(id: number | string) {
  return request({
    url: `/system/dicts/${id}`,
    method: 'delete'
  })
}

export function batchDeleteDicts(ids: (number | string)[]) {
  return request({
    url: '/system/dicts/batch-delete',
    method: 'post',
    data: { ids }
  })
}
