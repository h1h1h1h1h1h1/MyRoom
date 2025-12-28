<template>
  <div class="online-pay-page">
    <div class="page-header">
      <h1>在线购电</h1>
      <p class="subtitle">快速便捷地为您的电表账户充值</p>
    </div>

    <div class="content-wrapper">
      <!-- 左侧：购电表单 -->
      <div class="pay-form-section">
        <div class="card">
          <div class="card-header">
            <h2>选择购电金额</h2>
          </div>
          <div class="card-body">
            <!-- 账户余额信息 -->
            <div class="balance-info">
              <div class="balance-item">
                <span class="label">当前电费余额：</span>
                <span class="value amount">{{ formatCurrency(userStore.balance) }}</span>
              </div>
              <div class="balance-item">
                <span class="label">绑定客户编号：</span>
                <span class="value">{{ userStore.customerNumber || '未绑定' }}</span>
              </div>
              <div class="balance-item">
                <span class="label">上月电费：</span>
                <span class="value amount">{{ formatCurrency(lastMonthBill) }}</span>
              </div>
            </div>

            <!-- 金额选项 -->
            <div class="amount-options">
              <h3>选择充值金额</h3>
              <div class="amount-grid">
                <button
                  v-for="amount in presetAmounts"
                  :key="amount"
                  class="amount-btn"
                  :class="{ active: selectedAmount === amount }"
                  @click="selectAmount(amount)"
                >
                  {{ formatCurrency(amount) }}
                </button>
              </div>
              <div class="custom-amount">
                <label for="customAmount">自定义金额：</label>
                <div class="input-group">
                  <span class="input-prefix">¥</span>
                  <input
                    id="customAmount"
                    v-model="customAmount"
                    type="number"
                    min="10"
                    max="10000"
                    step="10"
                    placeholder="输入金额（10-10000元）"
                    @input="handleCustomAmountInput"
                  />
                </div>
              </div>
            </div>

            <!-- 支付方式 -->
            <div class="payment-methods">
              <h3>选择支付方式</h3>
              <div class="methods-grid">
                <div
                  v-for="method in paymentMethods"
                  :key="method.id"
                  class="method-card"
                  :class="{ active: selectedMethod === method.id }"
                  @click="selectPaymentMethod(method.id)"
                >
                  <div class="method-icon">
                    <i :class="method.icon"></i>
                  </div>
                  <div class="method-info">
                    <h4>{{ method.name }}</h4>
                    <p>{{ method.description }}</p>
                  </div>
                  <div class="method-check">
                    <i class="fas fa-check-circle"></i>
                  </div>
                </div>
              </div>
            </div>

            <!-- 确认支付 -->
            <div class="payment-summary">
              <div class="summary-item">
                <span>充值金额：</span>
                <span class="summary-value">{{ formatCurrency(finalAmount) }}</span>
              </div>
              <div class="summary-item">
                <span>支付方式：</span>
                <span class="summary-value">{{ getPaymentMethodName(selectedMethod) }}</span>
              </div>
              <div class="summary-item total">
                <span>实付金额：</span>
                <span class="summary-value total-amount">{{ formatCurrency(finalAmount) }}</span>
              </div>
              <button
                class="pay-btn"
                :disabled="isPaying || !canPay"
                @click="handlePayment"
              >
                <span v-if="isPaying">
                  <i class="fas fa-spinner fa-spin"></i> 支付处理中...
                </span>
                <span v-else>
                  确认支付 {{ formatCurrency(finalAmount) }}
                </span>
              </button>
              <p class="payment-tips">
                <i class="fas fa-info-circle"></i>
                支付成功后，电费余额将立即更新。如有问题请联系客服。
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：购电记录和提示 -->
      <div class="right-sidebar">
        <!-- 最近购电记录 -->
        <div class="card recent-purchases">
          <div class="card-header">
            <h2>最近购电记录</h2>
            <router-link to="/purchases" class="view-all">查看全部</router-link>
          </div>
          <div class="card-body">
            <div v-if="recentPurchases.length === 0" class="empty-state">
              <i class="fas fa-receipt"></i>
              <p>暂无购电记录</p>
            </div>
            <div v-else class="purchases-list">
              <div
                v-for="purchase in recentPurchases"
                :key="purchase.id"
                class="purchase-item"
              >
                <div class="purchase-info">
                  <div class="purchase-amount">{{ formatCurrency(purchase.amount) }}</div>
                  <div class="purchase-details">
                    <span class="purchase-date">{{ formatDate(purchase.createdAt) }}</span>
                    <span class="purchase-method">{{ purchase.paymentMethod }}</span>
                  </div>
                </div>
                <div class="purchase-status" :class="purchase.status">
                  {{ getStatusText(purchase.status) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 购电提示 -->
        <div class="card payment-tips-card">
          <div class="card-header">
            <h2>购电提示</h2>
          </div>
          <div class="card-body">
            <ul class="tips-list">
              <li>
                <i class="fas fa-bell"></i>
                <span>电费余额不足时，系统会自动发送通知</span>
              </li>
              <li>
                <i class="fas fa-bolt"></i>
                <span>建议保持电费余额充足，避免断电影响</span>
              </li>
              <li>
                <i class="fas fa-shield-alt"></i>
                <span>所有支付均通过安全加密通道处理</span>
              </li>
              <li>
                <i class="fas fa-clock"></i>
                <span>支付成功后，余额立即更新</span>
              </li>
              <li>
                <i class="fas fa-phone-alt"></i>
                <span>如有问题，请联系客服：95598</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- 支付成功弹窗 -->
    <div v-if="showSuccessModal" class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <i class="fas fa-check-circle success-icon"></i>
          <h2>支付成功！</h2>
        </div>
        <div class="modal-body">
          <p>您已成功充值 <strong>{{ formatCurrency(finalAmount) }}</strong></p>
          <p>新的电费余额：<strong>{{ formatCurrency(newBalance) }}</strong></p>
          <p>支付方式：{{ getPaymentMethodName(selectedMethod) }}</p>
          <p>交易时间：{{ formatDate(new Date()) }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn-primary" @click="closeSuccessModal">完成</button>
          <button class="btn-secondary" @click="viewReceipt">查看收据</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { formatCurrency, formatDate } from '@/utils/formatters'

interface PurchaseRecord {
  id: number
  amount: number
  paymentMethod: string
  status: 'success' | 'pending' | 'failed'
  createdAt: string
}

const router = useRouter()
const userStore = useAuthStore()

// 状态
const selectedAmount = ref<number>(100)
const customAmount = ref<string>('')
const selectedMethod = ref<string>('alipay')
const isPaying = ref<boolean>(false)
const showSuccessModal = ref<boolean>(false)
const newBalance = ref<number>(0)
const recentPurchases = ref<PurchaseRecord[]>([])
const lastMonthBill = ref<number>(156.78) // 模拟数据

// 预设金额选项
const presetAmounts = [50, 100, 200, 500, 1000]

// 支付方式
const paymentMethods = [
  {
    id: 'alipay',
    name: '支付宝',
    icon: 'fab fa-alipay',
    description: '推荐使用，即时到账'
  },
  {
    id: 'wechat',
    name: '微信支付',
    icon: 'fab fa-weixin',
    description: '扫码支付，安全便捷'
  },
  {
    id: 'bank',
    name: '银行卡支付',
    icon: 'fas fa-credit-card',
    description: '支持各大银行储蓄卡/信用卡'
  },
  {
    id: 'balance',
    name: '余额支付',
    icon: 'fas fa-wallet',
    description: '使用账户余额支付'
  }
]

// 计算最终金额
const finalAmount = computed(() => {
  if (customAmount.value && Number(customAmount.value) > 0) {
    return Number(customAmount.value)
  }
  return selectedAmount.value
})

// 检查是否可以支付
const canPay = computed(() => {
  return finalAmount.value >= 10 && finalAmount.value <= 10000
})

// 选择金额
const selectAmount = (amount: number) => {
  selectedAmount.value = amount
  customAmount.value = ''
}

// 处理自定义金额输入
const handleCustomAmountInput = () => {
  if (customAmount.value) {
    const amount = Number(customAmount.value)
    if (amount >= 10 && amount <= 10000) {
      selectedAmount.value = 0
    }
  }
}

// 选择支付方式
const selectPaymentMethod = (methodId: string) => {
  selectedMethod.value = methodId
}

// 获取支付方式名称
const getPaymentMethodName = (methodId: string) => {
  const method = paymentMethods.find(m => m.id === methodId)
  return method ? method.name : '未知方式'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    success: '成功',
    pending: '处理中',
    failed: '失败'
  }
  return statusMap[status] || status
}

// 模拟加载最近购电记录
const loadRecentPurchases = () => {
  // 模拟数据
  recentPurchases.value = [
    {
      id: 1,
      amount: 200,
      paymentMethod: '支付宝',
      status: 'success',
      createdAt: '2025-12-25T14:30:00'
    },
    {
      id: 2,
      amount: 100,
      paymentMethod: '微信支付',
      status: 'success',
      createdAt: '2025-12-20T10:15:00'
    },
    {
      id: 3,
      amount: 500,
      paymentMethod: '银行卡支付',
      status: 'success',
      createdAt: '2025-12-15T16:45:00'
    }
  ]
}

// 处理支付
const handlePayment = async () => {
  if (!canPay.value) {
    alert('请输入有效的金额（10-10000元）')
    return
  }

  isPaying.value = true

  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1500))

    // 模拟支付成功
    newBalance.value = userStore.balance + finalAmount.value
    showSuccessModal.value = true

    // 添加到最近记录
    const newPurchase: PurchaseRecord = {
      id: Date.now(),
      amount: finalAmount.value,
      paymentMethod: getPaymentMethodName(selectedMethod.value),
      status: 'success',
      createdAt: new Date().toISOString()
    }
    recentPurchases.value.unshift(newPurchase)

    // 更新用户余额（实际应用中应该调用API）
    userStore.updateBalance(newBalance.value)

  } catch (error) {
    console.error('支付失败:', error)
    alert('支付失败，请重试或联系客服')
  } finally {
    isPaying.value = false
  }
}

