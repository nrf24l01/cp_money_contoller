import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useTasksStore = defineStore('tasks', () => {
  // Интервал автообновления в секундах (по умолчанию 10 секунд)
  const autoRefreshInterval = ref(parseInt(localStorage.getItem('tasksAutoRefreshInterval')) || 10)
  
  // Включено ли автообновление
  const autoRefreshEnabled = ref(JSON.parse(localStorage.getItem('tasksAutoRefreshEnabled') || 'true'))

  function setAutoRefreshInterval(interval) {
    autoRefreshInterval.value = interval
    localStorage.setItem('tasksAutoRefreshInterval', interval.toString())
  }

  function setAutoRefreshEnabled(enabled) {
    autoRefreshEnabled.value = enabled
    localStorage.setItem('tasksAutoRefreshEnabled', JSON.stringify(enabled))
  }

  return {
    autoRefreshInterval,
    autoRefreshEnabled,
    setAutoRefreshInterval,
    setAutoRefreshEnabled
  }
})
