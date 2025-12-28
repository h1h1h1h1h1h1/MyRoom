<template>
  <div class="purchase-page">
    <NavBar />

    <div class="page-container">
      <div class="page-header">
        <h1>购电记录管理</h1>
        <div class="header-actions">
          <button class="btn btn-primary" @click="goToOnlinePay">
            <span class="icon">💳</span> 在线购电
          </button>
          <button class="btn btn-secondary" @click="refreshData">
            <span class="icon">🔄</span> 刷新记录
          </button>
        </div>
      </div>

      <div class="stats-cards">
        <div class="stat-card">
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <h3>累计购电金额</h3>
            <p class="stat-value">¥ {{ formatCurrency(totalPurchaseAmount) }}</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">📅</div>
          <div class="stat-content">
            <h3>本月购电次数</h3>
            <p class="stat-value">{{ currentMonthPurchases }} 次</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⚡</div>
          <div class="stat-content">
            <h3>累计购电量</h3>
            <p class="stat-value">{{ totalPurchasePower }} kWh</p>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-content">
            <h3>平均购电单价</h3>
            <p class="stat-value">¥ {{ formatCurrency(averagePrice) }}/kWh</p>
          </div>
        </div>
      </div>

      <div class="filter-section">
        <div class="filter-row">
          <div class="filter-group">
            <label for="startDate">开始日期</label>
            <input type="date" id="startDate" v-model="startDate" @change="filterData">
          </div>
          <div class="filter-group">
            <label for="endDate">结束日期</label>
            <input type="date" id="endDate" v-model="endDate" @change="filterData">
          </div>
          <div class="filter-group">
            <label for="paymentMethod">支付方式</label>
            <select id="paymentMethod" v-model="selectedPaymentMethod" @change="filterData">
              <option value="all">全部</option>
              <option value="alipay">支付宝</option>
              <option value="wechat">微信支付</option>
              <option value="bank">网银支付</option>
              <option value="cash">现金</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="status">状态</label>
            <select id="status" v-model="selectedStatus" @change="filterData">
              <option value="all">全部</option>
              <option value="success">成功</option>
              <option value="pending">处理中</option>
              <option value="failed">失败</option>
              <option value="refunded">已退款</option>
            </select>
          </div>
        </div>
      </div>

      <div class="data-table">
        <table>
          <thead>
            <tr>
              <th>订单号</th>
              <th>购电时间</th>
              <th>购电金额 (¥)</th>
              <th>购电量 (kWh)</th>
              <th>单价 (¥/kWh)</th>
              <th>支付方式</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in filteredRecords" :key="record.id">
              <td>{{ record.orderNumber }}</td>
              <td>{{ formatDateTime(record.purchaseTime) }}</td>
              <td>{{ formatCurrency(record.amount) }}</td>
              <td>{{ record.power }}</td>
              <td>{{ formatCurrency(record.unitPrice) }}</td>
              <td>
                <span class="payment-method" :class="record.paymentMethod">
                  {{ record.paymentMethodText }}
                </span>
              </td>
              <td>
                <span class="status-badge" :class="record.status">
                  {{ record.statusText }}
                </span>
              </td>
              <td>
                <button class="btn btn-small btn-primary" @click="viewDetails(record)">
                  详情
                </button>
                <button v-if="record.status === 'failed'" class="btn btn-small btn-success" @click="retryPayment(record)">
                  重试
                </button>
                <button v-if="record.status === 'success'" class="btn btn-small btn-secondary" @click="downloadInvoice(record)">
                  发票
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

      <div class="recent-purchases">
        <h3>最近购电趋势</h3>
        <div class="trend-chart">
          <div class="chart-placeholder">
            <p>📈 购电趋势图表区域</p>
            <p class="chart-hint">这里将显示最近12个月的购电趋势图表</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/NavBar.vue'

interface PurchaseRecord {
  id: number
  orderNumber: string
  purchaseTime: string
  amount: number
  power: number
  unitPrice: number
  paymentMethod: string
  paymentMethodText: string
  status: string
  statusText: string
}

const router = useRouter()
const authStore = useAuthStore()

// 筛选条件
const startDate = ref('2024-01-01')
const endDate = ref('2024-12-31')
const selectedPaymentMethod = ref('all')
const selectedStatus = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)

