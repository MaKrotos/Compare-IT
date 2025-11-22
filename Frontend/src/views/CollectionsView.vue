<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Заголовок -->
      <div class="text-center mb-12">
        <div class="flex items-center justify-center gap-4 mb-4">
          <div class="p-3 bg-blue-600 rounded-2xl">
            <Folder class="w-8 h-8 text-white" />
          </div>
        </div>
        <h1 class="text-4xl sm:text-5xl font-light tracking-tight text-gray-900 mb-3">
          Мои подборки
        </h1>
        <p class="text-lg text-gray-600 font-light">
          Управляйте своими коллекциями сравнений
        </p>
      </div>

      <!-- Кнопки управления -->
      <div class="mb-8 text-center flex flex-wrap gap-4 justify-center">
        <Button
          @click="createNewCollection"
          class="gap-2 px-6 py-3 text-base font-medium bg-blue-600 hover:bg-blue-700"
        >
          <Plus class="w-5 h-5" />
          Создать подборку
        </Button>
        <Button
          @click="goToComparison"
          variant="outline"
          class="gap-2 px-6 py-3 text-base font-medium"
        >
          <Scale class="w-5 h-5" />
          Сравнить товары
        </Button>
      </div>

      <!-- Список подборок -->
      <div v-if="isLoading" class="flex items-center justify-center py-20">
        <Loader2 class="w-8 h-8 animate-spin text-blue-600" />
      </div>
      <div v-else-if="!collections || collections.length === 0" class="text-center py-20">
        <div>
          <div class="inline-block bg-gray-100 rounded-full mb-4">
            <FolderOpen class="w-12 h-12 text-gray-400" />
          </div>
          <p class="text-gray-500 text-lg font-light">
            У вас пока нет подборок
          </p>
          <p class="text-gray-400 text-sm mt-2">
            Создайте первую подборку, чтобы начать организовывать сравнения
          </p>
        </div>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="collection in collections"
          :key="collection.id"
          class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow duration-200"
        >
          <div class="p-6">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <h3 class="text-lg font-medium text-gray-900 mb-1">
                  {{ collection.name || 'Без названия' }}
                </h3>
                <p class="text-sm text-gray-500">
                  Создана: {{ formatDate(collection.created_at) }}
                </p>
                <div v-if="collection.public_link" class="mt-2 flex items-center text-sm text-gray-500">
                  <Link class="w-4 h-4 mr-1" />
                  <span>Публичная ссылка создана</span>
                </div>
              </div>
              <div class="flex gap-1">
                <Button
                  variant="ghost"
                  size="icon"
                  @click="editCollection(collection)"
                  class="text-gray-500 hover:text-blue-600"
                >
                  <Edit3 class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  @click="togglePublicLink(collection)"
                  class="text-gray-500 hover:text-green-600"
                >
                  <component :is="collection.public_link ? 'Link' : 'Link2'" class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  @click="togglePinCollection(collection)"
                  :class="collection.is_pinned ? 'text-yellow-500 hover:text-yellow-600' : 'text-gray-500 hover:text-yellow-500'"
                >
                  <component :is="collection.is_pinned ? 'Bookmark' : 'BookmarkPlus'" class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  @click="deleteCollection(collection.id)"
                  class="text-gray-500 hover:text-red-600"
                >
                  <Trash2 class="w-4 h-4" />
                </Button>
              </div>
            </div>
            <div class="mt-4 flex items-center justify-between">
              <Button
                @click="openCollection(collection)"
                class="w-full bg-blue-50 hover:bg-blue-100 text-blue-700"
              >
                Открыть
              </Button>
            </div>
          </div>
        </div>
      </div>

      <!-- Модальное окно для создания/редактирования подборки -->
      <Dialog v-model:open="showCollectionDialog" @open-change="onDialogOpenChange">
        <DialogContent class="max-w-md">
          <DialogHeader>
            <DialogTitle>{{ editingCollection ? 'Редактировать подборку' : 'Создать подборку' }}</DialogTitle>
          </DialogHeader>
          <div class="py-4">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Название подборки
                </label>
                <Input
                  v-model="collectionForm.name"
                  placeholder="Введите название подборки"
                  class="w-full"
                  @focus="$event.target.select()"
                  ref="collectionNameInput"
                />
              </div>
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <Button variant="outline" @click="showCollectionDialog = false">Отмена</Button>
            <Button @click="saveCollection">
              {{ editingCollection ? 'Сохранить' : 'Создать' }}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { Folder, FolderOpen, Plus, Loader2, Edit3, Trash2, Bookmark, BookmarkPlus, Scale, Link, Link2 } from 'lucide-vue-next'
import { getComparisonCollections, createComparisonCollection, updateComparisonCollection, deleteComparisonCollection } from '@/utils/comparisons.js'
import { generatePublicLink, removePublicLink } from '@/utils/collections.js'

export default {
  name: 'CollectionsView',
  components: {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    Button,
    Input,
    Folder,
    FolderOpen,
    Plus,
    Loader2,
    Edit3,
    Trash2,
    Bookmark,
    BookmarkPlus,
    Scale,
    Link,
    Link2
  },
  setup() {
    const router = useRouter()
    const collections = ref([])
    const isLoading = ref(true)
    const showCollectionDialog = ref(false)
    const editingCollection = ref(null)
    const collectionForm = ref({
      name: ''
    })
    const collectionNameInput = ref(null)

    // Загрузка подборок при монтировании
    onMounted(async () => {
      await loadCollections()
    })

    // Обработчик открытия/закрытия диалога
    const onDialogOpenChange = (open) => {
      showCollectionDialog.value = open
      if (open) {
        // Фокус на поле ввода при открытии диалога
        setTimeout(() => {
          if (collectionNameInput.value) {
            collectionNameInput.value.focus()
          }
        }, 100)
      }
    }

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
      collectionForm.value = { name: '' }
      showCollectionDialog.value = true
    }

    // Редактирование подборки
    const editCollection = (collection) => {
      editingCollection.value = collection
      collectionForm.value = { name: collection.name }
      showCollectionDialog.value = true
    }

    // Сохранение подборки
    const saveCollection = async () => {
      try {
        if (editingCollection.value) {
          // Обновление существующей подборки
          const updatedCollection = {
            ...editingCollection.value,
            name: collectionForm.value.name
          }
          await updateComparisonCollection(updatedCollection)
        } else {
          // Создание новой подборки
          const newCollection = {
            name: collectionForm.value.name,
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

    // Форматирование даты
    const formatDate = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString)
      return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      })
    }

    return {
      collections,
      isLoading,
      showCollectionDialog,
      editingCollection,
      collectionForm,
      loadCollections,
      createNewCollection,
      editCollection,
      saveCollection,
      deleteCollection,
      togglePinCollection,
      togglePublicLink,
      openCollection,
      goToComparison,
      formatDate
    }
  }
}
</script>

<style scoped>
/* Дополнительные стили при необходимости */
</style>