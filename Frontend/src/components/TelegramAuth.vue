<template>
  <div v-if="!isAuthenticated" class="telegram-auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <div class="auth-icon">
          <Lock class="icon" />
        </div>
        <h2 class="auth-title">Авторизация</h2>
        <p class="auth-subtitle">Войдите через Telegram для доступа к подборкам</p>
      </div>
      
      <div class="auth-content">
        <Button
          @click="loginWithTelegram"
          class="auth-button"
          :disabled="isLoading"
        >
          <component :is="isLoading ? 'Loader2' : 'LogIn'" class="button-icon" />
          {{ isLoading ? 'Авторизация...' : 'Войти через Telegram' }}
        </Button>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <p class="auth-note">
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

<style scoped>
.telegram-auth-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f9fafb 0%, #eff6ff 30%, #f9fafb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-card {
  max-width: 28rem;
  width: 100%;
  background-color: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  padding: 2rem;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  background-color: #dbeafe;
  border-radius: 9999px;
  margin-bottom: 1rem;
}

.auth-icon .icon {
  width: 2rem;
  height: 2rem;
  color: #2563eb;
}

.auth-title {
  font-size: 1.5rem;
  font-weight: 300;
  color: #111827;
  margin-bottom: 0.5rem;
}

.auth-subtitle {
  color: #4b5563;
}

.auth-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.auth-button {
  width: 100%;
  gap: 0.5rem;
  background-color: #2563eb;
  color: white;
}

.auth-button:hover {
  background-color: #1d4ed8;
}

.button-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.error-message {
  font-size: 0.875rem;
  color: #ef4444;
  text-align: center;
}

.auth-note {
  font-size: 0.875rem;
  color: #6b7280;
  text-align: center;
}
</style>