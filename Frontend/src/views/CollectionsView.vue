<template>
  <div class="collections-view">
    <div class="collections-container">
      <!-- Заголовок -->
      <div class="collections-header">
        <div class="collections-logo">
          <div class="collections-logo-icon">
            <Folder class="collections-logo-folder" />
          </div>
        </div>
        <h1 class="collections-title">
          Мои подборки
        </h1>
        <p class="collections-subtitle">
          Управляйте своими коллекциями сравнений
        </p>
      </div>

      <!-- Кнопки управления -->
      <div class="collections-actions">
        <Button
          @click="createNewCollection"
          class="collections-action-button collections-action-create"
        >
          <Plus class="collections-action-icon" />
          Создать подборку
        </Button>
        <Button
          @click="goToComparison"
          variant="outline"
          class="collections-action-button collections-action-compare"
        >
          <Scale class="collections-action-icon" />
          Сравнить товары
        </Button>
      </div>

      <!-- Список подборок -->
      <div v-if="isLoading" class="collections-loading">
        <Loader2 class="collections-loading-icon" />
      </div>
      <div v-else-if="!collections || collections.length === 0" class="collections-empty">
        <div class="collections-empty-content">
          <div class="collections-empty-icon">
            <FolderOpen class="collections-empty-folder" />
          </div>
          <p class="collections-empty-title">
            У вас пока нет подборок
          </p>
          <p class="collections-empty-description">
            Создайте первую подборку, чтобы начать организовывать сравнения
          </p>
        </div>
      </div>
      <div v-else class="collections-grid">
        <CollectionCard
          v-for="collection in collections"
          :key="collection.id"
          :collection="collection"
          @edit="editCollection"
          @delete="deleteCollection"
          @toggle-public-link="togglePublicLink"
          @toggle-pin="togglePinCollection"
          @open="openCollection"
        />
      </div>

      <!-- Модальное окно для создания/редактирования подборки -->
      <CollectionDialog
        v-model:open="showCollectionDialog"
        :editing-collection="editingCollection"
        @save="saveCollection"
      />
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import CollectionCard from '@/components/collections/CollectionCard.vue'
import CollectionDialog from '@/components/collections/CollectionDialog.vue'
import Button from '@/components/ui/Button.vue'
import { Folder, FolderOpen, Plus, Loader2, Scale } from 'lucide-vue-next'
import { getComparisonCollections, createComparisonCollection, updateComparisonCollection, deleteComparisonCollection } from '@/utils/comparisons.js'
import { generatePublicLink, removePublicLink } from '@/utils/collections.js'

