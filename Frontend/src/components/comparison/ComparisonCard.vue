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
              <ChevronLeft class="w-5 h-5 text-gray-700" />
            </Button>
            
            <Button
              v-if="images.length > 1"
              variant="ghost"
              size="icon"
              @click="nextImage"
              class="comparison-card-nav-button comparison-card-nav-button-right"
            >
              <ChevronRight class="w-5 h-5 text-gray-700" />
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
            <Image class="w-12 h-12 text-gray-300" />
          </div>
          
          <!-- Рейтинг badge -->
          <div v-if="item.rating !== undefined" class="comparison-card-rating-badge">
            <Badge class="comparison-card-rating-badge-content">
              <Star class="w-4 h-4 fill-yellow-400 text-yellow-400 mr-1" />
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
              class="hidden"
              multiple
            />
            <Button
              @click="$refs.fileInput.click()"
              size="sm"
              variant="outline"
              class="flex-1"
            >
              <Plus class="w-4 h-4 mr-1" />
              Добавить фото
            </Button>
            <Button
              @click="saveImages"
              size="sm"
              class="flex-1 bg-blue-600 hover:bg-blue-700"
            >
              <Save class="w-4 h-4 mr-1" />
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
            >
              <Plus class="w-4 h-4" />
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
                class="h-8 w-8 text-red-500 hover:text-red-700 hover:bg-red-50"
              >
                <X class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
        
        <!-- Заголовок -->
        <div class="comparison-card-title-container">
          <div v-if="isEditingTitle" class="comparison-card-title-editing">
            <Input
              v-model="title"
              class="text-base font-medium"
              @keyup.enter="saveTitle"
            />
            <Button size="icon" @click="saveTitle" class="bg-blue-600 hover:bg-blue-700">
              <Save class="w-4 h-4" />
            </Button>
          </div>
          <div v-else class="comparison-card-title-display">
            <div class="comparison-card-title-buttons">
              <Button
                variant="ghost"
                size="icon"
                @click="isEditingTitle = true"
                class="comparison-card-title-button"
              >
                <Edit2 class="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                @click="onDelete"
                class="comparison-card-delete-button"
              >
                <Trash2 class="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                @click="startEditingImages"
                class="comparison-card-image-edit-button"
              >
                <Image class="w-4 h-4" />
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
            <ExternalLink class="w-3 h-3" />
            Открыть ссылку
          </a>

          <!-- Цена -->
          <div class="comparison-card-price-container">
            <div v-if="isEditingPrice" class="comparison-card-price-editing">
              <Input
                type="number"
                v-model="price"
                placeholder="Цена"
                class="text-base"
                @keyup.enter="savePrice"
              />
              <Button size="icon" @click="savePrice" class="bg-blue-600 hover:bg-blue-700 h-9 w-9">
                <Save class="w-4 h-4" />
              </Button>
            </div>
            <button
              v-else
              @click="isEditingPrice = true"
              class="comparison-card-price-display"
            >
              <DollarSign class="comparison-card-price-icon" />
              <span v-if="item.price">{{ item.price.toLocaleString() }} {{ item.currency || '₽' }}</span>
              <span v-else>Указать цену</span>
            </button>
          </div>
        </div>
      </CardHeader>

      <CardContent class="comparison-card-content">
        <ProConItem
          type="pro"
          :items="item.pros || []"
          :on-update="updatePros"
        />
        
        <ProConItem
          type="con"
          :items="item.cons || []"
          :on-update="updateCons"
        />
      </CardContent>

      <!-- Статистика -->
      <div class="comparison-card-stats">
        <div class="comparison-card-stats-grid">
          <div>
            <div class="comparison-card-stat-value comparison-card-stat-value-pro">
              {{ (item.pros || []).reduce((sum, p) => sum + (p.impact || 5), 0) }}
            </div>
            <div class="comparison-card-stat-label">{{ (item.pros || []).length }} плюсов</div>
          </div>
          <div>
            <div class="comparison-card-stat-value comparison-card-stat-value-con">
              {{ (item.cons || []).reduce((sum, c) => sum + (c.impact || 5), 0) }}
            </div>
            <div class="comparison-card-stat-label">{{ (item.cons || []).length }} минусов</div>
          </div>
          <div>
            <div class="comparison-card-stat-value comparison-card-stat-value-rating">
              {{ item.rating?.toFixed(1) || '—' }}
            </div>
            <div class="comparison-card-stat-label">рейтинг</div>
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
    X
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
      
      return { ...itemData, rating: Math.round(finalRating * 10) / 10 }
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
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.1s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
