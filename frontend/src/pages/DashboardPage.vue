<template>
  <div class="dashboard-page">
    <NavBar />

    <div class="dashboard-container">
      <aside class="sidebar">
        <div class="user-info">
          <div class="avatar">
            {{ userInitials }}
          </div>
          <div class="user-details">
            <h3>{{ user?.realName || user?.username }}</h3>
            <p>{{ user?.email }}</p>
            <p class="customer-id">客户编号: {{ user?.customerNumber || '未绑定' }}</p>
          </div>
        </div>

        <nav class="sidebar-nav">
          <router-link to="/dashboard" class="nav-item active">
            <span class="icon">📊</span>
            <span>仪表板</span>
          </router-link>
          <router-link to="/electricity" class="nav-item">
            <span class="icon">⚡</span>
            <span>电量电费</span>
          </router-link>
          <router-link to="/purchase" class="nav-item">
            <span class="icon">💰</span>
            <span>购电记录</span>
          </router-link>
          <router-link to="/online-pay" class="nav-item">
            <span class="icon">💳</span>
            <span>在线购电</span>
          </router-link>
          <router-link to="/usage-chart" class="nav-item">
            <span class="icon">📈</span>
            <span>用电曲线</span>
          </router-link>
          <router-link to="/notifications" class="nav-item">
            <span class="icon">🔔</span>
            <span>通知中心</span>
          </router-link>
          <router-link to="/information" class="nav-item">
            <span class="icon">📰</span>
            <span>信息发布</span>
          </router-link>
          <router-link to="/applications" class="nav-item">
            <span class="icon">📝</span>
            <span>用电申请</span>
          </router-link>
          <router-link to="/bind-customer" class="nav-item">
            <span class="icon">🔗</span>
            <span>绑定客户</span>
          </router-link>
          <router-link to="/profile" class="nav-item">
            <span class="icon">👤</span>
            <span>个人中心</span>
          </router-link>
        </nav>
      </aside>

      <main class="main-content">
        <div class="content-header">
          <h1>电费管理仪表板</h1>
          <div class="header-actions">
            <button class="btn btn-primary" @click="refreshData">
              <span class="icon">🔄</span> 刷新数据
            </button>
            <button class="btn btn-secondary" @click="handleLogout">
              <span class="icon">🚪</span> 退出登录
            </button>
          </div>
        </div>

        <div class="stats-grid">
          <div class="stat-card balance-card">
            <div class="stat-icon">💰</div>
            <div class="stat-content">
              <h3>电费余额</h3>
              <p class="stat-value">¥ {{ formatCurrency(balance) }}</p>
              <p class="stat-trend" :class="balanceTrend">
                {{ balanceChange >= 0 ? '+' : '' }}{{ formatCurrency(balanceChange) }} (本月)
              </p>
            </div>
          </div>

          <div class="stat-card usage-card">
            <div class="stat-icon">⚡</div>
            <div class="stat-content">
              <h3>本月用电量</h3>
              <p class="stat-value">{{ currentMonthUsage }} kWh</p>
              <p class="stat-trend" :class="usageTrend">
                {{ usageChange >= 0 ? '+' : '' }}{{ usageChange }}% (相比上月)
              </p>
            </div>
          </div>

          <div class="stat-card bill-card">
            <div class="stat-icon">📋</div>
            <div class="stat-content">
              <h3>上月电费</h3>
              <p class="stat-value">¥ {{ formatCurrency(lastMonthBill) }}</p>
              <p class="stat-due" :class="billStatusClass">
                {{ billStatusText }}
              </p>
            </div>
          </div>

          <div class="stat-card notification-card">
            <div class="stat-icon">🔔</div>
            <div class="stat-content">
              <h3>未读通知</h3>
              <p class="stat-value">{{ unreadNotifications }}</p>
              <p class="stat-link">
                <router-link to="/notifications">查看详情 →</router-link>
              </p>
            </div>
          </div>
        </div>

        <div class="content-grid">
          <div class="content-card recent-usage">
            <div class="card-header">
              <h3>最近用电记录</h3>
              <router-link to="/electricity" class="view-all">查看全部 →</router-link>
            </div>
            <div class="card-content">
              <table class="usage-table">
                <thead>
                  <tr>
                    <th>月份</th>
                    <th>用电量 (kWh)</th>
                    <th>电费 (¥)</th>
                    <th>状态</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="record in recentUsage" :key="record.month">
                    <td>{{ record.month }}</td>
                    <td>{{ record.usage }}</td>
                    <td>{{ formatCurrency(record.cost) }}</td>
                    <td>
                      <span class="status-badge" :class="record.status">
                        {{ record.statusText }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="content-card quick-actions">
            <div class="card-header">
              <h3>快捷操作</h3>
            </div>
            <div class="card-content">
              <div class="action-grid">
                <button class="action-btn" @click="goToOnlinePay">
                  <span class="action-icon">💳</span>
                  <span class="action-text">在线购电</span>
                </button>
                <button class="action-btn" @click="goToUsageChart">
                  <span class="action-icon">📈</span>
                  <span class="action-text">查看用电曲线</span>
                </button>
                <button class="action-btn" @click="goToBindCustomer">
                  <span class="action-icon">🔗</span>
                  <span class="action-text">绑定客户编号</span>
                </button>
                <button class="action-btn" @click="goToApplications">
                  <span class="action-icon">📝</span>
                  <span class="action-text">提交用电申请</span>
                </button>
                <button class="action-btn" @click="goToInformation">
                  <span class="action-icon">📰</span>
                  <span class="action-text">查看停电公告</span>
                </button>
                <button class="action-btn" @click="goToProfile">
                  <span class="action-icon">👤</span>
                  <span class="action-text">个人中心</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="content-card upcoming-bills">
          <div class="card-header">
            <h3>即将到期的账单</h3>
          </div>
          <div class="card-content">
            <div v-if="upcomingBills.length === 0" class="empty-state">
              <p>暂无即将到期的账单</p>
            </div>
            <div v-else class="bills-list">
              <div v-for="bill in upcomingBills" :key="bill.id" class="bill-item">
                <div class="bill-info">
                  <h4>{{ bill.title }}</h4>
                  <p>到期日: {{ bill.dueDate }}</p>
                  <p>金额: ¥ {{ formatCurrency(bill.amount) }}</p>
                </div>
                <div class="bill-actions">
                  <button class="btn btn-small btn-primary" @click="payBill(bill)">
                    立即支付
                  </button>
                  <button class="btn btn-small btn-secondary" @click="viewBillDetails(bill)">
                    查看详情
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/NavBar.vue'

interface Bill {
  id: number
  title: string
  dueDate: string
  amount: number
}

interface UsageRecord {
  month: string
  usage: number
  cost: number
  status: string
  statusText: string
}

const router = useRouter()
const authStore = useAuthStore()

// 用户信息
const user = computed(() => authStore.user)

// 用户姓名首字母
const userInitials = computed(() => {
  if (!user.value) return 'U'
  const name = user.value.realName || user.value.username
  return name.charAt(0).toUpperCase()
})

// 模拟数据
const balance = ref(356.78)
const balanceChange = ref(-45.23)
const currentMonthUsage = ref(245.6)
const usageChange = ref(12.5)
const lastMonthBill = ref(289.45)
const unreadNotifications = ref(3)
const recentUsage = ref<UsageRecord[]>([
  { month: '2024-11', usage: 256.3, cost: 289.45, status: 'paid', statusText: '已支付' },
  { month: '2024-10', usage: 228.7, cost: 258.12, status: 'paid', statusText: '已支付' },
  { month: '2024-09', usage: 210.5, cost: 237.89, status: 'paid', statusText: '已支付' },
  { month: '2024-08', usage: 195.8, cost: 221.45, status: 'paid', statusText: '已支付' },
  { month: '2024-07', usage: 189.2, cost: 213.78, status: 'paid', statusText: '已支付' },
])

const upcomingBills = ref<Bill[]>([
  { id: 1, title: '2024年12月电费', dueDate: '2025-01-10', amount: 320.50 },
  { id: 2, title: '2025年1月电费', dueDate: '2025-02-10', amount: 0 },
])

// 计算属性
const balanceTrend = computed(() => balanceChange.value >= 0 ? 'positive' : 'negative')
const usageTrend = computed(() => usageChange.value >= 0 ? 'positive' : 'negative')
const billStatusClass = computed(() => {
  const daysUntilDue = 5 // 模拟数据
  if (daysUntilDue <= 0) return 'overdue'
  if (daysUntilDue <= 3) return 'warning'
  return 'normal'
})
const billStatusText = computed(() => {
  const daysUntilDue = 5 // 模拟数据
  if (daysUntilDue <= 0) return '已逾期'
  if (daysUntilDue <= 3) return `${daysUntilDue}天后到期`
  return '正常'
})

// 方法
const formatCurrency = (amount: number) => {
  return amount.toFixed(2)
}

const refreshData = () => {
  // 这里应该调用API刷新数据
  alert('数据已刷新')
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const goToOnlinePay = () => {
  router.push('/online-pay')
}

const goToUsageChart = () => {
  router.push('/usage-chart')
}

const goToBindCustomer = () => {
  router.push('/bind-customer')
}

const goToApplications = () => {
  router.push('/applications')
}

const goToInformation = () => {
  router.push('/information')
}

const goToProfile = () => {
  router.push('/profile')
}

const payBill = (bill: Bill) => {
  alert(`支付账单: ${bill.title}`)
  // 这里应该调用支付API
}

const viewBillDetails = (bill: Bill) => {
  alert(`查看账单详情: ${bill.title}`)
  // 这里应该跳转到账单详情页面
}

onMounted(() => {
  // 检查用户是否已登录
  if (!authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<style scoped>
.dashboard-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.dashboard-container {
  display: flex;
  min-height: calc(100vh - 60px);
}

.sidebar {
  width: 250px;
  background: white;
  border-right: 1px solid #eaeaea;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eaeaea;
}

.avatar {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: bold;
}

.user-details h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.user-details p {
  margin: 5px 0 0;
  font-size: 14px;
  color: #666;
}

.customer-id {
  font-size: 12px;
  color: #999;
  margin-top: 3px;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 15px;
  color: #555;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s;
}

.nav-item:hover {
  background-color: #f0f2f5;
  color: #667eea;
}

.nav-item.active {
  background-color: #667eea;
  color: white;
}

.nav-item .icon {
  font-size: 18px;
}

.main-content {
  flex: 1;
  padding: 30px;
  overflow-y: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.content-header h1 {
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

.btn-small {
  padding: 6px 12px;
  font-size: 12px;
}

.stats-grid {
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
  margin: 0 0 5px;
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.stat-trend, .stat-due, .stat-link {
  margin: 0;
  font-size: 14px;
}

.stat-trend.positive {
  color: #10b981;
}

.stat-trend.negative {
  color: #ef4444;
}

.stat-due.normal {
  color: #10b981;
}

.stat-due.warning {
  color: #f59e0b;
}

.stat-due.overdue {
  color: #ef4444;
}

.stat-link a {
  color: #667eea;
  text-decoration: none;
}

.stat-link a:hover {
  text-decoration: underline;
}

.content-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.content-card {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.view-all {
  color: #667eea;
  text-decoration: none;
}

.view-all:hover {
  text-decoration: underline;
}

.usage-table {
  width: 100%;
  border-collapse: collapse;
}

.usage-table th,
.usage-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eaeaea;
}

.usage-table th {
  font-weight: 600;
  color: #666;
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

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px;
  background: white;
  border: 1px solid #eaeaea;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.1);
}

.action-icon {
  font-size: 24px;
}

.action-text {
  font-size: 14px;
  color: #333;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #666;
}

.bills-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.bill-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #eaeaea;
}

.bill-info h4 {
  margin: 0 0 5px;
  color: #333;
}

.bill-info p {
  margin: 3px 0;
  font-size: 14px;
  color: #666;
}

.bill-actions {
  display: flex;
  gap: 10px;
}

</style>
