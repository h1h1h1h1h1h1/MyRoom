<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <div class="error-code">404</div>
      <h1 class="error-title">页面未找到</h1>
      <p class="error-message">
        抱歉，您访问的页面不存在或已被移除。
      </p>
      <div class="action-buttons">
        <button @click="goBack" class="btn-back">
          ← 返回上一页
        </button>
        <button @click="goHome" class="btn-home">
          返回首页
        </button>
      </div>
      <div class="suggestions">
        <h3>您可以尝试：</h3>
        <ul>
          <li>检查网址是否正确</li>
          <li>使用导航菜单浏览其他页面</li>
          <li>返回仪表板查看您的用电信息</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.not-found-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.not-found-container {
  text-align: center;
  max-width: 600px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  color: #007bff;
  line-height: 1;
  margin-bottom: 20px;
  text-shadow: 3px 3px 0 rgba(0, 123, 255, 0.1);
}

.error-title {
  font-size: 32px;
  color: #333;
  margin-bottom: 15px;
}

.error-message {
  font-size: 18px;
  color: #666;
  margin-bottom: 30px;
  line-height: 1.6;
}

.action-buttons {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-bottom: 40px;
}

.btn-back, .btn-home {
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.btn-back {
  background: #6c757d;
  color: white;
}

.btn-back:hover {
  background: #5a6268;
  transform: translateY(-2px);
}

.btn-home {
  background: #007bff;
  color: white;
}

.btn-home:hover {
  background: #0056b3;
  transform: translateY(-2px);
}

.suggestions {
  text-align: left;
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-top: 30px;
}

.suggestions h3 {
  color: #333;
  margin-bottom: 15px;
  font-size: 18px;
}

.suggestions ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestions li {
  padding: 8px 0;
  color: #666;
  position: relative;
  padding-left: 20px;
}

.suggestions li:before {
  content: "•";
  color: #007bff;
  font-size: 20px;
  position: absolute;
  left: 0;
  top: 5px;
}
</style>
