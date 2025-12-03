<template>
  <div class="comparison-card-image-container">
    <!-- Контейнер для изображений -->
    <div v-if="images.length > 0" class="comparison-card-image-wrapper">
      <!-- Текущее изображение -->
      <transition name="comparison-card-fade" mode="out-in">
        <img
          :key="currentImage"
          :src="currentImage"
          :alt="title"
          class="comparison-card-image"
        />
      </transition>
      
      <!-- Кнопки навигации -->
      <Button
        v-if="images.length > 1"
        variant="ghost"
        size="icon"
        @click="prevImage"
        class="comparison-card-nav-button comparison-card-nav-button-left"
      >
        <ChevronLeft class="comparison-card-nav-icon" />
      </Button>
      
      <Button
        v-if="images.length > 1"
        variant="ghost"
        size="icon"
        @click="nextImage"
        class="comparison-card-nav-button comparison-card-nav-button-right"
      >
        <ChevronRight class="comparison-card-nav-icon" />
      </Button>
      
      <!-- Индикаторы изображений -->
      <div v-if="images.length > 1" class="comparison-card-indicators">
        <div
          v-for="(img, index) in images"
          :key="index"
          :class="[
            'comparison-card-indicator',
            index === currentImageIndex ? 'comparison-card-indicator-active' : 'comparison-card-indicator-inactive'
          ]"
        />
      </div>
    </div>
    
    <!-- Заглушка если нет изображений -->
    <div v-else class="comparison-card-image-placeholder">
      <Image class="comparison-card-placeholder-icon" />
    </div>
    
    <!-- Рейтинг badge -->
    <div v-if="item.rating !== undefined" class="comparison-card-rating-badge">
      <Badge class="comparison-card-rating-badge-content">
        <Star class="comparison-card-star-icon" />
        {{ item.rating }}
      </Badge>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'
import { ChevronLeft, ChevronRight, Image, Star } from 'lucide-vue-next'

export default {
  name: 'ImageGallery',
  components: {
    Button,
    Badge,
    ChevronLeft,
    ChevronRight,
    Image,
    Star
  },
  props: {
    item: {
      type: Object,
      required: true
    },
    title: {
      type: String,
      default: 'Без названия'
    }
  },
  setup(props) {
    const currentImageIndex = ref(0)
    
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
    
    // Текущее изображение
    const currentImage = computed(() => {
      return images.value[currentImageIndex.value] || ''
    })
    
    // Перейти к следующему изображению
    const nextImage = () => {
      if (images.value.length > 1) {
        currentImageIndex.value = (currentImageIndex.value + 1) % images.value.length
      }
    }
    
    // Перейти к предыдущему изображению
    const prevImage = () => {
      if (images.value.length > 1) {
        currentImageIndex.value = (currentImageIndex.value - 1 + images.value.length) % images.value.length
      }
    }
    
    // Наблюдение за изменением props.item для сброса индекса
    watch(() => props.item, () => {
      currentImageIndex.value = 0
    }, { deep: true })
    
    return {
      images,
      currentImage,
      currentImageIndex,
      nextImage,
      prevImage
    }
  }
}
</script>

<style scoped>
/* Стили для изображения */
.comparison-card-image-container {
  position: relative;
  aspect-ratio: 1 / 1;
  background-color: #f9fafb;
  overflow: hidden;
  border-top-left-radius: 0.75rem;
  border-top-right-radius: 0.75rem;
}

.comparison-card-image-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.comparison-card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Навигация по изображениям */
.comparison-card-nav-button {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  height: 2rem;
  width: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}

.comparison-card-nav-button:hover {
  background-color: #ffffff;
}

.comparison-card-nav-button-left {
  left: 0.5rem;
}

.comparison-card-nav-button-right {
  right: 0.5rem;
}

.comparison-card-nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #374151;
}

/* Индикаторы изображений */
.comparison-card-indicators {
  position: absolute;
  bottom: 0.75rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 0.375rem;
}

.comparison-card-indicator {
  height: 0.375rem;
  border-radius: 9999px;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-indicator-active {
  width: 1.5rem;
  background-color: #ffffff;
}

.comparison-card-indicator-inactive {
  width: 0.375rem;
  background-color: rgba(255, 255, 255, 0.6);
}

/* Заглушка для изображения */
.comparison-card-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
}

.comparison-card-placeholder-icon {
  width: 3rem;
  height: 3rem;
  color: #d1d5db;
}

/* Рейтинг badge */
.comparison-card-rating-badge {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
}

.comparison-card-rating-badge-content {
  background-color: rgba(255, 255, 255, 0.95);
  color: #111827;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1);
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 700;
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
}

.comparison-card-star-icon {
  width: 1rem;
  height: 1rem;
  fill: #fbbf24;
  color: #fbbf24;
  margin-right: 0.25rem;
}

/* Анимации */
.comparison-card-fade-enter-active {
  transition: opacity 100ms cubic-bezier(0.4, 0, 1, 1);
}

.comparison-card-fade-leave-active {
  transition: opacity 100ms cubic-bezier(0, 0, 0.2, 1);
}

.comparison-card-fade-enter-from,
.comparison-card-fade-leave-to {
  opacity: 0;
}
</style>