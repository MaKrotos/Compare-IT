<template>
  <div class="comparison-card-container">
    <Card class="comparison-card">
      <CardHeader class="comparison-card-header">
        <!-- Изображение -->
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
        
        <!-- Кнопки управления изображениями -->
        <div v-if="isEditingImages" class="comparison-card-editing-container">
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
        
        <!-- Заголовок -->
        <div class="comparison-card-title-container">
          <div v-if="isEditingTitle" class="comparison-card-title-editing">
            <Input
              v-model="title"
              class="comparison-card-title-input"
              @keyup.enter="saveTitle"
            />
            <Button size="icon" @click="saveTitle" class="comparison-card-save-title-button">
              <Save class="comparison-card-save-title-icon" />
            </Button>
          </div>
          <div v-else class="comparison-card-title-display">
            <div v-if="!readonly" class="comparison-card-title-buttons">
              <Button
                variant="ghost"
                size="icon"
                @click="isEditingTitle = true"
                class="comparison-card-title-button"
              >
                <Edit2 class="comparison-card-edit-icon" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                @click="onDelete"
                class="comparison-card-delete-button"
              >
                <Trash2 class="comparison-card-trash-icon" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                @click="startEditingImages"
                class="comparison-card-image-edit-button"
              >
                <Image class="comparison-card-image-icon" />
              </Button>
            </div>
            <h3 class="comparison-card-title-text">
              {{ title }}
            </h3>
          </div>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="comparison-card-link"
          >
            <ExternalLink class="comparison-card-link-icon" />
            Открыть ссылку
          </a>

          <!-- Цена -->
          <div class="comparison-card-price-container">
            <div v-if="isEditingPrice" class="comparison-card-price-editing">
              <Input
                type="number"
                v-model="price"
                placeholder="Цена"
                class="comparison-card-price-input"
                @keyup.enter="savePrice"
              />
              <Button size="icon" @click="savePrice" class="comparison-card-save-price-button">
                <Save class="comparison-card-save-price-icon" />
              </Button>
            </div>
            <button
              v-else
              @click="isEditingPrice = true"
              :disabled="readonly"
              class="comparison-card-price-display"
              :class="{ 'comparison-card-price-disabled': readonly }"
            >
              <DollarSign class="comparison-card-price-icon" />
              <span v-if="item.price" class="comparison-card-price-value">{{ item.price.toLocaleString() }} {{ item.currency || '₽' }}</span>
              <span v-else class="comparison-card-price-placeholder">Указать цену</span>
            </button>
          </div>
        </div>
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
      </CardContent>

      <!-- Статистика -->
      <div class="comparison-card-stats">
        <div class="comparison-card-stats-grid">
          <div class="comparison-card-stat-item">
            <div class="comparison-card-stat-value comparison-card-stat-value-pro">
              {{ (item.pros || []).reduce((sum, p) => sum + (p.impact || 5), 0) }}
            </div>
            <div class="comparison-card-stat-label">{{ (item.pros || []).length }} плюсов</div>
          </div>
          <div class="comparison-card-stat-item">
            <div class="comparison-card-stat-value comparison-card-stat-value-con">
              {{ (item.cons || []).reduce((sum, c) => sum + (c.impact || 5), 0) }}
            </div>
            <div class="comparison-card-stat-label">{{ (item.cons || []).length }} минусов</div>
          </div>
          <div class="comparison-card-stat-item">
            <div class="comparison-card-stat-value comparison-card-stat-value-rating">
              {{ item.rating?.toFixed(1) || '—' }}
            </div>
            <div class="comparison-card-stat-label">общий рейтинг</div>
          </div>
          <div class="comparison-card-stat-item">
            <div class="comparison-card-stat-value comparison-card-stat-value-price">
              {{ item.priceRating?.toFixed(1) || '—' }}
            </div>
            <div class="comparison-card-stat-label">рейтинг цены</div>
          </div>
          <div class="comparison-card-stat-item">
            <div class="comparison-card-stat-value comparison-card-stat-value-pros-cons">
              {{ item.prosConsRating?.toFixed(1) || '—' }}
            </div>
            <div class="comparison-card-stat-label">рейтинг плюсов/минусов</div>
          </div>
        </div>
      </div>
    </Card>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import Card from '@/components/ui/Card.vue'
