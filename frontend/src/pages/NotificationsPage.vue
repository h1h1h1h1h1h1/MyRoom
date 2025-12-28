<template>
  <div class="notifications-page">
    <div class="page-header">
      <h1 class="page-title">电费余额不足通知</h1>
      <p class="page-subtitle">管理您的电费通知和提醒设置</p>
    </div>

    <div class="notification-stats">
      <div class="stat-card">
        <div class="stat-icon">🔔</div>
        <div class="stat-content">
          <div class="stat-value">{{ unreadCount }}</div>
          <div class="stat-label">未读通知</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">⚠️</div>
        <div class="stat-content">
          <div class="stat-value">{{ warningCount }}</div>
          <div class="stat-label">余额不足提醒</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🔌</div>
        <div class="stat-content">
          <div class="stat-value">{{ criticalCount }}</div>
          <div class="stat-label">即将断电通知</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <div class="stat-value">{{ resolvedCount }}</div>
          <div class="stat-label">已处理通知</div>
        </div>
      </div>
    </div>

    <div class="notification-controls">
      <div class="control-group">
        <label for="notification-type">通知类型：</label>
        <select id="notification-type" v-model="selectedType" class="form-select">
          <option value="all">全部通知</option>
          <option value="warning">余额不足提醒</option>
          <option value="critical">即将断电通知</option>
          <option value="payment">缴费提醒</option>
          <option value="system">系统通知</option>
          <option value="info">信息公告</option>
        </select>
      </div>

      <div class="control-group">
        <label for="notification-status">状态：</label>
        <select id="notification-status" v-model="selectedStatus" class="form-select">
          <option value="all">全部状态</option>
          <option value="unread">未读</option>
          <option value="read">已读</option>
          <option value="resolved">已处理</option>
        </select>
      </div>

      <div class="control-group">
        <label for="date-range">时间范围：</label>
        <select id="date-range" v-model="selectedDateRange" class="form-select">
          <option value="7days">最近7天</option>
          <option value="30days">最近30天</option>
          <option value="90days">最近90天</option>
          <option value="all">全部时间</option>
        </select>
      </div>

      <button @click="markAllAsRead" class="btn btn-secondary" :disabled="unreadCount === 0">
        标记全部为已读
      </button>

      <button @click="clearAll" class="btn btn-danger" :disabled="notifications.length === 0">
        清空通知
      </button>
    </div>

    <div class="notifications-list">
      <div v-if="filteredNotifications.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <h3>暂无通知</h3>
        <p>您当前没有任何通知</p>
      </div>

      <div v-else class="notification-items">
        <div v-for="notification in filteredNotifications" :key="notification.id"
             :class="['notification-item', { unread: !notification.read, critical: notification.type === 'critical' }]"
             @click="toggleRead(notification)">
          <div class="notification-icon">
            <span v-if="notification.type === 'warning'">⚠️</span>
            <span v-else-if="notification.type === 'critical'">🔌</span>
            <span v-else-if="notification.type === 'payment'">💰</span>
            <span v-else-if="notification.type === 'system'">⚙️</span>
            <span v-else>📢</span>
          </div>

          <div class="notification-content">
            <div class="notification-header">
              <h4 class="notification-title">{{ notification.title }}</h4>
              <span class="notification-time">{{ formatDate(notification.createdAt) }}</span>
            </div>

            <p class="notification-message">{{ notification.message }}</p>

            <div class="notification-details">
              <span class="notification-customer" v-if="notification.customerNumber">
                客户编号：{{ formatCustomerNumber(notification.customerNumber) }}
              </span>
              <span class="notification-balance" v-if="notification.balance !== undefined">
                当前余额：{{ formatCurrency(notification.balance) }}
              </span>
              <span class="notification-amount" v-if="notification.amount !== undefined">
                欠费金额：{{ formatCurrency(notification.amount) }}
              </span>
            </div>

            <div class="notification-actions" v-if="!notification.read">
              <button @click.stop="markAsRead(notification)" class="btn btn-sm btn-primary">
                标记为已读
              </button>
              <button @click.stop="resolveNotification(notification)" class="btn btn-sm btn-success">
                已处理
              </button>
            </div>
          </div>

          <div class="notification-status">
            <span v-if="!notification.read" class="badge badge-unread">未读</span>
            <span v-else class="badge badge-read">已读</span>
          </div>
        </div>
      </div>
    </div>

    <div class="notification-settings">
      <h3>通知设置</h3>
      <div class="settings-grid">
        <div class="setting-item">
          <div class="setting-info">
            <h4>余额不足提醒</h4>
            <p>当电费余额低于上月电费时发送提醒</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.balanceWarning" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <h4>二次通知提醒</h4>
            <p>第一次通知后10天发送二次提醒</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.secondWarning" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <h4>断电预警通知</h4>
            <p>二次通知后10天内未缴费将发送断电预警</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.powerCutWarning" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <h4>缴费成功通知</h4>
            <p>缴费成功后发送确认通知</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.paymentSuccess" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <h4>系统维护通知</h4>
            <p>系统维护或升级时发送通知</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.systemMaintenance" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <h4>停电公告通知</h4>
            <p>计划停电时提前发送通知</p>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input type="checkbox" v-model="settings.powerOutage" @change="updateSettings">
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <div class="notification-methods">
        <h4>通知方式</h4>
        <div class="methods-grid">
          <div class="method-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="settings.methods.inApp" @change="updateSettings">
              <span>站内通知</span>
            </label>
          </div>
          <div class="method-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="settings.methods.email" @change="updateSettings">
              <span>电子邮件</span>
            </label>
          </div>
          <div class="method-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="settings.methods.sms" @change="updateSettings">
              <span>短信通知</span>
            </label>
          </div>
          <div class="method-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="settings.methods.push" @change="updateSettings">
              <span>推送通知</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { formatCurrency, formatCustomerNumber, formatDate } from '@/utils/formatters'

