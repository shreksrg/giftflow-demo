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
    if (props.user.role === 'CPRO_ADMIN') {
      items.push({ path: '/gift-management', label: 'Manage Gifts' })
    }
  }
  return items
})
</script>

<template>
  <nav class="bg-indigo-600">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between py-3">
        <div class="flex-shrink-0 text-white font-bold text-xl">
          GiftFlow
        </div>

        <div class="flex items-center gap-3 text-white text-sm">
          <span class="hidden sm:inline">{{ user.username }} ({{ user.role }})</span>
          <span class="sm:hidden">{{ user.username }}</span>
          <button 
            @click="emit('logout')"
            class="bg-indigo-800 hover:bg-indigo-900 px-3 py-2 rounded-md text-xs text-white"
          >
            Logout
          </button>
        </div>
      </div>

      <div class="border-t border-indigo-500/50">
        <div class="flex overflow-x-auto">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            active-class="border-white text-white"
            class="whitespace-nowrap flex-none px-4 py-3 text-sm font-medium text-indigo-200 border-b-2 border-transparent hover:text-white hover:border-indigo-200 transition-colors"
          >
            {{ item.label }}
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>
