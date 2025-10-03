<template>
  <form @submit.prevent="submit" class="space-y-4">
    <div>
      <label class="block text-gray-600 mb-1">Readable Name</label>
      <input v-model="newTask.task_name" required minlength="3" placeholder="e.g. Email Notification" class="border rounded px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
    </div>
    <div>
      <label class="block text-gray-600 mb-1">Internal Type ID</label>
      <input v-model="newTask.task_type" required pattern="^[a-z0-9_]+$" title="Lowercase, numbers and underscores only" placeholder="e.g. email_notification" class="border rounded px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
    </div>
    
    <!-- Toggle between Form and JSON modes -->
    <div class="flex gap-2 border-b">
      <button type="button" @click="editMode = 'form'" :class="['px-4 py-2 font-semibold', editMode === 'form' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-500']">
        Form Editor
      </button>
      <button type="button" @click="editMode = 'json'" :class="['px-4 py-2 font-semibold', editMode === 'json' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-500']">
        JSON Editor
      </button>
    </div>

    <!-- Form Mode -->
    <div v-if="editMode === 'form'">
      <label class="block text-gray-600 mb-2">Template Fields</label>
      <div v-for="(field, index) in newTask.task_template" :key="index" class="border rounded p-4 mb-3">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-gray-600 mb-1">Field Name</label>
            <input v-model="field.field_name" required placeholder="field_name" class="border rounded px-2 py-1 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
          <div>
            <label class="block text-gray-600 mb-1">Field Hint</label>
            <input v-model="field.field_hint" placeholder="field_hint" class="border rounded px-2 py-1 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
          <div>
            <label class="block text-gray-600 mb-1">Field Type</label>
            <select v-model="field.field_type" required class="border rounded px-2 py-1 w-full focus:outline-none focus:ring-2 focus:ring-blue-400">
              <option value="string">String</option>
              <option value="int">Integer</option>
              <option value="float">Float</option>
              <option value="bool">Boolean</option>
              <option value="json">JSON</option>
              <option value="array">Array</option>
            </select>
          </div>
          <div>
            <label class="block text-gray-600 mb-1">Required</label>
            <label class="flex items-center">
              <input v-model="field.required" type="checkbox" class="mr-2" />
              <span class="text-sm">Field is required</span>
            </label>
          </div>
          <div class="sm:col-span-2">
            <label class="block text-gray-600 mb-1">Field Description</label>
            <input v-model="field.field_description" placeholder="field_description" class="border rounded px-2 py-1 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
        </div>
        <button type="button" @click="removeField(index)" class="mt-2 text-red-500 hover:underline" @click.prevent="confirmRemoval(index)">Remove Field</button>
      </div>
      <button type="button" @click="addField" class="text-blue-500 hover:underline">+ Add Field</button>
    </div>

    <!-- JSON Mode -->
    <div v-if="editMode === 'json'" class="space-y-4">
      <div>
        <label class="block text-gray-600 mb-2">Template JSON</label>
        <textarea 
          v-model="jsonText" 
          @input="onJsonChange"
          rows="15" 
          class="border rounded px-3 py-2 w-full font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          :class="{ 'border-red-500': jsonError }"
          placeholder='[{"field_name": "example", "field_hint": "hint", ...}]'
        ></textarea>
        <p v-if="jsonError" class="text-red-500 text-sm mt-1">{{ jsonError }}</p>
      </div>
      
      <!-- JSON Format Help -->
      <div class="bg-blue-50 border border-blue-200 rounded p-4">
        <h3 class="font-semibold text-blue-900 mb-2">📘 JSON Format Guide</h3>
        <p class="text-sm text-gray-700 mb-2">Template should be an array of field objects. Each field object contains:</p>
        <ul class="text-sm text-gray-700 space-y-1 list-disc list-inside">
          <li><code class="bg-white px-1 rounded">field_name</code> (string, required) - Unique field identifier</li>
          <li><code class="bg-white px-1 rounded">field_hint</code> (string) - Short hint for the field</li>
          <li><code class="bg-white px-1 rounded">field_description</code> (string) - Detailed field description</li>
          <li><code class="bg-white px-1 rounded">field_type</code> (string, required) - One of: "string", "int", "float", "bool", "json", "array"</li>
          <li><code class="bg-white px-1 rounded">required</code> (boolean) - Whether field is required (default: false)</li>
        </ul>
        <details class="mt-3">
          <summary class="cursor-pointer text-blue-700 hover:text-blue-900 font-semibold">Show Example</summary>
          <pre class="bg-white rounded p-3 mt-2 text-xs overflow-x-auto"><code>[
  {
    "field_name": "email",
    "field_hint": "User email address",
    "field_description": "Email where notification will be sent",
    "field_type": "string",
    "required": true
  },
  {
    "field_name": "priority",
    "field_hint": "Task priority level",
    "field_description": "Priority from 1 (low) to 10 (high)",
    "field_type": "int",
    "required": false
  },
  {
    "field_name": "metadata",
    "field_hint": "Additional data",
    "field_description": "Optional metadata in JSON format",
    "field_type": "json",
    "required": false
  }
]</code></pre>
        </details>
      </div>
    </div>

    <div>
      <button type="submit" :disabled="creating || (editMode === 'json' && jsonError)" class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded px-4 py-2 disabled:opacity-50">
        <template v-if="!creating">Create Task Type</template>
        <svg v-else class="w-5 h-5 animate-spin inline-block" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
        </svg>
      </button>
    </div>
  </form>
