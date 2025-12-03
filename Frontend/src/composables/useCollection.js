import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getComparisonCollectionById } from '@/utils/comparisons.js'
import { getToken } from '@/utils/auth'

export function useCollection() {
  const route = useRoute()
  const items = ref([])
  const isLoading = ref(true)
  const collectionId = ref(null)
  const isCollectionLoaded = ref(false)
  const originalCollectionName = ref('')
  const saveTimeout = ref(null)
  const ratingWeights = ref({
    priceRatingWeight: 20,
    prosConsRatingWeight: 80
  })

  // Инициализация коллекции
  const initializeCollection = async () => {
    collectionId.value = route.query.collectionId
    if (collectionId.value) {
      isCollectionLoaded.value = true
      try {
        const collection = await getComparisonCollectionById(collectionId.value)
        items.value = collection.items || []
        // Инициализируем настройки рейтингов из данных коллекции
        ratingWeights.value = {
          priceRatingWeight: collection.price_rating_weight || 20,
          prosConsRatingWeight: collection.pros_cons_rating_weight || 80
        }
        // Сохраняем оригинальное название коллекции
        originalCollectionName.value = collection.name || `Коллекция ${new Date().toLocaleDateString()}`
      } catch (error) {
        console.error('Ошибка при загрузке коллекции:', error)
        // Устанавливаем имя по умолчанию в случае ошибки
        originalCollectionName.value = `Коллекция ${new Date().toLocaleDateString()}`
      }
    } else {
      // Для новой коллекции устанавливаем имя по умолчанию
      originalCollectionName.value = `Коллекция ${new Date().toLocaleDateString()}`
    }
    isLoading.value = false
  }

  // Функция для автоматического сохранения коллекции
  const autoSaveCollection = async (updatedItems, updatedRatingWeights) => {
    // Очищаем предыдущий таймер, если он есть
    if (saveTimeout.value) {
      clearTimeout(saveTimeout.value)
    }

    // Устанавливаем новый таймер на 1 секунду
    saveTimeout.value = setTimeout(async () => {
      try {
        if (collectionId.value && isCollectionLoaded.value) {
          const token = getToken()
          if (!token) {
            console.log('Пользователь не авторизован, автосохранение отключено')
            return
          }

          const collection = {
            id: parseInt(collectionId.value),
            name: originalCollectionName.value,
            items: updatedItems,
            price_rating_weight: updatedRatingWeights.priceRatingWeight,
            pros_cons_rating_weight: updatedRatingWeights.prosConsRatingWeight
          }

          const response = await fetch('/backend/comparisons/update', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(collection)
          })

          if (!response.ok) {
            console.error('Ошибка при автосохранении коллекции:', response.status)
          }
        }
      } catch (error) {
        console.error('Ошибка при автосохранении коллекции:', error)
      }
    }, 1000)
  }

  // Сохранение текущей коллекции
  const saveCurrentCollection = async (updatedItems, updatedRatingWeights) => {
    try {
      const token = getToken()
      if (!token) {
        alert('Для сохранения коллекции необходимо авторизоваться')
        return
      }

      let response
      if (collectionId.value) {
        // Обновляем существующую коллекцию
        const collection = {
          id: parseInt(collectionId.value),
          name: originalCollectionName.value,
          items: updatedItems,
          price_rating_weight: updatedRatingWeights.priceRatingWeight,
          pros_cons_rating_weight: updatedRatingWeights.prosConsRatingWeight
        }

        response = await fetch('/backend/comparisons/update', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(collection)
        })
      } else {
        // Создаем новую коллекцию
        const collection = {
          name: `Коллекция ${new Date().toLocaleDateString()}`,
          items: updatedItems,
          price_rating_weight: updatedRatingWeights.priceRatingWeight,
          pros_cons_rating_weight: updatedRatingWeights.prosConsRatingWeight
        }

        response = await fetch('/backend/comparisons', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(collection)
        })
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      alert('Коллекция успешно сохранена!')
    } catch (error) {
      console.error('Ошибка при сохранении коллекции:', error)
      alert('Не удалось сохранить коллекцию: ' + error.message)
    }
  }

  return {
    items,
    isLoading,
    collectionId,
    isCollectionLoaded,
    originalCollectionName,
    ratingWeights,
    initializeCollection,
    autoSaveCollection,
    saveCurrentCollection
  }
}