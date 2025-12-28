<template>
  <div class="usage-chart-page">
    <div class="page-header">
      <h1 class="page-title">用电信息查询</h1>
      <p class="page-subtitle">查看您的用电量趋势和统计分析</p>
    </div>

    <div class="controls-section">
      <div class="control-group">
        <label for="customer-select">选择客户编号：</label>
        <select id="customer-select" v-model="selectedCustomer" class="form-select">
          <option value="">全部客户</option>
          <option v-for="customer in customers" :key="customer.id" :value="customer.customerNumber">
            {{ formatCustomerNumber(customer.customerNumber) }} - {{ customer.address }}
          </option>
        </select>
      </div>

      <div class="control-group">
        <label for="period-select">查询周期：</label>
        <select id="period-select" v-model="selectedPeriod" class="form-select">
          <option value="monthly">月度</option>
          <option value="quarterly">季度</option>
          <option value="yearly">年度</option>
          <option value="custom">自定义</option>
        </select>
      </div>

      <div v-if="selectedPeriod === 'custom'" class="control-group">
        <label for="start-date">开始日期：</label>
        <input type="date" id="start-date" v-model="startDate" class="form-input">
        <label for="end-date">结束日期：</label>
        <input type="date" id="end-date" v-model="endDate" class="form-input">
      </div>

      <button @click="fetchData" class="btn btn-primary">
        <span v-if="loading">查询中...</span>
        <span v-else>查询</span>
      </button>
    </div>

    <div class="charts-section">
      <div class="chart-container">
        <h3>用电量趋势图</h3>
        <div class="chart-wrapper">
          <!-- 这里应该使用图表库，如Chart.js或ECharts -->
          <div class="mock-chart">
            <div class="chart-placeholder">
              <div class="chart-title">用电量趋势图（千瓦时）</div>
              <div class="chart-content">
                <div class="chart-bars">
                  <div v-for="(data, index) in chartData" :key="index" class="chart-bar-container">
                    <div class="chart-bar" :style="{ height: (data.usage / maxUsage * 100) + '%' }"></div>
                    <div class="chart-label">{{ data.month }}</div>
                  </div>
                </div>
                <div class="chart-y-axis">
                  <div v-for="tick in yAxisTicks" :key="tick" class="y-tick">{{ tick }} kWh</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="chart-container">
        <h3>电费趋势图</h3>
        <div class="chart-wrapper">
          <div class="mock-chart">
            <div class="chart-placeholder">
              <div class="chart-title">电费趋势图（元）</div>
              <div class="chart-content">
                <div class="chart-bars">
                  <div v-for="(data, index) in chartData" :key="index" class="chart-bar-container">
                    <div class="chart-bar bill-bar" :style="{ height: (data.bill / maxBill * 100) + '%' }"></div>
                    <div class="chart-label">{{ data.month }}</div>
                  </div>
                </div>
                <div class="chart-y-axis">
                  <div v-for="tick in billYAxisTicks" :key="tick" class="y-tick">{{ formatCurrency(tick) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="stats-section">
      <div class="stat-card">
        <div class="stat-icon">⚡</div>
        <div class="stat-content">
          <div class="stat-value">{{ formatElectricityUsage(totalUsage) }}</div>
          <div class="stat-label">总用电量</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">💰</div>
        <div class="stat-content">
          <div class="stat-value">{{ formatCurrency(totalBill) }}</div>
          <div class="stat-label">总电费</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📈</div>
        <div class="stat-content">
          <div class="stat-value">{{ formatElectricityUsage(averageUsage) }}</div>
          <div class="stat-label">月均用电量</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-value">{{ formatCurrency(averageBill) }}</div>
          <div class="stat-label">月均电费</div>
        </div>
      </div>
    </div>

    <div class="data-table-section">
      <h3>详细数据</h3>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>月份</th>
              <th>客户编号</th>
              <th>用电量 (kWh)</th>
              <th>电费 (元)</th>
              <th>单价 (元/kWh)</th>
              <th>用电状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in tableData" :key="index">
              <td>{{ record.month }}</td>
              <td>{{ formatCustomerNumber(record.customerNumber) }}</td>
              <td>{{ formatElectricityUsage(record.usage) }}</td>
              <td>{{ formatCurrency(record.bill) }}</td>
              <td>{{ record.unitPrice.toFixed(2) }}</td>
              <td :class="formatStatusColor(record.status)">
                {{ formatStatus(record.status) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="export-section">
      <button @click="exportData" class="btn btn-secondary">
        <span>导出数据 (CSV)</span>
      </button>
      <button @click="printChart" class="btn btn-secondary">
        <span>打印图表</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { formatCurrency, formatCustomerNumber, formatStatus, formatStatusColor, formatElectricityUsage } from '@/utils/formatters'

interface Customer {
  id: number
  customerNumber: string
  address: string
}

interface ChartData {
  month: string
  usage: number
  bill: number
  unitPrice: number
  status: string
  customerNumber: string
}

const selectedCustomer = ref('')
const selectedPeriod = ref('monthly')
const startDate = ref('')
const endDate = ref('')
const loading = ref(false)

// 模拟数据
const customers = ref<Customer[]>([
  { id: 1, customerNumber: '20230001', address: '北京市朝阳区建国路88号' },
  { id: 2, customerNumber: '20230002', address: '北京市海淀区中关村大街1号' },
  { id: 3, customerNumber: '20230003', address: '北京市东城区王府井大街100号' }
])

const chartData = ref<ChartData[]>([
  { month: '1月', usage: 320, bill: 192, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '2月', usage: 280, bill: 168, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '3月', usage: 350, bill: 210, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '4月', usage: 420, bill: 252, unitPrice: 0.6, status: 'warning', customerNumber: '20230001' },
  { month: '5月', usage: 380, bill: 228, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '6月', usage: 450, bill: 270, unitPrice: 0.6, status: 'critical', customerNumber: '20230001' },
  { month: '7月', usage: 520, bill: 312, unitPrice: 0.6, status: 'critical', customerNumber: '20230001' },
  { month: '8月', usage: 480, bill: 288, unitPrice: 0.6, status: 'warning', customerNumber: '20230001' },
  { month: '9月', usage: 400, bill: 240, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '10月', usage: 360, bill: 216, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '11月', usage: 320, bill: 192, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' },
  { month: '12月', usage: 300, bill: 180, unitPrice: 0.6, status: 'normal', customerNumber: '20230001' }
])

const tableData = computed(() => {
  if (selectedCustomer.value) {
    return chartData.value.filter(data => data.customerNumber === selectedCustomer.value)
  }
  return chartData.value
})

const maxUsage = computed(() => {
  return Math.max(...chartData.value.map(data => data.usage))
})

const maxBill = computed(() => {
  return Math.max(...chartData.value.map(data => data.bill))
})

const yAxisTicks = computed(() => {
  const max = maxUsage.value
  const ticks = []
  for (let i = 0; i <= 5; i++) {
    ticks.push(Math.round((max / 5) * i))
  }
  return ticks
})

const billYAxisTicks = computed(() => {
  const max = maxBill.value
  const ticks = []
  for (let i = 0; i <= 5; i++) {
    ticks.push(Math.round((max / 5) * i))
  }
  return ticks
})

const totalUsage = computed(() => {
  return tableData.value.reduce((sum, data) => sum + data.usage, 0)
})

const totalBill = computed(() => {
  return tableData.value.reduce((sum, data) => sum + data.bill, 0)
})

const averageUsage = computed(() => {
  return tableData.value.length > 0 ? totalUsage.value / tableData.value.length : 0
})

const averageBill = computed(() => {
  return tableData.value.length > 0 ? totalBill.value / tableData.value.length : 0
})

onMounted(() => {
  // 设置默认日期
  const today = new Date()
  const lastYear = new Date(today.getFullYear() - 1, today.getMonth(), today.getDate())
  startDate.value = lastYear.toISOString().split('T')[0] || ''
  endDate.value = today.toISOString().split('T')[0] || ''
})

const fetchData = async () => {
  loading.value = true
  // 模拟API调用
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
}

const exportData = () => {
  alert('数据导出功能将在实际项目中实现')
}

const printChart = () => {
  window.print()
}
</script>

<style scoped>
.usage-chart-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 16px;
  color: #6b7280;
}

.controls-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-end;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 200px;
}

.control-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-select,
.form-input {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-select:focus,
.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.btn {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  height: 40px;
}

.btn-primary {
  background-color: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background-color: #2563eb;
}

.btn-secondary {
  background-color: #6b7280;
  color: white;
}

.btn-secondary:hover {
  background-color: #4b5563;
}

.charts-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 24px;
}

.chart-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chart-container h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.chart-wrapper {
  height: 300px;
}

.mock-chart {
  height: 100%;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  width: 90%;
  height: 90%;
  display: flex;
  flex-direction: column;
}

.chart-title {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 16px;
  text-align: center;
}

.chart-content {
  flex: 1;
  display: flex;
  position: relative;
}

.chart-bars {
  flex: 1;
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  padding-bottom: 30px;
}

.chart-bar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  width: 30px;
}

.chart-bar {
  width: 20px;
  background-color: #3b82f6;
  border-radius: 4px 4px 0 0;
  min-height: 1px;
}

.bill-bar {
  background-color: #10b981;
}

.chart-label {
  font-size: 12px;
  color: #6b7280;
  margin-top: 8px;
  transform: rotate(-45deg);
  white-space: nowrap;
}

.chart-y-axis {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 30px;
  width: 60px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-end;
  padding-right: 8px;
}

.y-tick {
  font-size: 12px;
  color: #6b7280;
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 32px;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  border-radius: 50%;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}

.data-table-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.data-table-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  background-color: #f9fafb;
  padding: 12px 16px;
  text-align: left;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.data-table td {
  padding: 12px 16px;
  font-size: 14px;
  color: #1f2937;
  border-bottom: 1px solid #f3f4f6;
}

.data-table tbody tr:hover {
  background-color: #f9fafb;
}

.export-section {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .charts-section {
    grid-template-columns: 1fr;
  }

  .stats-section {
    grid-template-columns: repeat(2, 1fr);
  }

  .controls-section {
    flex-direction: column;
    align-items: stretch;
  }

  .control-group {
    min-width: auto;
  }
}

</style>
