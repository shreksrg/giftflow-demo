import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import GiftList from '../components/GiftList.vue'
import Dashboard from '../components/Dashboard.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/gifts',
    name: 'GiftList',
    component: GiftList,
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
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

export default router
