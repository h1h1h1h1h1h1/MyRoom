<template>
  <div class="electricity-page">
    <NavBar />

    <div class="page-container">
      <div class="page-header">
        <h1>电量电费管理</h1>
        <div class="header-actions">
          <button class="btn btn-primary" @click="refreshData">
            <span class="icon">🔄</span> 刷新数据
          </button>
          <button class="btn btn-secondary" @click="exportData">
            <span class="icon">📥</span> 导出数据
          </button>
        </div>
      </div>

      <div class="filter-section">
        <div class="filter-row">
          <div class="filter-group">
            <label for="year">年份</label>
            <select id="year" v-model="selectedYear" @change="filterData">
              <option v-for="year in years" :key="year" :value="year">{{ year }}年</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="month">月份</label>
            <select id="month" v-model="selectedMonth" @change="filterData">
              <option value="all">全部</option>
              <option v-for="month in months" :key="month.value" :value="month.value">
                {{ month.label }}
              </option>
            </select>
          </div>
          <div class="filter-group">
            <label for="status">状态</label>
            <select id="status" v-model="selectedStatus" @change="filterData">
              <option value="all">全部</option>
              <option value="paid">已支付</option>
              <option value="unpaid">未支付</option>
              <option value="overdue">已逾期</option>
            </select>
          </div>
        </div>
      </div>

      <div class="stats-summary">
        <div class="stat-item">
          <div class="stat-label">总用电量</div>
          <div class="stat-value">{{ totalUsage }} kWh</div>
        </div>
        <div class="stat-item">
          <div class="stat-label">总电费</div>
          <div class="stat-value">¥ {{ formatCurrency(totalCost) }}</div>
        </div>
        <div class="stat-item">
          <div class="stat-label">平均单价</div>
          <div class="stat-value">¥ {{ formatCurrency(averagePrice) }}/kWh</div>
        </div>
        <div class="stat-item">
          <div class="stat-label">未支付账单</div>
          <div class="stat-value">{{ unpaidCount }} 笔</div>
        </div>
      </div>

      <div class="data-table">
        <table>
          <thead>
            <tr>
              <th>月份</th>
              <th>用电量 (kWh)</th>
              <th>电费 (¥)</th>
              <th>单价 (¥/kWh)</th>
              <th>抄表日期</th>
              <th>支付状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in filteredRecords" :key="record.id">
              <td>{{ record.month }}</td>
              <td>{{ record.usage }}</td>
              <td>{{ formatCurrency(record.cost) }}</td>
              <td>{{ formatCurrency(record.unitPrice) }}</td>
              <td>{{ record.readingDate }}</td>
              <td>
                <span class="status-badge" :class="record.status">
                  {{ record.statusText }}
                </span>
              </td>
              <td>
                <button class="btn btn-small btn-primary" @click="viewDetails(record)">
                  详情
                </button>
                <button v-if="record.status === 'unpaid'" class="btn btn-small btn-success" @click="payBill(record)">
                  支付
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination" v-if="totalPages > 1">
        <button class="page-btn" :disabled="currentPage === 1" @click="prevPage">
          ← 上一页
        </button>
        <span class="page-info">第 {{ currentPage }} 页 / 共 {{ totalPages }} 页</span>
        <button class="page-btn" :disabled="currentPage === totalPages" @click="nextPage">
          下一页 →
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/NavBar.vue'

interface ElectricityRecord {
  id: number
  month: string
  usage: number
  cost: number
  unitPrice: number
  readingDate: string
  status: string
  statusText: string
}

const router = useRouter()
const authStore = useAuthStore()

// 筛选条件
const selectedYear = ref('2024')
const selectedMonth = ref('all')
const selectedStatus = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)

// 模拟数据
const electricityRecords = ref<ElectricityRecord[]>([
  { id: 1, month: '2024-12', usage: 245.6, cost: 320.50, unitPrice: 1.305, readingDate: '2025-01-01', status: 'unpaid', statusText: '未支付' },
  { id: 2, month: '2024-11', usage: 256.3, cost: 289.45, unitPrice: 1.129, readingDate: '2024-12-01', status: 'paid', statusText: '已支付' },
  { id: 3, month: '2024-10', usage: 228.7, cost: 258.12, unitPrice: 1.129, readingDate: '2024-11-01', status: 'paid', statusText: '已支付' },
  { id: 4, month: '2024-09', usage: 210.5, cost: 237.89, unitPrice: 1.130, readingDate: '2024-10-01', status: 'paid', statusText: '已支付' },
  { id: 5, month: '2024-08', usage: 195.8, cost: 221.45, unitPrice: 1.131, readingDate: '2024-09-01', status: 'paid', statusText: '已支付' },
  { id: 6, month: '2024-07', usage: 189.2, cost: 213.78, unitPrice: 1.130, readingDate: '2024-08-01', status: 'paid', statusText: '已支付' },
  { id: 7, month: '2024-06', usage: 175.4, cost: 198.20, unitPrice: 1.130, readingDate: '2024-07-01', status: 'paid', statusText: '已支付' },
  { id: 8, month: '2024-05', usage: 162.8, cost: 184.00, unitPrice: 1.130, readingDate: '2024-06-01', status: 'paid', statusText: '已支付' },
  { id: 9, month: '2024-04', usage: 158.3, cost: 178.90, unitPrice: 1.130, readingDate: '2024-05-01', status: 'paid', statusText: '已支付' },
  { id: 10, month: '2024-03', usage: 145.6, cost: 164.55, unitPrice: 1.130, readingDate: '2024-04-01', status: 'paid', statusText: '已支付' },
])

