<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <div class="max-w-7xl mx-auto">
      <!-- Header with title and create button -->
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-extrabold text-gray-800">Task Types Management</h1>
        <button @click="showModal = true" class="bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded px-4 py-2">
          Create New Task Type
        </button>
      </div>
      <!-- Alert -->
      <div v-if="alert.message" :class="['mb-4 p-4 rounded', alert.type === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800']">
        {{ alert.message }}
      </div>
      <!-- Search filter -->
      <div class="mb-4">
        <input v-model="searchQuery" type="text" placeholder="Search task types..." class="border rounded px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
      </div>
      <div class="grid grid-cols-1 gap-6">
        <!-- Existing Types -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-semibold mb-4 text-gray-700">Existing Task Types</h2>
          <div v-if="loading" class="flex justify-center items-center py-10">
            <svg class="w-10 h-10 animate-spin text-blue-600" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
            </svg>
          </div>
          <div v-else>
            <table v-if="filteredTaskTypes.length" class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-100">
                <tr>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type ID</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Template Fields</th>
                  <th class="px-4 py-2 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="type in filteredTaskTypes" :key="type.uuid">
                  <td class="px-4 py-2 text-gray-800 font-medium">{{ type.task_name }}</td>
                  <td class="px-4 py-2 text-gray-600">{{ type.task_type }}</td>
                  <td class="px-4 py-2">
                    <ul class="list-disc list-inside">
                      <li v-for="field in type.task_template" :key="field.field_name">
                        <span class="font-semibold">{{ field.field_name }}</span> - {{ field.field_hint }}
                      </li>
                    </ul>
                  </td>
                  <td class="px-4 py-2 text-center">
                    <button @click="deleteTaskType(type.uuid)" :disabled="deleting[type.uuid]" class="text-red-600 hover:text-red-800 disabled:opacity-50">
                      <span v-if="!deleting[type.uuid]">Delete</span>
                      <svg v-else class="w-5 h-5 animate-spin inline-block" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
                      </svg>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="text-center text-gray-500 py-6">No task types found.</div>
          </div>
        </div>
       </div>
      <!-- Modal for creating task type -->
      <transition name="fade">
        <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div class="bg-white rounded-lg w-full max-w-2xl p-6">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-semibold">New Task Type</h3>
              <button @click="showModal = false" class="text-gray-500 hover:text-gray-700 text-2xl leading-none">&times;</button>
            </div>
            <TaskTypeForm @success="handleSuccess" @error="handleError" />
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import api from '@/axios'
import TaskTypeForm from '@/components/TaskTypeForm.vue'
 
const taskTypes = ref([])
const alert = ref({ message: '', type: '' })
const showModal = ref(false)
const loading = ref(false)
const creating = ref(false)
const deleting = ref({})

async function fetchTaskTypes() {
  loading.value = true
  try {
    const response = await api.get('/task/type')
    taskTypes.value = response.data
  } catch (e) {
    console.error('Failed to fetch task types', e)
    alert.value = { message: 'Failed to load task types.', type: 'error' }
  } finally {
    loading.value = false
  }
}

const newTask = ref({
  task_name: '',
  task_type: '',
  task_template: [
    { field_name: '', field_hint: '', field_description: '' }
  ]
})

function addField() {
  newTask.value.task_template.push({ field_name: '', field_hint: '', field_description: '' })
}

function removeField(index) {
  newTask.value.task_template.splice(index, 1)
}

async function createTaskType() {
  creating.value = true
  try {
    await api.post('/task/type', newTask.value)
    alert.value = { message: 'Task type created successfully!', type: 'success' }
    newTask.value = { task_name: '', task_type: '', task_template: [ { field_name: '', field_hint: '', field_description: '' } ] }
    fetchTaskTypes()
  } catch (e) {
    console.error('Failed to create task type', e)
    alert.value = { message: 'Failed to create task type.', type: 'error' }
  } finally {
    creating.value = false
  }
}

// delete existing task type
async function deleteTaskType(uuid) {
  if (!confirm('Are you sure you want to delete this task type?')) return
  deleting.value = { ...deleting.value, [uuid]: true }
  try {
    await api.delete(`/task/type/${uuid}`)
    alert.value = { message: 'Task type deleted.', type: 'success' }
    fetchTaskTypes()
  } catch (e) {
    console.error('Delete failed', e)
    alert.value = { message: 'Failed to delete task type.', type: 'error' }
  } finally {
    deleting.value = { ...deleting.value, [uuid]: false }
  }
}

// clear alert after timeout
watch(alert, (val) => {
  if (val.message) setTimeout(() => alert.value = { message: '', type: '' }, 4000)
})

onMounted(fetchTaskTypes)
const searchQuery = ref('')
const filteredTaskTypes = computed(() => {
  if (!searchQuery.value) return taskTypes.value
  return taskTypes.value.filter(type =>
    type.task_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    type.task_type.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Handlers for modal form events
function handleSuccess(msg) {
  alert.value = { message: msg, type: 'success' }
  showModal.value = false
  fetchTaskTypes()
}

function handleError(msg) {
  alert.value = { message: msg, type: 'error' }
}
</script>

<style scoped>
/* You can add component-specific styles here */
</style>