import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/pages/RegisterPage.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/pages/DashboardPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/electricity',
      name: 'electricity',
      component: () => import('@/pages/ElectricityPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/purchase',
      name: 'purchase',
      component: () => import('@/pages/PurchasePage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/online-pay',
      name: 'online-pay',
      component: () => import('@/pages/OnlinePayPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/usage-chart',
      name: 'usage-chart',
      component: () => import('@/pages/UsageChartPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/notifications',
      name: 'notifications',
      component: () => import('@/pages/NotificationsPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/information',
      name: 'information',
      component: () => import('@/pages/InformationPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/applications',
      name: 'applications',
      component: () => import('@/pages/ApplicationsPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/bind-customer',
      name: 'bind-customer',
      component: () => import('@/pages/BindCustomerPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/pages/ProfilePage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFoundPage.vue'),
      meta: { requiresAuth: false }
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
