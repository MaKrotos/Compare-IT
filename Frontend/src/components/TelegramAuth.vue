<template>
  <div v-if="!isAuthenticated" class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50 flex items-center justify-center">
    <div class="max-w-md w-full bg-white rounded-xl shadow-sm border border-gray-200 p-8">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
          <Lock class="w-8 h-8 text-blue-600" />
        </div>
        <h2 class="text-2xl font-light text-gray-900 mb-2">Авторизация</h2>
        <p class="text-gray-600">Войдите через Telegram для доступа к подборкам</p>
      </div>
      
      <div class="space-y-4">
        <Button
          @click="loginWithTelegram"
          class="w-full gap-2 bg-blue-600 hover:bg-blue-700 text-white"
          :disabled="isLoading"
        >
          <component :is="isLoading ? 'Loader2' : 'LogIn'" class="w-5 h-5" />
          {{ isLoading ? 'Авторизация...' : 'Войти через Telegram' }}
        </Button>
        
        <div v-if="error" class="text-sm text-red-500 text-center">
          {{ error }}
        </div>
        
        <p class="text-sm text-gray-500 text-center">
          Для работы с подборками необходимо авторизоваться через Telegram
        </p>
      </div>
    </div>
  </div>
  
  <div v-else>
    <slot></slot>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { Lock, LogIn, Loader2 } from 'lucide-vue-next'

export default {
  name: 'TelegramAuth',
  components: {
    Button,
    Lock,
    LogIn,
    Loader2
  },
  setup() {
    const route = useRoute()
    const isAuthenticated = ref(false)
    const isLoading = ref(false)
    const error = ref(null)
    
    // Проверка, является ли текущий маршрут публичным
    const isPublicRoute = route.name === 'PublicCollection'
    
    // Проверка авторизации при монтировании
    onMounted(() => {
      // Для публичных маршрутов пропускаем авторизацию
      if (isPublicRoute) {
        isAuthenticated.value = true
        return
      }
      
      checkAuthStatus()
      // Попытка автоматической авторизации только если пользователь еще не авторизован
      if (window.Telegram?.WebApp && !isAuthenticated.value) {
        loginWithTelegram()
      }
    })
    
    // Проверка статуса авторизации
    const checkAuthStatus = () => {
      const token = localStorage.getItem('auth_token')
      isAuthenticated.value = !!token
    }
    
    // Авторизация через Telegram
    const loginWithTelegram = async () => {
      try {
        console.log('Попытка авторизации через Telegram...')
        isLoading.value = true
        error.value = null
        
        // Проверяем, что Telegram Web App доступен
        if (!window.Telegram || !window.Telegram.WebApp) {
          throw new Error('Telegram Web App не доступен. Откройте приложение через Telegram.')
        }
        
        const webApp = window.Telegram.WebApp
        webApp.ready()
        console.log('Telegram Web App готов')
        
        // Получаем initData
        const initData = webApp.initData
        console.log('Init data получены:', initData ? 'есть' : 'нет')
        
        // Проверяем, что initData не пустая
        if (!initData) {
          throw new Error('Нет данных для авторизации. Откройте приложение через Telegram.')
        }
        
        console.log('Отправка данных авторизации на сервер...')
        // Отправляем данные на сервер для проверки
        const response = await fetch('/api/auth/telegram', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ initData })
        })
        
        console.log('Ответ от сервера:', response.status, response.statusText)
        
        if (!response.ok) {
          const errorText = await response.text()
          console.error('Ошибка авторизации от сервера:', errorText)
          throw new Error(`Ошибка авторизации: ${response.status} - ${errorText}`)
        }
        
        const data = await response.json()
        console.log('Успешная авторизация:', data)
        
        // Сохраняем токен в localStorage
        localStorage.setItem('auth_token', data.token)
        
        // Обновляем статус авторизации
        isAuthenticated.value = true
        
        // Не перезагружаем страницу, чтобы избежать бесконечного перерендеринга
        // window.location.reload()
      } catch (error) {
        console.error('Ошибка авторизации:', error)
        error.value = error.message
        // Не показываем alert при автоматической авторизации
        if (!window.Telegram?.WebApp) {
          alert('Ошибка авторизации: ' + error.message)
        }
      } finally {
        isLoading.value = false
      }
    }
    
    return {
      isAuthenticated,
      isLoading,
      error,
      loginWithTelegram
    }
  }
}
</script>