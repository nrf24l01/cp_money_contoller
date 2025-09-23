<template>
    <nav class="flex justify-between items-center p-4 bg-gray-800 text-white">
        <div class="flex items-center">
            <img src="/cp_logo.png" alt="Logo" class="h-8 w-8 mr-2" />
            <span class="text-lg font-bold">CP money controller</span>
        </div>
        <div class="flex space-x-4">
            <router-link
                v-for="link in links"
                :key="link.name"
                :to="link.to"
                class="hover:underline"
            >
                {{ link.name }}
            </router-link>
        </div>
    </nav>
</template>

<script setup>
import { ref } from 'vue'

const links = ref([
    { name: 'Home', to: '/' },
    { name: 'About', to: '/about' },
    { name: 'Contact', to: '/contact' }
])
</script>