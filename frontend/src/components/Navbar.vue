<template>
    <nav class="flex justify-between items-center p-4 bg-gray-800 text-white">
            <router-link class="flex items-center hover:text-gray-300" :to="{ name: 'Index' }">
                <img src="/cp_logo.png" alt="Logo" class="h-8 w-8 mr-2" />
                <span class="text-lg font-bold">CP money controller</span>
            </router-link>
        <div class="flex space-x-4">
            <template v-if="authStore.isAuthenticated">
                <router-link
                    v-for="link in authenticatedLinks"
                    :key="link.name"
                    :to="link.to"
                    class="hover:text-gray-300"
                >
                    {{ link.name }}
                </router-link>
                <button @click="logout" class="hover:text-red-300">
                    <i class="pi pi-sign-out"></i>
                </button>
            </template>
            <template v-else>
                <router-link
                    v-for="link in guestLinks"
                    :key="link.name"
                    :to="link.to"
                    class="hover:underline"
                >
                    {{ link.name }}
                </router-link>
            </template>
        </div>
    </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const authenticatedLinks = ref([
    { name: 'Dashboard', to: '/' },
    { name: 'Tasks', to: '/tasks' },
    { name: 'Create Task', to: '/tasks/create' },
    { name: 'Task types', to: '/task_types' }
])

const guestLinks = ref([
    { name: 'Login', to: '/login' }
])

function logout() {
    authStore.logout()
    router.push('/login')
}
</script>