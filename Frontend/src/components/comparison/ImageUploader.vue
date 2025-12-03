<template>
  <div class="comparison-card-editing-container">
    <div class="comparison-card-editing-buttons">
      <input
        type="file"
        ref="fileInput"
        @change="handleFileUpload"
        accept="image/*"
        class="comparison-card-file-input"
        multiple
      />
      <Button
        @click="$refs.fileInput.click()"
        size="sm"
        variant="outline"
        class="comparison-card-upload-button"
      >
        <Plus class="comparison-card-button-icon" />
        Добавить фото
      </Button>
      <Button
        @click="saveImages"
        size="sm"
        class="comparison-card-save-button"
      >
        <Save class="comparison-card-button-icon" />
        Сохранить
      </Button>
    </div>
    
    <!-- Область для drag and drop -->
    <div
      @dragover.prevent
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      @paste="handlePaste"
      tabindex="0"
      class="comparison-card-drag-area"
      :class="{ 'comparison-card-drag-area-active': isDragOver }"
    >
      <div class="comparison-card-drag-content">
        <Upload class="comparison-card-drag-icon" />
        <p class="comparison-card-drag-text">Перетащите изображения сюда или нажмите для выбора</p>
        <p class="comparison-card-drag-subtext">Поддерживаются JPG, PNG, GIF</p>
      </div>
    </div>
    
    <!-- Поле для добавления изображений по ссылке -->
    <div class="comparison-card-url-input-container">
      <Input
        v-model="imageUrlInput"
        placeholder="Введите URL изображения"
        class="comparison-card-url-input"
        @keyup.enter="addImageFromUrl"
      />
      <Button
        @click="addImageFromUrl"
        size="sm"
        variant="outline"
        :disabled="!imageUrlInput"
        class="comparison-card-url-button"
      >
        <Plus class="comparison-card-button-icon" />
      </Button>
    </div>
    
    <!-- Список изображений с кнопками удаления -->
    <div class="comparison-card-images-list">
      <div
        v-for="(img, index) in editingImages"
        :key="index"
        class="comparison-card-image-item"
      >
        <div class="comparison-card-image-thumbnail">
          <img :src="img" class="comparison-card-image-thumbnail-img" />
        </div>
        <span class="comparison-card-image-text">Изображение {{ index + 1 }}</span>
        <Button
          @click="removeEditingImage(index)"
          size="icon"
          variant="ghost"
          class="comparison-card-remove-image-button"
        >
          <X class="comparison-card-remove-icon" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { Plus, Save, Upload, X } from 'lucide-vue-next'

