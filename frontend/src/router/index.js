import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Index',
      component: () => import('@/views/Index.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { guest: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    },
    {
      path: '/task',
      name: 'TasksList',
      component: () => import('@/views/TasksList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/task_types',
      name: 'TaskTypes',
      component: () => import('@/views/TaskTypesList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/tasks',
      name: 'Tasks',
      component: () => import('@/views/TasksList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/tasks/:uuid',
      name: 'TaskDetails',
      component: () => import('@/views/TaskDetails.vue'),
      meta: { requiresAuth: true }
    }
  ],
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Check if the route requires authentication
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!authStore.isAuthenticated) {
      // Пробуем обновить токен только один раз за сессию
      if (!authStore.tokenRefreshed && await authStore.refreshToken()) {
        next()
      } else {
        next({ name: 'Login' })
      }
    } else {
      next()
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    if (authStore.isAuthenticated) {
      next({ name: 'Index' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