// 关闭成功弹窗
const closeSuccessModal = () => {
  showSuccessModal.value = false
  // 重置表单
  selectedAmount.value = 100
  customAmount.value = ''
  selectedMethod.value = 'alipay'
}

// 查看收据
const viewReceipt = () => {
  closeSuccessModal()
  router.push('/purchases')
}

// 初始化
onMounted(() => {
  loadRecentPurchases()
})
</script>

<style scoped>
.online-pay-page {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #666;
}

.content-wrapper {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
}

.card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
  overflow: hidden;
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.view-all {
  color: #1890ff;
  text-decoration: none;
  font-size: 14px;
}

.view-all:hover {
  text-decoration: underline;
}

.card-body {
  padding: 24px;
}

/* 余额信息 */
.balance-info {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 24px;
}

.balance-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.balance-item:last-child {
  margin-bottom: 0;
}

.balance-item .label {
  color: #666;
  font-size: 14px;
}

.balance-item .value {
  font-size: 16px;
  font-weight: 500;
}

.balance-item .amount {
  color: #1890ff;
  font-weight: 600;
  font-size: 18px;
}

/* 金额选项 */
.amount-options h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #1a1a1a;
}

.amount-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.amount-btn {
  padding: 16px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  background: white;
  font-size: 16px;
  font-weight: 500;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
}

