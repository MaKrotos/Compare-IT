import { ref } from 'vue'
import { useRoute } from 'vue-router'

// Создаем реактивные переменные на уровне модуля
export const isAuthenticated = ref(false)
export const isLoading = ref(false)
export const error = ref(null)

// Получение токена из localStorage
export const getToken = () => {
  const token = localStorage.getItem('auth_token')
  
  // Если токена нет, возвращаем null
  if (!token) {
    return null
  }
  
  // Проверяем срок действия токена
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const now = Math.floor(Date.now() / 1000)
    
    // Если токен истек, удаляем его и возвращаем null
    if (payload.exp && payload.exp < now) {
      localStorage.removeItem('auth_token')
      isAuthenticated.value = false
      // Выполняем повторную инициализацию авторизации
      reinitAuth()
      return null
    }
    
    return token
  } catch (error) {
    // Если не удалось распарсить токен, удаляем его
      // Выполняем повторную инициализацию авторизации
      reinitAuth()
    localStorage.removeItem('auth_token')
    isAuthenticated.value = false
    return null
  }
}

// Проверка статуса авторизации
export const checkAuthStatus = () => {
  const token = getToken()
  isAuthenticated.value = !!token
}

// Авторизация через Telegram
export const loginWithTelegram = async () => {
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
  } catch (err) {
    console.error('Ошибка авторизации:', err)
    error.value = err.message
    // Не показываем alert при автоматической авторизации
    if (!window.Telegram?.WebApp) {
      alert('Ошибка авторизации: ' + err.message)
    }
  } finally {
    isLoading.value = false
  }
}

// Инициализация состояния авторизации
export const initAuth = (route) => {
  // Если route не передан, получаем его из useRoute
  const routeToUse = route || useRoute()
  
  // Проверка, является ли текущий маршрут публичным
  const isPublicRoute = routeToUse.name === 'PublicCollection'
  
  // Для публичных маршрутов пропускаем авторизацию
  if (isPublicRoute) {
    isAuthenticated.value = true
    return
  }
  
  checkAuthStatus()
  
  // Попытка автоматической авторизации только если пользователь еще не авторизован
  // и это не публичный маршрут
  if (window.Telegram?.WebApp && !isAuthenticated.value) {
    loginWithTelegram()
  }
}

// Повторная инициализация авторизации
export const reinitAuth = () => {
  // Сбрасываем состояние авторизации
  isAuthenticated.value = false
  
  // Повторная инициализация
  initAuth()
}
