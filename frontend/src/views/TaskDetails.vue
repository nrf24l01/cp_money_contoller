<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <div class="max-w-6xl mx-auto">
      <!-- Header with back button -->
      <div class="flex items-center mb-6">
        <button @click="$router.back()" class="mr-4 text-blue-600 hover:text-blue-800">
          ← Back to Tasks
        </button>
        <h1 class="text-2xl font-extrabold text-gray-800">Task Details</h1>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="flex justify-center items-center py-20">
        <svg class="w-12 h-12 animate-spin text-blue-600" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
        </svg>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {{ error }}
      </div>

      <!-- Task details -->
      <div v-else-if="task" class="space-y-6">
        <!-- Basic Info Card -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-4 text-gray-700">Basic Information</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-500">UUID</label>
              <p class="mt-1 text-sm font-mono text-gray-900">{{ task.uuid }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Type</label>
              <p class="mt-1 text-sm text-gray-900 capitalize">{{ task.type }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Status</label>
              <span :class="[
                'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                getStatusClass(task.status)
              ]">
                {{ task.status }}
              </span>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Success</label>
              <span v-if="task.payload_output && typeof task.payload_output.success !== 'undefined'" :class="[
                'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                task.payload_output.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              ]">
                {{ task.payload_output.success ? 'Success' : 'Failed' }}
              </span>
              <span v-else class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-800">
                Unknown
              </span>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500">Last Update</label>
              <p class="mt-1 text-sm text-gray-900">{{ formatDate(task.last_update) }}</p>
            </div>
          </div>
        </div>

        <!-- Input Payload Card -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-4 text-gray-700">Input Payload</h2>
          <div class="bg-gray-50 rounded p-4 overflow-auto">
            <pre class="text-sm text-gray-800">{{ JSON.stringify(task.payload_input, null, 2) }}</pre>
          </div>
        </div>

        <!-- Output Payload Card -->
        <div v-if="task.payload_output" class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-4 text-gray-700">Output Payload</h2>
          <div class="bg-gray-50 rounded p-4 overflow-auto">
            <pre class="text-sm text-gray-800">{{ JSON.stringify(task.payload_output, null, 2) }}</pre>
          </div>
        </div>

        <!-- Logs Card -->
        <div class="bg-white shadow rounded-lg p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-lg font-semibold text-gray-700">Logs</h2>
            <button
              @click="showLogsModal = true"
              class="bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded px-4 py-2 text-sm"
            >
              View All Logs
            </button>
          </div>
          
          <!-- Preview of first few logs -->
          <div v-if="task.logs && task.logs.length > 0" class="space-y-2">
            <div
              v-for="(log, index) in task.logs.slice(0, 5)"
              :key="index"
              :class="[
                'p-3 rounded font-mono text-sm border-l-4',
                getLogLevelClass(log)
              ]"
            >
              {{ log }}
            </div>
            <div v-if="task.logs.length > 5" class="text-center py-2">
              <span class="text-gray-500 text-sm">
                ... and {{ task.logs.length - 5 }} more entries
              </span>
            </div>
          </div>
          <div v-else class="text-gray-500 text-center py-4">
            No logs available
          </div>
        </div>
      </div>
    </div>

    <!-- Logs Modal -->
    <TaskLogsModal
      :visible="showLogsModal"
      :logs="task?.logs || []"
      :task-uuid="task?.uuid || ''"
      @close="showLogsModal = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/axios'
import TaskLogsModal from '@/components/TaskLogsModal.vue'

const route = useRoute()
const task = ref(null)
const loading = ref(false)
const error = ref('')
const showLogsModal = ref(false)

// Fetch task details
async function fetchTaskDetails() {
  loading.value = true
  error.value = ''
  
  try {
    const response = await api.get(`/task/${route.params.uuid}`)
    task.value = response.data
  } catch (e) {
    console.error('Failed to fetch task details', e)
    error.value = e.response?.data?.message || 'Failed to load task details'
  } finally {
    loading.value = false
  }
}

// Get status badge classes
function getStatusClass(status) {
  switch (status?.toLowerCase()) {
    case 'done':
      // Для завершенных задач смотрим на успешность
      if (task.value?.payload_output?.success === false) {
        return 'bg-orange-100 text-orange-800'
      }
      return 'bg-green-100 text-green-800'
    case 'running':
      return 'bg-blue-100 text-blue-800'
    case 'pending':
      return 'bg-yellow-100 text-yellow-800'
    case 'failed':
    case 'error':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// Extract log level from log string for styling
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

// Format Unix timestamp to readable string
function formatDate(timestamp) {
  return new Date(timestamp * 1000).toLocaleString()
}

onMounted(fetchTaskDetails)
</script>

<style scoped>
/* Component-specific styles */
</style>
