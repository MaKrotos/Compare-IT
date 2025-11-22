<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50/30 to-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Заголовок -->
      <div class="mb-8">
        <Button
          @click="$router.back()"
          variant="outline"
          class="mb-4 gap-2"
        >
          <ArrowLeft class="w-4 h-4" />
          Назад к подборкам
        </Button>
        
        <div class="flex items-center justify-between">
          <div>
            <h1 v-if="!isEditingName" class="text-3xl font-light tracking-tight text-gray-900 mb-2 flex items-center gap-2">
              {{ collection.name || 'Подборка' }}
              <Button @click="startEditingName" variant="ghost" size="icon" class="ml-2">
                <Edit3 class="w-4 h-4" />
              </Button>
            </h1>
            <div v-else class="flex items-center gap-2">
              <Input
                v-model="newCollectionName"
                placeholder="Введите новое название"
                class="text-3xl font-light"
                @keyup.enter="saveCollectionName"
                @blur="saveCollectionName"
                ref="nameInput"
              />
              <Button @click="saveCollectionName" variant="ghost" size="icon">
                <Check class="w-4 h-4" />
              </Button>
            </div>
            <p class="text-gray-600">
              {{ collection.items?.length || 0 }} товаров в подборке
            </p>
          </div>
          <div class="flex gap-2">
            <Button
              v-if="!collection.public_link"
              @click="generatePublicLink"
              variant="outline"
              class="gap-2"
            >
              <Link class="w-4 h-4" />
              Создать публичную ссылку
            </Button>
            <Button
              v-else
              @click="copyPublicLink"
              variant="outline"
              class="gap-2"
            >
              <Link class="w-4 h-4" />
              Копировать ссылку
            </Button>
            <Button
              v-if="collection.public_link"
              @click="removePublicLink"
              variant="outline"
              class="gap-2 text-red-600 hover:text-red-700 border-red-200 hover:border-red-300"
            >
              <Unlink class="w-4 h-4" />
              Удалить ссылку
            </Button>
            <Button
              @click="deleteCollection"
              variant="outline"
              class="gap-2 text-red-600 hover:text-red-700 border-red-200 hover:border-red-300"
            >
              <Trash2 class="w-4 h-4" />
              Удалить
            </Button>
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
          <p class="text-gray-400 text-sm mt-2">
            Добавьте товары для начала сравнения
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
import { ref, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import ComparisonCard from '@/components/comparison/ComparisonCard.vue'
import { ArrowLeft, Edit3, Trash2, Loader2, Package, Check } from 'lucide-vue-next'
import { getComparisonCollectionById, updateComparisonCollection, deleteComparisonCollection } from '@/utils/comparisons.js'
import { generatePublicLink, removePublicLink } from '@/utils/publicCollections.js'
import { Link, Unlink } from 'lucide-vue-next'

export default {
  name: 'CollectionDetailView',
  components: {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    Button,
    Input,
    ComparisonCard,
    ArrowLeft,
    Edit3,
    Trash2,
    Loader2,
    Package,
    Check,
    Link,
    Unlink
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const collection = ref({})
    const isLoading = ref(true)
    const showEditDialog = ref(false)
    const editForm = ref({
      name: ''
    })
    const isEditingName = ref(false)
    const newCollectionName = ref('')
    const nameInput = ref(null)

    // Загрузка подборки при монтировании
    onMounted(async () => {
      await loadCollection()
    })
    
    // Загрузка подборки при изменении параметров маршрута
    watch(() => route.params.id, async (newId, oldId) => {
      if (newId !== oldId) {
        await loadCollection()
      }
    })

    // Загрузка подборки
    const loadCollection = async () => {
      try {
        isLoading.value = true
        const collectionId = route.params.id
        const data = await getComparisonCollectionById(collectionId)
        collection.value = data
      } catch (error) {
        console.error('Ошибка при загрузке подборки:', error)
      } finally {
        isLoading.value = false
      }
    }

    // Начало редактирования названия
    const startEditingName = () => {
      isEditingName.value = true
      newCollectionName.value = collection.value.name || ''
      nextTick(() => {
        if (nameInput.value) {
          nameInput.value.focus()
        }
      })
    }

    // Сохранение названия подборки
    const saveCollectionName = async () => {
      if (!newCollectionName.value.trim()) {
        isEditingName.value = false
        return
      }

      try {
        const updatedCollection = {
          ...collection.value,
          name: newCollectionName.value
        }
        await updateComparisonCollection(updatedCollection)
        collection.value.name = newCollectionName.value
        isEditingName.value = false
      } catch (error) {
        console.error('Ошибка при сохранении названия подборки:', error)
        alert('Не удалось сохранить название подборки: ' + error.message)
      }
    }

    // Удаление подборки
    const deleteCollectionHandler = async () => {
      try {
        await deleteComparisonCollection(collection.value.id)
        router.push('/collections')
      } catch (error) {
        console.error('Ошибка при удалении подборки:', error)
        alert('Не удалось удалить подборку: ' + error.message)
      }
    }

    // Подтверждение удаления подборки
    const deleteCollection = () => {
    // Генерация публичной ссылки
    const generatePublicLinkHandler = async () => {
      try {
        const publicLink = await generatePublicLink(collection.value.id)
        collection.value.public_link = publicLink
        alert('Публичная ссылка создана успешно!')
      } catch (error) {
        console.error('Ошибка при создании публичной ссылки:', error)
        alert('Не удалось создать публичную ссылку: ' + error.message)
      }
    }

    // Копирование публичной ссылки
    const copyPublicLink = async () => {
      try {
        const publicUrl = `${window.location.origin}/public-collection/${collection.value.public_link}`
        await navigator.clipboard.writeText(publicUrl)
        alert('Ссылка скопирована в буфер обмена!')
      } catch (error) {
        console.error('Ошибка при копировании ссылки:', error)
        alert('Не удалось скопировать ссылку: ' + error.message)
      }
    }

    // Удаление публичной ссылки
    const removePublicLinkHandler = async () => {
      try {
        await removePublicLink(collection.value.id)
        collection.value.public_link = null
        alert('Публичная ссылка удалена успешно!')
      } catch (error) {
        console.error('Ошибка при удалении публичной ссылки:', error)
        alert('Не удалось удалить публичную ссылку: ' + error.message)
      }
    }

    return {
      collection,
      isLoading,
      showEditDialog,
      editForm,
      isEditingName,
      newCollectionName,
      nameInput,
      loadCollection,
      startEditingName,
      saveCollectionName,
      deleteCollection,
      handleUpdateItem,
      handleDeleteItem,
      generatePublicLink: generatePublicLinkHandler,
      copyPublicLink,
      removePublicLink: removePublicLinkHandler
    }
      if (confirm('Вы уверены, что хотите удалить эту подборку?')) {
        deleteCollectionHandler()
      }
    }

    // Обновление элемента
    const handleUpdateItem = (id, data) => {
      // В данной реализации мы не обновляем элементы в подборке напрямую
      // Это можно реализовать при необходимости
      console.log('Обновление элемента:', id, data)
    }

    // Удаление элемента
    const handleDeleteItem = (id) => {
      // В данной реализации мы не удаляем элементы из подборки напрямую
      // Это можно реализовать при необходимости
      console.log('Удаление элемента:', id)
    }

    return {
      collection,
      isLoading,
      showEditDialog,
      editForm,
      isEditingName,
      newCollectionName,
      nameInput,
      loadCollection,
      startEditingName,
      saveCollectionName,
      deleteCollection,
      handleUpdateItem,
      handleDeleteItem
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