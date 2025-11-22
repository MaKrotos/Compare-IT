<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Заголовок -->
      <div class="mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-light tracking-tight text-gray-900 mb-2">
              {{ collection.name || 'Подборка' }}
            </h1>
            <p class="text-gray-600">
              {{ collection.items?.length || 0 }} товаров в подборке
            </p>
          </div>
        </div>
      </div>

      <!-- Список товаров -->
      <div v-if="isLoading" class="flex items-center justify-center py-20">
        <Loader2 class="w-8 h-8 animate-spin text-blue-600" />
      </div>
      <div v-else-if="!collection.items || collection.items.length === 0" class="text-center py-20">
        <div>
          <div class="inline-block bg-gray-100 rounded-full mb-4">
            <Package class="w-12 h-12 text-gray-400" />
          </div>
          <p class="text-gray-500 text-lg font-light">
            В этой подборке пока нет товаров
          </p>
        </div>
      </div>
      <div v-else class="flex flex-wrap justify-center gap-6">
        <transition-group name="list" tag="div" class="flex flex-wrap justify-center gap-6">
          <div
            v-for="(item, index) in collection.items"
            :key="item.id"
            class="w-full md:w-80"
          >
            <ComparisonCard
              :item="item"
              :on-update="() => {}"
              :on-delete="() => {}"
              :readonly="true"
            />
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ComparisonCard from '@/components/comparison/ComparisonCard.vue'
import { Loader2, Package } from 'lucide-vue-next'
import { getPublicCollection } from '@/utils/publicCollections.js'

export default {
  name: 'PublicCollectionDetailView',
  components: {
    ComparisonCard,
    Loader2,
    Package
  },
  setup() {
    const route = useRoute()
    const collection = ref({})
    const isLoading = ref(true)

    // Загрузка публичной подборки при монтировании
    onMounted(async () => {
      await loadPublicCollection()
    })

    // Загрузка публичной подборки
    const loadPublicCollection = async () => {
      try {
        isLoading.value = true
        const publicLink = route.params.publicLink
        const data = await getPublicCollection(publicLink)
        collection.value = data
      } catch (error) {
        console.error('Ошибка при загрузке публичной подборки:', error)
      } finally {
        isLoading.value = false
      }
    }

    return {
      collection,
      isLoading,
      loadPublicCollection
    }
  }
}
</script>

<style scoped>
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