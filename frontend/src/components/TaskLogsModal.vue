<template>
  <transition name="fade">
    <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-lg w-full max-w-4xl max-h-[90vh] flex flex-col shadow-xl">
        <!-- Header -->
        <div class="flex justify-between items-center p-6 border-b">
          <h3 class="text-lg font-semibold">Task Logs - {{ taskUuid }}</h3>
          <button @click="close" class="text-gray-500 hover:text-gray-700 text-2xl leading-none">&times;</button>
        </div>
        
        <!-- Controls -->
        <div class="p-4 border-b bg-gray-50">
          <div class="flex flex-col space-y-3">
            <!-- Filter controls -->
            <div class="flex flex-col md:flex-row md:space-x-4 space-y-2 md:space-y-0">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search logs..."
                class="border rounded px-3 py-2 flex-grow focus:outline-none focus:ring-2 focus:ring-blue-400"
              />
              <select v-model="logLevelFilter" class="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400">
                <option value="">All Levels</option>
                <option value="INFO">INFO</option>
                <option value="ERROR">ERROR</option>
                <option value="DEBUG">DEBUG</option>
                <option value="WARNING">WARNING</option>
              </select>
              <select v-model="sortOrder" class="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400">
                <option value="asc">Oldest First</option>
                <option value="desc">Newest First</option>
              </select>
            </div>
            
            <!-- Auto-refresh controls -->
            <div class="flex flex-col sm:flex-row sm:items-center sm:space-x-4 space-y-2 sm:space-y-0 text-sm">
              <div class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="autoRefresh"
                  :checked="logsStore.autoRefreshEnabled"
                  @change="toggleAutoRefresh"
                  class="rounded focus:ring-2 focus:ring-blue-400"
                />
                <label for="autoRefresh" class="text-gray-700">Auto-refresh logs</label>
              </div>
              
              <div v-if="logsStore.autoRefreshEnabled" class="flex items-center space-x-2">
                <label for="refreshInterval" class="text-gray-600">Interval:</label>
                <select
                  id="refreshInterval"
                  :value="logsStore.autoRefreshInterval"
                  @change="updateRefreshInterval(parseInt($event.target.value))"
                  class="border rounded px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
                >
                  <option value="1">1 sec</option>
                  <option value="2">2 sec</option>
                  <option value="5">5 sec</option>
                  <option value="10">10 sec</option>
                  <option value="30">30 sec</option>
                  <option value="60">1 min</option>
                </select>
              </div>
              
              <button
                @click="fetchLogs"
                :disabled="loading"
                class="px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ loading ? 'Loading...' : 'Refresh Now' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Logs content -->
        <div class="flex-1 overflow-auto p-4">
          <div v-if="filteredLogs.length === 0" class="text-center text-gray-500 py-8">
            No logs found matching your criteria.
          </div>
          <div v-else class="space-y-1">
            <div
              v-for="(log, index) in filteredLogs"
              :key="index"
              :class="[
                'p-3 rounded font-mono text-sm border-l-4',
                getLogLevelClass(log)
              ]"
            >
              <div class="flex items-start space-x-2">
                <span class="text-xs text-gray-500 min-w-0 flex-shrink-0">{{ index + 1 }}</span>
                <span class="break-all">{{ log }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer with stats -->
        <div class="p-4 border-t bg-gray-50 text-sm text-gray-600">
          Showing {{ filteredLogs.length }} of {{ logs.length }} log entries
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useLogsStore } from '@/stores/logsStore'
import api from '@/axios'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  taskUuid: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const logsStore = useLogsStore()
const logs = ref([])
const loading = ref(false)
const error = ref('')

const searchQuery = ref('')
const logLevelFilter = ref('')
const sortOrder = ref('asc')

let refreshInterval = null

// Функция для получения логов с сервера
async function fetchLogs() {
  if (!props.taskUuid) return
  
  try {
    loading.value = true
    error.value = ''
    
    const response = await api.get(`/task/${props.taskUuid}/logs`)
    
    logs.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch logs:', err)
    error.value = err.response?.data?.message || 'Failed to fetch logs'
  } finally {
    loading.value = false
  }
}

// Функция для настройки автообновления
function setupAutoRefresh() {
  clearInterval(refreshInterval)
  
  if (logsStore.autoRefreshEnabled && props.visible && props.taskUuid) {
    refreshInterval = setInterval(() => {
      fetchLogs()
    }, logsStore.autoRefreshInterval * 1000)
  }
}

// Обработчики изменения настроек автообновления
function toggleAutoRefresh() {
  logsStore.setAutoRefreshEnabled(!logsStore.autoRefreshEnabled)
  setupAutoRefresh()
}

function updateRefreshInterval(interval) {
  logsStore.setAutoRefreshInterval(interval)
  setupAutoRefresh()
}

// Extract log level from log string
function extractLogLevel(log) {
  const match = log.match(/\[(INFO|ERROR|DEBUG|WARNING)\]/)
  return match ? match[1] : 'UNKNOWN'
}

// Get CSS classes based on log level
function getLogLevelClass(log) {
  const level = extractLogLevel(log)
  switch (level) {
    case 'ERROR':
      return 'bg-red-50 border-red-400 text-red-800'
    case 'WARNING':
      return 'bg-yellow-50 border-yellow-400 text-yellow-800'
    case 'INFO':
      return 'bg-blue-50 border-blue-400 text-blue-800'
    case 'DEBUG':
      return 'bg-gray-50 border-gray-400 text-gray-800'
    default:
      return 'bg-gray-50 border-gray-300 text-gray-700'
  }
}

// Filtered and sorted logs
const filteredLogs = computed(() => {
  let filtered = logs.value

  // Filter by search query
  if (searchQuery.value) {
    filtered = filtered.filter(log =>
      log.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  // Filter by log level
  if (logLevelFilter.value) {
    filtered = filtered.filter(log =>
      extractLogLevel(log) === logLevelFilter.value
    )
  }

  // Sort logs
  if (sortOrder.value === 'desc') {
    filtered = [...filtered].reverse()
  }

  return filtered
})

function close() {
  emit('close')
}

// Watchers для автоматического получения логов и настройки автообновления
watch(() => props.visible, (newValue) => {
  if (newValue && props.taskUuid) {
    // Сбрасываем фильтры при открытии модального окна
    searchQuery.value = ''
    logLevelFilter.value = ''
    sortOrder.value = 'asc'
    
    // Получаем логи
    fetchLogs()
    
    // Настраиваем автообновление
    setupAutoRefresh()
  } else {
    // Останавливаем автообновление при закрытии модального окна
    clearInterval(refreshInterval)
  }
})

watch(() => props.taskUuid, (newUuid) => {
  if (newUuid && props.visible) {
    fetchLogs()
    setupAutoRefresh()
  }
})

// Отслеживание изменений настроек автообновления
watch(() => logsStore.autoRefreshEnabled, () => {
  setupAutoRefresh()
})

watch(() => logsStore.autoRefreshInterval, () => {
  setupAutoRefresh()
})

// Lifecycle hooks
onMounted(() => {
  if (props.visible && props.taskUuid) {
    fetchLogs()
    setupAutoRefresh()
  }
})

onUnmounted(() => {
  clearInterval(refreshInterval)
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
