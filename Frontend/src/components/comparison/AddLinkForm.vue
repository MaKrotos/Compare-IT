<template>
  <Card class="add-link-form-container">
    <CardHeader class="add-link-form-header">
      <CardTitle class="add-link-form-title">
        <div class="add-link-form-icon-container">
          <Link2 class="w-5 h-5 text-blue-600" />
        </div>
        Добавить товары для сравнения
      </CardTitle>
      <p class="add-link-form-subtitle">
        Вставьте ссылки на товары (до 10 штук)
      </p>
    </CardHeader>
    <CardContent class="add-link-form-content">
      <div class="add-link-form-space-y-4">
        <transition-group name="list" tag="div" class="add-link-form-space-y-4">
          <div v-for="(link, index) in links" :key="index" class="space-y-3">
            <div class="flex gap-3 items-start">
              <div class="flex-1 space-y-3">
                <!-- URL -->
                <div class="relative add-link-form-input-container">
                  <Input
                    :value="link.url"
                    @input="updateLinkUrl(index, $event.target.value)"
                    :placeholder="`Ссылка на товар ${index + 1}`"
                    :disabled="isLoading"
                  />
                  <Link2 class="w-4 h-4 text-gray-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                </div>

                <!-- Название товара -->
                <div class="relative add-link-form-input-container">
                  <Input
                    :value="link.title"
                    @input="updateLinkTitle(index, $event.target.value)"
                    placeholder="Название товара"
                    :disabled="isLoading"
                  />
                  <Tag class="w-4 h-4 text-gray-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                </div>

                <!-- Цена -->
                <div class="relative add-link-form-price-container">
                  <Input
                    type="number"
                    :value="link.price"
                    @input="updateLinkPrice(index, $event.target.value)"
                    placeholder="Цена, ₽"
                    :disabled="isLoading"
                  />
                  <DollarSign class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
                </div>
                
                <!-- Ссылки на изображения -->
                <div class="relative add-link-form-textarea-container">
                  <Textarea
                    :value="link.imageUrls"
                    @input="updateLinkImageUrls(index, $event.target.value)"
                    @paste="handleImagePaste(index, $event)"
                    placeholder="Ссылки на изображения (по одной на строку) или вставьте изображения (Ctrl+V)"
                    :disabled="isLoading"
                  />
                  <Image class="w-4 h-4 text-gray-400 absolute left-3 top-3" />
                </div>
                
                <!-- Область для drag and drop -->
                <div
                  @dragover.prevent
                  @dragenter="handleDragEnter"
                  @dragleave="handleDragLeave"
                  @drop="event => handleDrop(event, index)"
                  class="add-link-form-drag-area"
                  :class="{ 'add-link-form-drag-area-active': isDragOver }"
                >
                  <div class="add-link-form-drag-content">
                    <Upload class="add-link-form-drag-icon" />
                    <p class="add-link-form-drag-text">Перетащите изображения сюда</p>
                    <p class="add-link-form-drag-subtext">Поддерживаются JPG, PNG, GIF</p>
                  </div>
                </div>
              </div>

              <!-- Кнопка удаления -->
              <Button
                v-if="links.length > 1"
                type="button"
                variant="ghost"
                size="icon"
                @click="removeLink(index)"
                :disabled="isLoading"
              >
                <X class="w-5 h-5" />
              </Button>
            </div>
          </div>
        </transition-group>

        <div class="add-link-form-button-container">
          <Button
            v-if="links.length < 10"
            type="button"
            variant="outline"
            @click="addLinkField"
            class="add-link-form-add-button"
            :disabled="isLoading"
          >
            <Plus class="add-link-form-button-icon" />
            Добавить ссылку
          </Button>
          <Button
            type="button"
            @click="handleAddItems"
            :disabled="isLoading || !links.some(l => l.url.trim())"
            class="add-link-form-submit-button"
          >
            <Loader2 v-if="isLoading" class="add-link-form-button-icon animate-spin" />
            <span v-if="isLoading">Загрузка...</span>
            <span v-else>Загрузить товары</span>
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<script>
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Card from '@/components/ui/Card.vue'
import CardContent from '@/components/ui/CardContent.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardTitle from '@/components/ui/CardTitle.vue'
import { Plus, Link2, Loader2, DollarSign, X, Image, Tag, Upload } from 'lucide-vue-next'
import { extractProductName, getCachedPreviewData } from '@/utils/preview.js'

