<template>
  <div class="comparison-view">
    <div class="comparison-container">
      <!-- Заголовок -->
      <div class="comparison-header">
        <div class="comparison-logo">
          <div class="comparison-logo-icon">
            <Scale class="comparison-logo-scale" />
          </div>
        </div>
        <h1 class="comparison-title">
          Сравнение товаров
        </h1>
        <p class="comparison-subtitle">
          Добавьте ссылки на товары и сравните их плюсы и минусы
        </p>
      </div>
    
    <!-- Кнопки экспорта/импорта -->
    <div class="comparison-actions">
      <div class="comparison-action-buttons">
        <Button variant="outline" class="comparison-action-button" @click="handleExport">
          <Download class="comparison-action-icon" />
          Экспорт
        </Button>
        
        <Button variant="outline" class="comparison-action-button" @click="triggerImport">
          <Upload class="comparison-action-icon" />
          Импорт
        </Button>
        
        <!-- Кнопки для работы с коллекциями -->
        <Button variant="outline" class="comparison-action-button" @click="goToCollections">
          <FolderOpen class="comparison-action-icon" />
          Открыть
        </Button>
        
        <!-- Кнопки для работы с публичными ссылками -->
        <Button
          v-if="collectionId && isCollectionLoaded"
          variant="outline"
          class="comparison-action-button"
          @click="togglePublicLink"
        >
          <component :is="publicLink ? 'Link' : 'Link2'" class="comparison-action-icon" />
          {{ publicLink ? 'Удалить ссылку' : 'Создать ссылку' }}
        </Button>
        
        <Button
          v-if="publicLink"
          variant="outline"
          class="comparison-action-button"
          @click="copyPublicLink"
        >
          <Copy class="comparison-action-icon" />
          Копировать ссылку
        </Button>
      </div>
      
      <!-- Отображение публичной ссылки -->
      <div v-if="publicLink" class="comparison-public-link">
        <div class="comparison-public-link-content">
          <div class="comparison-public-link-info">
            <p class="comparison-public-link-title">Публичная ссылка для просмотра:</p>
            <p class="comparison-public-link-url">{{ fullPublicLink }}</p>
          </div>
          <Button variant="ghost" size="icon" @click="copyPublicLink" class="comparison-public-link-copy">
            <Copy class="comparison-public-link-copy-icon" />
          </Button>
        </div>
      </div>
      
      <!-- Поле для ввода данных импорта -->
      <div v-if="showImportField" class="comparison-import-field">
        <Textarea
          v-model="importData"
          placeholder="Вставьте JSON данные для импорта..."
          class="comparison-import-textarea"
        />
        <div class="comparison-import-buttons">
          <Button @click="handleImport" class="comparison-import-button">
            Импортировать
          </Button>
          <Button variant="outline" @click="showImportField = false" class="comparison-import-cancel">
            Отмена
          </Button>
        </div>
      </div>
      
      <!-- Модальное окно для выбора коллекции -->
      <Dialog v-model:open="showCollectionsDialog">
        <DialogContent class="comparison-collections-dialog">
          <DialogHeader>
            <DialogTitle class="comparison-collections-title">Сохраненные коллекции</DialogTitle>
          </DialogHeader>
          <div class="comparison-collections-content">
            <div class="comparison-collections-empty">
              Для работы с коллекциями перейдите на страницу коллекций
            </div>
          </div>
          <div class="comparison-collections-footer">
            <Button variant="outline" @click="showCollectionsDialog = false" class="comparison-collections-close">
              Закрыть
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
    
    <!-- Форма добавления -->
    <div class="comparison-add-form-container">
      <AddLinkForm :on-add-links="handleAddLinks" :is-loading="isLoadingItems" />
    </div>
    
    <!-- Элементы управления рейтингами -->
    <div class="comparison-rating-controls">
      <RatingControls
        :initial-price-weight="ratingWeights.priceRatingWeight"
        :initial-pros-cons-weight="ratingWeights.prosConsRatingWeight"
        :on-save="handleRatingWeightsSave"
      />
    </div>

    <!-- Список товаров -->
    <div v-if="isLoading" class="comparison-loading">
      <Loader2 class="comparison-loading-icon" />
    </div>
    <div v-else-if="items.length === 0" class="comparison-empty">
      <div class="comparison-empty-content">
        <div class="comparison-empty-icon">
          <Scale class="comparison-empty-scale" />
        </div>
        <p class="comparison-empty-text">
          Добавьте товары для начала сравнения
        </p>
      </div>
    </div>
    <div v-else class="comparison-items">
      <transition-group name="comparison-list" tag="div" class="comparison-items-grid">
        <div
          v-for="(item, index) in items"
          :key="item.id"
          class="comparison-item-container"
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
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'
import Button from '@/components/ui/Button.vue'
import Textarea from '@/components/ui/Textarea.vue'
import AddLinkForm from '@/components/comparison/AddLinkForm.vue'
import ComparisonCard from '@/components/comparison/ComparisonCard.vue'
import { Scale, Loader2, Download, Upload, FolderOpen, Link, Link2, Copy } from 'lucide-vue-next'
import { getComparisonCollectionById } from '@/utils/comparisons.js'
import { getToken } from '@/utils/auth'
import { generatePublicLink, removePublicLink } from '@/utils/collections.js'
import RatingControls from '@/components/comparison/RatingControls.vue'

