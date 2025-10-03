import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLogsStore = defineStore('logs', () => {
  // Интервал автообновления в секундах (по умолчанию 5 секунд)
  const autoRefreshInterval = ref(parseInt(localStorage.getItem('logsAutoRefreshInterval')) || 5)
  
  // Включено ли автообновление
  const autoRefreshEnabled = ref(JSON.parse(localStorage.getItem('logsAutoRefreshEnabled') || 'true'))

  function setAutoRefreshInterval(interval) {
    autoRefreshInterval.value = interval
    localStorage.setItem('logsAutoRefreshInterval', interval.toString())
  }

  function setAutoRefreshEnabled(enabled) {
    autoRefreshEnabled.value = enabled
    localStorage.setItem('logsAutoRefreshEnabled', JSON.stringify(enabled))
  }

  return {
    autoRefreshInterval,
    autoRefreshEnabled,
    setAutoRefreshInterval,
    setAutoRefreshEnabled
  }
})