export default {
  name: 'CollectionsView',
  components: {
    CollectionCard,
    CollectionDialog,
    Button,
    Folder,
    FolderOpen,
    Plus,
    Loader2,
    Scale
  },
  setup() {
    const router = useRouter()
    const collections = ref([])
    const isLoading = ref(true)
    const showCollectionDialog = ref(false)
    const editingCollection = ref(null)

    // Загрузка подборок при монтировании
    onMounted(async () => {
      await loadCollections()
    })

    // Загрузка подборок
    const loadCollections = async () => {
      try {
        isLoading.value = true
        const data = await getComparisonCollections()
        // Сортируем коллекции: закрепленные первыми
        collections.value = data.sort((a, b) => {
          // Если обе закреплены или обе не закреплены, сортируем по дате создания (новые первыми)
          if (a.is_pinned === b.is_pinned) {
            return new Date(b.created_at) - new Date(a.created_at);
          }
          // Закрепленные коллекции идут первыми
          return b.is_pinned - a.is_pinned;
        });
      } catch (error) {
        console.error('Ошибка при загрузке подборок:', error)
      } finally {
        isLoading.value = false
      }
    }

    // Создание новой подборки
    const createNewCollection = () => {
      editingCollection.value = null
      showCollectionDialog.value = true
    }

    // Редактирование подборки
    const editCollection = (collection) => {
      editingCollection.value = collection
      showCollectionDialog.value = true
    }

    // Сохранение подборки
    const saveCollection = async (collectionData) => {
      try {
        if (editingCollection.value) {
          // Обновление существующей подборки
          const updatedCollection = {
            ...editingCollection.value,
            name: collectionData.name
          }
          await updateComparisonCollection(updatedCollection)
        } else {
          // Создание новой подборки
          const newCollection = {
            name: collectionData.name,
            items: []
          }
          await createComparisonCollection(newCollection)
        }
        
        // Закрываем диалог и перезагружаем список
        showCollectionDialog.value = false
        await loadCollections()
      } catch (error) {
        console.error('Ошибка при сохранении подборки:', error)
        alert('Не удалось сохранить подборку: ' + error.message)
      }
    }

    // Удаление подборки
    const removeCollection = async (collectionId) => {
      try {
        await deleteComparisonCollection(collectionId)
        await loadCollections()
      } catch (error) {
        console.error('Ошибка при удалении подборки:', error)
        alert('Не удалось удалить подборку: ' + error.message)
      }
    }

    // Подтверждение удаления подборки
    const deleteCollection = (collectionId) => {
      if (confirm('Вы уверены, что хотите удалить эту подборку?')) {
        removeCollection(collectionId)
      }
    }
    
    // Закрепление/открепление подборки
    const togglePinCollection = async (collection) => {
      try {
        const updatedCollection = {
          ...collection,
          is_pinned: !collection.is_pinned
        }
        await updateComparisonCollection(updatedCollection)
        // Обновляем состояние в списке
        const index = collections.value.findIndex(c => c.id === collection.id)
        if (index !== -1) {
          collections.value[index] = updatedCollection
        }
      } catch (error) {
        console.error('Ошибка при закреплении подборки:', error)
        alert('Не удалось изменить статус закрепления подборки: ' + error.message)
      }
    }

    // Включение/выключение публичной ссылки
    const togglePublicLink = async (collection) => {
      try {
        let updatedCollection;
        if (collection.public_link) {
          // Удаляем публичную ссылку
          await removePublicLink(collection.id);
          updatedCollection = {
            ...collection,
            public_link: null
          };
          // Обновляем состояние в списке
          const index = collections.value.findIndex(c => c.id === collection.id);
          if (index !== -1) {
            collections.value[index] = updatedCollection;
          }
          alert('Публичная ссылка удалена');
        } else {
          // Создаем публичную ссылку
          const publicLink = await generatePublicLink(collection.id);
          updatedCollection = {
            ...collection,
            public_link: publicLink
          };
          // Обновляем состояние в списке
          const index = collections.value.findIndex(c => c.id === collection.id);
          if (index !== -1) {
            collections.value[index] = updatedCollection;
          }
          alert('Публичная ссылка создана: ' + publicLink);
        }
      } catch (error) {
        console.error('Ошибка при работе с публичной ссылкой:', error);
        alert('Не удалось изменить статус публичной ссылки: ' + error.message);
      }
    };

    // Открытие подборки
    const openCollection = (collection) => {
      // Переход к странице сравнений этой коллекции
      router.push(`/compare?collectionId=${collection.id}`)
    }

    // Переход к сравнению товаров
    const goToComparison = () => {
      router.push('/compare')
    }

    return {
      collections,
      isLoading,
      showCollectionDialog,
      editingCollection,
      loadCollections,
      createNewCollection,
      editCollection,
      saveCollection,
      deleteCollection,
      togglePinCollection,
      togglePublicLink,
      openCollection,
      goToComparison
    }
  }
}
</script>

<style scoped>
/* Стили для компонента CollectionsView */
.collections-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f9fafb 0%, #eff6ff 30%, #f9fafb 100%);
}

.collections-container {
  max-width: 80rem;
  margin: 0 auto;
  padding: 3rem 1rem;
}

/* Заголовок */
.collections-header {
  text-align: center;
  margin-bottom: 3rem;
}

.collections-logo {
  margin-bottom: 1rem;
}

.collections-logo-icon {
  display: inline-flex;
  padding: 0.75rem;
  background-color: #2563eb;
  border-radius: 1rem;
}

.collections-logo-folder {
  width: 2rem;
  height: 2rem;
  color: white;
}

.collections-title {
  font-size: 3rem;
  font-weight: 300;
  letter-spacing: -0.025em;
  color: #111827;
  margin-bottom: 0.75rem;
}

