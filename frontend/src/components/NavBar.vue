<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

export default defineComponent({
  name: 'NavBar',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const isMenuOpen = ref(false)

    const menuItems = [
      { name: '仪表盘', path: '/dashboard', icon: 'dashboard' },
      { name: '电量电费', path: '/electricity', icon: 'bolt' },
      { name: '购电记录', path: '/purchase', icon: 'receipt' },
      { name: '在线购电', path: '/online-pay', icon: 'shopping_cart' },
      { name: '用电查询', path: '/usage-chart', icon: 'query_stats' },
      { name: '通知', path: '/notifications', icon: 'notifications' },
      { name: '信息发布', path: '/information', icon: 'info' },
      { name: '用电申请', path: '/applications', icon: 'assignment' },
      { name: '用户绑定', path: '/bind-customer', icon: 'link' },
    ]

    const logout = () => {
      authStore.logout()
      router.push('/login')
    }

    const toggleMenu = () => {
      isMenuOpen.value = !isMenuOpen.value
    }

    return {
      authStore,
      isMenuOpen,
      menuItems,
      logout,
      toggleMenu,
    }
  }
})
</script>

<template>
  <nav class="navbar">
    <div class="navbar-container">
      <div class="navbar-brand">
        <button class="menu-toggle" @click="toggleMenu">
          <span class="material-icons">menu</span>
        </button>
        <h1 class="navbar-title">电费管理系统</h1>
      </div>

      <div class="navbar-menu" :class="{ 'is-open': isMenuOpen }">
        <div class="navbar-items">
          <router-link
            v-for="item in menuItems"
            :key="item.path"
            :to="item.path"
            class="navbar-item"
            active-class="active"
            @click="isMenuOpen = false"
          >
            <span class="material-icons">{{ item.icon }}</span>
            <span>{{ item.name }}</span>
          </router-link>
        </div>

        <div class="navbar-user">
          <div class="user-info">
            <span class="material-icons">account_circle</span>
            <span>{{ authStore.user?.username || '用户' }}</span>
          </div>
          <button class="logout-btn" @click="logout">
            <span class="material-icons">logout</span>
            <span>退出</span>
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background-color: #1976d2;
  color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  flex-direction: column;
}

.navbar-brand {
  display: flex;
  align-items: center;
  padding: 15px 0;
}

.menu-toggle {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  margin-right: 15px;
  display: none;
}

.navbar-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.navbar-menu {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-items {
  display: flex;
  gap: 10px;
}

.navbar-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 10px 15px;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.navbar-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.navbar-item.active {
  background-color: rgba(255, 255, 255, 0.2);
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 5px;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

@media (max-width: 768px) {
  .menu-toggle {
    display: block;
  }

  .navbar-menu {
    display: none;
    flex-direction: column;
    align-items: stretch;
    padding: 10px 0;
  }

  .navbar-menu.is-open {
    display: flex;
  }

  .navbar-items {
    flex-direction: column;
    gap: 5px;
  }

  .navbar-user {
    flex-direction: column;
    gap: 10px;
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }
}
</style>
