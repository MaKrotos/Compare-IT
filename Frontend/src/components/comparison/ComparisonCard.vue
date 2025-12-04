<template>
  <div class="comparison-card-container">
    <Card class="comparison-card">
      <CardHeader class="comparison-card-header">
        <!-- Изображение -->
        <ImageGallery :item="item" :title="item.title || 'Без названия'" :item-number="itemNumber" />
        
        <!-- Кнопки управления изображениями -->
        <ImageUploader
          v-if="isEditingImages"
          :initial-images="editingImages"
          @save="saveImages"
        />
        
        <!-- Заголовок -->
        <ProductHeader
          :item="item"
          :readonly="readonly"
          @update="handleUpdate"
          @delete="onDelete"
          @edit-images="startEditingImages"
        />
      </CardHeader>

      <CardContent class="comparison-card-content">
        <ProConItem
          type="pro"
          :items="item.pros || []"
          :on-update="updatePros"
          :readonly="readonly"
        />
        
        <ProConItem
          type="con"
          :items="item.cons || []"
          :on-update="updateCons"
          :readonly="readonly"
        />
        
        <ProConItem
          type="note"
          :items="item.notes || []"
          :on-update="updateNotes"
          :readonly="readonly"
        />
      </CardContent>

      <!-- Статистика -->
      <ProductStats :item="item" />
    </Card>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import CardContent from '@/components/ui/CardContent.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import ProConItem from './ProConItem.vue'
import ImageGallery from './ImageGallery.vue'
import ImageUploader from './ImageUploader.vue'
import ProductHeader from './ProductHeader.vue'
import ProductStats from './ProductStats.vue'

export default {
  name: 'ComparisonCard',
  components: {
    Card,
    CardContent,
    CardHeader,
    ProConItem,
    ImageGallery,
    ImageUploader,
    ProductHeader,
    ProductStats
  },
  props: {
    item: {
      type: Object,
      required: true
    },
    onUpdate: {
      type: Function,
      required: true
    },
    onDelete: {
      type: Function,
      required: true
    },
    readonly: {
      type: Boolean,
      default: false
    },
    itemNumber: {
      type: Number,
      default: 0
    }
  },
  setup(props) {
    const isEditingImages = ref(false)
    const editingImages = ref([])
    
    // Вычисляем массив изображений из item.images или item.image
    const images = computed(() => {
      if (Array.isArray(props.item.images) && props.item.images.length > 0) {
        return props.item.images
      }
      if (props.item.image) {
        return [props.item.image]
      }
      return []
    })
    
    // Начать редактирование изображений
    const startEditingImages = () => {
      isEditingImages.value = true
      // Копируем текущие изображения для редактирования
      editingImages.value = [...images.value]
    }
    
    // Сохранить изменения изображений
    const saveImages = (newImages) => {
      props.onUpdate({ ...props.item, images: newImages })
      isEditingImages.value = false
    }
    
    // Обработчик обновления элемента
    const handleUpdate = (updatedItem) => {
      props.onUpdate(updatedItem)
    }
    
    // Обновить плюсы
    const updatePros = (pros) => {
      props.onUpdate({ ...props.item, pros })
    }
    
    // Обновить минусы
    const updateCons = (cons) => {
      props.onUpdate({ ...props.item, cons })
    }
    
    // Обновить заметки
    const updateNotes = (notes) => {
      props.onUpdate({ ...props.item, notes })
    }
    
    return {
      isEditingImages,
      editingImages,
      images,
      startEditingImages,
      saveImages,
      handleUpdate,
      updatePros,
      updateCons,
      updateNotes
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ComparisonCard */
.comparison-card-container {
  height: 100%;
  max-width: 100%;
}

.comparison-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  max-width: 100%;
}

.comparison-card-header {
  padding: 0;
  margin-bottom: 0;
}

/* Ограничение ширины карточек на мобильных устройствах */
@media (max-width: 768px), (orientation: portrait) and (max-aspect-ratio: 1/1) {
  .comparison-card-container {
    max-width: 100vw;
  }
  
  .comparison-card {
    max-width: 100vw;
  }
}

/* Контент карточки */
.comparison-card-content {
  padding: 1rem;
  margin-bottom: 1rem;
  flex: 1;
}
</style>