<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Заголовок -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-4 mb-4">
          <div class="p-3 bg-blue-600 rounded-2xl">
            <Scale class="w-8 h-8 text-white" />
          </div>
        </div>
        <h1 class="text-4xl sm:text-5xl font-light tracking-tight text-gray-900 mb-3">
          {{ collectionName || 'Публичная коллекция' }}
        </h1>
        <p class="text-lg text-gray-600 font-light">
          Публичный доступ к коллекции товаров
        </p>
      </div>

      <!-- Сообщение о публичном доступе -->
      <div class="mb-8 p-4 bg-blue-50 rounded-lg border border-blue-200">
        <div class="flex items-start">
          <Info class="w-5 h-5 text-blue-600 mt-0.5 mr-2 flex-shrink-0" />
          <div>
            <p class="text-sm text-blue-800 font-medium">Публичный доступ</p>
            <p class="text-sm text-blue-600">
              Эта коллекция доступна для просмотра без авторизации. 
              Для редактирования коллекции необходимо авторизоваться.
            </p>
          </div>
        </div>
      </div>

      <!-- Список товаров -->
      <div v-if="isLoading" class="flex items-center justify-center py-20">
        <Loader2 class="w-8 h-8 animate-spin text-blue-600" />
      </div>
      <div v-else-if="items.length === 0" class="text-center py-20">
        <div>
          <div class="inline-block bg-gray-100 rounded-full mb-4">
            <Scale class="w-12 h-12 text-gray-400" />
          </div>
          <p class="text-gray-500 text-lg font-light">
            В этой коллекции пока нет товаров
          </p>
        </div>
      </div>
      <div v-else class="flex flex-wrap justify-center gap-6">
        <transition-group name="list" tag="div" class="flex flex-wrap justify-center gap-6">
          <div
            v-for="(item, index) in items"
            :key="item.id"
            class="w-full md:w-80"
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
        isLoading.value = true
        error.value = null
        
        const publicLink = route.params.publicLink
        if (!publicLink) {
          throw new Error('Не указана публичная ссылка')
        }
        
        const collection = await getPublicCollection(publicLink)
        items.value = collection.items || []
        collectionName.value = collection.name || 'Публичная коллекция'
      } catch (err) {
        console.error('Ошибка при загрузке публичной коллекции:', err)
        error.value = err.message || 'Не удалось загрузить коллекцию'
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
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

* {
  font-family: 'Inter', sans-serif;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>