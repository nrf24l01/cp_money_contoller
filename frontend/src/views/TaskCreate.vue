<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <div class="max-w-4xl mx-auto">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div>
          <h1 class="text-3xl font-extrabold text-gray-800">Create New Task</h1>
          <p class="text-gray-600 mt-2">Fill in the form below to create a new task</p>
        </div>
        <router-link 
          to="/tasks" 
          class="bg-gray-600 hover:bg-gray-700 text-white font-semibold rounded px-4 py-2"
        >
          Back to Tasks
        </router-link>
      </div>

      <!-- Alert -->
      <div v-if="alert.message" :class="['mb-6 p-4 rounded-lg', alert.type === 'success' ? 'bg-green-100 text-green-800 border border-green-200' : 'bg-red-100 text-red-800 border border-red-200']">
        <div class="flex items-center">
          <svg v-if="alert.type === 'success'" class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <svg v-else class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
          </svg>
          {{ alert.message }}
        </div>
      </div>

      <!-- Main Form Card -->
      <div class="bg-white rounded-lg shadow-lg">
        <div class="p-8">
          <form @submit.prevent="createTask" class="space-y-8">
            
            <!-- Task Type Selection -->
            <div class="space-y-3">
              <label class="block text-lg font-semibold text-gray-900">
                Select Task Type
                <span class="text-red-500">*</span>
              </label>
              <select 
                v-model="selectedTaskType" 
                @change="onTaskTypeChange"
                required 
                class="w-full border border-gray-300 rounded-lg px-4 py-3 text-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">Choose a task type...</option>
                <option v-for="taskType in taskTypes" :key="taskType.uuid" :value="taskType">
                  {{ taskType.task_name }} ({{ taskType.task_type }})
                </option>
              </select>
              <p class="text-sm text-gray-500">Select the type of task you want to create</p>
            </div>

            <!-- Task Type Info Card -->
            <div v-if="selectedTaskType" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h3 class="font-semibold text-blue-900 mb-2">{{ selectedTaskType.task_name }}</h3>
              <p class="text-blue-800 text-sm">
                <strong>Type:</strong> {{ selectedTaskType.task_type }}
              </p>
              <p v-if="selectedTaskType.task_template && selectedTaskType.task_template.length" class="text-blue-800 text-sm mt-1">
                <strong>Required fields:</strong> {{ selectedTaskType.task_template.length }}
              </p>
            </div>

            <!-- Dynamic Form Fields -->
            <div v-if="selectedTaskType && selectedTaskType.task_template && selectedTaskType.task_template.length" class="space-y-6">
              <div class="border-t border-gray-200 pt-6">
                <h3 class="text-xl font-semibold text-gray-900 mb-4">Task Parameters</h3>
                <div class="grid gap-6">
                  <div 
                    v-for="field in selectedTaskType.task_template" 
                    :key="field.field_name"
                    class="space-y-2"
                  >
                    <label class="block text-base font-medium text-gray-700">
                      {{ field.field_hint || field.field_name }}
                      <span v-if="field.required !== false" class="text-red-500">*</span>
                      <span class="text-sm text-gray-500 font-normal ml-2">({{ field.field_type }})</span>
                    </label>
                    
                    <!-- String input -->
                    <input 
                      v-if="field.field_type === 'string'"
                      v-model="taskPayload[field.field_name]"
                      :placeholder="field.field_description || field.field_hint"
                      :required="field.required !== false"
                      type="text"
                      class="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                    
                    <!-- Integer input -->
                    <input 
                      v-else-if="field.field_type === 'int'"
                      v-model.number="taskPayload[field.field_name]"
                      :placeholder="field.field_description || field.field_hint"
                      :required="field.required !== false"
                      type="number"
                      step="1"
                      class="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                    
                    <!-- Float input -->
                    <input 
                      v-else-if="field.field_type === 'float'"
                      v-model.number="taskPayload[field.field_name]"
                      :placeholder="field.field_description || field.field_hint"
                      :required="field.required !== false"
                      type="number"
                      step="any"
                      class="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                    
                    <!-- Boolean input -->
                    <div v-else-if="field.field_type === 'bool'" class="flex items-center space-x-3">
                      <input 
                        v-model="taskPayload[field.field_name]"
                        type="checkbox"
                        class="h-4 w-4 rounded border-gray-300 focus:ring-2 focus:ring-blue-500"
                      />
                      <span class="text-base text-gray-700">{{ field.field_description || field.field_hint || 'Enable this option' }}</span>
                    </div>
                    
                    <!-- JSON input -->
                    <div v-else-if="field.field_type === 'json'" class="space-y-2">
                      <textarea 
                        v-model="taskPayload[field.field_name]"
                        :placeholder="field.field_description || field.field_hint || 'Enter valid JSON'"
                        :required="field.required !== false"
                        rows="6"
                        class="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
                        @blur="validateJSON(field.field_name)"
                      ></textarea>
                      <p v-if="jsonErrors[field.field_name]" class="text-red-500 text-sm flex items-center">
                        <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                        {{ jsonErrors[field.field_name] }}
                      </p>
                    </div>
                    
                    <!-- Array input -->
                    <div v-else-if="field.field_type === 'array'" class="space-y-2">
                      <textarea 
                        v-model="taskPayload[field.field_name]"
                        :placeholder="field.field_description || field.field_hint || 'Enter JSON array'"
                        :required="field.required !== false"
                        rows="4"
                        class="w-full border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
                        @blur="validateArray(field.field_name)"
                      ></textarea>
                      <p v-if="arrayErrors[field.field_name]" class="text-red-500 text-sm flex items-center">
                        <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                        </svg>
                        {{ arrayErrors[field.field_name] }}
                      </p>
                    </div>
                    
                    <p v-if="field.field_description" class="text-sm text-gray-500">
                      {{ field.field_description }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- No fields message -->
            <div v-else-if="selectedTaskType && (!selectedTaskType.task_template || selectedTaskType.task_template.length === 0)" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <p class="text-yellow-800">This task type doesn't require any additional parameters.</p>
            </div>

            <!-- JSON Preview -->
            <div v-if="selectedTaskType" class="space-y-3">
              <label class="block text-lg font-semibold text-gray-900">Payload Preview</label>
              <div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
                <pre class="text-sm text-gray-800 whitespace-pre-wrap overflow-x-auto">{{ JSON.stringify({ type: selectedTaskType.task_type, payload: taskPayload }, null, 2) }}</pre>
              </div>
              <p class="text-sm text-gray-500">This is how your task data will be sent to the system</p>
            </div>

            <!-- Submit Button -->
            <div class="border-t border-gray-200 pt-6">
              <div class="flex justify-end space-x-4">
                <router-link 
                  to="/tasks"
                  class="px-6 py-3 text-gray-700 bg-gray-100 border border-gray-300 rounded-lg hover:bg-gray-200 font-medium"
                >
                  Cancel
                </router-link>
                <button 
                  type="submit" 
                  :disabled="creating || !selectedTaskType"
                  class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium text-lg"
                >
                  <span v-if="!creating">Create Task</span>
                  <span v-else class="flex items-center">
                    <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
                    </svg>
                    Creating Task...
                  </span>
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/axios'

const router = useRouter()
const taskTypes = ref([])
const selectedTaskType = ref('')
const taskPayload = ref({})
const creating = ref(false)
const alert = ref({ message: '', type: '' })
const jsonErrors = ref({})
const arrayErrors = ref({})

// Load task types on component mount
onMounted(async () => {
  try {
    const response = await api.get('/task/type')
    taskTypes.value = response.data
  } catch (error) {
    console.error('Failed to load task types:', error)
    alert.value = { message: 'Failed to load task types. Please try refreshing the page.', type: 'error' }
  }
})

// Handle task type selection change
function onTaskTypeChange() {
  taskPayload.value = {}
  jsonErrors.value = {}
  arrayErrors.value = {}
  alert.value = { message: '', type: '' }
  
  if (selectedTaskType.value && selectedTaskType.value.task_template) {
    // Initialize payload with appropriate default values for each field type
    selectedTaskType.value.task_template.forEach(field => {
      switch (field.field_type) {
        case 'string':
          taskPayload.value[field.field_name] = ''
          break
        case 'int':
          taskPayload.value[field.field_name] = 0
          break
        case 'float':
          taskPayload.value[field.field_name] = 0.0
          break
        case 'bool':
          taskPayload.value[field.field_name] = false
          break
        case 'json':
          taskPayload.value[field.field_name] = '{\n  \n}'
          break
        case 'array':
          taskPayload.value[field.field_name] = '[\n  \n]'
          break
        default:
          taskPayload.value[field.field_name] = ''
      }
    })
  }
}

// Validate JSON input
function validateJSON(fieldName) {
  try {
    const value = taskPayload.value[fieldName]
    if (value && value.trim()) {
      JSON.parse(value)
      delete jsonErrors.value[fieldName]
    }
  } catch (error) {
    jsonErrors.value[fieldName] = 'Invalid JSON format'
  }
}

// Validate Array input
function validateArray(fieldName) {
  try {
    const value = taskPayload.value[fieldName]
    if (value && value.trim()) {
      const parsed = JSON.parse(value)
      if (!Array.isArray(parsed)) {
        arrayErrors.value[fieldName] = 'Value must be a valid JSON array'
      } else {
        delete arrayErrors.value[fieldName]
      }
    }
  } catch (error) {
    arrayErrors.value[fieldName] = 'Invalid JSON array format'
  }
}

// Create task
async function createTask() {
  // Validate JSON and Array fields before submission
  let hasErrors = false
  
  if (selectedTaskType.value && selectedTaskType.value.task_template) {
    selectedTaskType.value.task_template.forEach(field => {
      if (field.field_type === 'json') {
        validateJSON(field.field_name)
        if (jsonErrors.value[field.field_name]) hasErrors = true
      } else if (field.field_type === 'array') {
        validateArray(field.field_name)
        if (arrayErrors.value[field.field_name]) hasErrors = true
      }
    })
  }
  
  if (hasErrors) {
    alert.value = { message: 'Please fix validation errors before submitting', type: 'error' }
    return
  }
  
  creating.value = true
  alert.value = { message: '', type: '' }
  
  try {
    // Process payload to convert string JSON/Array values to actual objects
    const processedPayload = { ...taskPayload.value }
    
    if (selectedTaskType.value && selectedTaskType.value.task_template) {
      selectedTaskType.value.task_template.forEach(field => {
        const value = processedPayload[field.field_name]
        
        if (field.field_type === 'json' && typeof value === 'string' && value.trim()) {
          try {
            processedPayload[field.field_name] = JSON.parse(value)
          } catch (error) {
            // Keep as string if parsing fails
          }
        } else if (field.field_type === 'array' && typeof value === 'string' && value.trim()) {
          try {
            processedPayload[field.field_name] = JSON.parse(value)
          } catch (error) {
            // Keep as string if parsing fails
          }
        }
      })
    }
    
    const requestData = {
      type: selectedTaskType.value.task_type,
      payload: processedPayload
    }
    
    const response = await api.post('/task', requestData)
    
    alert.value = { message: `Task created successfully! UUID: ${response.data.uuid}`, type: 'success' }
    
    // Redirect to task details after a short delay
    setTimeout(() => {
      router.push({ name: 'TaskDetails', params: { uuid: response.data.uuid } })
    }, 2000)
    
  } catch (error) {
    console.error('Failed to create task:', error)
    let errorMessage = 'Failed to create task'
    
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.response?.data?.error) {
      errorMessage = error.response.data.error
    }
    
    alert.value = { message: errorMessage, type: 'error' }
  } finally {
    creating.value = false
  }
}
</script>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
