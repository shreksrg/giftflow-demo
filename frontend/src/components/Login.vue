<script setup>
import { ref } from 'vue'
import axios from 'axios'

const emit = defineEmits(['login'])
const loading = ref(false)
const error = ref('')
const heroUrl = import.meta.env.VITE_LOGIN_HERO || '/login-hero.jpg'

const mockUsers = [
  { username: 'dept_admin', role: 'DEPT_ADMIN', label: 'Dept Admin (Applicant)' },
  { username: 'dept_head', role: 'DEPT_HEAD', label: 'Dept Head (Approver)' },
  { username: 'cpro_admin', role: 'CPRO_ADMIN', label: 'CPRO Admin (VIP Approver)' }
]

const loginAs = async (mockUser) => {
  loading.value = true
  error.value = ''
  try {
    // In a real app, we would authenticate here.
    // For mock, we just verify we can talk to backend with this user.
    // We set a global axios interceptor to inject the user header
    axios.defaults.headers.common['X-Mock-User'] = mockUser.username
    
    const response = await axios.get('/api/me')
    emit('login', response.data)
  } catch (err) {
    console.error(err)
    error.value = 'Failed to login. Is backend running?'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="w-full max-w-4xl mx-auto bg-white rounded-lg shadow overflow-hidden">
    <div class="grid grid-cols-1 md:grid-cols-2">
      <div class="relative h-40 md:h-auto">
        <img :src="heroUrl" alt="Login Visual" class="absolute inset-0 w-full h-full object-cover">
        <div class="absolute inset-0 bg-indigo-900/40"></div>
        <div class="relative z-10 p-6 md:p-8 text-white">
          <div class="text-2xl font-bold">GiftFlow</div>
          <div class="mt-1 text-sm text-gray-300">Internal Gifting Workflow</div>
        </div>
      </div>
      <div class="p-6 md:p-8">
        <h1 class="text-2xl font-bold mb-4 text-gray-800 text-center md:text-left">Login</h1>
        <p class="mb-4 text-gray-600 text-center md:text-left">Select a role to simulate login:</p>
        <div class="space-y-3">
          <button 
            v-for="user in mockUsers" 
            :key="user.username"
            @click="loginAs(user)"
            :disabled="loading"
            class="w-full py-2.5 px-4 rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-colors"
          >
            {{ user.label }}
          </button>
        </div>
        <div v-if="error" class="mt-4 text-red-600 text-sm text-center md:text-left">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>