@media (max-width: 640px) {
  .collections-title {
    font-size: 2.25rem;
  }
}

.collections-subtitle {
  font-size: 1.125rem;
  color: #4b5563;
  font-weight: 300;
}

/* Кнопки управления */
.collections-actions {
  margin-bottom: 2rem;
  text-align: center;
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  justify-content: center;
}

.collections-action-button {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.collections-action-create {
  background-color: #2563eb;
  color: white;
  border: none;
}

.collections-action-create:hover {
  background-color: #1d4ed8;
}

.collections-action-compare {
  background-color: white;
  color: #374151;
  border: 1px solid #d1d5db;
}

.collections-action-compare:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
}

.collections-action-icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* Состояние загрузки */
.collections-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5rem 0;
}

.collections-loading-icon {
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
.collections-empty {
  text-align: center;
  padding: 5rem 0;
}

.collections-empty-content {
  display: inline-block;
}

.collections-empty-icon {
  display: inline-flex;
  background-color: #f3f4f6;
  border-radius: 50%;
  padding: 1rem;
  margin-bottom: 1rem;
}

.collections-empty-folder {
  width: 3rem;
  height: 3rem;
  color: #9ca3af;
}

.collections-empty-title {
  color: #6b7280;
  font-size: 1.125rem;
  font-weight: 300;
}

.collections-empty-description {
  color: #9ca3af;
  font-size: 0.875rem;
  margin-top: 0.5rem;
}

/* Сетка подборок */
.collections-grid {
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 1.5rem;
}

@media (min-width: 768px) {
  .collections-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .collections-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* Карточка подборки */
.collection-card {
  background-color: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  transition: box-shadow 200ms ease;
}

.collection-card:hover {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1);
}

.collection-card-content {
  padding: 1.5rem;
}

.collection-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.collection-card-info {
  flex: 1;
}

.collection-card-title {
  font-size: 1.125rem;
  line-height: 1.75rem;
  font-weight: 500;
  color: #111827;
  margin-bottom: 0.25rem;
}

.collection-card-date {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #6b7280;
}

.collection-card-link-indicator {
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #6b7280;
}

.collection-card-link-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.25rem;
  color: #6b7280;
}

.collection-card-link-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.collection-card-actions {
  display: flex;
  gap: 0.25rem;
}

.collection-card-action {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
  border: none;
  cursor: pointer;
  transition: background-color 150ms ease;
}

.collection-card-edit {
  color: #6b7280;
}

.collection-card-edit:hover {
  color: #2563eb;
  background-color: #dbeafe;
}

.collection-card-link {
  color: #6b7280;
}

.collection-card-link:hover {
  color: #16a34a;
  background-color: #dcfce7;
}

.collection-card-pin {
  color: #6b7280;
}

.collection-card-pinned {
  color: #f59e0b;
}

.collection-card-unpinned:hover {
  color: #f59e0b;
  background-color: #fef3c7;
}

.collection-card-pinned:hover {
  color: #d97706;
  background-color: #fef3c7;
}

.collection-card-delete {
  color: #6b7280;
}

.collection-card-delete:hover {
  color: #dc2626;
  background-color: #fef2f2;
}

.collection-card-action-icon {
  width: 1rem;
  height: 1rem;
}

.collection-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.collection-card-open-button {
  width: 100%;
  background-color: #dbeafe;
  color: #1e40af;
  border: none;
  border-radius: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 150ms ease;
}

.collection-card-open-button:hover {
  background-color: #bfdbfe;
}

/* Модальное окно */
.collection-dialog {
  max-width: 28rem;
}

.collection-dialog-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  color: #111827;
}

.collection-dialog-content {
  padding: 1rem 0;
}

.collection-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.collection-form-label {
  display: block;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.25rem;
}

.collection-form-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.collection-form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.collection-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.collection-dialog-cancel {
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

.collection-dialog-cancel:hover {
  background-color: #f9fafb;
}

.collection-dialog-save {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  background-color: #2563eb;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
}

.collection-dialog-save:hover {
  background-color: #1d4ed8;
}
</style>