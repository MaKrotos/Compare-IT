<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Заголовок -->
      <div
        class="text-center mb-12"
      >
        <div class="flex items-center justify-center gap-4 mb-4">
          <div class="p-3 bg-blue-600 rounded-2xl">
            <Scale class="w-8 h-8 text-white" />
          </div>
        </div>
        <h1 class="text-4xl sm:text-5xl font-light tracking-tight text-gray-900 mb-3">
          Сравнение товаров
        </h1>
        <p class="text-lg text-gray-600 font-light">
          Добавьте ссылки на товары и сравните их плюсы и минусы
        </p>
      </div>
    
    <!-- Кнопки экспорта/импорта -->
    <div class="import-export-buttons mb-8">
      <div class="flex flex-wrap gap-3 justify-center">
        <Button variant="outline" class="gap-2 px-4 py-2 text-base font-medium" @click="handleExport">
          <Download class="w-5 h-5" />
          Экспорт
        </Button>
        
        <Button variant="outline" class="gap-2 px-4 py-2 text-base font-medium" @click="triggerImport">
          <Upload class="w-5 h-5" />
          Импорт
        </Button>
        
        <!-- Кнопки для работы с коллекциями -->
        <Button variant="outline" class="gap-2 px-4 py-2 text-base font-medium" @click="goToCollections">
          <FolderOpen class="w-5 h-5" />
          Открыть
        </Button>
        
      </div>
      
      <!-- Поле для ввода данных импорта -->
      <div v-if="showImportField" class="mt-4 p-4 bg-gray-50 rounded-lg">
        <Textarea
          v-model="importData"
          placeholder="Вставьте JSON данные для импорта..."
          class="h-32 font-mono text-xs w-full"
        />
        <div class="flex gap-2 mt-3">
          <Button @click="handleImport" class="flex-1">
            Импортировать
          </Button>
          <Button variant="outline" @click="showImportField = false" class="flex-1">
            Отмена
          </Button>
        </div>
      </div>
      
      <!-- Модальное окно для выбора коллекции -->
      <Dialog v-model:open="showCollectionsDialog">
        <DialogContent class="max-w-md">
          <DialogHeader>
            <DialogTitle>Сохраненные коллекции</DialogTitle>
          </DialogHeader>
          <div class="py-4">
            <div class="text-center py-4 text-gray-500">
              Для работы с коллекциями перейдите на страницу коллекций
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <Button variant="outline" @click="showCollectionsDialog = false">Закрыть</Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
    
    <!-- Форма добавления -->
    <div
      class="mb-12"
    >
      <AddLinkForm :on-add-links="handleAddLinks" :is-loading="isLoadingItems" />
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
          Добавьте товары для начала сравнения
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
            :on-update="(data) => handleUpdateItem(item.id, data)"
            :on-delete="() => handleDeleteItem(item.id)"
          />
        </div>
      </transition-group>
    </div>
  </div>
</div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'
import DialogTrigger from '@/components/ui/DialogTrigger.vue'
import Button from '@/components/ui/Button.vue'
import Textarea from '@/components/ui/Textarea.vue'
import AddLinkForm from '@/components/comparison/AddLinkForm.vue'
import ComparisonCard from '@/components/comparison/ComparisonCard.vue'
import { Scale, Loader2, Download, Upload, FolderOpen, Save } from 'lucide-vue-next'
import { getComparisonCollectionById } from '@/utils/comparisons.js'

