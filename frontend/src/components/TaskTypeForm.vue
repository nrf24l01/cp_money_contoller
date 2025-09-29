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
    <div>
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
          <div class="sm:col-span-2">
            <label class="block text-gray-600 mb-1">Field Description</label>
            <input v-model="field.field_description" placeholder="field_description" class="border rounded px-2 py-1 w-full focus:outline-none focus:ring-2 focus:ring-blue-400" />
          </div>
        </div>
        <button type="button" @click="removeField(index)" class="mt-2 text-red-500 hover:underline" @click.prevent="confirmRemoval(index)">Remove Field</button>
      </div>
      <button type="button" @click="addField" class="text-blue-500 hover:underline">+ Add Field</button>
    </div>
    <div>
      <button type="submit" :disabled="creating" class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded px-4 py-2 disabled:opacity-50">
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
import { ref } from 'vue'
import api from '@/axios'
const emit = defineEmits({
  success: (msg) => typeof msg === 'string',
  error: (msg) => typeof msg === 'string'
})

const newTask = ref({
  task_name: '',
  task_type: '',
  task_template: [{ field_name: '', field_hint: '', field_description: '' }]
})
const creating = ref(false)

function addField() {
  newTask.value.task_template.push({ field_name: '', field_hint: '', field_description: '' })
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
    newTask.value = { task_name: '', task_type: '', task_template: [{ field_name: '', field_hint: '', field_description: '' }] }
  } catch (e) {
    console.error('Create failed', e)
    emit('error', 'Failed to create task type.')
  } finally {
    creating.value = false
  }
}
</script>
