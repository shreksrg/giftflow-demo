<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const props = defineProps(['user'])
const applications = ref([])
const loading = ref(true)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 10

const fetchApplications = async (page = 1) => {
  loading.value = true
  try {
    const pageNum = typeof page === 'number' ? page : currentPage.value
    const response = await axios.get('http://localhost:8080/api/applications', {
      params: { page: pageNum, limit }
    })
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

onMounted(fetchApplications)

const canApprove = (app) => {
  if (props.user.role === 'DEPT_HEAD' && app.status === 'PENDING_HEAD') return true
  if (props.user.role === 'CPRO_ADMIN' && app.status === 'PENDING_CPRO') return true
  return false
}

const updateStatus = async (app, action) => {
  const comment = prompt(`Enter comment for ${action}:`)
  if (comment === null) return // Cancelled

  try {
    await axios.put(`http://localhost:8080/api/applications/${app.id}/status`, {
      action,
      comment
    })
    await fetchApplications(currentPage.value) // Refresh list
  } catch (err) {
    alert('Failed to update status: ' + (err.response?.data?.error || err.message))
  }
}

const statusClass = (status) => {
  switch (status) {
    case 'PENDING_HEAD': return 'bg-yellow-100 text-yellow-800'
    case 'PENDING_CPRO': return 'bg-orange-100 text-orange-800'
    case 'APPROVED_HEAD': return 'bg-blue-100 text-blue-800' // Intermediate
    case 'COMPLETED': return 'bg-green-100 text-green-800'
    case 'REJECTED': return 'bg-red-100 text-red-800'
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
                    Gift
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
                    <div class="text-sm font-medium text-gray-900">{{ app.giftName }}</div>
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
                      {{ app.status === 'COMPLETED' ? 'Done' : 'Pending' }}
                    </div>
                  </td>
                </tr>
                <tr v-if="applications.length === 0">
                  <td colspan="6" class="px-6 py-4 text-center text-gray-500">
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
  </div>
</template>