interface Notification {
  id: number
  title: string
  message: string
  type: 'warning' | 'critical' | 'payment' | 'system' | 'info'
  read: boolean
  resolved: boolean
  createdAt: string
  customerNumber?: string
  balance?: number
  amount?: number
}

interface NotificationSettings {
  balanceWarning: boolean
  secondWarning: boolean
  powerCutWarning: boolean
  paymentSuccess: boolean
  systemMaintenance: boolean
  powerOutage: boolean
  methods: {
    inApp: boolean
    email: boolean
    sms: boolean
    push: boolean
  }
}

const selectedType = ref('all')
const selectedStatus = ref('all')
const selectedDateRange = ref('30days')

// 模拟通知数据
const notifications = ref<Notification[]>([
  {
    id: 1,
    title: '电费余额不足提醒',
    message: '您的电费余额已低于上月电费，请及时充值以避免影响正常用电。',
    type: 'warning',
    read: false,
    resolved: false,
    createdAt: '2025-12-25T10:30:00',
    customerNumber: '20230001',
    balance: 50.00,
    amount: 192.00
  },
  {
    id: 2,
    title: '二次提醒通知',
    message: '距离第一次余额不足提醒已过去10天，请尽快缴纳电费。',
    type: 'warning',
    read: false,
    resolved: false,
    createdAt: '2025-12-20T14:15:00',
    customerNumber: '20230001',
    balance: 30.00,
    amount: 192.00
  },
  {
    id: 3,
    title: '即将断电预警',
    message: '您的电费已严重不足，如10天内仍未缴费，系统将自动断电。',
    type: 'critical',
    read: false,
    resolved: false,
    createdAt: '2025-12-15T09:45:00',
    customerNumber: '20230001',
    balance: 10.00,
    amount: 192.00
  },
  {
    id: 4,
    title: '缴费成功通知',
    message: '您已成功缴纳电费200元，当前余额为210元。',
    type: 'payment',
    read: true,
    resolved: true,
    createdAt: '2025-12-10T16:20:00',
    customerNumber: '20230001',
    balance: 210.00
  },
  {
    id: 5,
    title: '系统维护通知',
    message: '系统将于2025年12月28日00:00-06:00进行维护，期间可能无法正常使用。',
    type: 'system',
    read: true,
    resolved: true,
    createdAt: '2025-12-05T11:00:00'
  },
  {
    id: 6,
    title: '停电公告',
    message: '计划停电通知：您所在的区域将于2025年12月30日09:00-17:00停电检修。',
    type: 'info',
    read: true,
    resolved: true,
    createdAt: '2025-12-01T08:30:00',
    customerNumber: '20230001'
  },
  {
    id: 7,
    title: '电费账单生成',
    message: '2025年11月电费账单已生成，用电量320kWh，电费192元。',
    type: 'info',
    read: true,
    resolved: true,
    createdAt: '2025-11-28T00:00:00',
    customerNumber: '20230001',
    amount: 192.00
  },
  {
    id: 8,
    title: '余额不足提醒',
    message: '客户编号20230002的电费余额不足，请及时充值。',
    type: 'warning',
    read: false,
    resolved: false,
    createdAt: '2025-11-25T15:30:00',
    customerNumber: '20230002',
    balance: 80.00,
    amount: 168.00
  }
])