</template>

<script setup>
import { ref, watch } from 'vue'
import api from '@/axios'

const emit = defineEmits({
  success: (msg) => typeof msg === 'string',
  error: (msg) => typeof msg === 'string'
})

const newTask = ref({
  task_name: '',
  task_type: '',
  task_template: [{ field_name: '', field_hint: '', field_description: '', field_type: 'string', required: true }]
})
const creating = ref(false)
const editMode = ref('form') // 'form' or 'json'
const jsonText = ref('')
const jsonError = ref('')

// Initialize JSON text from task_template
watch(() => newTask.value.task_template, (newTemplate) => {
  if (editMode.value === 'form') {
    try {
      jsonText.value = JSON.stringify(newTemplate, null, 2)
      jsonError.value = ''
    } catch (e) {
      jsonError.value = 'Error serializing template'
    }
  }
}, { deep: true, immediate: true })

// Switch to JSON mode - sync form to JSON
watch(editMode, (newMode) => {
  if (newMode === 'json') {
    try {
      jsonText.value = JSON.stringify(newTask.value.task_template, null, 2)
      jsonError.value = ''
    } catch (e) {
      jsonError.value = 'Error serializing template'
    }
  }
})

// Handle JSON input changes
function onJsonChange() {
  try {
    const parsed = JSON.parse(jsonText.value)
    
    // Validate that it's an array
    if (!Array.isArray(parsed)) {
      jsonError.value = 'Template must be an array of field objects'
      return
    }
    
    // Validate each field
    for (let i = 0; i < parsed.length; i++) {
      const field = parsed[i]
      if (!field.field_name || typeof field.field_name !== 'string') {
        jsonError.value = `Field ${i + 1}: 'field_name' is required and must be a string`
        return
      }
      if (!field.field_type || typeof field.field_type !== 'string') {
        jsonError.value = `Field ${i + 1}: 'field_type' is required and must be a string`
        return
      }
      const validTypes = ['string', 'int', 'float', 'bool', 'json', 'array']
      if (!validTypes.includes(field.field_type)) {
        jsonError.value = `Field ${i + 1}: 'field_type' must be one of: ${validTypes.join(', ')}`
        return
      }
      
      // Ensure required fields have default values
      if (!field.hasOwnProperty('field_hint')) field.field_hint = ''
      if (!field.hasOwnProperty('field_description')) field.field_description = ''
      if (!field.hasOwnProperty('required')) field.required = false
    }
    
    // Update the form data
    newTask.value.task_template = parsed
    jsonError.value = ''
  } catch (e) {
    jsonError.value = 'Invalid JSON: ' + e.message
  }
}

function addField() {
  newTask.value.task_template.push({ 
    field_name: '', 
    field_hint: '', 
    field_description: '', 
    field_type: 'string', 
    required: true 
  })
}

function removeField(index) {
  newTask.value.task_template.splice(index, 1)
}

function confirmRemoval(index) {
  removeField(index)
}

async function submit() {
  creating.value = true
  try {
    await api.post('/task/type', newTask.value)
    emit('success', 'Task type created successfully!')
    newTask.value = { 
      task_name: '', 
      task_type: '', 
      task_template: [{ 
        field_name: '', 
        field_hint: '', 
        field_description: '', 
        field_type: 'string', 
        required: true 
      }] 
    }
    editMode.value = 'form'
  } catch (e) {
    console.error('Create failed', e)
    emit('error', 'Failed to create task type.')
  } finally {
    creating.value = false
  }
}
</script>
