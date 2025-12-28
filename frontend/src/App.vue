<script setup lang="ts">
import { RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="app-container">
    <NavBar v-if="authStore.isAuthenticated" />
    <main class="main-content">
      <RouterView />
    </main>
    <footer class="app-footer">
      <p>电费管理系统 &copy; {{ new Date().getFullYear() }}</p>
    </footer>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.app-footer {
  background-color: #f8f9fa;
  padding: 20px;
  text-align: center;
  border-top: 1px solid #dee2e6;
  margin-top: auto;
}
</style>
