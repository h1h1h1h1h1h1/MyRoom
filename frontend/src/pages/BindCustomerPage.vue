<template>
  <div class="bind-customer-page">
    <h1>绑定客户号</h1>
    <div class="bind-form">
      <form @submit.prevent="bindCustomer">
        <div class="form-group">
          <label for="customerNumber">客户号</label>
          <input
            type="text"
            id="customerNumber"
            v-model="customerNumber"
            placeholder="请输入您的客户号"
            required
          />
        </div>
        <div class="form-group">
          <label for="customerName">客户姓名</label>
          <input
            type="text"
            id="customerName"
            v-model="customerName"
            placeholder="请输入客户姓名"
            required
          />
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? '绑定中...' : '绑定客户号' }}
        </button>
      </form>
    </div>

    <div v-if="boundCustomers.length > 0" class="bound-customers">
      <h2>已绑定的客户号</h2>
      <ul>
        <li v-for="customer in boundCustomers" :key="customer.id">
          {{ customer.customer_number }} - {{ customer.customer_name }}
          <button @click="unbindCustomer(customer.id)" class="unbind-btn">解绑</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/services/api'

interface BoundCustomer {
  id: number
  customer_number: string
  customer_name: string
}

const authStore = useAuthStore()
const customerNumber = ref('')
const customerName = ref('')
const loading = ref(false)
const boundCustomers = ref<BoundCustomer[]>([])

const fetchBoundCustomers = async () => {
  try {
    const response = await api.get('/customer-numbers')
    boundCustomers.value = response.data
  } catch (error) {
    console.error('获取已绑定客户号失败:', error)
  }
}

const bindCustomer = async () => {
  if (!customerNumber.value.trim() || !customerName.value.trim()) {
    alert('请填写完整的客户信息')
    return
  }

  loading.value = true
  try {
    await api.post('/customer-numbers', {
      customer_number: customerNumber.value,
      customer_name: customerName.value
    })

    alert('绑定成功！')
    customerNumber.value = ''
    customerName.value = ''
    await fetchBoundCustomers()
  } catch (error: unknown) {
    console.error('绑定失败:', error)
    const err = error as { response?: { data?: { message?: string } } }
    alert(err.response?.data?.message || '绑定失败，请重试')
  } finally {
    loading.value = false
  }
}

const unbindCustomer = async (id: number) => {
  if (!confirm('确定要解绑此客户号吗？')) {
    return
  }

  try {
    await api.delete(`/customer-numbers/${id}`)
    alert('解绑成功！')
    await fetchBoundCustomers()
  } catch (error: unknown) {
    console.error('解绑失败:', error)
    const err = error as { response?: { data?: { message?: string } } }
    alert(err.response?.data?.message || '解绑失败，请重试')
  }
}

onMounted(() => {
  fetchBoundCustomers()
})
</script>

<style scoped>
.bind-customer-page {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.bind-form {
  background: #f8f9fa;
  padding: 30px;
  border-radius: 8px;
  margin-bottom: 30px;
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

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

button {
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

button:hover:not(:disabled) {
  background: #0056b3;
}

button:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.bound-customers {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
}

.bound-customers h2 {
  margin-bottom: 15px;
  color: #333;
}

.bound-customers ul {
  list-style: none;
  padding: 0;
}

.bound-customers li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
}

.bound-customers li:last-child {
  border-bottom: none;
}

.unbind-btn {
  width: auto;
  padding: 5px 15px;
  background: #dc3545;
  font-size: 14px;
}

.unbind-btn:hover {
  background: #c82333;
}
</style>
