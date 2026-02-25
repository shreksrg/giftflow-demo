<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps(['user'])
const gifts = ref([])
const loading = ref(true)
const applying = ref(null)
const successMsg = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 8

const fetchGifts = async (page = 1) => {
  loading.value = true
  try {
    const pageNum = typeof page === 'number' ? page : 1
    const response = await axios.get('http://127.0.0.1:8080/api/gifts', {
      params: { page: pageNum, limit }
    })
    gifts.value = response.data.data || []
    currentPage.value = response.data.page || 1
    totalPages.value = response.data.totalPages || 1
  } catch (err) {
    console.error(err)
    // gifts.value = [] // Optional: clear list on error
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchGifts(1))

const applyForGift = async (gift) => {
  if (confirm(`Apply for ${gift.name}?`)) {
    applying.value = gift.id
    try {
      await axios.post('http://127.0.0.1:8080/api/applications', {
        giftId: gift.id
      })
      
      // Update local stock immediately
      gift.quantity--
      
      successMsg.value = `Successfully applied for ${gift.name}!`
      setTimeout(() => successMsg.value = '', 3000)
    } catch (err) {
      alert('Failed to apply: ' + (err.response?.data?.error || err.message))
    } finally {
      applying.value = null
    }
  }
}
</script>

<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">Available Gifts</h2>
      <div v-if="successMsg" class="bg-green-100 text-green-800 px-4 py-2 rounded-md transition-all">
        {{ successMsg }}
      </div>
    </div>

    <div v-if="loading" class="text-center py-10">
      Loading gifts...
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="gift in gifts" :key="gift.id" class="bg-white overflow-hidden shadow rounded-lg border border-gray-200 flex flex-col">
        <div class="h-48 w-full bg-gray-200 relative">
          <img :src="gift.imageUrl" :alt="gift.name" class="w-full h-full object-cover">
          <div class="absolute top-2 right-2">
            <span 
              :class="[
                gift.type === 'VIP' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800',
                'px-2 py-1 rounded-full text-xs font-semibold uppercase tracking-wide'
              ]"
            >
              {{ gift.type }}
            </span>
          </div>
        </div>
        <div class="p-4 flex-1 flex flex-col">
          <h3 class="text-lg font-medium text-gray-900 truncate">{{ gift.name }}</h3>
          <p class="mt-1 text-sm text-gray-500 flex-1">{{ gift.description }}</p>
          <div class="mt-4 flex items-center justify-between">
            <span class="text-sm text-gray-500">Stock: {{ gift.quantity }}</span>
            <button 
              @click="applyForGift(gift)"
              :disabled="applying === gift.id"
              class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              {{ applying === gift.id ? 'Applying...' : 'Apply' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="!loading && totalPages > 1" class="mt-6 flex justify-center items-center space-x-2">
      <button 
        @click="fetchGifts(currentPage - 1)" 
        :disabled="currentPage === 1"
        class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Previous
      </button>
      <span class="text-sm text-gray-700">
        Page {{ currentPage }} of {{ totalPages }}
      </span>
      <button 
        @click="fetchGifts(currentPage + 1)" 
        :disabled="currentPage === totalPages"
        class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next
      </button>
    </div>
  </div>
</template>
