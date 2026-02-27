<script setup>
import { computed, ref } from 'vue'

const props = defineProps(['user'])
const emit = defineEmits(['logout'])

const isOpen = ref(false)

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
              class="bg-indigo-800 hover:bg-indigo-900 px-3 py-2 rounded-md text-xs text-white"
            >
              Logout
            </button>
          </div>
        </div>
        <div class="-mr-2 flex md:hidden">
          <!-- Mobile menu button -->
          <button @click="isOpen = !isOpen" type="button" class="bg-indigo-600 inline-flex items-center justify-center p-2 rounded-md text-indigo-200 hover:text-white hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-600 focus:ring-white" aria-controls="mobile-menu" aria-expanded="false">
            <span class="sr-only">Open main menu</span>
            <svg :class="{'hidden': isOpen, 'block': !isOpen }" class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg :class="{'block': isOpen, 'hidden': !isOpen }" class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile menu, show/hide based on menu state. -->
    <div v-show="isOpen" class="md:hidden" id="mobile-menu">
      <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          active-class="bg-indigo-700 text-white"
          class="text-indigo-200 hover:bg-indigo-500 hover:text-white block px-3 py-2 rounded-md text-base font-medium"
          @click="isOpen = false"
        >
          {{ item.label }}
        </router-link>
      </div>
      <div class="pt-4 pb-3 border-t border-indigo-700">
        <div class="flex items-center px-5">
          <div class="ml-3">
            <div class="text-base font-medium leading-none text-white">{{ user.username }}</div>
            <div class="text-sm font-medium leading-none text-indigo-200">{{ user.role }}</div>
          </div>
          <button 
            @click="emit('logout')"
            class="ml-auto bg-indigo-800 hover:bg-indigo-900 px-3 py-2 rounded-md text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-800 focus:ring-white"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>
