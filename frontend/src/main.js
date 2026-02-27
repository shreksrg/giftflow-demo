import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'

axios.defaults.baseURL =
  import.meta.env.VITE_API_BASE_URL || `http://${window.location.hostname}:8080`

createApp(App).use(router).mount('#app')
