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
import RatingControls from '@/components/comparison/RatingControls.vue'
import { useCollection } from '@/composables/useCollection.js'
import { useImportExport } from '@/composables/useImportExport.js'
import { usePublicLink } from '@/composables/usePublicLink.js'
import { useRating } from '@/composables/useRating.js'

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
    
    // Используем composable для управления коллекцией
    const {
      items,
      isLoading,
      collectionId,
      isCollectionLoaded,
      ratingWeights,
      initializeCollection,
      autoSaveCollection,
      saveCurrentCollection
    } = useCollection()
    
    // Используем composable для управления импортом/экспортом
    const { handleExport: exportItems, handleImport: importItems } = useImportExport()
    
    // Используем composable для управления публичными ссылками
    const { 
      publicLink, 
      togglePublicLink: toggleLink, 
      copyPublicLink: copyLink, 
      getFullPublicLink 
    } = usePublicLink()
    
    // Используем composable для управления рейтингами
    const { 
      handleRatingWeightsSave: saveRatingWeights,
      createItem,
      updateItem,
      deleteItem,
      recalculateAllRatings
    } = useRating()
    
    const isLoadingItems = ref(false)
    const importData = ref('')
    const showImportField = ref(false)
    const showCollectionsDialog = ref(false)
    
    // Инициализация
    onMounted(async () => {
      await initializeCollection()
    })
    
    // Отслеживаем изменения в items и запускаем автосохранение
    watch(items, () => {
      // Проверяем, есть ли ID коллекции в URL и была ли коллекция загружена
      if (collectionId.value && isCollectionLoaded.value) {
        autoSaveCollection(items.value, ratingWeights.value)
      }
    }, { deep: true })
    
    // Включение/выключение публичной ссылки
    const togglePublicLink = async () => {
      await toggleLink(collectionId.value)
    }
    
    // Копирование публичной ссылки в буфер обмена
    const copyPublicLink = async () => {
      await copyLink(fullPublicLink.value)
    }
    
    // Полная публичная ссылка
    const fullPublicLink = computed(() => {
      return getFullPublicLink(publicLink.value)
    })
    
    // Функция для экспорта данных
    const handleExport = () => {
      exportItems(items.value)
    }
    
    // Функция для импорта данных
    const handleImport = () => {
      importItems(
        importData.value,
        items.value,
        (newItems) => { items.value = newItems },
        (show) => { showImportField.value = show },
        (data) => { importData.value = data }
      )
    }
    
    const triggerImport = () => {
      showImportField.value = true
    }
    
    // Обработчик сохранения весов рейтингов
    const handleRatingWeightsSave = (weights) => {
      saveRatingWeights(weights, (newWeights) => {
        ratingWeights.value = newWeights
      })
      // Пересчитываем рейтинги с новыми весами
      const updatedItems = recalculateAllRatings(items.value, weights)
      items.value = updatedItems
      // Сохраняем коллекцию с новыми настройками рейтингов
      saveCurrentCollection(items.value, weights)
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
          }
          
          // Добавляем новый элемент и пересчитываем рейтинги
          const { itemsWithRatings } = createItem(newItem, items.value, ratingWeights.value)
          items.value = itemsWithRatings
        }
      } catch (error) {
        console.error('Ошибка при загрузке товаров:', error)
      } finally {
        isLoadingItems.value = false
      }
    }
    
    const handleUpdateItem = (id, data) => {
      const updatedItems = updateItem(id, data, items.value, ratingWeights.value)
      items.value = updatedItems
    }
    
    const handleDeleteItem = (id) => {
      const updatedItems = deleteItem(id, items.value)
      items.value = updatedItems
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