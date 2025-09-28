import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { refreshAccessToken } from '@/axios'

function isTokenExpired(token) {
  if (!token) return true
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    if (!payload.exp) return false
    return Date.now() >= payload.exp * 1000
  } catch (e) {
    return true
  }
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(null)
  const user_id = ref(null)
  const username = ref(null)
  const tokenRefreshed = ref(false) // Флаг для отслеживания попытки обновления
  const rememberMe = ref(false)
  const isAuthLoading = ref(false)

  const isAuthenticated = computed(() =>
    !!accessToken.value && !isTokenExpired(accessToken.value)
  )
  const authHeader = computed(() =>
    accessToken.value ? { Authorization: 'Bearer ' + accessToken.value } : {}
  )

  function setToken(token, remember = false) {
    accessToken.value = token
    rememberMe.value = remember
    
    // Extract user_id and username from JWT
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        user_id.value = payload.user_id || null
        username.value = payload.username || null
      } catch (e) {
        user_id.value = null
        username.value = null
      }
    } else {
      user_id.value = null
      username.value = null
    }
  }

  function logout() {
    accessToken.value = null
    user_id.value = null
    username.value = null
    rememberMe.value = false
  }
  
  // Add method to refresh token
  async function refreshToken() {
    try {
      const newToken = await refreshAccessToken()
      setToken(newToken)
      return true
    } catch (error) {
      logout()
      return false
    } finally {
      tokenRefreshed.value = true // Отмечаем, что попытка была
    }
  }

  return {
    accessToken,
    user_id,
    username,
    isAuthenticated,
    authHeader,
    setToken,
    logout,
    refreshToken,
    tokenRefreshed
  }
}, {
  persist: {
    afterRestore: (ctx) => {
      // Сбрасываем флаг после восстановления состояния,
      // чтобы при первой загрузке проверка всегда была
      ctx.store.tokenRefreshed = false
    }
  }
})