export default {
  name: 'ComparisonView',
  components: {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    Button,
    Textarea,
    AddLinkForm,
    ComparisonCard,
    RatingControls,
    Scale,
    Loader2,
    Download,
    Upload,
    FolderOpen,
    Link,
    Link2,
    Copy
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
    const collectionId = ref(null)
    const publicLink = ref(null)
    const ratingWeights = ref({
      priceRatingWeight: 20,
      prosConsRatingWeight: 80
    })
    let saveTimeout = null;
    let isCollectionLoaded = false;
    let originalCollectionName = '';

    // Инициализация
    onMounted(async () => {
      // Проверяем, есть ли параметр collectionId в URL
      collectionId.value = route.query.collectionId;
      if (collectionId.value) {
        // Устанавливаем флаг до загрузки коллекции, чтобы избежать автосохранения
        isCollectionLoaded = true;
        // Загружаем коллекцию с бэкенда
        try {
          const collection = await getComparisonCollectionById(collectionId.value);
          items.value = collection.items || [];
          publicLink.value = collection.public_link || null;
          // Инициализируем настройки рейтингов из данных коллекции
          ratingWeights.value = {
            priceRatingWeight: collection.price_rating_weight || 20,
            prosConsRatingWeight: collection.pros_cons_rating_weight || 80
          };
          // Сохраняем оригинальное название коллекции
          originalCollectionName = collection.name || `Коллекция ${new Date().toLocaleDateString()}`;
        } catch (error) {
          console.error('Ошибка при загрузке коллекции:', error);
          // Устанавливаем имя по умолчанию в случае ошибки
          originalCollectionName = `Коллекция ${new Date().toLocaleDateString()}`;
        }
      } else {
        // Для новой коллекции устанавливаем имя по умолчанию
        originalCollectionName = `Коллекция ${new Date().toLocaleDateString()}`;
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
          
          if (collectionId && isCollectionLoaded) {
            // Отправляем коллекцию на сервер
            const token = getToken();
            if (!token) {
              console.log('Пользователь не авторизован, автосохранение отключено');
              return;
            }
            
            const collection = {
              id: parseInt(collectionId),
              name: originalCollectionName,
              items: items.value,
              price_rating_weight: ratingWeights.value.priceRatingWeight,
              pros_cons_rating_weight: ratingWeights.value.prosConsRatingWeight
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
      // Проверяем, есть ли ID коллекции в URL и была ли коллекция загружена
      if (collectionId.value && isCollectionLoaded) {
        autoSaveCollection();
      }
    }, { deep: true })

    // Включение/выключение публичной ссылки
    const togglePublicLink = async () => {
      try {
        if (publicLink.value) {
          // Удаляем публичную ссылку
          await removePublicLink(collectionId.value);
          publicLink.value = null;
          alert('Публичная ссылка удалена');
        } else {
          // Создаем публичную ссылку
          const link = await generatePublicLink(collectionId.value);
          publicLink.value = link;
          alert('Публичная ссылка создана');
        }
      } catch (error) {
        console.error('Ошибка при работе с публичной ссылкой:', error);
        alert('Не удалось изменить статус публичной ссылки: ' + error.message);
      }
    };

    // Копирование публичной ссылки в буфер обмена
    const copyPublicLink = async () => {
      try {
        await navigator.clipboard.writeText(fullPublicLink.value);
        alert('Ссылка скопирована в буфер обмена!');
      } catch (error) {
        console.error('Ошибка при копировании ссылки:', error);
        // Альтернативный способ копирования
        const textArea = document.createElement('textarea');
        textArea.value = fullPublicLink.value;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
        alert('Ссылка скопирована в буфер обмена!');
      }
    };

    // Полная публичная ссылка
    const fullPublicLink = computed(() => {
      if (!publicLink.value) return '';
      // В реальном приложении домен будет браться из конфигурации
      const domain = window.location.origin;
      return `${domain}/public/${publicLink.value}`;
    });

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
    
    // Обработчик сохранения весов рейтингов
    const handleRatingWeightsSave = (weights) => {
      ratingWeights.value = weights;
      // Пересчитываем рейтинги с новыми весами
      recalculateAllRatings();
      // Сохраняем коллекцию с новыми настройками рейтингов
      saveCurrentCollection();
    };

    // Пересчет всех рейтингов
    const recalculateAllRatings = () => {
      const updatedItems = items.value.map(item =>
        calculateRating(item, items.value)
      );
      items.value = updatedItems;
    };

    const createItem = (itemData) => {
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

    const calculateRating = (itemData, allItems) => {
      // Конфигурационные параметры для настройки веса рейтинга
      const RATING_CONFIG = {
        // Максимальный рейтинг
        MAX_RATING: 100,
        // Вес рейтинга по плюсам/минусам (в процентах от максимального рейтинга)
        PROS_CONS_WEIGHT: ratingWeights.value.prosConsRatingWeight,
        // Вес рейтинга по цене (в процентах от максимального рейтинга)
        PRICE_WEIGHT: ratingWeights.value.priceRatingWeight
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

    const updateItem = (id, data) => {
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
          const newItem = {
            url: linkData.url,
            title: linkData.title || 'Товар ' + (items.value.length + 1),
            description: '',
            images: linkData.images || [],
            price: linkData.price ? parseFloat(linkData.price) : 0,
            currency: '₽',
            pros: [],
            cons: [],
            rating: 50
          };
          
          // Добавляем новый элемент и пересчитываем рейтинги
          createItem(newItem);
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
    
    // Сохранение текущей коллекции
    const saveCurrentCollection = async () => {
      try {
        // Получаем ID текущей коллекции из URL
        const collectionId = route.query.collectionId;
        
        // Отправляем коллекцию на сервер
        const token = getToken();
        if (!token) {
          alert('Для сохранения коллекции необходимо авторизоваться');
          return;
        }
        
        let response;
        if (collectionId) {
          // Обновляем существующую коллекцию
          const collection = {
            id: parseInt(collectionId),
            name: originalCollectionName,
            items: items.value,
            price_rating_weight: ratingWeights.value.priceRatingWeight,
            pros_cons_rating_weight: ratingWeights.value.prosConsRatingWeight
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
            items: items.value,
            price_rating_weight: ratingWeights.value.priceRatingWeight,
            pros_cons_rating_weight: ratingWeights.value.prosConsRatingWeight
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
      saveCurrentCollection,
      goToCollections,
      togglePublicLink,
      copyPublicLink,
      fullPublicLink,
      collectionId,
      publicLink,
      ratingWeights,
      handleRatingWeightsSave
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ComparisonView */
.comparison-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f9fafb 0%, #eff6ff 30%, #f9fafb 100%);
}

.comparison-container {
  margin: 0 auto;
  padding: 3rem 1rem;
}

/* Заголовок */
.comparison-header {
  text-align: center;
  margin-bottom: 3rem;
}

.comparison-logo {
  margin-bottom: 1rem;
}

.comparison-logo-icon {
  display: inline-flex;
  padding: 0.75rem;
  background-color: #2563eb;
  border-radius: 1rem;
}

.comparison-logo-scale {
  width: 2rem;
  height: 2rem;
  color: white;
}

.comparison-title {
  font-size: 3rem;
  font-weight: 300;
  letter-spacing: -0.025em;
  color: #111827;
  margin-bottom: 0.75rem;
}

@media (max-width: 640px) {
  .comparison-title {
    font-size: 2.25rem;
  }
}

.comparison-subtitle {
  font-size: 1.125rem;
  color: #4b5563;
  font-weight: 300;
}

/* Кнопки действий */
.comparison-actions {
  margin-bottom: 3rem;
}

.comparison-action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  justify-content: center;
}

.comparison-action-button {
  padding: 0.5rem 1rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  background-color: white;
  color: #374151;
  border: 1px solid #d1d5db;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-action-button:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

.comparison-action-icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* Публичная ссылка */
.comparison-public-link {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #dbeafe;
  border-radius: 0.5rem;
  border: 1px solid #bfdbfe;
}

.comparison-public-link-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.comparison-public-link-info {
  flex: 1;
}

.comparison-public-link-title {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #1e40af;
  margin-bottom: 0.25rem;
}

.comparison-public-link-url {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #2563eb;
  word-break: break-all;
}

.comparison-public-link-copy {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  color: #2563eb;
  background-color: transparent;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 150ms ease;
}

.comparison-public-link-copy:hover {
  background-color: rgba(37, 99, 235, 0.1);
}

.comparison-public-link-copy-icon {
  width: 1rem;
  height: 1rem;
}

/* Поле импорта */
.comparison-import-field {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f3f4f6;
  border-radius: 0.5rem;
}

.comparison-import-textarea {
  height: 8rem;
  font-family: monospace;
  font-size: 0.75rem;
  line-height: 1rem;
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  background-color: white;
}

.comparison-import-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.comparison-import-buttons {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.comparison-import-button,
.comparison-import-cancel {
  flex: 1;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  border-radius: 0.375rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.comparison-import-button {
  background-color: #2563eb;
  color: white;
  border: none;
}

.comparison-import-button:hover {
  background-color: #1d4ed8;
}

.comparison-import-cancel {
  background-color: white;
  color: #374151;
  border: 1px solid #d1d5db;
}

.comparison-import-cancel:hover {
  background-color: #f9fafb;
}

/* Диалог коллекций */
.comparison-collections-dialog {
  max-width: 28rem;
}

.comparison-collections-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  color: #111827;
}

.comparison-collections-content {
  padding: 1rem 0;
}

.comparison-collections-empty {
  text-align: center;
  padding: 1rem 0;
  color: #6b7280;
}

.comparison-collections-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.comparison-collections-close {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  background-color: white;
  color: #374151;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  cursor: pointer;
}

.comparison-collections-close:hover {
  background-color: #f9fafb;
}

/* Форма добавления */
.comparison-add-form-container {
  margin-bottom: 3rem;
}

/* Элементы управления рейтингами */
.comparison-rating-controls {
  margin-bottom: 2rem;
}

/* Состояние загрузки */
.comparison-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5rem 0;
}

.comparison-loading-icon {
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
.comparison-empty {
  text-align: center;
  padding: 5rem 0;
}

.comparison-empty-content {
  display: inline-block;
}

.comparison-empty-icon {
  display: inline-flex;
  background-color: #f3f4f6;
  border-radius: 50%;
  padding: 1rem;
  margin-bottom: 1rem;
}

.comparison-empty-scale {
  width: 3rem;
  height: 3rem;
  color: #9ca3af;
}

.comparison-empty-text {
  color: #6b7280;
  font-size: 1.125rem;
  font-weight: 300;
}

/* Список товаров */
.comparison-items-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1.5rem;
}

.comparison-item-container {
  width: 100%;
  max-width: 20rem;
}

/* Анимации */
.comparison-list-enter-active,
.comparison-list-leave-active {
  transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.comparison-list-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.comparison-list-leave-active {
  position: absolute;
}
</style>