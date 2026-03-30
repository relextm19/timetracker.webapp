import { createApp } from 'vue'
import App from './App.vue'
import './styles/tailwind.css'
import router from './router.ts'

import { setupCalendar } from 'v-calendar'

createApp(App).use(router).use(setupCalendar, {}).mount('#app');