export default {
  name: 'ComparisonView',
  components: {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    Button,
    Textarea,
    AddLinkForm,
    ComparisonCard,
    Scale,
    Loader2,
    Download,
    Upload,
    Save,
    FolderOpen
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const items = ref([])
    const isLoadingItems = ref(false)
    const isLoading = ref(true)
    const importData = ref('')
    const showImportField = ref(false)
    const showCollectionsDialog = ref(false)
    let saveTimeout = null;

    // Инициализация
    onMounted(async () => {
      // Проверяем, есть ли параметр collectionId в URL
      const collectionId = route.query.collectionId;
      if (collectionId) {
        // Загружаем коллекцию с бэкенда
        try {
          const collection = await getComparisonCollectionById(collectionId);
          items.value = collection.items || [];
        } catch (error) {
          console.error('Ошибка при загрузке коллекции:', error);
        }
      }
      isLoading.value = false
    })
    
    // Функция для автоматического сохранения коллекции
    const autoSaveCollection = async () => {
      // Очищаем предыдущий таймер, если он есть
      if (saveTimeout) {
        clearTimeout(saveTimeout);
      }
      
      // Устанавливаем новый таймер на 1 секунду
      saveTimeout = setTimeout(async () => {
        try {
          // Получаем ID текущей коллекции из URL
          const collectionId = route.query.collectionId;
          
          if (collectionId) {
            // Отправляем коллекцию на сервер
            const token = localStorage.getItem('auth_token');
            if (!token) {
              console.log('Пользователь не авторизован, автосохранение отключено');
              return;
            }
            
            const collection = {
              id: parseInt(collectionId),
              name: `Коллекция ${new Date().toLocaleDateString()}`,
              items: items.value
            };
            
            const response = await fetch('/backend/comparisons/update', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
              },
              body: JSON.stringify(collection)
            });
            
            if (!response.ok) {
              console.error('Ошибка при автосохранении коллекции:', response.status);
            }
          }
        } catch (error) {
          console.error('Ошибка при автосохранении коллекции:', error)
        }
      }, 1000); // Сохраняем через 1 секунду после последнего изменения
    }
    
    // Отслеживаем изменения в items и запускаем автосохранение
    watch(items, () => {
      // Проверяем, есть ли ID коллекции в URL
      const collectionId = route.query.collectionId;
      if (collectionId) {
        autoSaveCollection();
      }
    }, { deep: true })

    // Функция для экспорта данных
    const handleExport = () => {
      try {
        const dataStr = JSON.stringify(items.value)
        const base64 = btoa(encodeURIComponent(dataStr))
        navigator.clipboard.writeText(base64).catch(err => {
          console.error('Не удалось скопировать в буфер обмена:', err)
        })
        alert('Данные скопированы в буфер обмена!')
      } catch (error) {
        console.error('Ошибка экспорта:', error)
      }
    }
    
    // Функция для импорта данных
    const handleImport = () => {
      try {
        if (!importData.value.trim()) {
          alert('Нет данных для импорта')
          return
        }
        
        // Декодируем данные из base64
        const decoded = decodeURIComponent(atob(importData.value.trim()))
        const importedItems = JSON.parse(decoded)
        
        if (Array.isArray(importedItems)) {
          // Добавляем уникальные ID для каждого импортированного элемента
          const itemsWithIds = importedItems.map(item => ({
            ...item,
            id: item.id || Date.now().toString() + Math.random().toString(36).substr(2, 9),
            created_date: item.created_date || new Date().toISOString()
          }))
          
          // Добавляем импортированные элементы к существующим
          items.value = [...items.value, ...itemsWithIds]
          showImportField.value = false
          importData.value = ''
          alert('Данные успешно импортированы!')
        }
      } catch (error) {
        console.error('Ошибка импорта:', error)
        alert('Ошибка импорта: ' + error.message)
      }
    }
    
    const triggerImport = () => {
      showImportField.value = true
    }
    
    const createItem = (itemData) => {
      // Импортируем функцию calculateRating из ComparisonCard
      const calculateRating = (itemData, allItems) => {
        // Конфигурационные параметры для настройки веса рейтинга
        const RATING_CONFIG = {
          // Максимальный рейтинг
          MAX_RATING: 100,
          // Вес рейтинга по плюсам/минусам (в процентах от максимального рейтинга)
          PROS_CONS_WEIGHT: 80,
          // Вес рейтинга по цене (в процентах от максимального рейтинга)
          PRICE_WEIGHT: 20
        };

        const pros = itemData.pros || []
        const cons = itemData.cons || []
        const itemPrice = itemData.price || 0

        // Сумма положительных баллов
        const prosScore = pros.reduce((sum, pro) => sum + (pro.impact || 5), 0)
        
        // Сумма отрицательных баллов
        const consScore = cons.reduce((sum, con) => sum + (con.impact || 5), 0)
        
        // Базовый рейтинг (0-100)
        const maxPossibleScore = Math.max((pros.length + cons.length) * 10, 1)
        const baseRating = ((prosScore - consScore) / maxPossibleScore) * 50 + 50
        
        // Корректировка на цену (чем выше цена, тем меньше рейтинг)
        // Находим максимальную и минимальную цену среди всех товаров для нормализации
        let priceImpact = 0;
        if (allItems.length > 0 && itemPrice > 0) {
          const maxPrice = Math.max(...allItems.map(item => item.price || 0));
          const minPrice = Math.min(...allItems.map(item => item.price || 0));
          if (maxPrice > minPrice) {
            // Чем выше цена по сравнению с минимальной, тем меньше рейтинг
            // Чем ниже цена по сравнению с максимальной, тем выше рейтинг
            const priceRange = maxPrice - minPrice;
            const priceRatio = (maxPrice - itemPrice) / priceRange;
            // Применяем вес цены к максимальному бонусу/штрафу
            priceImpact = priceRatio * RATING_CONFIG.PRICE_WEIGHT;
          }
        }
        
        // Учет относительности плюсов/минусов
        let prosConsImpact = 0;
        if (allItems.length > 0) {
          // Средние значения по всем карточкам
          const totalProsScore = allItems.reduce((sum, item) => sum + (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0), 0);
          const totalConsScore = allItems.reduce((sum, item) => sum + (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0), 0);
          const avgProsScore = totalProsScore / allItems.length;
          const avgConsScore = totalConsScore / allItems.length;
          
          // Сравнение с лучшей карточкой
          const maxProsScore = Math.max(...allItems.map(item => (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0)));
          const maxConsScore = Math.max(...allItems.map(item => (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0)));
          
          // Нормализованные значения (0-1)
          const normalizedPros = maxProsScore > 0 ? prosScore / maxProsScore : 0;
          const normalizedCons = maxConsScore > 0 ? consScore / maxConsScore : 0;
          
          // Влияние плюсов (чем больше плюсов по сравнению со средним и чем больше по сравнению с лучшей карточкой, тем выше рейтинг)
          const prosVsAvg = avgProsScore > 0 ? prosScore / avgProsScore : 1;
          const prosBonus = Math.min(1.5, Math.max(0.5, prosVsAvg)); // Ограничиваем влияние
          
          // Влияние минусов (чем меньше минусов по сравнению со средним и чем меньше по сравнению с худшей карточкой, тем выше рейтинг)
          const consVsAvg = avgConsScore > 0 ? consScore / avgConsScore : 0;
          const consBonus = Math.min(1.5, Math.max(0.5, 1 - consVsAvg)); // Ограничиваем влияние
          
          // Комбинированное влияние с учетом веса плюсов/минусов
          prosConsImpact = ((normalizedPros * prosBonus) - (normalizedCons * consBonus)) * RATING_CONFIG.PROS_CONS_WEIGHT * 0.15;
        }
        
        const finalRating = Math.max(0, Math.min(RATING_CONFIG.MAX_RATING, baseRating + priceImpact + prosConsImpact))
        
        return { ...itemData, rating: Math.round(finalRating * 10) / 10 }
      }
      
      const newItem = {
        ...itemData,
        id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
        created_date: new Date().toISOString(),
      }
      
      // Добавляем новый элемент
      const updatedItems = [newItem, ...items.value]
      
      // Пересчитываем рейтинги для всех элементов
      const itemsWithRatings = updatedItems.map(item =>
        calculateRating(item, updatedItems)
      )
      
      // Обновляем состояние
      items.value = itemsWithRatings
      return newItem
    }

    const updateItem = (id, data) => {
      // Импортируем функцию calculateRating из ComparisonCard
      const calculateRating = (itemData, allItems) => {
        const pros = itemData.pros || []
        const cons = itemData.cons || []
        const itemPrice = itemData.price || 0

        // Сумма положительных баллов
        const prosScore = pros.reduce((sum, pro) => sum + (pro.impact || 5), 0)
        
        // Сумма отрицательных баллов
        const consScore = cons.reduce((sum, con) => sum + (con.impact || 5), 0)
        
        // Базовый рейтинг (0-100)
        const maxPossibleScore = Math.max((pros.length + cons.length) * 10, 1)
        const baseRating = ((prosScore - consScore) / maxPossibleScore) * 50 + 50
        
        // Корректировка на цену (чем выше цена, тем меньше рейтинг)
        // Находим максимальную и минимальную цену среди всех товаров для нормализации
        let priceImpact = 0;
        if (allItems.length > 0 && itemPrice > 0) {
          const maxPrice = Math.max(...allItems.map(item => item.price || 0));
          const minPrice = Math.min(...allItems.map(item => item.price || 0));
          if (maxPrice > minPrice) {
            // Чем выше цена по сравнению с минимальной, тем меньше рейтинг
            // Чем ниже цена по сравнению с максимальной, тем выше рейтинг
            const priceRange = maxPrice - minPrice;
            const priceRatio = (maxPrice - itemPrice) / priceRange;
            priceImpact = priceRatio * 20; // Максимальный бонус/штраф 20 баллов
          }
        }
        
        // Учет относительности плюсов/минусов
        let prosConsImpact = 0;
        if (allItems.length > 0) {
          // Средние значения по всем карточкам
          const totalProsScore = allItems.reduce((sum, item) => sum + (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0), 0);
          const totalConsScore = allItems.reduce((sum, item) => sum + (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0), 0);
          const avgProsScore = totalProsScore / allItems.length;
          const avgConsScore = totalConsScore / allItems.length;
          
          // Сравнение с лучшей карточкой
          const maxProsScore = Math.max(...allItems.map(item => (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0)));
          const maxConsScore = Math.max(...allItems.map(item => (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0)));
          
          // Нормализованные значения (0-1)
          const normalizedPros = maxProsScore > 0 ? prosScore / maxProsScore : 0;
          const normalizedCons = maxConsScore > 0 ? consScore / maxConsScore : 0;
          
          // Влияние плюсов (чем больше плюсов по сравнению со средним и чем больше по сравнению с лучшей карточкой, тем выше рейтинг)
          const prosVsAvg = avgProsScore > 0 ? prosScore / avgProsScore : 1;
          const prosBonus = Math.min(1.5, Math.max(0.5, prosVsAvg)); // Ограничиваем влияние
          
          // Влияние минусов (чем меньше минусов по сравнению со средним и чем меньше по сравнению с худшей карточкой, тем выше рейтинг)
          const consVsAvg = avgConsScore > 0 ? consScore / avgConsScore : 0;
          const consBonus = Math.min(1.5, Math.max(0.5, 1 - consVsAvg)); // Ограничиваем влияние
          
          // Комбинированное влияние
          prosConsImpact = ((normalizedPros * prosBonus) - (normalizedCons * consBonus)) * 15;
        }
        
        const finalRating = Math.max(0, Math.min(100, baseRating + priceImpact + prosConsImpact))
        
        return { ...itemData, rating: Math.round(finalRating * 10) / 10 }
      }
      
      // Обновляем элемент
      const updatedItems = items.value.map(item =>
        item.id === id ? { ...item, ...data } : item
      )
      
      // Пересчитываем рейтинги для всех элементов
      const itemsWithRatings = updatedItems.map(item =>
        calculateRating(item, updatedItems)
      )
      
      // Обновляем состояние
      items.value = itemsWithRatings
    }

    const deleteItem = (id) => {
      items.value = items.value.filter(item => item.id !== id)
    }

    const handleAddLinks = async (linksData) => {
      isLoadingItems.value = true
      
      try {
        for (const linkData of linksData) {
          // Создаем товар с данными из формы
          createItem({
            url: linkData.url,
            title: linkData.title || 'Товар ' + (items.value.length + 1),
            description: '',
            images: linkData.images || [],
            price: linkData.price ? parseFloat(linkData.price) : 0,
            currency: '₽',
            pros: [],
            cons: [],
            rating: 50
          })
        }
      } catch (error) {
        console.error('Ошибка при загрузке товаров:', error)
      } finally {
        isLoadingItems.value = false
      }
    }

    const handleUpdateItem = (id, data) => {
      updateItem(id, data)
    }

    const handleDeleteItem = (id) => {
      deleteItem(id)
    }
    
    const selectText = (event) => {
      event.target.select()
    }
    
    // Сохранение текущей коллекции
    const saveCurrentCollection = async () => {
      try {
        // Получаем ID текущей коллекции из URL
        const collectionId = route.query.collectionId;
        
        // Отправляем коллекцию на сервер
        const token = localStorage.getItem('auth_token');
        if (!token) {
          alert('Для сохранения коллекции необходимо авторизоваться');
          return;
        }
        
        let response;
        if (collectionId) {
          // Обновляем существующую коллекцию
          const collection = {
            id: parseInt(collectionId),
            name: `Коллекция ${new Date().toLocaleDateString()}`,
            items: items.value
          };
          
          response = await fetch('/backend/comparisons/update', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(collection)
          });
        } else {
          // Создаем новую коллекцию
          const collection = {
            name: `Коллекция ${new Date().toLocaleDateString()}`,
            items: items.value
          };
          
          response = await fetch('/backend/comparisons', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(collection)
          });
        }
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        alert('Коллекция успешно сохранена!')
      } catch (error) {
        console.error('Ошибка при сохранении коллекции:', error)
        alert('Не удалось сохранить коллекцию: ' + error.message)
      }
    }
    
    // Переход к странице коллекций
    const goToCollections = () => {
      router.push('/')
    }

    return {
      items,
      isLoadingItems,
      isLoading,
      importData,
      showImportField,
      showCollectionsDialog,
      handleAddLinks,
      handleUpdateItem,
      handleDeleteItem,
      handleExport,
      handleImport,
      triggerImport,
      selectText,
      saveCurrentCollection,
      goToCollections
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