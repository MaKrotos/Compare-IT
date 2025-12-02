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
import { onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import { Lock, LogIn, Loader2 } from 'lucide-vue-next'
import { isAuthenticated, isLoading, error, loginWithTelegram, initAuth } from '@/utils/auth'

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
    
    // Проверка, является ли текущий маршрут публичным
    const isPublicRoute = computed(() => route.name === 'PublicCollection')
    
    // Для публичных маршрутов не показываем форму авторизации
    if (isPublicRoute.value) {
      isAuthenticated.value = true
    }
    
    // Следим за изменением маршрута
    watch(route, (newRoute) => {
      const isPublic = newRoute.name === 'PublicCollection'
      if (isPublic) {
        isAuthenticated.value = true
      } else {
        initAuth(newRoute)
      }
    })
    
    // Проверка авторизации при монтировании для не публичных маршрутов
    onMounted(() => {
      if (!isPublicRoute.value) {
        initAuth(route)
      }
    })
    
    return {
      isAuthenticated,
      isLoading,
      error,
      loginWithTelegram
    }
  }
}
</script>