<template>
  <div class="product-link-field">
    <div class="product-link-field-row">
      <div class="product-link-field-column">
        <!-- URL -->
        <div class="product-link-input-container product-link-url-input">
          <Input :value="link.url" @input="updateUrl" :placeholder="`Ссылка на товар ${index + 1}`"
            :disabled="isLoading" />
          <Link2 class="product-link-input-icon product-link-url-icon" />
        </div>

        <!-- Название товара -->
        <div class="product-link-input-container product-link-title-input">
          <Input :value="link.title" @input="updateTitle" placeholder="Название товара" :disabled="isLoading" />
          <Tag class="product-link-input-icon product-link-title-icon" />
        </div>

        <!-- Цена -->
        <div class="product-link-input-container product-link-price-input">
          <Input type="number" :value="link.price" @input="updatePrice" placeholder="Цена, ₽" :disabled="isLoading" />
          <DollarSign class="product-link-input-icon product-link-price-icon" />
        </div>

        <!-- Ссылки на изображения -->
        <ProductImageField :imageUrls="link.imageUrls" @update:imageUrls="updateImageUrls" @paste="handleImagePaste"
          :isLoading="isLoading" @dragenter="handleDragEnter" @dragleave="handleDragLeave" @drop="handleDrop"
          :isDragOver="isDragOver" />
      </div>

      <!-- Кнопка удаления -->
      <Button v-if="showRemoveButton" type="button" variant="ghost" size="icon" @click="remove" :disabled="isLoading"
        class="product-link-remove-button">
        <X class="product-link-remove-icon" />
      </Button>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import ProductImageField from './ProductImageField.vue'
import { Link2, DollarSign, X, Tag } from 'lucide-vue-next'
import { getCachedPreviewData } from '@/utils/preview.js'

export default {
  name: 'ProductLinkField',
  components: {
    Button,
    Input,
    ProductImageField,
    Link2,
    DollarSign,
    X,
    Tag
  },
  props: {
    link: {
      type: Object,
      required: true
    },
    index: {
      type: Number,
      required: true
    },
    isLoading: {
      type: Boolean,
      default: false
    },
    showRemoveButton: {
      type: Boolean,
      default: true
    },
    isDragOver: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:link', 'remove', 'dragenter', 'dragleave', 'drop'],
  setup(props, { emit }) {
    const updateUrl = (event) => {
      const newUrl = event.target.value
      console.log('updateUrl вызван с новым URL:', newUrl);
      emit('update:link', { ...props.link, url: newUrl })

      // Если URL валиден и заголовок еще не установлен (или содержит значение по умолчанию), пытаемся получить имя товара
      if (newUrl && (!props.link.title || props.link.title === 'Загрузка...' || props.link.title === 'Без названия')) {
        fetchProductName(newUrl)
      }
    }

    const updateTitle = (event) => {
      emit('update:link', { ...props.link, title: event.target.value })
    }

    const updatePrice = (event) => {
      emit('update:link', { ...props.link, price: event.target.value })
    }

    const updateImageUrls = (value) => {
      emit('update:link', { ...props.link, imageUrls: value })
    }

    const handleImagePaste = (event) => {
      // Просто пробрасываем событие выше
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

    const remove = () => {
      emit('remove')
    }

    /**
     * Получает имя товара с сайта
     * @param {string} url - URL сайта
     */
    const fetchProductName = async (url) => {
      console.log('fetchProductName вызван с URL:', url);
      try {
        // Проверяем, что URL валиден
        new URL(url)

        // Показываем индикатор загрузки
        const updatedLink = { ...props.link, url: url, title: 'Загрузка...' };
        emit('update:link', updatedLink)

        // Получаем данные предварительного просмотра с кэшированием
        const previewData = await getCachedPreviewData(url)
        console.log('Получены данные предварительного просмотра:', previewData);

        // Обновляем название товара
        // Используем текущее значение props.link для получения актуальных данных
        const updatedLink = {
          ...props.link,
          url: url,
          title: previewData.title || 'Без названия',
          imageUrls: props.link.imageUrls || previewData.image || '',
          description: props.link.description || previewData.description || ''
        };
        emit('update:link', updatedLink)
      } catch (error) {
        console.error('Ошибка при получении данных товара:', error)
        const updatedLink = { ...props.link, url: url, title: 'Без названия' };
        emit('update:link', updatedLink)
      }
    }

    return {
      updateUrl,
      updateTitle,
      updatePrice,
      updateImageUrls,
      handleImagePaste,
      handleDragEnter,
      handleDragLeave,
      handleDrop,
      remove
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ProductLinkField */
.product-link-field-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.product-link-field-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.product-link-input-container {
  position: relative;
}

.product-link-input-icon {
  position: absolute;
  left: 0.875rem;
  top: 50%;
  transform: translateY(-50%);
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
  pointer-events: none;
}

.product-link-url-input {
  padding-left: 2.75rem;
}

.product-link-title-input {
  padding-left: 2.75rem;
}

.product-link-price-input {
  padding-left: 2.5rem;
}

.product-link-price-icon {
  left: 0.75rem;
}

.product-link-remove-button {
  margin-top: 0.25rem;
}

.product-link-remove-icon {
  width: 1.25rem;
  height: 1.25rem;
}
</style>