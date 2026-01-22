import {defineStore} from 'pinia'
import {ref} from 'vue'
import type {User} from '../types/user'

export const useUserStore = defineStore('users', () => {
  const users = ref<User[]>([
    {
      id: '1',
      name: '张三',
      phone: '13800138000',
      level: 'gold',
      discount: 0.8,
      balance: 5000,
      points: 1200,
      totalConsumed: 10000,
      tags: ['常客', '大方'],
      createdAt: '2023-01-01T10:00:00Z',
      lastVisit: '2023-10-25T14:30:00Z'
    },
    {
      id: '2',
      name: '李四',
      phone: '13900139000',
      level: 'silver',
      discount: 0.9,
      balance: 200,
      points: 300,
      totalConsumed: 1000,
      tags: ['偶尔'],
      createdAt: '2023-05-15T11:00:00Z',
      lastVisit: '2023-10-20T18:00:00Z'
    },
    {
      id: '3',
      name: '王五',
      phone: '13700137000',
      level: 'normal',
      discount: 1.0,
      balance: 0,
      points: 50,
      totalConsumed: 200,
      tags: [],
      createdAt: '2023-09-01T09:00:00Z',
      lastVisit: '2023-09-01T09:00:00Z'
    }
  ])

  const addUser = (user: Omit<User, 'id' | 'createdAt' | 'lastVisit'>) => {
    const newUser: User = {
      ...user,
      id: Math.random().toString(36).substr(2, 9),
      createdAt: new Date().toISOString(),
      lastVisit: new Date().toISOString()
    }
    users.value.push(newUser)
  }

  const updateUser = (id: string, updates: Partial<Omit<User, 'id'>>) => {
    const index = users.value.findIndex(u => u.id === id)
    if (index !== -1) {
      const current = users.value[index]!
      users.value[index] = { ...current, ...updates, id: current.id } as User
    }
  }

  const deleteUser = (id: string) => {
    const index = users.value.findIndex(u => u.id === id)
    if (index !== -1) {
      users.value.splice(index, 1)
    }
  }

  return {
    users,
    addUser,
    updateUser,
    deleteUser
  }
})
