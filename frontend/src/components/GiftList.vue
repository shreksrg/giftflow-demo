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

const showConfirmModal = ref(false)
const pendingGift = ref(null)

const fetchGifts = async (page = 1) => {
  loading.value = true
  try {
    const pageNum = typeof page === 'number' ? page : 1
    const response = await axios.get('/api/gifts', {
      params: { page: pageNum, limit, published: true }
    })
    gifts.value = (response.data.data || []).map(g => ({ ...g, selectedQuantity: 1 }))
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
  if (gift.selectedQuantity < 1) {
      alert("Quantity must be at least 1")
      return
  }
  if (gift.selectedQuantity > gift.quantity) {
      alert("Not enough stock")
      return
  }

  pendingGift.value = gift
  showConfirmModal.value = true
}

const confirmApplication = async () => {
  if (!pendingGift.value) return
  
  const gift = pendingGift.value
  applying.value = gift.id
  showConfirmModal.value = false
  
  try {
    await axios.post('/api/applications', {
      giftId: gift.id,
      quantity: gift.selectedQuantity
    })
    
    // Update local stock immediately
    gift.quantity -= gift.selectedQuantity
    
    successMsg.value = `Successfully applied for ${gift.selectedQuantity} x ${gift.name}!`
    setTimeout(() => successMsg.value = '', 3000)
  } catch (err) {
    alert('Failed to apply: ' + (err.response?.data?.error || err.message))
  } finally {
    applying.value = null
    pendingGift.value = null
  }
}
</script>

<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="flex flex-col gap-3 mb-6 sm:flex-row sm:justify-between sm:items-center">
      <h2 class="text-2xl font-bold text-gray-900">Available Gifts</h2>
      <div v-if="successMsg" class="bg-green-100 text-green-800 px-4 py-2 rounded-md transition-all sm:self-auto self-start">
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
          <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-0">
            <span class="text-sm text-gray-500">Stock: {{ gift.quantity }}</span>
            <div class="flex items-center gap-2">
              <input 
                type="number" 
                v-model.number="gift.selectedQuantity" 
                min="1" 
                :max="gift.quantity"
                class="w-20 h-10 border-gray-300 rounded-md text-sm focus:ring-indigo-500 focus:border-indigo-500 px-2"
              >
              <button 
                @click="applyForGift(gift)"
                :disabled="applying === gift.id || gift.quantity <= 0"
                class="flex-1 sm:flex-none inline-flex justify-center items-center px-4 h-10 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                {{ applying === gift.id ? 'Applying...' : 'Apply' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="!loading && totalPages > 1" class="mt-6 flex flex-col gap-3 sm:flex-row sm:justify-center sm:items-center sm:gap-0 sm:space-x-2">
      <button 
        @click="fetchGifts(currentPage - 1)" 
        :disabled="currentPage === 1"
        class="w-full sm:w-auto px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Previous
      </button>
      <span class="text-sm text-gray-700">
        Page {{ currentPage }} of {{ totalPages }}
      </span>
      <button 
        @click="fetchGifts(currentPage + 1)" 
        :disabled="currentPage === totalPages"
        class="w-full sm:w-auto px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next
      </button>
    </div>

    <!-- Confirmation Modal -->
    <div v-if="showConfirmModal" class="fixed z-20 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showConfirmModal = false"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-indigo-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  Confirm Application
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    Are you sure you want to apply for <span class="font-bold text-indigo-600">{{ pendingGift?.selectedQuantity }} x {{ pendingGift?.name }}</span>?
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button type="button" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm" @click="confirmApplication">
              Confirm
            </button>
            <button type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" @click="showConfirmModal = false">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
