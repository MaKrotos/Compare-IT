<template>
  <div class="public-collection-view">
    <div class="public-collection-container">
      <!-- Заголовок -->
      <div class="public-collection-header">
        <div class="public-collection-logo">
          <div class="public-collection-logo-icon">
            <Scale class="public-collection-logo-scale" />
          </div>
        </div>
        <h1 class="public-collection-title">
          {{ collectionName || 'Публичная коллекция' }}
        </h1>
        <p class="public-collection-subtitle">
          Публичный доступ к коллекции товаров
        </p>
      </div>

      <!-- Сообщение о публичном доступе -->
      <div class="public-collection-message">
        <div class="public-collection-message-content">
          <Info class="public-collection-message-icon" />
          <div class="public-collection-message-text">
            <p class="public-collection-message-title">Публичный доступ</p>
            <p class="public-collection-message-description">
              Эта коллекция доступна для просмотра без авторизации.
              Для редактирования коллекции необходимо авторизоваться.
            </p>
          </div>
        </div>
      </div>

      <!-- Список товаров -->
      <div v-if="isLoading" class="public-collection-loading">
        <Loader2 class="public-collection-loading-icon" />
      </div>
      <div v-else-if="items.length === 0" class="public-collection-empty">
        <div class="public-collection-empty-content">
          <div class="public-collection-empty-icon">
            <Scale class="public-collection-empty-scale" />
          </div>
          <p class="public-collection-empty-text">
            В этой коллекции пока нет товаров
          </p>
        </div>
      </div>
      <div v-else class="public-collection-items">
        <transition-group
          name="public-collection-list"
          tag="div"
          class="public-collection-items-grid"
        >
          <div
            v-for="item in items"
            :key="`item-${item.id}`"
            class="public-collection-item-container"
          >
            <ComparisonCard
              :item="item"
              :readonly="true"
            />
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import ComparisonCard from '@/components/comparison/ComparisonCard.vue'
import { Scale, Loader2, Info } from 'lucide-vue-next'
import { getPublicCollection } from '@/utils/publicCollections.js'

export default {
  name: 'PublicCollection',
  components: {
    ComparisonCard,
    Scale,
    Loader2,
    Info
  },
  setup() {
    const route = useRoute()
    const items = ref([])
    const collectionName = ref('')
    const isLoading = ref(true)
    const error = ref(null)
    const refreshInterval = ref(null)

    // Загрузка публичной коллекции
    const loadPublicCollection = async () => {
      try {
        // При обновлении не показываем индикатор загрузки
        const wasLoading = isLoading.value;
        if (items.value.length > 0) {
          isLoading.value = false
        } else {
          isLoading.value = true
        }
        error.value = null
        
        const publicLink = route.params.publicLink
        if (!publicLink) {
          throw new Error('Не указана публичная ссылка')
        }
        
        const collection = await getPublicCollection(publicLink)
        
        // Сравниваем новые и старые значения поэлементно для более точного определения изменений
        const oldItems = items.value
        const newItems = collection.items || []
        
        // Проверяем, изменились ли данные
        const hasChanges = JSON.stringify(oldItems) !== JSON.stringify(newItems);
        
        // Обновляем только если данные действительно изменились
        if (hasChanges) {
          items.value = newItems
        }
        
        collectionName.value = collection.name || 'Публичная коллекция'
      } catch (err) {
        console.error('Ошибка при загрузке публичной коллекции:', err)
        error.value = err.message || 'Не удалось загрузить коллекцию'
        
        // Показываем ошибку пользователю (в реальном приложении можно использовать toast-уведомления)
        // Например: showToast('error', error.value)
      } finally {
        isLoading.value = false
      }
    }

    onMounted(() => {
      loadPublicCollection()
      // Установка интервала для автоматического обновления каждые 30 секунд
      refreshInterval.value = setInterval(() => {
        loadPublicCollection()
      }, 30000)
    })

    // Очистка интервала при размонтировании компонента
    onUnmounted(() => {
      if (refreshInterval.value) {
        clearInterval(refreshInterval.value)
      }
    })

    return {
      items,
      collectionName,
      isLoading,
      error
    }
  }
}
</script>

<style scoped>
/* Стили для компонента PublicCollection */
.public-collection-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f9fafb 0%, #eff6ff 30%, #f9fafb 100%);
}

.public-collection-container {
  margin: 0 auto;
  padding: 3rem 1rem;
}

/* Заголовок */
.public-collection-header {
  text-align: center;
  margin-bottom: 3rem;
}

.public-collection-logo {
  margin-bottom: 1rem;
}

.public-collection-logo-icon {
  display: inline-flex;
  padding: 0.75rem;
  background-color: #2563eb;
  border-radius: 1rem;
}

.public-collection-logo-scale {
  width: 2rem;
  height: 2rem;
  color: white;
}

.public-collection-title {
  font-size: 3rem;
  font-weight: 300;
  letter-spacing: -0.025em;
  color: #111827;
  margin-bottom: 0.75rem;
}

@media (max-width: 640px) {
  .public-collection-title {
    font-size: 2.25rem;
  }
}

.public-collection-subtitle {
  font-size: 1.125rem;
  color: #4b5563;
  font-weight: 300;
}

/* Сообщение о публичном доступе */
.public-collection-message {
  margin-bottom: 2rem;
  padding: 1rem;
  background-color: #dbeafe;
  border-radius: 0.5rem;
  border: 1px solid #bfdbfe;
}

.public-collection-message-content {
  display: flex;
  align-items: flex-start;
}

.public-collection-message-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563eb;
  margin-top: 0.125rem;
  margin-right: 0.5rem;
  flex-shrink: 0;
}

.public-collection-message-text {
  flex: 1;
}

.public-collection-message-title {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #1e40af;
  margin-bottom: 0.125rem;
}

.public-collection-message-description {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #2563eb;
}

/* Состояние загрузки */
.public-collection-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5rem 0;
}

.public-collection-loading-icon {
  width: 2rem;
  height: 2rem;
  animation: spin 1s linear infinite;
  color: #2563eb;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Пустое состояние */
.public-collection-empty {
  text-align: center;
  padding: 5rem 0;
}

.public-collection-empty-content {
  display: inline-block;
}

.public-collection-empty-icon {
  display: inline-flex;
  background-color: #f3f4f6;
  border-radius: 50%;
  padding: 1rem;
  margin-bottom: 1rem;
}

.public-collection-empty-scale {
  width: 3rem;
  height: 3rem;
  color: #9ca3af;
}

.public-collection-empty-text {
  color: #6b7280;
  font-size: 1.125rem;
  font-weight: 300;
}

/* Список товаров */
.public-collection-items-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1.5rem;
}

.public-collection-item-container {
  width: 100%;
  max-width: 20rem;
}

/* Анимации */
.public-collection-list-enter-active,
.public-collection-list-leave-active {
  transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
}

.public-collection-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.public-collection-list-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.public-collection-list-leave-active {
  position: absolute;
}
</style>
