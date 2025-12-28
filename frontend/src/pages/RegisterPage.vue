<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-header">
        <h1>电费管理系统</h1>
        <p>创建新账户</p>
      </div>

      <div class="register-card">
        <h2>用户注册</h2>

        <form @submit.prevent="handleRegister" class="register-form">
          <div class="form-row">
            <div class="form-group">
              <label for="username">用户名 *</label>
              <input
                id="username"
                v-model="registerForm.username"
                type="text"
                placeholder="请输入用户名"
                required
              />
            </div>

            <div class="form-group">
              <label for="email">邮箱 *</label>
              <input
                id="email"
                v-model="registerForm.email"
                type="email"
                placeholder="请输入邮箱"
                required
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="password">密码 *</label>
              <input
                id="password"
                v-model="registerForm.password"
                type="password"
                placeholder="请输入密码"
                required
              />
            </div>

            <div class="form-group">
              <label for="confirmPassword">确认密码 *</label>
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                required
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="phone">手机号 *</label>
              <input
                id="phone"
                v-model="registerForm.phone"
                type="tel"
                placeholder="请输入手机号"
                required
              />
            </div>

            <div class="form-group">
              <label for="realName">真实姓名</label>
              <input
                id="realName"
                v-model="registerForm.realName"
                type="text"
                placeholder="请输入真实姓名"
              />
            </div>
          </div>

          <div class="form-group">
            <label for="address">地址</label>
            <input
              id="address"
              v-model="registerForm.address"
              type="text"
              placeholder="请输入详细地址"
            />
          </div>

          <div class="form-group">
            <label class="terms-checkbox">
              <input type="checkbox" v-model="agreeTerms" required />
              <span>我已阅读并同意 <a href="#">服务条款</a> 和 <a href="#">隐私政策</a></span>
            </label>
          </div>

          <button type="submit" class="register-btn" :disabled="loading">
            {{ loading ? '注册中...' : '注册' }}
          </button>

          <div class="login-link">
            已有账号？
            <router-link to="/login">立即登录</router-link>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const registerForm = ref({
  username: '',
  password: '',
  email: '',
  phone: '',
  realName: '',
  address: ''
})

const confirmPassword = ref('')
const agreeTerms = ref(false)
const loading = ref(false)

const handleRegister = async () => {
  // 验证表单
  if (!registerForm.value.username || !registerForm.value.password ||
      !registerForm.value.email || !registerForm.value.phone) {
    alert('请填写所有必填字段')
    return
  }

  if (registerForm.value.password !== confirmPassword.value) {
    alert('两次输入的密码不一致')
    return
  }

  if (registerForm.value.password.length < 6) {
    alert('密码长度至少为6位')
    return
  }

  if (!agreeTerms.value) {
    alert('请同意服务条款和隐私政策')
    return
  }

  loading.value = true
  try {
    const result = await authStore.register({
      username: registerForm.value.username,
      password: registerForm.value.password,
      email: registerForm.value.email,
      phone: registerForm.value.phone
    })

    if (result.success) {
      alert('注册成功！')
      router.push('/dashboard')
    } else {
      alert(result.error || '注册失败，请稍后重试')
    }
  } catch (error) {
    console.error('注册失败:', error)
    alert('注册失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-container {
  width: 100%;
  max-width: 800px;
}

.register-header {
  text-align: center;
  color: white;
  margin-bottom: 30px;
}

.register-header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.register-header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.register-card {
  background: white;
  border-radius: 10px;
  padding: 40px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.register-card h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.terms-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #555;
}

.terms-checkbox a {
  color: #667eea;
  text-decoration: none;
}

.terms-checkbox a:hover {
  text-decoration: underline;
}

.register-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.3s;
  margin-top: 10px;
}

.register-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.register-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-link {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.login-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.login-link a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
