import {defineStore} from 'pinia'
import {ref} from 'vue'
import type {Category, Product} from '../types/product'

export const useProductStore = defineStore('products', () => {
  const categories = ref<Category[]>([
    { id: '1', name: '棋牌室', type: 'chess' },
    { id: '2', name: '台球桌', type: 'billiards' },
    { id: '3', name: '麻将机', type: 'mahjong' },
    { id: '4', name: '饮料零食', type: 'other' },
  ])

  const products = ref<Product[]>([
    {
      id: '1',
      name: '豪华棋牌室 A01',
      categoryId: '1',
      price: 88,
      stock: 1,
      pricingType: 'hourly',
      status: 'active',
      createdAt: new Date().toISOString()
    },
    {
      id: '2',
      name: '标准台球桌 B01',
      categoryId: '2',
      price: 45,
      stock: 1,
      pricingType: 'hourly',
      status: 'active',
      createdAt: new Date().toISOString()
    },
    {
      id: '3',
      name: '可乐',
      categoryId: '4',
      price: 5,
      stock: 100,
      pricingType: 'fixed',
      status: 'active',
      createdAt: new Date().toISOString()
    }
  ])

  const getCategoryName = (id: string) => {
    const category = categories.value.find(c => c.id === id)
    return category ? category.name : 'Unknown'
  }

  const addProduct = (product: Omit<Product, 'id' | 'createdAt'>) => {
    const newProduct: Product = {
      ...product,
      id: Math.random().toString(36).substr(2, 9),
      createdAt: new Date().toISOString()
    }
    products.value.push(newProduct)
  }

  const updateProduct = (id: string, updates: Partial<Omit<Product, 'id'>>) => {
    const index = products.value.findIndex(p => p.id === id)
    if (index !== -1) {
      const current = products.value[index]!
      products.value[index] = { ...current, ...updates, id: current.id } as Product
    }
  }

  const deleteProduct = (id: string) => {
    const index = products.value.findIndex(p => p.id === id)
    if (index !== -1) {
      products.value.splice(index, 1)
    }
  }

  return {
    categories,
    products,
    getCategoryName,
    addProduct,
    updateProduct,
    deleteProduct
  }
})