.amount-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.amount-btn.active {
  border-color: #1890ff;
  background: #e6f7ff;
  color: #1890ff;
}

.custom-amount {
  margin-top: 20px;
}

.custom-amount label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.input-group {
  display: flex;
  align-items: center;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
}

.input-prefix {
  padding: 0 12px;
  background: #fafafa;
  border-right: 1px solid #d9d9d9;
  height: 40px;
  display: flex;
  align-items: center;
  color: #666;
}

.input-group input {
  flex: 1;
  border: none;
  padding: 0 12px;
  height: 40px;
  font-size: 16px;
  outline: none;
}

.input-group input:focus {
  border-color: #1890ff;
}

/* 支付方式 */
.payment-methods h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 32px 0 16px;
  color: #1a1a1a;
}

.methods-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.method-card {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.method-card:hover {
  border-color: #1890ff;
  background: #f8f9fa;
}

.method-card.active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.method-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
  border-radius: 8px;
  margin-right: 16px;
  font-size: 24px;
  color: #1890ff;
}

.method-info {
  flex: 1;
}

.method-info h4 {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 4px;
  color: #1a1a1a;
}

.method-info p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.method-check {
  color: #52c41a;
  font-size: 20px;
  opacity: 0;
  transition: opacity 0.2s;
}

