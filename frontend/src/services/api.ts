import axios from 'axios'

export const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 清除本地存储的token
      localStorage.removeItem('token')
      // 重定向到登录页面
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const getApplications = (userId: number) => api.get(`/user/applications?user_id=${userId}`)
export const submitApplication = (data: any) => api.post('/user/application', data)
export const getAnnouncements = () => api.get('/info/announcements')
export const getNotifications = (userId: number) => api.get(`/user/notifications?user_id=${userId}`)
export const getCustomers = (userId: number) => api.get(`/user/customers?user_id=${userId}`)
export const bindCustomer = (data: any) => api.post('/user/bind', data)
export const getUsage = (customerId: number) => api.get(`/electricity/usage?customer_id=${customerId}`)
export const getPayments = (customerId: number) => api.get(`/electricity/payments?customer_id=${customerId}`)
export const pay = (data: any) => api.post('/electricity/pay', data)
export const checkSystem = () => api.post('/system/check')

export default {
  getApplications,
  submitApplication,
  getAnnouncements,
  getNotifications,
  getCustomers,
  bindCustomer,
  getUsage,
  getPayments,
  pay,
  checkSystem
}

// Force update
