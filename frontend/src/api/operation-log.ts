import request from '@/utils/request'

export interface OperationLog {
  id: number
  userId: number
  username: string
  method: string
  path: string
  query: string
  body: string
  ip: string
  userAgent: string
  status: number
  latency: number
  errorMessage: string
  operationTime: string
  createdAt: string
  updatedAt: string
}

export interface GetOperationLogsParams {
  page?: number
  pageSize?: number
  keyword?: string
  status?: number
  method?: string
  startDate?: string
  endDate?: string
}

export interface OperationLogListResult {
  items: OperationLog[]
  total: number
}

export function getOperationLogs(params: GetOperationLogsParams) {
  return request<OperationLogListResult>({
    url: '/system/operation-logs',
    method: 'get',
    params,
  })
}

export function getOperationLog(id: number) {
  return request<OperationLog>({
    url: `/system/operation-logs/${id}`,
    method: 'get',
  })
}
