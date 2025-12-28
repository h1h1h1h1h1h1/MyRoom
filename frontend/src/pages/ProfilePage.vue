<template>
  <div class="profile-page">
    <h1>个人资料</h1>

    <div class="profile-card">
      <div class="profile-header">
        <div class="avatar">
          <span>{{ userInitials }}</span>
        </div>
        <div class="user-info">
          <h2>{{ user.name || '未设置姓名' }}</h2>
          <p>{{ user.email }}</p>
        </div>
      </div>

      <div class="profile-form">
        <form @submit.prevent="updateProfile">
          <div class="form-group">
            <label for="name">姓名</label>
            <input
              type="text"
              id="name"
              v-model="form.name"
              placeholder="请输入您的姓名"
            />
          </div>

          <div class="form-group">
            <label for="phone">手机号</label>
            <input
              type="tel"
              id="phone"
              v-model="form.phone"
              placeholder="请输入手机号"
            />
          </div>

          <div class="form-group">
            <label for="address">地址</label>
            <textarea
              id="address"
              v-model="form.address"
              placeholder="请输入地址"
              rows="3"
            ></textarea>
          </div>

          <button type="submit" :disabled="loading">
            {{ loading ? '保存中...' : '保存更改' }}
          </button>
        </form>
      </div>
    </div>

    <div class="account-actions">
      <h3>账户操作</h3>
      <div class="actions-grid">
        <button class="action-btn" @click="changePassword">
          <span class="icon">🔒</span>
          <span>修改密码</span>
        </button>
        <button class="action-btn" @click="logout">
          <span class="icon">🚪</span>
          <span>退出登录</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { api } from '@/services/api'

interface UserProfile {
  id: number
  email: string
  name: string
  phone: string
  address: string
  created_at: string
}

const authStore = useAuthStore()
const router = useRouter()

const user = ref<UserProfile>({
  id: 0,
  email: '',
  name: '',
  phone: '',
  address: '',
  created_at: ''
})

const form = ref({
  name: '',
  phone: '',
  address: ''
})

const loading = ref(false)

const userInitials = computed(() => {
  if (!user.value.name) return 'U'
  return user.value.name.charAt(0).toUpperCase()
})

const fetchProfile = async () => {
  try {
    const response = await api.get('/profile')
    user.value = response.data
    form.value = {
      name: response.data.name || '',
      phone: response.data.phone || '',
      address: response.data.address || ''
    }
  } catch (error) {
    console.error('获取个人资料失败:', error)
  }
}

const updateProfile = async () => {
  loading.value = true
  try {
    const response = await api.put('/profile', form.value)
    user.value = { ...user.value, ...response.data }
    alert('个人资料更新成功！')
  } catch (error: unknown) {
    console.error('更新个人资料失败:', error)
    const err = error as { response?: { data?: { message?: string } } }
    alert(err.response?.data?.message || '更新失败，请重试')
  } finally {
    loading.value = false
  }
}

const changePassword = () => {
  const newPassword = prompt('请输入新密码:')
  if (newPassword) {
    alert('密码修改功能正在开发中...')
  }
}

const logout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.profile-card {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 30px;
  margin-bottom: 30px;
}

.profile-header {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #dee2e6;
}

.avatar {
  width: 80px;
  height: 80px;
  background: #007bff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
}

.avatar span {
  font-size: 32px;
  font-weight: bold;
  color: white;
}

.user-info h2 {
  margin: 0 0 5px 0;
  color: #333;
}

.user-info p {
  margin: 0;
  color: #6c757d;
}

.profile-form {
  background: white;
  padding: 20px;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #555;
}

input, textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  font-family: inherit;
}

input:focus, textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

button[type="submit"] {
  width: 100%;
  padding: 12px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.3s;
}

button[type="submit"]:hover:not(:disabled) {
  background: #0056b3;
}

button[type="submit"]:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.account-actions h3 {
  margin-bottom: 20px;
  color: #333;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.action-btn {
  display: flex;
  align-items: center;
  padding: 15px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  background: #f8f9fa;
  border-color: #007bff;
  transform: translateY(-2px);
}

.action-btn .icon {
  font-size: 20px;
  margin-right: 10px;
}

.action-btn span:last-child {
  font-size: 16px;
  color: #333;
}
</style>
