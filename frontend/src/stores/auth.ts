import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '../types/user'
import { api } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const balance = ref<number>(0)

  const isAuthenticated = computed(() => !!token.value)
  const customerNumber = computed(() => user.value?.customerNumber || '')

  const setAuth = (newToken: string, userData: User) => {
    token.value = newToken
    user.value = userData
    // 初始化余额（实际应用中应该从API获取）
    balance.value = userData.balance || 0
    localStorage.setItem('token', newToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  const clearAuth = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    delete api.defaults.headers.common['Authorization']
  }

  const login = async (username: string, password: string) => {
    try {
      const response = await api.post('/login', { username, password })
      const { token: authToken, user: userData } = response.data
      setAuth(authToken, userData)
      return { success: true }
    } catch (error: unknown) {
      console.error('Login failed:', error)
      return { success: false, error: '登录失败，请检查用户名和密码' }
    }
  }

  const register = async (userData: {
    username: string
    password: string
    email: string
    phone: string
  }) => {
    try {
      const response = await api.post('/register', userData)
      const { token: authToken, user: newUser } = response.data
      setAuth(authToken, newUser)
      return { success: true }
    } catch (error: unknown) {
      console.error('Registration failed:', error)
      return { success: false, error: '注册失败，请稍后重试' }
    }
  }

  const logout = () => {
    clearAuth()
  }

  const checkAuth = async () => {
    if (!token.value) return false

    try {
      const response = await api.get('/profile')
      user.value = response.data
      // 更新余额
      balance.value = response.data.balance || 0
      return true
    } catch (error: unknown) {
      clearAuth()
      return false
    }
  }

  const updateBalance = (newBalance: number) => {
    balance.value = newBalance
  }

  // 初始化时设置Authorization头
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    user,
    token,
    balance,
    customerNumber,
    isAuthenticated,
    login,
    register,
    logout,
    checkAuth,
    setAuth,
    clearAuth,
    updateBalance,
  }
})
