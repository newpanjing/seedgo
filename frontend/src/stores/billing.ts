import {defineStore} from 'pinia'
import {ref} from 'vue'
import type {Order, Room, Session} from '../types/billing'
import {useProductStore} from './products'

export const useBillingStore = defineStore('billing', () => {
  const rooms = ref<Room[]>([
    { id: '101', name: '棋牌室 101', type: 'chess', status: 'available', hourlyRate: 88 },
    { id: '102', name: '棋牌室 102', type: 'chess', status: 'occupied', currentSessionId: 's1', hourlyRate: 88 },
    { id: '103', name: '麻将房 201', type: 'mahjong', status: 'available', hourlyRate: 68 },
    { id: '104', name: '台球桌 A1', type: 'billiards', status: 'available', hourlyRate: 45 },
    { id: '105', name: '台球桌 A2', type: 'billiards', status: 'available', hourlyRate: 45 },
    { id: '106', name: 'VIP包厢', type: 'chess', status: 'maintenance', hourlyRate: 128 },
  ])

  const activeSessions = ref<Session[]>([
    {
      id: 's1',
      roomId: '102',
      startTime: new Date(Date.now() - 1000 * 60 * 45).toISOString(), // started 45 mins ago
      status: 'active',
      totalAmount: 0,
      orders: [],
      durationMinutes: 45
    }
  ])

  const productStore = useProductStore()

  const getRoomById = (id: string) => rooms.value.find(r => r.id === id)
  
  const getSessionByRoomId = (roomId: string) => {
    return activeSessions.value.find(s => s.roomId === roomId && s.status === 'active')
  }

  const addRoom = (room: Omit<Room, 'id' | 'status'>) => {
    const newRoom: Room = {
      ...room,
      id: Math.random().toString(36).substr(2, 9),
      status: 'available'
    }
    rooms.value.push(newRoom)
  }

  const updateRoom = (id: string, updates: Partial<Omit<Room, 'id'>>) => {
    const index = rooms.value.findIndex(r => r.id === id)
    if (index !== -1) {
      const current = rooms.value[index]
      if (current) {
        rooms.value[index] = { 
          ...current, 
          ...updates, 
          id: current.id,
          // Ensure required fields are not undefined if updates has them as undefined
          name: updates.name ?? current.name,
          type: updates.type ?? current.type,
          status: updates.status ?? current.status,
          hourlyRate: updates.hourlyRate ?? current.hourlyRate
        }
      }
    }
  }

  const deleteRoom = (id: string) => {
    const index = rooms.value.findIndex(r => r.id === id)
    if (index !== -1) {
      rooms.value.splice(index, 1)
    }
  }

  const startSession = (roomId: string) => {
    const room = getRoomById(roomId)
    if (!room || room.status !== 'available') return

    const newSession: Session = {
      id: Math.random().toString(36).substr(2, 9),
      roomId,
      startTime: new Date().toISOString(),
      status: 'active',
      totalAmount: 0,
      orders: [],
      durationMinutes: 0
    }

    activeSessions.value.push(newSession)
    room.status = 'occupied'
    room.currentSessionId = newSession.id
  }

  const endSession = (roomId: string, settlement?: { totalAmount: number, discount: number, rounding: number }) => {
    const room = getRoomById(roomId)
    if (!room || !room.currentSessionId) return

    const session = activeSessions.value.find(s => s.id === room.currentSessionId)
    if (session) {
      session.endTime = new Date().toISOString()
      session.status = 'completed'
      
      if (settlement) {
        session.totalAmount = settlement.totalAmount
        session.discount = settlement.discount
        session.rounding = settlement.rounding
      } else {
        // Fallback logic if no settlement details provided
        const start = new Date(session.startTime).getTime()
        const end = new Date(session.endTime).getTime()
        const hours = (end - start) / (1000 * 60 * 60)
        const roomCharge = Math.ceil(hours * room.hourlyRate)
        const ordersTotal = session.orders.reduce((sum, order) => sum + order.total, 0)
        session.totalAmount = roomCharge + ordersTotal
      }
    }

    room.status = 'cleaning'
    room.currentSessionId = undefined
  }

  const setRoomAvailable = (roomId: string) => {
    const room = getRoomById(roomId)
    if (room) {
      room.status = 'available'
    }
  }

  const addOrderToSession = (sessionId: string, productId: string, quantity: number) => {
    const session = activeSessions.value.find(s => s.id === sessionId)
    const product = productStore.products.find(p => p.id === productId)
    
    if (session && product) {
      const order: Order = {
        id: Math.random().toString(36).substr(2, 9),
        productId,
        productName: product.name,
        quantity,
        price: product.price,
        total: product.price * quantity,
        timestamp: new Date().toISOString()
      }
      session.orders.push(order)
    }
  }

  // Timer to update duration
  setInterval(() => {
    const now = new Date().getTime()
    activeSessions.value.forEach(session => {
      if (session.status === 'active') {
        const start = new Date(session.startTime).getTime()
        session.durationMinutes = Math.floor((now - start) / (1000 * 60))
      }
    })
  }, 60000)

  return {
    rooms,
    activeSessions,
    addRoom,
    updateRoom,
    deleteRoom,
    startSession,
    endSession,
    setRoomAvailable,
    addOrderToSession,
    getSessionByRoomId
  }
})
