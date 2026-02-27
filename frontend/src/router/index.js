import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import GiftList from '../components/GiftList.vue'
import Dashboard from '../components/Dashboard.vue'
import GiftManagement from '../components/GiftManagement.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: 'Login' }
  },
  {
    path: '/gifts',
    name: 'GiftList',
    component: GiftList,
    meta: { requiresAuth: true, title: 'Gifts' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true, title: 'Dashboard' }
  },
  {
    path: '/gift-management',
    name: 'GiftManagement',
    component: GiftManagement,
    meta: { requiresAuth: true, title: 'Gift Management' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : null
  
  if (to.meta.requiresAuth && !user) {
    next('/login')
  } else if (to.path === '/login' && user) {
    // Already logged in
    if (user.role === 'DEPT_ADMIN') {
      next('/gifts')
    } else {
      next('/dashboard')
    }
  } else {
    next()
  }
})

router.afterEach((to) => {
  const baseTitle = 'GiftFlow'
  let suffix = to.meta?.title || ''

  if (to.name === 'Dashboard') {
    const userStr = localStorage.getItem('user')
    const user = userStr ? JSON.parse(userStr) : null
    if (user?.role === 'DEPT_ADMIN') suffix = 'My Applications'
    else if (user) suffix = 'Approvals'
  }

  document.title = suffix ? `${baseTitle}-${suffix}` : baseTitle
})

export default router
