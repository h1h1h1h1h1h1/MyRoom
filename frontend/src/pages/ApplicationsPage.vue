<template>
  <div class="applications-page">
    <NavBar />
    <div class="page-container">
      <div class="page-header">
        <h1>�õ�����</h1>
        <p>������װ�û����롢�����������롢���ܱ�У���ҵ��</p>
      </div>

      <div class="applications-content">
        <!-- ��������ѡ�� -->
        <div class="application-types-grid">
          <div class="type-card" @click="openApplicationModal('new_install')">
            <div class="icon-wrapper new-install">
              <span class="icon"></span>
            </div>
            <h3>��װ�û�����</h3>
            <p>Ϊ�·����³������밲װ���</p>
            <button class="btn-apply">��������</button>
          </div>

          <div class="type-card" @click="openApplicationModal('rename')">
            <div class="icon-wrapper rename">
              <span class="icon"></span>
            </div>
            <h3>������������</h3>
            <p>�����õ继���������</p>
            <button class="btn-apply">��������</button>
          </div>

          <div class="type-card" @click="openApplicationModal('check')">
            <div class="icon-wrapper check">
              <span class="icon"></span>
            </div>
            <h3>���ܱ�У��</h3>
            <p>����Ե������У����</p>
            <button class="btn-apply">��������</button>
          </div>
        </div>

        <!-- �����¼�б� -->
        <div class="applications-history">
          <h2>�ҵ������¼</h2>
          <div v-if="loading" class="loading-state">������...</div>
          <div v-else-if="applications.length === 0" class="empty-state">
            <p>���������¼</p>
          </div>
          <div v-else class="history-list">
            <div v-for="app in applications" :key="app.id" class="history-item">
              <div class="history-icon" :class="app.type">
                {{ getIconForType(app.type) }}
              </div>
              <div class="history-details">
                <h4>{{ getTypeName(app.type) }}</h4>
                <p class="app-date">{{ formatDate(app.created_at) }}</p>
                <p class="app-content">{{ app.content }}</p>
              </div>
              <div class="history-status">
                <span class="status-badge" :class="app.status">
                  {{ getStatusName(app.status) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ���뵯�� -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ getTypeName(currentType) }}</h3>
          <button class="close-btn" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitApplication">
            <div class="form-group">
              <label>����˵��</label>
              <textarea
                v-model="applicationForm.content"
                placeholder="����ϸ�������������������ַ����ϵ�˵���Ϣ..."
                rows="5"
                required
              ></textarea>
            </div>
            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closeModal">ȡ��</button>
              <button type="submit" class="btn-submit" :disabled="submitting">
                {{ submitting ? '�ύ��...' : '�ύ����' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NavBar from '@/components/NavBar.vue'
import { api } from '@/services/api'
import { formatDate } from '@/utils/formatters'

const router = useRouter()
const authStore = useAuthStore()

interface Application {
  id: number
  type: string
  content: string
  status: string
  created_at: string
}

const applications = ref<Application[]>([])
const loading = ref(false)
const showModal = ref(false)
const currentType = ref('')
const submitting = ref(false)

const applicationForm = reactive({
  content: ''
})

const getTypeName = (type: string) => {
  const map: Record<string, string> = {
    'new_install': '��װ�û�����',
    'rename': '������������',
    'check': '���ܱ�У��'
  }
  return map[type] || type
}

const getIconForType = (type: string) => {
  const map: Record<string, string> = {
    'new_install': '',
    'rename': '',
    'check': ''
  }
  return map[type] || ''
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    'pending': '�����',
    'approved': '��ͨ��',
    'rejected': '�Ѳ���',
    'completed': '�����'
  }
  return map[status] || status
}

const fetchApplications = async () => {
  if (!authStore.user) return

  loading.value = true
  try {
    const response = await api.get(`/user/applications?user_id=${authStore.user.id}`)
    applications.value = response.data
  } catch (error) {
    console.error('Failed to fetch applications:', error)
  } finally {
    loading.value = false
  }
}

const openApplicationModal = (type: string) => {
  currentType.value = type
  applicationForm.content = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const submitApplication = async () => {
  if (!authStore.user) return

  submitting.value = true
  try {
    await api.post('/user/application', {
      user_id: authStore.user.id,
      type: currentType.value,
      content: applicationForm.content
    })
    showModal.value = false
    fetchApplications()
    alert('�����ύ�ɹ���')
  } catch (error) {
    console.error('Failed to submit application:', error)
    alert('�ύʧ�ܣ�������')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  fetchApplications()
})
</script>

<style scoped>
.applications-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 20px;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 24px;
  color: #333;
  margin-bottom: 10px;
}

.page-header p {
  color: #666;
}

.application-types-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.type-card {
  background: white;
  border-radius: 12px;
  padding: 30px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.type-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 24px;
}

.icon-wrapper.new-install {
  background-color: #e6f7ff;
  color: #1890ff;
}

.icon-wrapper.rename {
  background-color: #f6ffed;
  color: #52c41a;
}

.icon-wrapper.check {
  background-color: #fff7e6;
  color: #fa8c16;
}

.type-card h3 {
  margin: 0 0 10px;
  color: #333;
}

.type-card p {
  color: #666;
  margin-bottom: 20px;
  font-size: 14px;
}

.btn-apply {
  background-color: #1890ff;
  color: white;
  border: none;
  padding: 8px 24px;
  border-radius: 20px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-apply:hover {
  background-color: #40a9ff;
}

.applications-history {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.applications-history h2 {
  font-size: 20px;
  margin-bottom: 20px;
  color: #333;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.history-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 1px solid #eaeaea;
  border-radius: 8px;
  gap: 15px;
}

.history-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  background-color: #f5f5f5;
}

.history-details {
  flex: 1;
}

.history-details h4 {
  margin: 0 0 5px;
  color: #333;
}

.app-date {
  font-size: 12px;
  color: #999;
  margin: 0 0 5px;
}

.app-content {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.pending {
  background-color: #fff7e6;
  color: #fa8c16;
}

.status-badge.approved {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-badge.rejected {
  background-color: #fff1f0;
  color: #f5222d;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
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
  padding: 20px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.modal-header h3 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
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

.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  resize: vertical;
  font-family: inherit;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-cancel {
  background: white;
  border: 1px solid #ddd;
  padding: 8px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.btn-submit {
  background: #1890ff;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 6px;
  cursor: pointer;
}

.btn-submit:disabled {
  background: #ccc;
  cursor: not-allowed;
}
</style>
