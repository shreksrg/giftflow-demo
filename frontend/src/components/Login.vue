<script setup>
import { ref } from 'vue'
import axios from 'axios'

const emit = defineEmits(['login'])
const loading = ref(false)
const error = ref('')

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
    
    const response = await axios.get('http://127.0.0.1:8080/api/me')
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
  <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
    <h1 class="text-2xl font-bold mb-6 text-center text-gray-800">GiftFlow Login</h1>
    <p class="mb-4 text-gray-600 text-center">Select a role to simulate login:</p>
    
    <div class="space-y-3">
      <button 
        v-for="user in mockUsers" 
        :key="user.username"
        @click="loginAs(user)"
        :disabled="loading"
        class="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-colors"
      >
        {{ user.label }}
      </button>
    </div>

    <div v-if="error" class="mt-4 text-red-600 text-sm text-center">
      {{ error }}
    </div>
  </div>
</template>