// 模拟数据
const purchaseRecords = ref<PurchaseRecord[]>([
  { id: 1, orderNumber: 'DD20241227001', purchaseTime: '2024-12-27 14:30:25', amount: 200.00, power: 153.2, unitPrice: 1.305, paymentMethod: 'alipay', paymentMethodText: '支付宝', status: 'success', statusText: '成功' },
  { id: 2, orderNumber: 'DD20241215002', purchaseTime: '2024-12-15 09:15:42', amount: 300.00, power: 229.8, unitPrice: 1.305, paymentMethod: 'wechat', paymentMethodText: '微信支付', status: 'success', statusText: '成功' },
  { id: 3, orderNumber: 'DD20241128003', purchaseTime: '2024-11-28 16:45:18', amount: 150.00, power: 114.9, unitPrice: 1.305, paymentMethod: 'alipay', paymentMethodText: '支付宝', status: 'success', statusText: '成功' },
  { id: 4, orderNumber: 'DD20241110004', purchaseTime: '2024-11-10 11:20:33', amount: 500.00, power: 383.0, unitPrice: 1.305, paymentMethod: 'bank', paymentMethodText: '网银支付', status: 'success', statusText: '成功' },
  { id: 5, orderNumber: 'DD20241022005', purchaseTime: '2024-10-22 13:55:47', amount: 100.00, power: 76.6, unitPrice: 1.305, paymentMethod: 'alipay', paymentMethodText: '支付宝', status: 'success', statusText: '成功' },
  { id: 6, orderNumber: 'DD20241005006', purchaseTime: '2024-10-05 10:10:10', amount: 250.00, power: 191.5, unitPrice: 1.305, paymentMethod: 'wechat', paymentMethodText: '微信支付', status: 'success', statusText: '成功' },
  { id: 7, orderNumber: 'DD20240918007', purchaseTime: '2024-09-18 15:40:22', amount: 400.00, power: 306.4, unitPrice: 1.305, paymentMethod: 'alipay', paymentMethodText: '支付宝', status: 'success', statusText: '成功' },
  { id: 8, orderNumber: 'DD20240901008', purchaseTime: '2024-09-01 08:25:19', amount: 350.00, power: 268.1, unitPrice: 1.305, paymentMethod: 'bank', paymentMethodText: '网银支付', status: 'success', statusText: '成功' },
  { id: 9, orderNumber: 'DD20240815009', purchaseTime: '2024-08-15 17:30:44', amount: 180.00, power: 137.9, unitPrice: 1.305, paymentMethod: 'alipay', paymentMethodText: '支付宝', status: 'success', statusText: '成功' },
  { id: 10, orderNumber: 'DD20240801010', purchaseTime: '2024-08-01 12:05:28', amount: 220.00, power: 168.5, unitPrice: 1.305, paymentMethod: 'wechat', paymentMethodText: '微信支付', status: 'success', statusText: '成功' },
])

// 计算属性
const filteredRecords = computed(() => {
  let filtered = [...purchaseRecords.value]

  // 按日期筛选
  if (startDate.value) {
    filtered = filtered.filter(record => record.purchaseTime >= startDate.value)
  }
  if (endDate.value) {
    filtered = filtered.filter(record => record.purchaseTime <= endDate.value + ' 23:59:59')
  }

  // 按支付方式筛选
  if (selectedPaymentMethod.value !== 'all') {
    filtered = filtered.filter(record => record.paymentMethod === selectedPaymentMethod.value)
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

const totalPurchaseAmount = computed(() => {
  return purchaseRecords.value.reduce((sum, record) => sum + record.amount, 0)
})

const currentMonthPurchases = computed(() => {
  const currentMonth = new Date().getMonth() + 1
  const currentYear = new Date().getFullYear()
  return purchaseRecords.value.filter(record => {
    const recordDate = new Date(record.purchaseTime)
    return recordDate.getMonth() + 1 === currentMonth && recordDate.getFullYear() === currentYear
  }).length
})

const totalPurchasePower = computed(() => {
  return purchaseRecords.value.reduce((sum, record) => sum + record.power, 0)
})

const averagePrice = computed(() => {
  return totalPurchasePower.value > 0 ? totalPurchaseAmount.value / totalPurchasePower.value : 0
})

const totalPages = computed(() => {
  return Math.ceil(purchaseRecords.value.length / pageSize.value)
})

// 方法
const formatCurrency = (amount: number) => {
  return amount.toFixed(2)
}

const formatDateTime = (datetime: string) => {
  return datetime.replace(' ', ' ')
}

const goToOnlinePay = () => {
  router.push('/online-pay')
}

const refreshData = () => {
  // 这里应该调用API刷新数据
  alert('购电记录已刷新')
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

const viewDetails = (record: PurchaseRecord) => {
  alert(`查看订单详情: ${record.orderNumber}`)
}

const retryPayment = (record: PurchaseRecord) => {
  alert(`重试支付: ${record.orderNumber}`)
}

const downloadInvoice = (record: PurchaseRecord) => {
  alert(`下载发票: ${record.orderNumber}`)
}

onMounted(() => {
  // 检查用户是否已登录
  if (!authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<style scoped>
.purchase-page {
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

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 10px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  font-size: 40px;
}

.stat-content h3 {
  margin: 0 0 10px;
  font-size: 16px;
  color: #666;
}

.stat-value {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
  color: #333;
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

.filter-group input,
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

.filter-group input:focus,
.filter-group select:focus {
  outline: none;
  border-color: #667eea;
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

.payment-method {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.payment-method.alipay {
  background-color: #e6f7ff;
  color: #1890ff;
}

.payment-method.wechat {
  background-color: #f6ffed;
  color: #52c41a;
}

.payment-method.bank {
  background-color: #f9f0ff;
  color: #722ed1;
}

.payment-method.cash {
  background-color: #fff7e6;
  color: #fa8c16;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.success {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge.pending {
  background-color: #fef3c7;
  color: #92400e;
}

.status-badge.failed {
  background-color: #fee2e2;
  color: #991b1b;
}

.status-badge.refunded {
  background-color: #e0e7ff;
  color: #3730a3;
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
  margin-bottom: 20px;
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

.recent-purchases {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.recent-purchases h3 {
  margin: 0 0 20px;
  color: #333;
  font-size: 18px;
}

.trend-chart {
  height: 300px;
  background: #f9fafb;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #666;
}

.chart-placeholder p {
  margin: 0;
}

.chart-hint {
  font-size: 14px;
  color: #999;
  margin-top: 10px;
}

</style>