const settings = ref<NotificationSettings>({
  balanceWarning: true,
  secondWarning: true,
  powerCutWarning: true,
  paymentSuccess: true,
  systemMaintenance: true,
  powerOutage: true,
  methods: {
    inApp: true,
    email: true,
    sms: false,
    push: true
  }
})

const filteredNotifications = computed(() => {
  let filtered = notifications.value

  // 按类型过滤
  if (selectedType.value !== 'all') {
    filtered = filtered.filter(n => n.type === selectedType.value)
  }

  // 按状态过滤
  if (selectedStatus.value === 'unread') {
    filtered = filtered.filter(n => !n.read)
  } else if (selectedStatus.value === 'read') {
    filtered = filtered.filter(n => n.read)
  } else if (selectedStatus.value === 'resolved') {
    filtered = filtered.filter(n => n.resolved)
  }

  // 按时间范围过滤
  const now = new Date()
  let cutoffDate = new Date()

  switch (selectedDateRange.value) {
    case '7days':
      cutoffDate.setDate(now.getDate() - 7)
      break
    case '30days':
      cutoffDate.setDate(now.getDate() - 30)
      break
    case '90days':
      cutoffDate.setDate(now.getDate() - 90)
      break
    case 'all':
    default:
      cutoffDate = new Date(0) // 最早日期
  }

  filtered = filtered.filter(n => new Date(n.createdAt) >= cutoffDate)

  // 按时间倒序排序
  return filtered.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
})

const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.read).length
})

const warningCount = computed(() => {
  return notifications.value.filter(n => n.type === 'warning' && !n.resolved).length
})

const criticalCount = computed(() => {
  return notifications.value.filter(n => n.type === 'critical' && !n.resolved).length
})

const resolvedCount = computed(() => {
  return notifications.value.filter(n => n.resolved).length
})

onMounted(() => {
  // 从本地存储加载设置
  const savedSettings = localStorage.getItem('notificationSettings')
  if (savedSettings) {
    settings.value = JSON.parse(savedSettings)
  }
})

const toggleRead = (notification: Notification) => {
  notification.read = !notification.read
}

const markAsRead = (notification: Notification) => {
  notification.read = true
}

const markAllAsRead = () => {
  notifications.value.forEach(n => {
    n.read = true
  })
}

const resolveNotification = (notification: Notification) => {
  notification.resolved = true
  notification.read = true
}

const clearAll = () => {
  if (confirm('确定要清空所有通知吗？此操作不可撤销。')) {
    notifications.value = []
  }
}

const updateSettings = () => {
  localStorage.setItem('notificationSettings', JSON.stringify(settings.value))
}
</script>

<style scoped>
.notifications-page {
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

.notification-stats {
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

.notification-controls {
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

.form-select {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-select:focus {
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #6b7280;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #4b5563;
}

.btn-danger {
  background-color: #ef4444;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background-color: #dc2626;
}

.btn-primary {
  background-color: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background-color: #2563eb;
}

.btn-success {
  background-color: #10b981;
  color: white;
}

.btn-success:hover {
  background-color: #059669;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
  height: auto;
}

.notifications-list {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 14px;
  color: #6b7280;
}

.notification-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.2s;
}

.notification-item:hover {
  background-color: #f9fafb;
}

.notification-item.unread {
  background-color: #eff6ff;
  border-color: #bfdbfe;
}

.notification-item.critical {
  background-color: #fef2f2;
  border-color: #fecaca;
}

.notification-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  border-radius: 8px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.notification-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.notification-time {
  font-size: 12px;
  color: #6b7280;
  white-space: nowrap;
}

.notification-message {
  font-size: 14px;
  color: #4b5563;
  margin-bottom: 12px;
  line-height: 1.5;
}

.notification-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
}

.notification-customer,
.notification-balance,
.notification-amount {
  font-size: 12px;
  padding: 4px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  color: #374151;
}

.notification-actions {
  display: flex;
  gap: 8px;
}

.notification-status {
  flex-shrink: 0;
}

.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.badge-unread {
  background-color: #dbeafe;
  color: #1e40af;
}

.badge-read {
  background-color: #d1fae5;
  color: #065f46;
}

.notification-settings {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.notification-settings h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 24px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.setting-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.setting-info p {
  font-size: 12px;
  color: #6b7280;
  margin: 0;
}

.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 34px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #3b82f6;
}

input:checked + .slider:before {
  transform: translateX(26px);
}

.notification-methods h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.methods-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.method-item {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .notification-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .settings-grid {
    grid-template-columns: 1fr;
  }

  .methods-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .notification-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .control-group {
    min-width: auto;
  }

  .notification-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}

</style>
