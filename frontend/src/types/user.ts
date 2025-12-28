export interface User {
  id: number
  username: string
  email: string
  phone: string
  role: 'user' | 'admin'
  realName?: string
  address?: string
  customerNumber?: string
  balance?: number
  created_at: string
  updated_at: string
}

export interface CustomerNumber {
  id: number
  user_id: number
  customer_number: string
  address: string
  meter_number: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  email: string
  phone: string
}

export interface AuthResponse {
  token: string
  user: User
}