import CardContent from '@/components/ui/CardContent.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import  Button  from '@/components/ui/Button.vue'
import  Input  from '@/components/ui/Input.vue'
import  Badge  from '@/components/ui/Badge.vue'
import ProConItem from './ProConItem.vue'
import {
  ExternalLink, Trash2, Edit2, Save, DollarSign, Star, ChevronLeft, ChevronRight, Image, Plus, X, Upload
} from 'lucide-vue-next'

export default {
  name: 'ComparisonCard',
  components: {
    Card,
    CardContent,
    CardHeader,
    Button,
    Input,
    Badge,
    ProConItem,
    ExternalLink,
    Trash2,
    Edit2,
    Save,
    DollarSign,
    Star,
    ChevronLeft,
    ChevronRight,
    Image,
    Plus,
    X,
    Upload
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
    }
  },
  setup(props) {
    const isEditingTitle = ref(false)
    const title = ref(props.item.title || 'Без названия')
    const isEditingPrice = ref(false)
    const price = ref(props.item.price || 0)
    const isEditingImages = ref(false)
    const editingImages = ref([])
    const fileInput = ref(null)
    const isDragOver = ref(false)
    const imageUrlInput = ref('')
    
    // Индекс текущего изображения
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
    
    // Начать редактирование изображений
    const startEditingImages = () => {
      isEditingImages.value = true
      // Копируем текущие изображения для редактирования
      editingImages.value = [...images.value]
    }
    
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
      props.onUpdate({ ...props.item, images: [...editingImages.value] })
      isEditingImages.value = false
    }

    // Наблюдение за изменением props.item для обновления локальных данных
    watch(() => props.item, (newItem) => {
      title.value = newItem.title || 'Без названия'
      price.value = newItem.price || 0
      // Сброс индекса при изменении изображений
      currentImageIndex.value = 0
    }, { deep: true })

    const saveTitle = () => {
      props.onUpdate({ ...props.item, title: title.value })
      isEditingTitle.value = false
    }

    const savePrice = () => {
      const updatedItem = { ...props.item, price: parseFloat(price.value) || 0 }
      props.onUpdate(calculateRating(updatedItem, [])) // Передаем пустой массив, так как пересчет будет в App.vue
      isEditingPrice.value = false
    }

    const updatePros = (pros) => {
      const updatedItem = { ...props.item, pros }
      props.onUpdate(calculateRating(updatedItem, [])) // Передаем пустой массив, так как пересчет будет в App.vue
    }

    const updateCons = (cons) => {
      const updatedItem = { ...props.item, cons }
      props.onUpdate(calculateRating(updatedItem, [])) // Передаем пустой массив, так как пересчет будет в App.vue
    }

    // Расчет рейтинга с учетом относительности плюсов/минусов
    const calculateRating = (itemData, allItems = []) => {
      // Конфигурационные параметры для настройки веса рейтинга
      const RATING_CONFIG = {
        // Максимальный рейтинг
        MAX_RATING: 100,
        // Вес рейтинга по плюсам/минусам (в процентах от максимального рейтинга)
        PROS_CONS_WEIGHT: 80,
        // Вес рейтинга по цене (в процентах от максимального рейтинга)
        PRICE_WEIGHT: 20
      };

      const pros = itemData.pros || []
      const cons = itemData.cons || []
      const itemPrice = itemData.price || 0

      // Сумма положительных баллов
      const prosScore = pros.reduce((sum, pro) => sum + (pro.impact || 5), 0)
      
      // Сумма отрицательных баллов
      const consScore = cons.reduce((sum, con) => sum + (con.impact || 5), 0)
      
      // Базовый рейтинг (0-100)
      const maxPossibleScore = Math.max((pros.length + cons.length) * 10, 1)
      const baseRating = ((prosScore - consScore) / maxPossibleScore) * 50 + 50
      
      // Корректировка на цену (чем выше цена, тем меньше рейтинг)
      // Находим максимальную и минимальную цену среди всех товаров для нормализации
      let priceImpact = 0;
      if (allItems.length > 0 && itemPrice > 0) {
        const maxPrice = Math.max(...allItems.map(item => item.price || 0));
        const minPrice = Math.min(...allItems.map(item => item.price || 0));
        if (maxPrice > minPrice) {
          // Чем выше цена по сравнению с минимальной, тем меньше рейтинг
          // Чем ниже цена по сравнению с максимальной, тем выше рейтинг
          const priceRange = maxPrice - minPrice;
          const priceRatio = (maxPrice - itemPrice) / priceRange;
          // Применяем вес цены к максимальному бонусу/штрафу
          priceImpact = priceRatio * RATING_CONFIG.PRICE_WEIGHT;
        }
      }
      
      // Учет относительности плюсов/минусов
      let prosConsImpact = 0;
      if (allItems.length > 0) {
        // Средние значения по всем карточкам
        const totalProsScore = allItems.reduce((sum, item) => sum + (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0), 0);
        const totalConsScore = allItems.reduce((sum, item) => sum + (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0), 0);
        const avgProsScore = totalProsScore / allItems.length;
        const avgConsScore = totalConsScore / allItems.length;
        
        // Сравнение с лучшей карточкой
        const maxProsScore = Math.max(...allItems.map(item => (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0)));
        const maxConsScore = Math.max(...allItems.map(item => (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0)));
        
        // Нормализованные значения (0-1)
        const normalizedPros = maxProsScore > 0 ? prosScore / maxProsScore : 0;
        const normalizedCons = maxConsScore > 0 ? consScore / maxConsScore : 0;
        
        // Влияние плюсов (чем больше плюсов по сравнению со средним и чем больше по сравнению с лучшей карточкой, тем выше рейтинг)
        const prosVsAvg = avgProsScore > 0 ? prosScore / avgProsScore : 1;
        const prosBonus = Math.min(1.5, Math.max(0.5, prosVsAvg)); // Ограничиваем влияние
        
        // Влияние минусов (чем меньше минусов по сравнению со средним и чем меньше по сравнению с худшей карточкой, тем выше рейтинг)
        const consVsAvg = avgConsScore > 0 ? consScore / avgConsScore : 0;
        const consBonus = Math.min(1.5, Math.max(0.5, 1 - consVsAvg)); // Ограничиваем влияние
        
        // Комбинированное влияние с учетом веса плюсов/минусов
        prosConsImpact = ((normalizedPros * prosBonus) - (normalizedCons * consBonus)) * RATING_CONFIG.PROS_CONS_WEIGHT * 0.15;
      }
      
      const finalRating = Math.max(0, Math.min(RATING_CONFIG.MAX_RATING, baseRating + priceImpact + prosConsImpact))
      
      // Добавляем отдельные рейтинги
      const priceRating = Math.max(0, Math.min(100, 50 + priceImpact));
      const prosConsRating = Math.max(0, Math.min(100, 50 + prosConsImpact));
      
      return { 
        ...itemData, 
        rating: Math.round(finalRating * 10) / 10,
        priceRating: Math.round(priceRating * 10) / 10,
        prosConsRating: Math.round(prosConsRating * 10) / 10
      }
    }

    return {
      isEditingTitle,
      title,
      isEditingPrice,
      price,
      isEditingImages,
      editingImages,
      fileInput,
      isDragOver,
      saveTitle,
      savePrice,
      updatePros,
      updateCons,
      calculateRating,
      images,
      currentImage,
      currentImageIndex,
      nextImage,
      prevImage,
      startEditingImages,
      handleFileUpload,
      handleDrop,
      handleDragEnter,
      handleDragLeave,
      handlePaste,
      removeEditingImage,
      saveImages,
      addImageFromUrl,
      imageUrlInput
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

/* Заголовок */
.comparison-card-title-container {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.comparison-card-title-editing {
  display: flex;
  gap: 0.5rem;
}

.comparison-card-title-input {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  flex: 1;
}

.comparison-card-save-title-button {
  height: 2.5rem;
  width: 2.5rem;
  background-color: #2563eb;
}

.comparison-card-save-title-button:hover {
  background-color: #1d4ed8;
}

.comparison-card-save-title-icon {
  width: 1rem;
  height: 1rem;
  color: #ffffff;
}

.comparison-card-title-display {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.comparison-card-title-text {
  font-size: 1.125rem;
  line-height: 1.375;
  font-weight: 500;
  color: #111827;
  flex: 1;
  word-wrap: break-word;
  overflow-wrap: break-word;
  word-break: break-word;
}

.comparison-card-title-buttons {
  display: flex;
  gap: 0.25rem;
}

.comparison-card-title-button,
.comparison-card-delete-button,
.comparison-card-image-edit-button {
  height: 2rem;
  width: 2rem;
  color: #9ca3af;
  border-radius: 9999px;
}

.comparison-card-title-button:hover {
  color: #2563eb;
  background-color: #dbeafe;
}

.comparison-card-delete-button:hover {
  color: #dc2626;
  background-color: #fef2f2;
}

.comparison-card-image-edit-button:hover {
  color: #16a34a;
  background-color: #dcfce7;
}

.comparison-card-edit-icon,
.comparison-card-trash-icon,
.comparison-card-image-icon {
  width: 1rem;
  height: 1rem;
}

/* Ссылка на товар */
.comparison-card-link {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #2563eb;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-top: 0.5rem;
  text-decoration: none;
}

.comparison-card-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.comparison-card-link-icon {
  width: 0.75rem;
  height: 0.75rem;
}

/* Цена */
.comparison-card-price-container {
  margin-top: 0.75rem;
}

.comparison-card-price-editing {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.comparison-card-price-input {
  font-size: 1.125rem;
  line-height: 1.75rem;
  flex: 1;
}

.comparison-card-save-price-button {
  height: 2.25rem;
  width: 2.25rem;
  background-color: #2563eb;
}

.comparison-card-save-price-button:hover {
  background-color: #1d4ed8;
}

.comparison-card-save-price-icon {
  width: 1rem;
  height: 1rem;
  color: #ffffff;
}

.comparison-card-price-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.125rem;
  line-height: 1.75rem;
  font-weight: 600;
  color: #2563eb;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  transition: color 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-price-display:hover:not(.comparison-card-price-disabled) {
  color: #1d4ed8;
}

.comparison-card-price-disabled {
  cursor: default;
}

.comparison-card-price-icon {
  width: 1.25rem;
  height: 1.25rem;
  transition: transform 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-price-display:hover:not(.comparison-card-price-disabled) .comparison-card-price-icon {
  transform: scale(1.1);
}

.comparison-card-price-value,
.comparison-card-price-placeholder {
  font-weight: 600;
}

/* Контент карточки */
.comparison-card-content {
  padding: 1rem;
  margin-bottom: 1rem;
  flex: 1;
}

/* Статистика */
.comparison-card-stats {
  padding: 1rem;
  border-top: 1px solid #f3f4f6;
  background-color: rgba(249, 250, 251, 0.5);
}

.comparison-card-stats-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 0.75rem;
  text-align: center;
}

.comparison-card-stat-item {
  display: flex;
  flex-direction: column;
}

.comparison-card-stat-value {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
}

.comparison-card-stat-value-pro {
  color: #16a34a;
}

.comparison-card-stat-value-con {
  color: #dc2626;
}

.comparison-card-stat-value-rating {
  color: #2563eb;
}

.comparison-card-stat-value-price {
  color: #f59e0b;
}

.comparison-card-stat-value-pros-cons {
  color: #8b5cf6;
}

.comparison-card-stat-label {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
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