// 计算属性
const years = computed(() => ['2024', '2023', '2022'])
const months = computed(() => [
  { value: '01', label: '1月' }, { value: '02', label: '2月' }, { value: '03', label: '3月' },
  { value: '04', label: '4月' }, { value: '05', label: '5月' }, { value: '06', label: '6月' },
  { value: '07', label: '7月' }, { value: '08', label: '8月' }, { value: '09', label: '9月' },
  { value: '10', label: '10月' }, { value: '11', label: '11月' }, { value: '12', label: '12月' },
])

const filteredRecords = computed(() => {
  let filtered = [...electricityRecords.value]

  // 按年份筛选
  if (selectedYear.value !== 'all') {
    filtered = filtered.filter(record => record.month.startsWith(selectedYear.value))
  }

  // 按月份筛选
  if (selectedMonth.value !== 'all') {
    filtered = filtered.filter(record => record.month.endsWith(`-${selectedMonth.value}`))
  }

  // 按状态筛选
  if (selectedStatus.value !== 'all') {
    filtered = filtered.filter(record => record.status === selectedStatus.value)
  }

  // 分页
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filtered.slice(start, end)
})

const totalUsage = computed(() => {
  return electricityRecords.value.reduce((sum, record) => sum + record.usage, 0)
})

const totalCost = computed(() => {
  return electricityRecords.value.reduce((sum, record) => sum + record.cost, 0)
})

const averagePrice = computed(() => {
  return totalUsage.value > 0 ? totalCost.value / totalUsage.value : 0
})

const unpaidCount = computed(() => {
  return electricityRecords.value.filter(record => record.status === 'unpaid').length
})

const totalPages = computed(() => {
  return Math.ceil(electricityRecords.value.length / pageSize.value)
})

// 方法
const formatCurrency = (amount: number) => {
  return amount.toFixed(2)
}

const refreshData = () => {
  // 这里应该调用API刷新数据
  alert('数据已刷新')
}

const exportData = () => {
  // 这里应该实现数据导出功能
  alert('数据导出功能')
}

const filterData = () => {
  currentPage.value = 1
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const viewDetails = (record: ElectricityRecord) => {
  alert(`查看详情: ${record.month} 用电记录`)
}

const payBill = (record: ElectricityRecord) => {
  alert(`支付账单: ${record.month} 电费 ¥${formatCurrency(record.cost)}`)
  // 这里应该调用支付API
}

onMounted(() => {
  // 检查用户是否已登录
  if (!authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<style scoped>
.electricity-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.page-container {
  padding: 30px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 0;
  color: #333;
  font-size: 28px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: opacity 0.3s;
}

.btn:hover {
  opacity: 0.9;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-secondary {
  background: #f0f2f5;
  color: #555;
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-small {
  padding: 6px 12px;
  font-size: 12px;
}

.filter-section {
  background: white;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.filter-row {
  display: flex;
  gap: 20px;
  align-items: flex-end;
}

.filter-group {
  flex: 1;
}

.filter-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #555;
  font-weight: 500;
}

.filter-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  color: #333;
  background-color: white;
  transition: border-color 0.3s;
}

.filter-group select:focus {
  outline: none;
  border-color: #667eea;
}

.stats-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-item {
  background: white;
  border-radius: 10px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.data-table {
  background: white;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.data-table table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 15px;
  text-align: left;
  border-bottom: 1px solid #eaeaea;
}

.data-table th {
  font-weight: 600;
  color: #666;
  background-color: #f9fafb;
}

.data-table tbody tr:hover {
  background-color: #f9fafb;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.paid {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge.unpaid {
  background-color: #fee2e2;
  color: #991b1b;
}

.status-badge.overdue {
  background-color: #fef3c7;
  color: #92400e;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  padding: 20px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  color: #555;
  cursor: pointer;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}
</style>
