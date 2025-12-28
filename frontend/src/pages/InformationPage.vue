<template>
  <div class="information-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>信息发布</h1>
      <p class="page-description">查看服务网点、停电公告等信息</p>
    </div>

    <!-- 信息统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-store"></i>
        </div>
        <div class="stat-content">
          <h3>{{ servicePointsCount }}</h3>
          <p>服务网点</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-bolt"></i>
        </div>
        <div class="stat-content">
          <h3>{{ outageNoticesCount }}</h3>
          <p>停电公告</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-newspaper"></i>
        </div>
        <div class="stat-content">
          <h3>{{ newsCount }}</h3>
          <p>最新动态</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-calendar-check"></i>
        </div>
        <div class="stat-content">
          <h3>{{ activeNoticesCount }}</h3>
          <p>有效公告</p>
        </div>
      </div>
    </div>

    <!-- 信息筛选 -->
    <div class="filter-section">
      <div class="filter-controls">
        <div class="filter-group">
          <label for="infoType">信息类型</label>
          <select id="infoType" v-model="selectedType" @change="filterInformation">
            <option value="all">全部类型</option>
            <option value="service_point">服务网点</option>
            <option value="outage_notice">停电公告</option>
            <option value="news">最新动态</option>
            <option value="policy">政策法规</option>
          </select>
        </div>
        <div class="filter-group">
          <label for="status">状态</label>
          <select id="status" v-model="selectedStatus" @change="filterInformation">
            <option value="all">全部状态</option>
            <option value="active">有效</option>
            <option value="expired">已过期</option>
            <option value="upcoming">即将生效</option>
          </select>
        </div>
        <div class="filter-group">
          <label for="dateRange">时间范围</label>
          <select id="dateRange" v-model="selectedDateRange" @change="filterInformation">
            <option value="all">全部时间</option>
            <option value="today">今天</option>
            <option value="week">本周</option>
            <option value="month">本月</option>
            <option value="year">本年</option>
          </select>
        </div>
        <button class="btn btn-primary" @click="resetFilters">
          <i class="fas fa-redo"></i> 重置筛选
        </button>
      </div>
    </div>

    <!-- 信息列表 -->
    <div class="information-list">
      <div class="list-header">
        <h2>信息列表</h2>
        <div class="list-actions">
          <button class="btn btn-secondary" @click="refreshList">
            <i class="fas fa-sync"></i> 刷新
          </button>
          <button class="btn btn-primary" @click="showAddModal = true">
            <i class="fas fa-plus"></i> 发布新信息
          </button>
        </div>
      </div>

      <div class="table-container">
        <table class="information-table">
          <thead>
            <tr>
              <th>标题</th>
              <th>类型</th>
              <th>发布时间</th>
              <th>有效期</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="info in filteredInformation" :key="info.id">
              <td>
                <div class="info-title">
                  <strong>{{ info.title }}</strong>
                  <span v-if="info.isImportant" class="important-badge">重要</span>
                </div>
                <div class="info-summary">{{ info.summary }}</div>
              </td>
              <td>
                <span :class="`type-badge type-${info.type}`">
                  {{ formatInfoType(info.type) }}
                </span>
              </td>
              <td>{{ formatDate(info.publishDate) }}</td>
              <td>{{ formatDate(info.expiryDate) }}</td>
              <td>
                <span :class="`status-badge status-${info.status}`">
                  {{ formatStatus(info.status) }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button class="btn-icon" @click="viewDetails(info)" title="查看详情">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button class="btn-icon" @click="editInfo(info)" title="编辑">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button class="btn-icon btn-danger" @click="deleteInfo(info)" title="删除">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="totalPages > 1">
        <button
          class="pagination-btn"
          :disabled="currentPage === 1"
          @click="changePage(currentPage - 1)"
        >
          <i class="fas fa-chevron-left"></i>
        </button>

        <button
          v-for="page in visiblePages"
          :key="page"
          class="pagination-btn"
          :class="{ active: page === currentPage }"
          @click="changePage(page)"
        >
          {{ page }}
        </button>

        <button
          class="pagination-btn"
          :disabled="currentPage === totalPages"
          @click="changePage(currentPage + 1)"
        >
          <i class="fas fa-chevron-right"></i>
        </button>
      </div>
    </div>

    <!-- 详情模态框 -->
    <div v-if="showDetailModal" class="modal-overlay" @click.self="closeDetailModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ selectedInfo?.title }}</h2>
          <button class="modal-close" @click="closeDetailModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="info-details">
            <div class="detail-row">
              <span class="detail-label">信息类型：</span>
              <span :class="`type-badge type-${selectedInfo?.type}`">
                {{ formatInfoType(selectedInfo?.type || '') }}
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">发布时间：</span>
              <span>{{ formatDate(selectedInfo?.publishDate || '') }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">有效期至：</span>
              <span>{{ formatDate(selectedInfo?.expiryDate || '') }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">状态：</span>
              <span :class="`status-badge status-${selectedInfo?.status}`">
                {{ formatStatus(selectedInfo?.status || '') }}
              </span>
            </div>
            <div class="detail-row full-width">
              <span class="detail-label">详细内容：</span>
              <div class="detail-content">
                {{ selectedInfo?.content }}
              </div>
            </div>
            <div v-if="selectedInfo?.type === 'service_point'" class="detail-row full-width">
              <span class="detail-label">服务网点信息：</span>
              <div class="service-point-details">
                <p><strong>地址：</strong>{{ selectedInfo?.address }}</p>
                <p><strong>联系电话：</strong>{{ selectedInfo?.phone }}</p>
                <p><strong>营业时间：</strong>{{ selectedInfo?.businessHours }}</p>
                <p><strong>服务范围：</strong>{{ selectedInfo?.serviceScope }}</p>
              </div>
            </div>
            <div v-if="selectedInfo?.type === 'outage_notice'" class="detail-row full-width">
              <span class="detail-label">停电信息：</span>
              <div class="outage-details">
                <p><strong>停电区域：</strong>{{ selectedInfo?.affectedArea }}</p>
                <p><strong>停电时间：</strong>{{ selectedInfo?.outageTime }}</p>
                <p><strong>预计恢复时间：</strong>{{ selectedInfo?.estimatedRestoreTime }}</p>
                <p><strong>影响用户数：</strong>{{ selectedInfo?.affectedUsers }} 户</p>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeDetailModal">关闭</button>
          <button class="btn btn-primary" @click="printInfo">打印</button>
        </div>
      </div>
    </div>

    <!-- 添加/编辑模态框 -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="closeAddModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEditing ? '编辑信息' : '发布新信息' }}</h2>
          <button class="modal-close" @click="closeAddModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveInfo">
            <div class="form-group">
              <label for="title">标题 *</label>
              <input
                type="text"
                id="title"
                v-model="formData.title"
                required
                placeholder="请输入信息标题"
              />
            </div>
            <div class="form-group">
              <label for="type">信息类型 *</label>
              <select id="type" v-model="formData.type" required>
                <option value="service_point">服务网点</option>
                <option value="outage_notice">停电公告</option>
                <option value="news">最新动态</option>
                <option value="policy">政策法规</option>
              </select>
            </div>
            <div class="form-group">
              <label for="summary">摘要</label>
              <textarea
                id="summary"
                v-model="formData.summary"
                rows="2"
                placeholder="请输入信息摘要"
              ></textarea>
            </div>
            <div class="form-group">
              <label for="content">详细内容 *</label>
              <textarea
                id="content"
                v-model="formData.content"
                rows="5"
                required
                placeholder="请输入详细内容"
              ></textarea>
            </div>
            <div v-if="formData.type === 'service_point'" class="form-group">
              <label for="address">地址</label>
              <input
                type="text"
                id="address"
                v-model="formData.address"
                placeholder="请输入服务网点地址"
              />
            </div>
            <div v-if="formData.type === 'service_point'" class="form-group">
              <label for="phone">联系电话</label>
              <input
                type="tel"
                id="phone"
                v-model="formData.phone"
                placeholder="请输入联系电话"
              />
            </div>
            <div v-if="formData.type === 'outage_notice'" class="form-group">
              <label for="affectedArea">停电区域</label>
              <input
                type="text"
                id="affectedArea"
                v-model="formData.affectedArea"
                placeholder="请输入停电区域"
              />
            </div>
            <div class="form-group">
              <label>
                <input type="checkbox" v-model="formData.isImportant" />
                标记为重要信息
              </label>
            </div>
            <div class="form-group">
              <label for="expiryDate">有效期至</label>
              <input
                type="date"
                id="expiryDate"
                v-model="formData.expiryDate"
              />
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeAddModal">取消</button>
          <button class="btn btn-primary" @click="saveInfo">
            {{ isEditing ? '更新' : '发布' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { formatDate } from '@/utils/formatters'

// 类型定义
interface Information {
  id: number
  title: string
  type: string
  summary: string
  content: string
  publishDate: string
  expiryDate: string
  status: string
  isImportant: boolean
  address?: string
  phone?: string
  businessHours?: string
  serviceScope?: string
  affectedArea?: string
  outageTime?: string
  estimatedRestoreTime?: string
  affectedUsers?: number
}

// 响应式数据
const servicePointsCount = ref(15)
const outageNoticesCount = ref(8)
const newsCount = ref(23)
const activeNoticesCount = ref(31)

const selectedType = ref('all')
const selectedStatus = ref('all')
const selectedDateRange = ref('all')
const currentPage = ref(1)
const itemsPerPage = ref(10)

const informationList = ref<Information[]>([
  {
    id: 1,
    title: '城东供电营业厅搬迁通知',
    type: 'service_point',
    summary: '城东供电营业厅将于2024年1月1日搬迁至新址',
    content: '为提供更好的服务环境，城东供电营业厅将于2024年1月1日搬迁至新址：人民路123号。新营业厅面积更大，设施更完善，欢迎广大用户前来办理业务。',
    publishDate: '2024-12-20',
    expiryDate: '2024-12-31',
    status: 'active',
    isImportant: true,
    address: '人民路123号',
    phone: '0571-88888888',
    businessHours: '周一至周五 8:30-17:00',
    serviceScope: '电费缴纳、业务咨询、用电申请'
  },
  {
    id: 2,
    title: '12月25日计划停电通知',
    type: 'outage_notice',
    summary: '因线路检修，部分区域将于12月25日停电',
    content: '为保障电网安全运行，计划于2024年12月25日对10kV城东线进行检修，届时将影响部分区域供电。',
    publishDate: '2024-12-18',
    expiryDate: '2024-12-25',
    status: 'active',
    isImportant: true,
    affectedArea: '城东区人民路、建设路周边',
    outageTime: '2024-12-25 08:00-17:00',
    estimatedRestoreTime: '2024-12-25 17:00',
    affectedUsers: 1500
  },
  {
    id: 3,
    title: '冬季用电安全提示',
    type: 'news',
    summary: '冬季用电高峰，请注意用电安全',
    content: '随着冬季来临，用电负荷增加，提醒广大用户注意用电安全，避免超负荷用电，及时检查家中线路。',
    publishDate: '2024-12-15',
    expiryDate: '2025-01-15',
    status: 'active',
    isImportant: false
  },
  {
    id: 4,
    title: '城南供电营业厅',
    type: 'service_point',
    summary: '城南区主要供电服务网点',
    content: '城南供电营业厅提供全面的供电服务，包括电费缴纳、业务咨询、用电申请等。',
    publishDate: '2024-12-10',
    expiryDate: '2025-12-10',
    status: 'active',
    isImportant: false,
    address: '城南路456号',
    phone: '0571-88888889',
    businessHours: '周一至周日 8:30-17:00',
    serviceScope: '电费缴纳、业务咨询、用电申请、故障报修'
  },
  {
    id: 5,
    title: '新电价政策解读',
    type: 'policy',
    summary: '2024年新电价政策详细解读',
    content: '根据国家发改委通知，2024年将实施新的电价政策，本次调整主要涉及峰谷电价时段和价格。',
    publishDate: '2024-12-05',
    expiryDate: '2025-12-05',
    status: 'active',
    isImportant: true
  }
])

const showDetailModal = ref(false)
const showAddModal = ref(false)
const selectedInfo = ref<Information | null>(null)
const isEditing = ref(false)

const formData = ref({
  title: '',
  type: 'service_point',
  summary: '',
  content: '',
  address: '',
  phone: '',
  affectedArea: '',
  isImportant: false,
  expiryDate: ''
})

// 计算属性
const filteredInformation = computed(() => {
  let filtered = informationList.value

  if (selectedType.value !== 'all') {
    filtered = filtered.filter(info => info.type === selectedType.value)
  }

  if (selectedStatus.value !== 'all') {
    filtered = filtered.filter(info => info.status === selectedStatus.value)
  }

  if (selectedDateRange.value !== 'all') {
    const now = new Date()
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())

    switch (selectedDateRange.value) {
      case 'today':
        filtered = filtered.filter(info => {
          const publishDate = new Date(info.publishDate)
          return publishDate >= today
        })
        break
      case 'week': {
        const weekAgo = new Date(today)
        weekAgo.setDate(weekAgo.getDate() - 7)
        filtered = filtered.filter(info => {
          const publishDate = new Date(info.publishDate)
          return publishDate >= weekAgo
        })
        break
      }
      case 'month': {
        const monthAgo = new Date(today)
        monthAgo.setMonth(monthAgo.getMonth() - 1)
        filtered = filtered.filter(info => {
          const publishDate = new Date(info.publishDate)
          return publishDate >= monthAgo
        })
        break
      }
      case 'year': {
        const yearAgo = new Date(today)
        yearAgo.setFullYear(yearAgo.getFullYear() - 1)
        filtered = filtered.filter(info => {
          const publishDate = new Date(info.publishDate)
          return publishDate >= yearAgo
        })
        break
      }
    }
  }

  return filtered
})

const totalPages = computed(() => {
  return Math.ceil(filteredInformation.value.length / itemsPerPage.value)
})

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  const end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

// 方法函数
const filterInformation = () => {
  currentPage.value = 1
}

const resetFilters = () => {
  selectedType.value = 'all'
  selectedStatus.value = 'all'
  selectedDateRange.value = 'all'
  currentPage.value = 1
}

const refreshList = () => {
  // 这里应该调用API获取最新数据
  console.log('刷新信息列表')
}

const viewDetails = (info: Information) => {
  selectedInfo.value = info
  showDetailModal.value = true
}

const editInfo = (info: Information) => {
  isEditing.value = true
  formData.value = {
    title: info.title,
    type: info.type,
    summary: info.summary,
    content: info.content,
    address: info.address || '',
    phone: info.phone || '',
    affectedArea: info.affectedArea || '',
    isImportant: info.isImportant,
    expiryDate: info.expiryDate
  }
  showAddModal.value = true
}

const deleteInfo = (info: Information) => {
  if (confirm(`确定要删除"${info.title}"吗？`)) {
    const index = informationList.value.findIndex(item => item.id === info.id)
    if (index !== -1) {
      informationList.value.splice(index, 1)
    }
  }
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedInfo.value = null
}

const printInfo = () => {
  window.print()
}

const closeAddModal = () => {
  showAddModal.value = false
  isEditing.value = false
  formData.value = {
    title: '',
    type: 'service_point',
    summary: '',
    content: '',
    address: '',
    phone: '',
    affectedArea: '',
    isImportant: false,
    expiryDate: ''
  }
}

const saveInfo = () => {
  if (!formData.value.title.trim()) {
    alert('请输入标题')
    return
  }

  if (!formData.value.content.trim()) {
    alert('请输入详细内容')
    return
  }

  if (isEditing.value && selectedInfo.value) {
    // 更新现有信息
    const index = informationList.value.findIndex(item => item.id === selectedInfo.value!.id)
    if (index !== -1) {
      informationList.value[index] = {
        ...informationList.value[index],
        ...formData.value,
        id: selectedInfo.value!.id,
        publishDate: new Date().toISOString().split('T')[0] || '',
        status: 'active'
      }
    }
  } else {
    // 添加新信息
    const newInfo: Information = {
      id: informationList.value.length + 1,
      title: formData.value.title,
      type: formData.value.type,
      summary: formData.value.summary,
      content: formData.value.content,
      publishDate: new Date().toISOString().split('T')[0] || '',
      expiryDate: formData.value.expiryDate || new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0] || '',
      status: 'active',
      isImportant: formData.value.isImportant,
      address: formData.value.address || undefined,
      phone: formData.value.phone || undefined,
      affectedArea: formData.value.affectedArea || undefined
    }
    informationList.value.unshift(newInfo)
  }

  closeAddModal()
}

const changePage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const formatInfoType = (type: string | undefined): string => {
  if (!type) return '未知类型'
  const typeMap: Record<string, string> = {
    'service_point': '服务网点',
    'outage_notice': '停电公告',
    'news': '最新动态',
    'policy': '政策法规'
  }
  return typeMap[type] || type
}

const formatStatus = (status: string | undefined): string => {
  if (!status) return '未知状态'
  const statusMap: Record<string, string> = {
    'active': '有效',
    'expired': '已过期',
    'upcoming': '即将生效'
  }
  return statusMap[status] || status
}

// 生命周期
onMounted(() => {
  // 初始化时可以加载数据
  console.log('信息发布页面已加载')
})
</script>

<style scoped>
.information-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 28px;
  color: #333;
  margin-bottom: 8px;
}

.page-description {
  color: #666;
  font-size: 16px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
}

.stat-content h3 {
  font-size: 24px;
  color: #333;
  margin-bottom: 5px;
}

.stat-content p {
  color: #666;
  font-size: 14px;
}

.filter-section {
  background: white;
  border-radius: 10px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filter-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  align-items: flex-end;
}

.filter-group {
  flex: 1;
  min-width: 150px;
}

.filter-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.filter-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
  background: white;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  font-size: 14px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

.btn-secondary {
  background: #f8f9fa;
  color: #333;
  border: 1px solid #ddd;
}

.btn-secondary:hover {
  background: #e9ecef;
}

.information-list {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.list-header h2 {
  font-size: 20px;
  color: #333;
}

.list-actions {
  display: flex;
  gap: 10px;
}

.table-container {
  overflow-x: auto;
}

.information-table {
  width: 100%;
  border-collapse: collapse;
}

.information-table th {
  background: #f8f9fa;
  padding: 12px 15px;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #dee2e6;
}

.information-table td {
  padding: 15px;
  border-bottom: 1px solid #eee;
  vertical-align: top;
}

.info-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 5px;
}

.important-badge {
  background: #ff6b6b;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.info-summary {
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.type-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 12px;
  font-weight: 500;
}

.type-service_point {
  background: #e3f2fd;
  color: #1976d2;
}

.type-outage_notice {
  background: #fff3e0;
  color: #f57c00;
}

.type-news {
  background: #e8f5e9;
  color: #388e3c;
}

.type-policy {
  background: #f3e5f5;
  color: #7b1fa2;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 12px;
  font-weight: 500;
}

.status-active {
  background: #e8f5e9;
  color: #388e3c;
}

.status-expired {
  background: #ffebee;
  color: #d32f2f;
}

.status-upcoming {
  background: #fff3e0;
  color: #f57c00;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 5px;
  background: #f8f9fa;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.btn-icon:hover {
  background: #e9ecef;
  color: #333;
}

.btn-danger {
  color: #dc3545;
}

.btn-danger:hover {
  background: #dc3545;
  color: white;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 5px;
  margin-top: 20px;
}

.pagination-btn {
  width: 36px;
  height: 36px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pagination-btn:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
}

.pagination-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

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
  border-radius: 10px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  font-size: 20px;
  color: #333;
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 5px;
  transition: all 0.3s ease;
}

.modal-close:hover {
  background: #f8f9fa;
  color: #333;
}

.modal-body {
  padding: 20px;
}

.info-details {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.detail-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.detail-row.full-width {
  flex-direction: column;
  align-items: stretch;
}

.detail-label {
  font-weight: 600;
  color: #333;
  min-width: 100px;
}

.detail-content {
  flex: 1;
  color: #666;
  line-height: 1.6;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 5px;
}

.service-point-details,
.outage-details {
  padding: 10px;
  background: #f8f9fa;
  border-radius: 5px;
  margin-top: 5px;
}

.service-point-details p,
.outage-details p {
  margin: 5px 0;
  color: #666;
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
}

.form-group input[type="text"],
.form-group input[type="tel"],
.form-group input[type="date"],
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
  background: white;
}

.form-group textarea {
  resize: vertical;
  min-height: 60px;
}

.form-group input[type="checkbox"] {
  margin-right: 8px;
}

.form-group label input[type="checkbox"] {
  display: inline-block;
  width: auto;
}

@media (max-width: 768px) {
  .information-page {
    padding: 10px;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .filter-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-group {
    min-width: 100%;
  }

  .list-header {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }

  .list-actions {
    justify-content: flex-start;
  }

  .modal-content {
    width: 95%;
    margin: 10px;
  }

  .detail-row {
    flex-direction: column;
    align-items: stretch;
    gap: 5px;
  }

  .detail-label {
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    flex-direction: column;
    gap: 5px;
  }

  .btn-icon {
    width: 28px;
    height: 28px;
  }
}
</style>
