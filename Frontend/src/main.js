import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Добавляем глобальные стили
import './assets/styles.css'
import './assets/theme.css'

const app = createApp(App)
app.use(router)
app.mount('#app')