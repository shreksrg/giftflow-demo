<script setup>
import { computed } from 'vue'

const props = defineProps(['user'])
const emit = defineEmits(['logout'])

const navItems = computed(() => {
  const items = []
  
  if (props.user.role === 'DEPT_ADMIN') {
    items.push({ path: '/gifts', label: 'Browse Gifts' })
    items.push({ path: '/dashboard', label: 'My Applications' })
  } else {
    items.push({ path: '/dashboard', label: 'Pending Approvals' })
  }
  return items
})
</script>

<template>
  <nav class="bg-indigo-600">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0 text-white font-bold text-xl">
            GiftFlow
          </div>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <router-link
                v-for="item in navItems"
                :key="item.path"
                :to="item.path"
                active-class="bg-indigo-700 text-white"
                class="text-indigo-200 hover:bg-indigo-500 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                {{ item.label }}
              </router-link>
            </div>
          </div>
        </div>
        <div class="hidden md:block">
          <div class="ml-4 flex items-center md:ml-6 text-white text-sm">
            <span class="mr-4">{{ user.username }} ({{ user.role }})</span>
            <button 
              @click="emit('logout')"
              class="bg-indigo-800 hover:bg-indigo-900 px-3 py-2 rounded-md text-xs"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>
