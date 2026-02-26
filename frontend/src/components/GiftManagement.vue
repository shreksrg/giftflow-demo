<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps(['user'])
const gifts = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const form = ref({
  id: null,
  name: '',
  description: '',
  type: 'NORMAL',
  quantity: 0,
  imageUrl: '',
  isPublished: false
})
const selectedFile = ref(null)

const fetchGifts = async () => {
  loading.value = true
  try {
    const res = await axios.get('http://127.0.0.1:8080/api/gifts?limit=100')
    gifts.value = res.data.data
  } catch (e) {
    alert('Failed to fetch gifts: ' + e.message)
  } finally {
    loading.value = false
  }
}

onMounted(fetchGifts)

const openAddModal = () => {
  isEditing.value = false
  form.value = { name: '', description: '', type: 'NORMAL', quantity: 0, imageUrl: '', isPublished: false }
  selectedFile.value = null
  showModal.value = true
}

const openEditModal = (gift) => {
  isEditing.value = true
  form.value = { ...gift, isPublished: gift.isPublished || false }
  selectedFile.value = null
  showModal.value = true
}

const handleFileChange = (e) => {
  selectedFile.value = e.target.files[0]
}

const submitForm = async () => {
  try {
    let imageUrl = form.value.imageUrl
    
    if (selectedFile.value) {
      const formData = new FormData()
      formData.append('file', selectedFile.value)
      const res = await axios.post('http://127.0.0.1:8080/api/upload', formData)
      imageUrl = res.data.url
    }

    const payload = {
      name: form.value.name,
      description: form.value.description,
      type: form.value.type,
      quantity: parseInt(form.value.quantity),
      imageUrl,
      isPublished: form.value.isPublished
    }

    if (isEditing.value) {
      await axios.put(`http://127.0.0.1:8080/api/gifts/${form.value.id}`, payload)
    } else {
      await axios.post('http://127.0.0.1:8080/api/gifts', payload)
    }
    
    showModal.value = false
    fetchGifts()
  } catch (e) {
    alert('Operation failed: ' + (e.response?.data?.error || e.message))
  }
}

const deleteGift = async (id) => {
  if (!confirm('Are you sure you want to delete this gift?')) return
  try {
    await axios.delete(`http://127.0.0.1:8080/api/gifts/${id}`)
    fetchGifts()
  } catch (e) {
    alert('Delete failed: ' + (e.response?.data?.error || e.message))
  }
}
</script>

<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">Gift Management</h2>
      <button @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded text-sm font-medium">
        Add Gift
      </button>
    </div>

    <div v-if="loading" class="text-center py-10">Loading...</div>

    <div v-else class="flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Image</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Qty</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="gift in gifts" :key="gift.id">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <img :src="gift.imageUrl" class="h-10 w-10 object-cover rounded bg-gray-100" alt="">
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ gift.name }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span :class="gift.type === 'VIP' ? 'bg-yellow-100 text-yellow-800' : 'bg-green-100 text-green-800'" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                      {{ gift.type }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <span :class="gift.isPublished ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                      {{ gift.isPublished ? 'Published' : 'Draft' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ gift.quantity }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-3">
                    <button @click="openEditModal(gift)" class="text-indigo-600 hover:text-indigo-900">Edit</button>
                    <button @click="deleteGift(gift.id)" class="text-red-600 hover:text-red-900">Delete</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="showModal" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showModal = false"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
              {{ isEditing ? 'Edit' : 'Add' }} Gift
            </h3>
            <form @submit.prevent="submitForm" class="mt-4 space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Name</label>
                <input v-model="form.name" required class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Description</label>
                <textarea v-model="form.description" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"></textarea>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Type</label>
                <select v-model="form.type" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                  <option value="NORMAL">Normal</option>
                  <option value="VIP">VIP</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Quantity</label>
                <input type="number" v-model="form.quantity" required min="0" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
              </div>
              <div class="flex items-center">
                <button 
                  type="button" 
                  class="relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" 
                  :class="form.isPublished ? 'bg-indigo-600' : 'bg-gray-200'"
                  @click="form.isPublished = !form.isPublished"
                >
                  <span class="sr-only">Published</span>
                  <span 
                    aria-hidden="true" 
                    class="pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200"
                    :class="form.isPublished ? 'translate-x-5' : 'translate-x-0'"
                  ></span>
                </button>
                <span class="ml-3 cursor-pointer" @click="form.isPublished = !form.isPublished">
                  <span class="text-sm font-medium text-gray-900">Published</span>
                </span>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Image</label>
                <input type="file" @change="handleFileChange" accept="image/*" class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100">
                <div v-if="form.imageUrl" class="mt-2">
                  <p class="text-xs text-gray-500 mb-1">Current Image:</p>
                  <img :src="form.imageUrl" class="h-20 w-auto rounded object-cover">
                </div>
              </div>
              <div class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense">
                <button type="submit" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:col-start-2 sm:text-sm">
                  Save
                </button>
                <button type="button" @click="showModal = false" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
