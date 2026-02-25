<script setup>
import { ref, computed } from 'vue'
import Login from './components/Login.vue'
import GiftList from './components/GiftList.vue'
import Dashboard from './components/Dashboard.vue'
import Navbar from './components/Navbar.vue'

const user = ref(null)
const currentView = ref('gifts')

const setUser = (newUser) => {
  user.value = newUser
  if (newUser) {
    // Default view based on role
    if (newUser.role === 'DEPT_ADMIN') {
      currentView.value = 'gifts'
    } else {
      currentView.value = 'dashboard'
    }
  } else {
    currentView.value = 'login' // This shouldn't happen as we unmount main app
  }
}

const setView = (view) => {
  currentView.value = view
}

const logout = () => {
  user.value = null
  currentView.value = 'login'
}
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div v-if="!user" class="flex items-center justify-center min-h-screen">
      <Login @login="setUser" />
    </div>
    
    <div v-else>
      <Navbar :user="user" :currentView="currentView" @navigate="setView" @logout="logout" />
      
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <GiftList v-if="currentView === 'gifts'" :user="user" />
        <Dashboard v-if="currentView === 'dashboard'" :user="user" />
      </main>
    </div>
  </div>
</template>
