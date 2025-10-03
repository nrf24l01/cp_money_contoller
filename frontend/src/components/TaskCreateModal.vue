<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] flex flex-col">
      <div class="p-6 overflow-y-auto flex-1">
        <!-- Header -->
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800">Create New Task</h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600 text-2xl">
            &times;
          </button>
        </div>

        <!-- Alert -->
        <div v-if="alert.message" :class="['mb-4 p-3 rounded', alert.type === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800']">
          {{ alert.message }}
        </div>

        <!-- Form -->
        <form @submit.prevent="createTask" class="space-y-4">
          <!-- Task Type Selection -->
          <div>
            <label class="block text-gray-700 font-medium mb-2">Task Type</label>
            <select 
              v-model="selectedTaskType" 
              @change="onTaskTypeChange"
              required 
              class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select a task type...</option>
              <option v-for="taskType in taskTypes" :key="taskType.uuid" :value="taskType">
                {{ taskType.task_name }}
              </option>
            </select>
          </div>

          <!-- Dynamic Form Fields -->
          <div v-if="selectedTaskType && selectedTaskType.task_template" class="space-y-4">
            <h3 class="text-lg font-medium text-gray-800">Task Parameters</h3>
            <div 
              v-for="field in selectedTaskType.task_template" 
              :key="field.field_name"
              class="space-y-2"
            >
              <label class="block text-gray-700 font-medium">
                {{ field.field_hint || field.field_name }}
                <span v-if="field.required !== false" class="text-red-500">*</span>
                <span class="text-sm text-gray-500 font-normal">({{ field.field_type }})</span>
              </label>
              
              <!-- String input -->
              <input 
                v-if="field.field_type === 'string'"
                v-model="taskPayload[field.field_name]"
                :placeholder="field.field_description || field.field_hint"
                :required="field.required !== false"
                type="text"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              
              <!-- Integer input -->
              <input 
                v-else-if="field.field_type === 'int'"
                v-model.number="taskPayload[field.field_name]"
                :placeholder="field.field_description || field.field_hint"
                :required="field.required !== false"
                type="number"
                step="1"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              
              <!-- Float input -->
              <input 
                v-else-if="field.field_type === 'float'"
                v-model.number="taskPayload[field.field_name]"
                :placeholder="field.field_description || field.field_hint"
                :required="field.required !== false"
                type="number"
                step="any"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              
              <!-- Boolean input -->
              <label v-else-if="field.field_type === 'bool'" class="flex items-center space-x-2">
                <input 
                  v-model="taskPayload[field.field_name]"
                  type="checkbox"
                  class="rounded border-gray-300 focus:ring-2 focus:ring-blue-500"
                />
                <span class="text-sm">{{ field.field_description || field.field_hint || 'Enable option' }}</span>
              </label>
              
              <!-- JSON input -->
              <div v-else-if="field.field_type === 'json'" class="space-y-2">
                <textarea 
                  v-model="taskPayload[field.field_name]"
                  :placeholder="field.field_description || field.field_hint || 'Enter valid JSON'"
                  :required="field.required !== false"
                  rows="4"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                  @blur="validateJSON(field.field_name)"
                ></textarea>
                <p v-if="jsonErrors[field.field_name]" class="text-red-500 text-sm">{{ jsonErrors[field.field_name] }}</p>
              </div>
              
              <!-- Array input -->
              <div v-else-if="field.field_type === 'array'" class="space-y-2">
                <textarea 
                  v-model="taskPayload[field.field_name]"
                  :placeholder="field.field_description || field.field_hint || 'Enter JSON array'"
                  :required="field.required !== false"
                  rows="3"
                  class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                  @blur="validateArray(field.field_name)"
                ></textarea>
                <p v-if="arrayErrors[field.field_name]" class="text-red-500 text-sm">{{ arrayErrors[field.field_name] }}</p>
              </div>
              
              <p v-if="field.field_description" class="text-sm text-gray-500">
                {{ field.field_description }}
              </p>
            </div>
          </div>

          <!-- JSON Preview -->
          <div v-if="selectedTaskType" class="space-y-2">
            <label class="block text-gray-700 font-medium">JSON Payload Preview</label>
            <pre class="bg-gray-100 border rounded-lg p-3 text-sm overflow-x-auto">{{ JSON.stringify(taskPayload, null, 2) }}</pre>
          </div>

          <!-- Buttons -->
          <div class="flex justify-end space-x-3 pt-4">
            <button 
              type="button" 
              @click="$emit('close')"
              class="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50"
            >
              Cancel
            </button>
            <button 
              type="submit" 
              :disabled="creating || !selectedTaskType"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="!creating">Create Task</span>
              <span v-else class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
                </svg>
                Creating...
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '@/axios'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'success'])

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
    alert.value = { message: 'Failed to load task types', type: 'error' }
  }
})

// Watch for modal visibility changes to clear form
watch(() => props.visible, (newVal) => {
  if (newVal) {
    // Reset form when modal opens
    selectedTaskType.value = ''
    taskPayload.value = {}
    jsonErrors.value = {}
    arrayErrors.value = {}
    alert.value = { message: '', type: '' }
  }
})

// Handle task type selection change
function onTaskTypeChange() {
  taskPayload.value = {}
  jsonErrors.value = {}
  arrayErrors.value = {}
  
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
          taskPayload.value[field.field_name] = '{}'
          break
        case 'array':
          taskPayload.value[field.field_name] = '[]'
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
    
    alert.value = { message: 'Task created successfully!', type: 'success' }
    
    // Emit success event with task UUID
    emit('success', response.data.uuid)
    
    // Close modal after a short delay
    setTimeout(() => {
      emit('close')
    }, 1500)
    
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
/* Custom scrollbar for the modal */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
