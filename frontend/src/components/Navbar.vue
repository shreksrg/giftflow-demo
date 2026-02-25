<script setup>
import { computed } from 'vue'

const props = defineProps(['user', 'currentView'])
const emit = defineEmits(['navigate', 'logout'])

const navItems = computed(() => {
  const items = []
  
  if (props.user.role === 'DEPT_ADMIN') {
    items.push({ id: 'gifts', label: 'Browse Gifts' })
    items.push({ id: 'dashboard', label: 'My Applications' })
  } else {
    items.push({ id: 'dashboard', label: 'Pending Approvals' })
    // Approvers can't browse/apply in this version
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
              <button
                v-for="item in navItems"
                :key="item.id"
                @click="emit('navigate', item.id)"
                :class="[
                  currentView === item.id 
                    ? 'bg-indigo-700 text-white' 
                    : 'text-indigo-200 hover:bg-indigo-500 hover:text-white',
                  'px-3 py-2 rounded-md text-sm font-medium transition-colors'
                ]"
              >
                {{ item.label }}
              </button>
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
