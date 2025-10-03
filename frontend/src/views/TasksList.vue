<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <div class="max-w-7xl mx-auto">
      <!-- Header with title and buttons -->
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-extrabold text-gray-800">Tasks List</h1>
        <div class="flex space-x-3">
          <router-link 
            to="/tasks/create" 
            class="bg-green-600 hover:bg-green-700 text-white font-semibold rounded px-4 py-2"
          >
            Create Task
          </router-link>
          <button 
            @click="showCreateModal = true" 
            class="bg-green-600 hover:bg-green-700 text-white font-semibold rounded px-4 py-2"
          >
            Quick Create
          </button>
          <button @click="fetchTasks" class="bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded px-4 py-2">
            Refresh
          </button>
        </div>
      </div>

      <!-- Alert notification -->
      <div v-if="alert.message" :class="['mb-4 p-4 rounded', alert.type === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800']">
        {{ alert.message }}
      </div>

      <!-- Filters: search, status, and sorting -->
      <div class="flex flex-col md:flex-row md:space-x-4 mb-4">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by UUID or Type..."
          class="border rounded px-3 py-2 flex-grow focus:outline-none focus:ring-2 focus:ring-blue-400 mb-2 md:mb-0"
        />
        <select v-model="statusFilter" class="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400 mb-2 md:mb-0">
          <option value="">All Statuses</option>
          <option v-for="status in statuses" :key="status" :value="status">
            {{ status }}
          </option>
        </select>
        <select v-model="sortOrder" class="border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400">
          <option value="desc">Newest First</option>
          <option value="asc">Oldest First</option>
        </select>
      </div>

      <!-- Tasks table -->
      <div v-if="loading" class="flex justify-center items-center py-10">
        <svg class="w-10 h-10 animate-spin text-blue-600" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
        </svg>
      </div>
      <div v-else>
        <table v-if="filteredTasks.length" class="min-w-full divide-y divide-gray-200 bg-white shadow rounded-lg">
          <thead class="bg-gray-100">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">UUID</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Success</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Last Update</th>
              <th class="px-4 py-2 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="task in filteredTasks" :key="task.uuid">
              <td class="px-4 py-2 text-gray-800 font-mono text-sm">{{ task.uuid.substring(0, 8) }}...</td>
              <td class="px-4 py-2 text-gray-700 capitalize">{{ task.type }}</td>
              <td class="px-4 py-2">
                <span :class="[
                  'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                  getStatusClass(task.status, task.payload_output)
                ]">
                  {{ task.status }}
                </span>
              </td>
              <td class="px-4 py-2">
                <span v-if="task.payload_output && typeof task.payload_output.success !== 'undefined'" :class="[
                  'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                  task.payload_output.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]">
                  {{ task.payload_output.success ? 'Success' : 'Failed' }}
                </span>
                <span v-else class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-800">
                  —
                </span>
              </td>
              <td class="px-4 py-2 text-gray-600 text-sm">{{ formatDate(task.last_update) }}</td>
              <td class="px-4 py-2 text-center space-x-2">
                <button 
                  @click="viewDetails(task.uuid)" 
                  class="text-blue-600 hover:text-blue-800 font-semibold"
                >
                  View
                </button>
                <button 
                  v-if="task.logs && task.logs.length > 0"
                  @click="viewLogs(task)" 
                  class="text-green-600 hover:text-green-800 font-semibold"
                >
                  Logs
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="text-center text-gray-500 py-6">No tasks found.</div>
      </div>
    </div>

    <!-- Logs Modal -->
    <TaskLogsModal
      :visible="showLogsModal"
      :logs="selectedTask?.logs || []"
      :task-uuid="selectedTask?.uuid || ''"
      @close="showLogsModal = false"
    />

    <!-- Create Task Modal -->
    <TaskCreateModal
      :visible="showCreateModal"
      @close="showCreateModal = false"
      @success="onTaskCreated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/axios'
import { useRouter } from 'vue-router'
import TaskLogsModal from '@/components/TaskLogsModal.vue'
import TaskCreateModal from '@/components/TaskCreateModal.vue'

const tasks = ref([])
const loading = ref(false)
const alert = ref({ message: '', type: '' })
const searchQuery = ref('')
const statusFilter = ref('')
const sortOrder = ref('desc')
const showLogsModal = ref(false)
const showCreateModal = ref(false)
const selectedTask = ref(null)
const router = useRouter()

// Fetch tasks from backend
async function fetchTasks() {
  loading.value = true
  try {
    const res = await api.get('/task')
    if (!res.data || res.data.length === 0) {
      console.log('No tasks found')
      alert.value = { message: 'No tasks found. Create a new task to get started!', type: 'info' }
    } else {
      tasks.value = res.data
      alert.value = { message: '', type: '' }
    }
  } catch (e) {
    console.error('Failed to fetch tasks', e)
    alert.value = { message: 'Failed to load tasks.', type: 'error' }
    // clear alert after timeout
    setTimeout(() => (alert.value = { message: '', type: '' }), 4000)
  } finally {
    loading.value = false
  }
}

// Derived list of statuses for filter dropdown
const statuses = computed(() => {
  const unique = new Set(tasks.value.map(t => t.status))
  return Array.from(unique)
})

// Filtered and sorted tasks based on search, status, and order
const filteredTasks = computed(() => {
  let filtered = tasks.value.filter(t => {
    const matchSearch =
      !searchQuery.value ||
      t.uuid.includes(searchQuery.value) ||
      t.type.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchStatus = !statusFilter.value || t.status === statusFilter.value
    return matchSearch && matchStatus
  })

  // Sort by last_update timestamp
  filtered.sort((a, b) => {
    if (sortOrder.value === 'desc') {
      return b.last_update - a.last_update
    } else {
      return a.last_update - b.last_update
    }
  })

  return filtered
})

// Navigate to detail view
function viewDetails(uuid) {
  router.push({ name: 'TaskDetails', params: { uuid } })
}

// Show logs modal
function viewLogs(task) {
  selectedTask.value = task
  showLogsModal.value = true
}

// Get status badge classes
function getStatusClass(status, payloadOutput = null) {
  switch (status?.toLowerCase()) {
    case 'done':
      // Для завершенных задач смотрим на успешность
      if (payloadOutput?.success === false) {
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

// Format Unix timestamp to readable string
function formatDate(ts) {
  return new Date(ts * 1000).toLocaleString()
}

// Handle task creation success
function onTaskCreated(taskUuid) {
  // Refresh tasks list
  fetchTasks()
  // Show success message
  alert.value = { message: `Task created successfully! UUID: ${taskUuid}`, type: 'success' }
  // Clear alert after timeout
  setTimeout(() => (alert.value = { message: '', type: '' }), 5000)
}

onMounted(fetchTasks)
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