export default {
  name: 'ImageUploader',
  components: {
    Button,
    Input,
    Plus,
    Save,
    Upload,
    X
  },
  props: {
    initialImages: {
      type: Array,
      default: () => []
    }
  },
  emits: ['save'],
  setup(props, { emit }) {
    const fileInput = ref(null)
    const isDragOver = ref(false)
    const imageUrlInput = ref('')
    const editingImages = ref([...props.initialImages])
    
    // Обработка загрузки файлов
    const handleFileUpload = (event) => {
      const files = event.target.files
      if (files && files.length > 0) {
        for (let i = 0; i < files.length; i++) {
          const file = files[i]
          if (file.type.startsWith('image/')) {
            const reader = new FileReader()
            reader.onload = (e) => {
              editingImages.value.push(e.target.result)
            }
            reader.readAsDataURL(file)
          }
        }
        // Очищаем input
        event.target.value = ''
      }
    }
    
    // Обработка drag and drop
    const handleDrop = (event) => {
      event.preventDefault()
      isDragOver.value = false
      
      const files = event.dataTransfer.files
      if (files && files.length > 0) {
        for (let i = 0; i < files.length; i++) {
          const file = files[i]
          if (file.type.startsWith('image/')) {
            const reader = new FileReader()
            reader.onload = (e) => {
              editingImages.value.push(e.target.result)
            }
            reader.readAsDataURL(file)
          }
        }
      }
    }
    
    // Обработка вставки из буфера обмена
    const handlePaste = (event) => {
      const items = event.clipboardData.items;
      
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        
        if (item.type.indexOf('image') !== -1) {
          const blob = item.getAsFile();
          const reader = new FileReader();
          
          reader.onload = (e) => {
            editingImages.value.push(e.target.result);
          };
          
          reader.readAsDataURL(blob);
        }
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
    
    // Удалить изображение из редактируемого списка
    const removeEditingImage = (index) => {
      editingImages.value.splice(index, 1)
    }
    
    // Добавить изображение по URL
    const addImageFromUrl = () => {
      if (imageUrlInput.value) {
        // Простая проверка, что это URL изображения
        if (imageUrlInput.value.match(/\.(jpeg|jpg|gif|png|webp)$/i)) {
          editingImages.value.push(imageUrlInput.value)
          imageUrlInput.value = ''
        } else {
          // Если URL не заканчивается расширением, все равно добавляем (может быть URL с параметрами)
          editingImages.value.push(imageUrlInput.value)
          imageUrlInput.value = ''
        }
      }
    }
    
    // Сохранить изменения изображений
    const saveImages = () => {
      emit('save', [...editingImages.value])
    }
    
    return {
      fileInput,
      isDragOver,
      imageUrlInput,
      editingImages,
      handleFileUpload,
      handleDrop,
      handleDragEnter,
      handleDragLeave,
      handlePaste,
      removeEditingImage,
      addImageFromUrl,
      saveImages
    }
  }
}
</script>

<style scoped>
/* Редактирование изображений */
.comparison-card-editing-container {
  padding: 0.75rem;
  background-color: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.comparison-card-editing-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.comparison-card-file-input {
  display: none;
}

.comparison-card-upload-button,
.comparison-card-save-button {
  flex: 1;
}

.comparison-card-save-button {
  background-color: #2563eb;
}

.comparison-card-save-button:hover {
  background-color: #1d4ed8;
}

.comparison-card-button-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.25rem;
}

/* Drag and drop область */
.comparison-card-drag-area {
  margin-top: 0.75rem;
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: border-color 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-drag-area:hover {
  border-color: #3b82f6;
}

.comparison-card-drag-area:focus {
  outline: 2px solid transparent;
  outline-offset: 2px;
}

.comparison-card-drag-area:focus-visible {
  box-shadow: 0 0 0 2px #3b82f6;
  border-color: transparent;
}

.comparison-card-drag-area-active {
  border-color: #3b82f6;
  background-color: #dbeafe;
}

.comparison-card-drag-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.comparison-card-drag-icon {
  width: 1.5rem;
  height: 1.5rem;
  color: #9ca3af;
}

.comparison-card-drag-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
}

.comparison-card-drag-subtext {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
  margin-top: 0.25rem;
}

/* URL ввод для изображений */
.comparison-card-url-input-container {
  margin-top: 0.75rem;
  display: flex;
  gap: 0.5rem;
}

.comparison-card-url-input {
  flex: 1;
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.comparison-card-url-button {
  padding: 0 0.75rem;
}

/* Список изображений */
.comparison-card-images-list {
  margin-top: 0.75rem;
  max-height: 10rem;
  overflow-y: auto;
}

.comparison-card-image-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background-color: #ffffff;
  border-radius: 0.25rem;
  border: 1px solid #e5e7eb;
  margin-bottom: 0.5rem;
}

.comparison-card-image-thumbnail {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.25rem;
  overflow: hidden;
  flex-shrink: 0;
}

.comparison-card-image-thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.comparison-card-image-text {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.comparison-card-remove-image-button {
  height: 2rem;
  width: 2rem;
  color: #ef4444;
}

.comparison-card-remove-image-button:hover {
  color: #dc2626;
  background-color: #fef2f2;
}

.comparison-card-remove-icon {
  width: 1rem;
  height: 1rem;
}
</style>