<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const props = defineProps(['user'])
const applications = ref([])
const loading = ref(true)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 10
const currentView = ref('pending') // pending or history
const selectedApp = ref(null)

const selectedImage = ref(null)

const showCommentModal = ref(false)
const pendingAction = ref(null)
const commentForm = ref({
  comment: ''
})

const searchFilters = ref({
  status: '',
  giftName: '',
  startDate: '',
  endDate: ''
})

const fetchApplications = async (page = 1) => {
  loading.value = true
  try {
    const pageNum = typeof page === 'number' ? page : currentPage.value
    const params = { page: pageNum, limit, view: currentView.value }

    if (props.user.role === 'DEPT_ADMIN') {
      if (searchFilters.value.status) params.status = searchFilters.value.status
      if (searchFilters.value.giftName) params.giftName = searchFilters.value.giftName
      if (searchFilters.value.startDate) params.startDate = searchFilters.value.startDate
      if (searchFilters.value.endDate) params.endDate = searchFilters.value.endDate
    }

    const response = await axios.get('/api/applications', { params })
    applications.value = response.data.data || []
    currentPage.value = response.data.page || 1
    totalPages.value = response.data.totalPages || 1
  } catch (err) {
    console.error(err)
    alert('Failed to fetch applications: ' + (err.response?.data?.error || err.message))
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchFilters.value = {
    status: '',
    giftName: '',
    startDate: '',
    endDate: ''
  }
  fetchApplications(1)
}

const setView = (view) => {
  currentView.value = view
  currentPage.value = 1
  fetchApplications(1)
}

onMounted(fetchApplications)

const canApprove = (app) => {
  if (props.user.role === 'DEPT_HEAD' && app.status === 'PENDING_HEAD') return true
  if (props.user.role === 'CPRO_ADMIN' && app.status === 'PENDING_CPRO') return true
  return false
}

const updateStatus = async (app, action) => {
  pendingAction.value = { app, action }
  commentForm.value.comment = ''
  showCommentModal.value = true
}

const confirmStatusUpdate = async () => {
  if (!pendingAction.value) return

  const { app, action } = pendingAction.value
  const comment = commentForm.value.comment || ''
  
  if (!comment.trim()) {
    alert('Comment is required.')
    return
  }

  try {
    await axios.put(`/api/applications/${app.id}/status`, {
      status: action === 'REJECT' ? 'REJECTED' : 'APPROVED',
      comment
    })
    showCommentModal.value = false
    await fetchApplications(currentPage.value) // Refresh list
  } catch (err) {
    alert('Failed to update status: ' + (err.response?.data?.error || err.message))
  } finally {
    pendingAction.value = null
  }
}

const statusClass = (status) => {
  switch (status) {
    case 'PENDING_HEAD': return 'bg-yellow-100 text-yellow-800'
    case 'PENDING_CPRO': return 'bg-orange-100 text-orange-800'
    case 'APPROVED_HEAD': return 'bg-blue-100 text-blue-800' // Intermediate
    case 'APPROVED_CPRO': return 'bg-green-100 text-green-800'
    case 'COMPLETED': return 'bg-green-100 text-green-800'
    case 'REJECT_HEAD': return 'bg-red-100 text-red-800'
    case 'REJECT_CPRO': return 'bg-red-200 text-red-900'
    default: return 'bg-gray-100 text-gray-800'
  }
}
</script>

<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-900">
        {{ user.role === 'DEPT_ADMIN' ? 'My Applications' : 'Approvals Dashboard' }}
      </h2>
      <button @click="fetchApplications(currentPage)" class="text-sm text-indigo-600 hover:text-indigo-900">
        Refresh
      </button>
    </div>

    <!-- Search Filters for Dept Admin -->
    <div v-if="user.role === 'DEPT_ADMIN'" class="mb-6 bg-white p-4 rounded-lg shadow sm:p-6">
      <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
        <!-- Status -->
        <div class="sm:col-span-2">
          <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
          <select v-model="searchFilters.status" id="status" class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md">
            <option value="">All Statuses</option>
            <option value="PENDING_HEAD">Pending Head</option>
            <option value="APPROVED_HEAD">Approved Head</option>
            <option value="PENDING_CPRO">Pending CPRO</option>
            <option value="APPROVED_CPRO">Approved CPRO</option>
            <option value="REJECT_HEAD">Rejected Head</option>
            <option value="REJECT_CPRO">Rejected CPRO</option>
          </select>
        </div>

        <!-- Product Name -->
        <div class="sm:col-span-2">
          <label for="giftName" class="block text-sm font-medium text-gray-700">Product Name</label>
          <input type="text" v-model="searchFilters.giftName" id="giftName" class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md" placeholder="Search by name">
        </div>

        <!-- Date Range -->
        <div class="sm:col-span-2">
          <label class="block text-sm font-medium text-gray-700">Date Range</label>
          <div class="flex space-x-2 mt-1">
            <input type="date" v-model="searchFilters.startDate" class="focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md">
            <span class="text-gray-500 self-center">-</span>
            <input type="date" v-model="searchFilters.endDate" class="focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md">
          </div>
        </div>
      </div>
      <div class="mt-4 flex justify-end">
        <button @click="fetchApplications(1)" class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          Search
        </button>
        <button @click="resetSearch" class="ml-3 inline-flex justify-center py-2 px-4 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          Reset
        </button>
      </div>
    </div>

    <!-- View Toggle for Dept Head and CPRO Admin -->
    <div v-if="user.role === 'DEPT_HEAD' || user.role === 'CPRO_ADMIN'" class="mb-4 border-b border-gray-200">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button
          @click="setView('pending')"
          :class="[
            currentView === 'pending'
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm'
          ]"
        >
          {{ user.role === 'CPRO_ADMIN' ? 'Active Applications' : 'Pending Approvals' }}
        </button>
        <button
          @click="setView(user.role === 'CPRO_ADMIN' ? 'rejected' : 'history')"
          :class="[
            (user.role === 'CPRO_ADMIN' ? currentView === 'rejected' : currentView === 'history')
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
            'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm'
          ]"
        >
          {{ user.role === 'CPRO_ADMIN' ? 'Rejected Records' : 'Approval History' }}
        </button>
      </nav>
    </div>

    <div v-if="loading" class="text-center py-10">Loading...</div>
    
    <div v-else class="flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Application ID
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Image
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Gift
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Quantity
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th v-if="user.role !== 'DEPT_ADMIN'" scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Applicant
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Date
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Last Comment
                  </th>
                  <th scope="col" class="relative px-6 py-3">
                    <span class="sr-only">Actions</span>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="app in applications" :key="app.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ app.id.substring(0, 8) }}...
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <img 
                      v-if="app.imageUrl" 
                      :src="app.imageUrl" 
                      alt="Gift Image" 
                      class="h-10 w-10 rounded-full object-cover cursor-pointer hover:opacity-75"
                      @click="selectedImage = app.imageUrl"
                    >
                    <span v-else class="text-gray-400 text-xs">No Img</span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ app.giftName }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ app.quantity || 1 }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span :class="['px-2 inline-flex text-xs leading-5 font-semibold rounded-full', statusClass(app.status)]">
                      {{ app.status }}
                    </span>
                  </td>
                  <td v-if="user.role !== 'DEPT_ADMIN'" class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ app.applicantId }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ new Date(app.createdAt).toLocaleDateString() }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 max-w-xs truncate" :title="app.history && app.history.length > 0 ? app.history[app.history.length - 1].comment : ''">
                    {{ (app.history && app.history.length > 0 && app.history[app.history.length - 1].comment) ? app.history[app.history.length - 1].comment : '--' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div v-if="canApprove(app)" class="flex justify-end space-x-2">
                      <button 
                        @click="updateStatus(app, 'APPROVE')"
                        class="text-green-600 hover:text-green-900"
                      >
                        Approve
                      </button>
                      <button 
                        @click="updateStatus(app, 'REJECT')"
                        class="text-red-600 hover:text-red-900"
                      >
                        Reject
                      </button>
                    </div>
                    <div v-else class="text-gray-400 italic text-xs">
                      <button @click="selectedApp = app" class="hover:text-indigo-600 hover:underline focus:outline-none">
                        {{ ['COMPLETED', 'REJECT_HEAD', 'REJECT_CPRO', 'APPROVED_CPRO'].includes(app.status) ? 'Done' : 'Pending' }}
                      </button>
                    </div>
                  </td>
                </tr>
                <tr v-if="applications.length === 0">
                  <td colspan="7" class="px-6 py-4 text-center text-gray-500">
                    No applications found.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      
      <!-- Pagination -->
      <div v-if="!loading && totalPages > 1" class="mt-4 flex justify-between items-center">
        <span class="text-sm text-gray-700">
          Page {{ currentPage }} of {{ totalPages }}
        </span>
        <div class="space-x-2">
          <button 
            @click="fetchApplications(currentPage - 1)" 
            :disabled="currentPage === 1"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            Previous
          </button>
          <button 
            @click="fetchApplications(currentPage + 1)" 
            :disabled="currentPage === totalPages"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <!-- History Modal -->
    <div v-if="selectedApp" class="fixed z-10 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="selectedApp = null"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  Approval History: {{ selectedApp.giftName }}
                </h3>
                <div class="mt-4">
                  <ul class="divide-y divide-gray-200">
                    <li v-for="(item, index) in selectedApp.history" :key="index" class="py-4">
                      <div class="flex space-x-3">
                        <div class="flex-1 space-y-1">
                          <div class="flex items-center justify-between">
                            <h3 class="text-sm font-medium" :class="statusClass(item.status).split(' ')[1]">{{ item.status }}</h3>
                            <p class="text-sm text-gray-500">{{ new Date(item.timestamp).toLocaleString() }}</p>
                          </div>
                          <p class="text-sm text-gray-500">User: {{ item.changedBy }}</p>
                          <p class="text-sm text-gray-700 mt-1">Comment: {{ item.comment }}</p>
                        </div>
                      </div>
                    </li>
                  </ul>
                  <div v-if="!selectedApp.history || selectedApp.history.length === 0" class="text-center text-gray-500 py-4">
                    No history available.
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" @click="selectedApp = null">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comment Modal -->
    <div v-if="showCommentModal" class="fixed z-20 inset-0 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showCommentModal = false"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
              {{ pendingAction?.action === 'APPROVE' ? 'Approve' : 'Reject' }} Application
            </h3>
            <div class="mt-4">
              <label for="comment" class="block text-sm font-medium text-gray-700">Comment <span class="text-red-500">*</span></label>
              <textarea
                id="comment"
                v-model="commentForm.comment"
                rows="3"
                required
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 mt-1 block w-full sm:text-sm border border-gray-300 rounded-md"
                placeholder="Enter your comment here (required)..."
              ></textarea>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              type="button" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm"
              :class="pendingAction?.action === 'APPROVE' ? 'bg-green-600 hover:bg-green-700 focus:ring-green-500' : 'bg-red-600 hover:bg-red-700 focus:ring-red-500'"
              @click="confirmStatusUpdate"
            >
              Confirm
            </button>
            <button 
              type="button" 
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" 
              @click="showCommentModal = false"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Image Modal -->
    <div v-if="selectedImage" class="fixed z-20 inset-0 overflow-y-auto" aria-labelledby="image-modal-title" role="dialog" aria-modal="true">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-black bg-opacity-75 transition-opacity" aria-hidden="true" @click="selectedImage = null"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
          <div class="bg-white p-2">
            <img :src="selectedImage" class="w-full h-auto rounded" alt="Full size gift image">
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" @click="selectedImage = null">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
