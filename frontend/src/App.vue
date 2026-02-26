<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from './components/Navbar.vue'
import axios from 'axios'

const router = useRouter()
const user = ref(null)

// Restore user session on app load
onMounted(() => {
  const storedUser = localStorage.getItem('user')
  if (storedUser) {
    try {
      const parsedUser = JSON.parse(storedUser)
      user.value = parsedUser
      // Set axios default header for all subsequent requests
      axios.defaults.headers.common['X-Mock-User'] = parsedUser.username
    } catch (e) {
      console.error('Failed to parse user session', e)
      localStorage.removeItem('user')
    }
  }
})

const handleLogin = (newUser) => {
  user.value = newUser
  localStorage.setItem('user', JSON.stringify(newUser))
  axios.defaults.headers.common['X-Mock-User'] = newUser.username
  
  if (newUser.role === 'DEPT_ADMIN') {
    router.push('/gifts')
  } else {
    router.push('/dashboard')
  }
}

const handleLogout = () => {
  user.value = null
  localStorage.removeItem('user')
  delete axios.defaults.headers.common['X-Mock-User']
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Navbar shows only when logged in -->
    <Navbar v-if="user" :user="user" @logout="handleLogout" />
    
    <main class="flex-grow w-full max-w-7xl mx-auto py-6 sm:px-6 lg:px-8" :class="{ 'flex items-center justify-center': !user }">
      <!-- 
        We use router-view to render the current page component.
        We pass 'user' prop to all components (Dashboard, GiftList need it).
        We listen for 'login' event from Login component.
      -->
      <router-view 
        :user="user" 
        @login="handleLogin" 
      />
    </main>
  </div>
</template>
