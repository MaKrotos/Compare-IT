<template>
  <div class="product-image-field">
    <!-- Ссылки на изображения -->
    <div class="product-image-textarea-container">
      <Textarea
        :value="imageUrls"
        @input="updateImageUrls"
        @paste="handlePaste"
        placeholder="Ссылки на изображения (по одной на строку) или вставьте изображения (Ctrl+V)"
        :disabled="isLoading"
      />
      <Image class="product-image-textarea-icon" />
    </div>
    
    <!-- Область для drag and drop -->
    <div
      @dragover.prevent
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      class="product-image-drag-area"
      :class="{ 'product-image-drag-area-active': isDragOver }"
    >
      <div class="product-image-drag-content">
        <Upload class="product-image-drag-icon" />
        <p class="product-image-drag-text">Перетащите изображения сюда</p>
        <p class="product-image-drag-subtext">Поддерживаются JPG, PNG, GIF</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import Textarea from '@/components/ui/Textarea.vue'
import { Image, Upload } from 'lucide-vue-next'

export default {
  name: 'ProductImageField',
  components: {
    Textarea,
    Image,
    Upload
  },
  props: {
    imageUrls: {
      type: String,
      default: ''
    },
    isLoading: {
      type: Boolean,
      default: false
    },
    isDragOver: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:imageUrls', 'paste', 'dragenter', 'dragleave', 'drop'],
  setup(props, { emit }) {
    const updateImageUrls = (event) => {
      emit('update:imageUrls', event.target.value)
    }

    const handlePaste = (event) => {
      emit('paste', event)
    }

    const handleDragEnter = (event) => {
      emit('dragenter', event)
    }

    const handleDragLeave = (event) => {
      emit('dragleave', event)
    }

    const handleDrop = (event) => {
      emit('drop', event)
    }

    return {
      updateImageUrls,
      handlePaste,
      handleDragEnter,
      handleDragLeave,
      handleDrop
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ProductImageField */
.product-image-textarea-container {
  position: relative;
}

.product-image-textarea-icon {
  position: absolute;
  left: 0.75rem;
  top: 0.75rem;
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
  pointer-events: none;
}

.product-image-drag-area {
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  padding: 1rem;
  text-align: center;
  cursor: pointer;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.product-image-drag-area:hover {
  border-color: #3b82f6;
}

.product-image-drag-area-active {
  border-color: #3b82f6;
  background-color: #dbeafe;
}

.product-image-drag-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.product-image-drag-icon {
  width: 1.5rem;
  height: 1.5rem;
  color: #9ca3af;
}

.product-image-drag-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
}

.product-image-drag-subtext {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
  margin-top: 0.25rem;
}
</style>