.method-card.active .method-check {
  opacity: 1;
}

/* 支付摘要 */
.payment-summary {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 16px;
}

.summary-item.total {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  font-size: 18px;
  font-weight: 600;
}

.summary-value {
  font-weight: 500;
  color: #1a1a1a;
}

.total-amount {
  color: #1890ff;
  font-size: 24px;
}

.pay-btn {
  width: 100%;
  padding: 16px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
  margin-top: 24px;
}

.pay-btn:hover:not(:disabled) {
  background: #096dd9;
}

.pay-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
}

.payment-tips {
  margin-top: 16px;
  font-size: 14px;
  color: #666;
  text-align: center;
}

.payment-tips i {
  margin-right: 8px;
  color: #1890ff;
}

/* 最近购电记录 */
.recent-purchases .card-body {
  padding: 0;
}

.empty-state {
  padding: 40px 24px;
  text-align: center;
  color: #999;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.purchases-list {
  max-height: 300px;
  overflow-y: auto;
}

.purchase-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.purchase-item:last-child {
  border-bottom: none;
}

.purchase-info {
  flex: 1;
}

.purchase-amount {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.purchase-details {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color: #666;
}

.purchase-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.purchase-status.success {
  background: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.purchase-status.pending {
  background: #fff7e6;
  color: #fa8c16;
  border: 1px solid #ffd591;
}

.purchase-status.failed {
  background: #fff1f0;
  color: #ff4d4f;
  border: 1px solid #ffa39e;
}

/* 购电提示 */
.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tips-list li {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.tips-list li:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.tips-list i {
  color: #1890ff;
  margin-right: 12px;
  margin-top: 2px;
  font-size: 16px;
}

.tips-list span {
  flex: 1;
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  overflow: hidden;
  animation: modalSlideIn 0.3s ease;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  padding: 32px 24px 24px;
  text-align: center;
  background: #f6ffed;
}

.success-icon {
  font-size: 64px;
  color: #52c41a;
  margin-bottom: 16px;
}

.modal-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.modal-body {
  padding: 24px;
  text-align: center;
}

.modal-body p {
  margin: 0 0 12px;
  font-size: 16px;
  color: #666;
}

.modal-body strong {
  color: #1a1a1a;
}

.modal-footer {
  padding: 24px;
  display: flex;
  gap: 12px;
  border-top: 1px solid #f0f0f0;
}

.btn-primary, .btn-secondary {
  flex: 1;
  padding: 12px;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #1890ff;
  color: white;
  border: none;
}

.btn-primary:hover {
  background: #096dd9;
}

.btn-secondary {
  background: white;
  color: #1890ff;
  border: 1px solid #1890ff;
}

.btn-secondary:hover {
  background: #e6f7ff;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .content-wrapper {
    grid-template-columns: 1fr;
  }

  .right-sidebar {
    margin-top: 24px;
  }
}

@media (max-width: 768px) {
  .online-pay-page {
    padding: 16px;
  }

  .amount-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .modal-footer {
    flex-direction: column;
  }
}
</style>
