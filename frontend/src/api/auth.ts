import request from '@/utils/request'

export interface LoginParams {
  username: string // Can be phone or username
  password: string
}

export interface LoginResult {
  token: string
  user: {
    id: number
    username: string
    realName?: string
    name?: string
    avatar?: string
    role?: string
    tenantId?: string
    isSuper?: boolean
  }
}

export function login(data: LoginParams) {
  return request<any, LoginResult>({
    url: '/auth/login',
    method: 'post',
    data,
  })
}

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post',
  })
}

export function getUserInfo() {
  return request({
    url: '/common/user/profile',
    method: 'get',
  })
}

export function updateUserInfo(data: { name?: string; phone?: string; email?: string }) {
  return request({
    url: '/common/user/profile',
    method: 'post',
    data,
  })
}

export function changePassword(data: any) {
  return request({
    url: '/common/user/change-password',
    method: 'post',
    data,
  })
}