export default {
  name: 'AddLinkForm',
  components: {
    Button,
    Input,
    Textarea,
    Card,
    CardContent,
    CardHeader,
    CardTitle,
    Plus,
    Link2,
    Loader2,
    DollarSign,
    X,
    Image,
    Tag,
    Upload
  },
  props: {
    onAddLinks: {
      type: Function,
      required: true
    },
    isLoading: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const links = ref([{ url: '', price: '', imageUrls: '', title: '', description: '' }])
    const isDragOver = ref(false)

    /**
     * Обновляет URL в ссылке по индексу
     * @param {number} index - Индекс ссылки
     * @param {string} value - Новое значение URL
     */
    const updateLinkUrl = (index, value) => {
      links.value[index].url = value
      
      // Если URL валиден и заголовок еще не установлен, пытаемся получить имя товара
      if (value && !links.value[index].title) {
        fetchProductName(index, value)
      }
    }

    /**
     * Получает имя товара с сайта
     * @param {number} index - Индекс ссылки
     * @param {string} url - URL сайта
     */
    const fetchProductName = async (index, url) => {
      try {
        // Проверяем, что URL валиден
        new URL(url)
        
        // Показываем индикатор загрузки
        const originalTitle = links.value[index].title
        links.value[index].title = 'Загрузка...'
        
        // Получаем данные предварительного просмотра с кэшированием
        const previewData = await getCachedPreviewData(url)
        
        // Обновляем название товара
        links.value[index].title = previewData.title || 'Без названия'
        
        // Если изображения еще не установлены, устанавливаем изображение из preview
        if (!links.value[index].imageUrls && previewData.image) {
          links.value[index].imageUrls = previewData.image
        }
        
        // Если описание еще не установлено, устанавливаем описание из preview
        if (!links.value[index].description && previewData.description) {
          links.value[index].description = previewData.description
        }
      } catch (error) {
        console.error('Ошибка при получении данных товара:', error)
        links.value[index].title = 'Без названия'
      }
    }

    /**
     * Обновляет цену в ссылке по индексу
     * @param {number} index - Индекс ссылки
     * @param {string} value - Новое значение цены
     */
    const updateLinkPrice = (index, value) => {
      links.value[index].price = value
    }

    /**
     * Обновляет ссылки на изображения в ссылке по индексу
     * @param {number} index - Индекс ссылки
     * @param {string} value - Новое значение ссылок на изображения
     */
    const updateLinkImageUrls = (index, value) => {
      links.value[index].imageUrls = value
    }

    /**
     * Обновляет заголовок в ссылке по индексу
     * @param {number} index - Индекс ссылки
     * @param {string} value - Новое значение заголовка
     */
    const updateLinkTitle = (index, value) => {
      links.value[index].title = value
    }

    /**
     * Обработка вставки изображений через Ctrl+V
     * @param {number} index - Индекс ссылки
     * @param {ClipboardEvent} event - Событие вставки
     */
    const handleImagePaste = (index, event) => {
      const items = (event.clipboardData || event.originalEvent.clipboardData).items;
      
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        
        // Если вставлен файл изображения
        if (item.type.indexOf('image') !== -1) {
          const blob = item.getAsFile();
          const reader = new FileReader();
          
          reader.onload = (e) => {
            // Создаем data URL изображения
            const dataURL = e.target.result;
            
            // Добавляем к существующим ссылкам
            const currentUrls = links.value[index].imageUrls || '';
            const separator = currentUrls ? '\n' : '';
            links.value[index].imageUrls = currentUrls + separator + dataURL;
          };
          
          reader.readAsDataURL(blob);
        }
        
        // Если вставлена ссылка
        if (item.type === 'text/plain') {
          item.getAsString((text) => {
            // Проверяем, является ли текст URL
            if (text.match(/^https?:\/\/.+\.(jpg|jpeg|png|gif|webp)/i)) {
              // Добавляем к существующим ссылкам
              const currentUrls = links.value[index].imageUrls || '';
              const separator = currentUrls ? '\n' : '';
              links.value[index].imageUrls = currentUrls + separator + text;
            }
          });
        }
      }
    };

    /**
     * Добавляет новое поле для ввода ссылки и цены
     */
    const addLinkField = () => {
      if (links.value.length < 10) {
        links.value.push({ url: '', price: '', imageUrls: '', title: '', description: '' })
      }
    }

    /**
     * Удаляет поле по индексу
     * @param {number} index - Индекс поля для удаления
     */
    const removeLink = (index) => {
      const newLinks = links.value.filter((_, i) => i !== index)
      links.value = newLinks.length ? newLinks : [{ url: '', price: '', imageUrls: '', title: '', description: '' }]
    }

    /**
     * Обработчик добавления товаров
     */
    const handleAddItems = () => {
      console.log('Ссылки перед фильтрацией:', links.value)
      const validLinks = links.value.filter(link => link.url.trim())
      console.log('Валидные ссылки:', validLinks)
      
      if (validLinks.length > 0) {
        // Преобразуем imageUrls в массив
        const processedLinks = validLinks.map(link => {
          let images = []
          if (link.imageUrls) {
            images = link.imageUrls.split('\n').map(url => url.trim()).filter(url => url)
          }
          return {
            ...link,
            images: images,
            // Используем заголовок из формы или извлекаем из URL
            title: link.title && link.title !== 'Загрузка...' ? link.title : 'Без названия'
          }
        })
        
        console.log('Добавляем товары:', processedLinks)
        props.onAddLinks(processedLinks)
        // Очищаем форму после добавления
        links.value = [{ url: '', price: '', imageUrls: '', title: '', description: '' }]
      }
    }
    
    // Обработка событий drag
    const handleDragEnter = (event) => {
      event.preventDefault()
      isDragOver.value = true
    }
    
    const handleDragLeave = (event) => {
      event.preventDefault()
      isDragOver.value = false
    }
    
    // Обработка drag and drop
    const handleDrop = (event, index) => {
      event.preventDefault()
      isDragOver.value = false
      
      const files = event.dataTransfer.files
      if (files && files.length > 0) {
        for (let i = 0; i < files.length; i++) {
          const file = files[i]
          if (file.type.startsWith('image/')) {
            const reader = new FileReader()
            reader.onload = (e) => {
              // Создаем data URL изображения
              const dataURL = e.target.result
              
              // Добавляем к существующим ссылкам
              const currentUrls = links.value[index].imageUrls || ''
              const separator = currentUrls ? '\n' : ''
              links.value[index].imageUrls = currentUrls + separator + dataURL
            }
            reader.readAsDataURL(file)
          }
        }
      }
    }

    return {
      links,
      updateLinkUrl,
      updateLinkPrice,
      updateLinkImageUrls,
      updateLinkTitle,
      handleImagePaste,
      addLinkField,
      removeLink,
      handleAddItems,
      isDragOver,
      handleDragEnter,
      handleDragLeave,
      handleDrop
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
  transform: translateX(-10px);
}
</style>
