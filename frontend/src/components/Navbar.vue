<template>
    <nav class="flex justify-between items-center p-4 bg-gray-800 text-white">
        <div class="flex items-center">
            <img src="/cp_logo.png" alt="Logo" class="h-8 w-8 mr-2" />
            <span class="text-lg font-bold">CP money controller</span>
        </div>
        <div class="flex space-x-4">
            <template v-if="authStore.isAuthenticated">
                <router-link
                    v-for="link in authenticatedLinks"
                    :key="link.name"
                    :to="link.to"
                    class="hover:underline"
                >
                    {{ link.name }}
                </router-link>
                <button @click="logout" class="hover:underline text-red-300">Logout</button>
                <span class="text-sm text-gray-300 ml-2">{{ authStore.username }}</span>
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
    { name: 'Tasks', to: '/tasks' }
])

const guestLinks = ref([
    { name: 'Login', to: '/login' }
])

function logout() {
    authStore.logout()
    router.push('/login')
}
</script>