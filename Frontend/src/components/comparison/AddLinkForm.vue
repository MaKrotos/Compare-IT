<template>
  <Card class="add-link-form-container">
    <CardHeader class="add-link-form-header">
      <CardTitle class="add-link-form-title">
        <div class="add-link-form-icon-container">
          <Link2 class="add-link-form-icon" />
        </div>
        Добавить товары для сравнения
      </CardTitle>
      <p class="add-link-form-subtitle">
        Вставьте ссылки на товары (до 10 штук)
      </p>
    </CardHeader>
    <CardContent class="add-link-form-content">
      <div class="add-link-form-fields">
        <transition-group name="list" tag="div" class="add-link-form-fields">
          <div v-for="(link, index) in links" :key="index" class="add-link-form-field-group">
            <div class="add-link-form-field-row">
              <div class="add-link-form-field-column">
                <!-- URL -->
                <div class="add-link-form-input-container add-link-form-url-input">
                  <Input
                    :value="link.url"
                    @input="updateLinkUrl(index, $event.target.value)"
                    :placeholder="`Ссылка на товар ${index + 1}`"
                    :disabled="isLoading"
                  />
                  <Link2 class="add-link-form-input-icon add-link-form-url-icon" />
                </div>

                <!-- Название товара -->
                <div class="add-link-form-input-container add-link-form-title-input">
                  <Input
                    :value="link.title"
                    @input="updateLinkTitle(index, $event.target.value)"
                    placeholder="Название товара"
                    :disabled="isLoading"
                  />
                  <Tag class="add-link-form-input-icon add-link-form-title-icon" />
                </div>

                <!-- Цена -->
                <div class="add-link-form-input-container add-link-form-price-input">
                  <Input
                    type="number"
                    :value="link.price"
                    @input="updateLinkPrice(index, $event.target.value)"
                    placeholder="Цена, ₽"
                    :disabled="isLoading"
                  />
                  <DollarSign class="add-link-form-input-icon add-link-form-price-icon" />
                </div>
                
                <!-- Ссылки на изображения -->
                <div class="add-link-form-textarea-container">
                  <Textarea
                    :value="link.imageUrls"
                    @input="updateLinkImageUrls(index, $event.target.value)"
                    @paste="handleImagePaste(index, $event)"
                    placeholder="Ссылки на изображения (по одной на строку) или вставьте изображения (Ctrl+V)"
                    :disabled="isLoading"
                  />
                  <Image class="add-link-form-textarea-icon" />
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
                class="add-link-form-remove-button"
              >
                <X class="add-link-form-remove-icon" />
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
            <Loader2 v-if="isLoading" class="add-link-form-loading-icon" />
            <span v-if="isLoading" class="add-link-form-loading-text">Загрузка...</span>
            <span v-else class="add-link-form-submit-text">Загрузить товары</span>
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
import { getCachedPreviewData } from '@/utils/preview.js'

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
/* Стили для компонента AddLinkForm */
.add-link-form-container {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  border-radius: 1rem;
  overflow: hidden;
}

.add-link-form-header {
  border-bottom: 1px solid #f3f4f6;
  padding-bottom: 1.25rem;
  background-image: linear-gradient(to right, rgba(219, 234, 254, 0.8), rgba(224, 231, 255, 0.8));
}

.add-link-form-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #1f2937;
}

.add-link-form-icon-container {
  padding: 0.5rem;
  background-color: #dbeafe;
  border-radius: 50%;
}

.add-link-form-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563eb;
}

.add-link-form-subtitle {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
  margin-top: 0.375rem;
}

.add-link-form-content {
  padding-top: 1.5rem;
}

.add-link-form-fields {
  margin-bottom: 1rem;
}

.add-link-form-fields > :not(:last-child) {
  margin-bottom: 1rem;
}

.add-link-form-field-group {
  margin-bottom: 1.5rem;
}

.add-link-form-field-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.add-link-form-field-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.add-link-form-input-container {
  position: relative;
}

.add-link-form-input-icon {
  position: absolute;
  left: 0.875rem;
  top: 50%;
  transform: translateY(-50%);
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
  pointer-events: none;
}

.add-link-form-url-input {
  padding-left: 2.75rem;
}

.add-link-form-title-input {
  padding-left: 2.75rem;
}

.add-link-form-price-input {
  padding-left: 2.5rem;
}

.add-link-form-price-icon {
  left: 0.75rem;
}

.add-link-form-textarea-container {
  position: relative;
}

.add-link-form-textarea-icon {
  position: absolute;
  left: 0.75rem;
  top: 0.75rem;
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
  pointer-events: none;
}

.add-link-form-drag-area {
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.add-link-form-drag-area:hover {
  border-color: #3b82f6;
}

.add-link-form-drag-area-active {
  border-color: #3b82f6;
  background-color: #dbeafe;
}

.add-link-form-drag-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.add-link-form-drag-icon {
  width: 1.5rem;
  height: 1.5rem;
  color: #9ca3af;
}

.add-link-form-drag-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
}

.add-link-form-drag-subtext {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
  margin-top: 0.25rem;
}

.add-link-form-remove-button {
  margin-top: 0.25rem;
}

.add-link-form-remove-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.add-link-form-button-container {
  display: flex;
  gap: 0.75rem;
  padding-top: 0.5rem;
}

.add-link-form-add-button {
  flex: 1;
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  height: 3rem;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.add-link-form-add-button:hover {
  border-color: #3b82f6;
  background-color: #dbeafe;
}

.add-link-form-submit-button {
  flex: 1;
  height: 3rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.add-link-form-submit-button:hover {
  background-color: #1d4ed8;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}

.add-link-form-button-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
}

.add-link-form-loading-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
  animation: spin 1s linear infinite;
}

.add-link-form-loading-text {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
}

.add-link-form-submit-text {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

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