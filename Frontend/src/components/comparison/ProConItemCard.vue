<template>
  <div
    :class="cardClasses"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
  >
    <div class="pro-con-item-content">
      <div :class="bulletClasses" />
      <div class="pro-con-item-text-container">
        <div class="pro-con-item-text-content">
          <p class="pro-con-item-text">{{ item.text }}</p>
          <Badge
            v-if="showImpactBadge"
            variant="outline"
            :class="impactBadgeClasses"
          >
            {{ item.impact || 5 }}/10
          </Badge>
        </div>
        <div v-if="showImpactSlider" class="pro-con-item-slider-container" :class="sliderContainerClasses">
          <div class="pro-con-item-slider-content">
            <div class="pro-con-item-slider-header">
              <div class="pro-con-item-slider-label">
                <TrendingUp class="pro-con-item-slider-icon" />
                Влияние
              </div>
              <Slider
                :model-value="[item.impact || 5]"
                @update:model-value="onUpdateImpact"
                :min="1"
                :max="10"
                :step="1"
                class="pro-con-item-slider"
              />
            </div>
          </div>
        </div>
      </div>
      <Button
        v-if="showDeleteButton"
        variant="ghost"
        size="icon"
        @click="onRemove"
        class="pro-con-item-delete-button"
      >
        <Trash2 class="pro-con-item-delete-icon" />
      </Button>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import Button from '@/components/ui/Button.vue'
import Badge from '@/components/ui/Badge.vue'
import Slider from '@/components/ui/Slider.vue'
import { Trash2, TrendingUp } from 'lucide-vue-next'

export default {
  name: 'ProConItemCard',
  components: {
    Button,
    Badge,
    Slider,
    Trash2,
    TrendingUp
  },
  props: {
    item: {
      type: Object,
      required: true
    },
    isPro: {
      type: Boolean,
      required: true
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  emits: ['remove', 'update:impact', 'mouseenter', 'mouseleave'],
  setup(props, { emit }) {
    const isHovered = ref(false)

    // Вычисляемые свойства для классов
    const cardClasses = computed(() => [
      'pro-con-item-card',
      props.isPro ? 'pro-con-item-card-pro' : (props.item.type === 'note' ? 'pro-con-item-card-note' : 'pro-con-item-card-con')
    ])

    const bulletClasses = computed(() => [
      'pro-con-item-bullet',
      props.item.type === 'note' ? 'pro-con-item-bullet-note' : (props.isPro ? 'pro-con-item-bullet-pro' : 'pro-con-item-bullet-con')
    ])

    const impactBadgeClasses = computed(() => [
      'pro-con-item-impact-badge',
      getImpactClass(props.item.impact || 5)
    ])

    const sliderContainerClasses = computed(() => ({
      'pro-con-item-slider-container': true,
      'pro-con-item-slider-container-open': isHovered.value
    }))

    // Условия отображения элементов
    const showImpactBadge = computed(() => props.item.impact !== undefined)
    const showImpactSlider = computed(() => props.item.impact !== undefined && !props.readonly)
    const showDeleteButton = computed(() => !props.readonly)

    // Функция для определения класса влияния
    const getImpactClass = (impact) => {
      if (impact <= 3) {
        return props.isPro ? 'pro-con-item-impact-badge-low-pro' : 'pro-con-item-impact-badge-low-con'
      }
      if (impact <= 6) {
        return props.isPro ? 'pro-con-item-impact-badge-medium-pro' : 'pro-con-item-impact-badge-medium-con'
      }
      return props.isPro ? 'pro-con-item-impact-badge-high-pro' : 'pro-con-item-impact-badge-high-con'
    }

    // Обработчики событий
    const onRemove = () => {
      emit('remove')
    }

    const onUpdateImpact = (value) => {
      emit('update:impact', value[0])
    }

    const onMouseEnter = () => {
      isHovered.value = true
      emit('mouseenter')
    }

    const onMouseLeave = () => {
      isHovered.value = false
      emit('mouseleave')
    }

    return {
      // Классы
      cardClasses,
      bulletClasses,
      impactBadgeClasses,
      sliderContainerClasses,
      // Условия отображения
      showImpactBadge,
      showImpactSlider,
      showDeleteButton,
      // Обработчики
      onRemove,
      onUpdateImpact,
      onMouseEnter,
      onMouseLeave
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ProConItemCard */
.pro-con-item-card {
  padding: 0.5rem;
  border-radius: 0.375rem;
  border: 1px solid;
}

.pro-con-item-card-note {
  background-color: rgba(224, 242, 254, 0.3) !important; /* sky-50 with 30% opacity */
  border-color: #e0f2fe !important; /* sky-100 */
}

.pro-con-item-bullet-note {
  background-color: #0ea5e9 !important; /* sky-500 */
}

.pro-con-item-card-wrapper.pro-con-item-card-note {
  background-color: rgba(224, 242, 254, 0.3) !important; /* sky-50 with 30% opacity */
  border-color: #e0f2fe !important; /* sky-100 */
}

.pro-con-item-card-wrapper .pro-con-item-bullet-note {
  background-color: #0ea5e9 !important; /* sky-500 */
}

div.pro-con-item-card-note {
  background-color: rgba(224, 242, 254, 0.3) !important; /* sky-50 with 30% opacity */
  border-color: #e0f2fe !important; /* sky-100 */
}

div.pro-con-item-bullet-note {
  background-color: #0ea5e9 !important; /* sky-500 */
}

.pro-con-item-card.pro-con-item-card-note {
  background-color: rgba(224, 242, 254, 0.3) !important; /* sky-50 with 30% opacity */
  border-color: #e0f2fe !important; /* sky-100 */
}

.pro-con-item-card .pro-con-item-bullet-note {
  background-color: #0ea5e9 !important; /* sky-500 */
}

.pro-con-item-card-pro {
  background-color: rgba(240, 253, 244, 0.3); /* green-50 with 30% opacity */
  border-color: #dcfce7; /* green-100 */
}

.pro-con-item-card-con {
  background-color: rgba(254, 242, 242, 0.3); /* red-50 with 30% opacity */
  border-color: #fee2e2; /* red-100 */
}

.pro-con-item-content {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.pro-con-item-bullet {
  margin-top: 0.375rem;
  width: 0.375rem;
  height: 0.375rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.pro-con-item-bullet-pro {
  background-color: #22c55e; /* green-500 */
}

.pro-con-item-bullet-con {
  background-color: #ef4444; /* red-500 */
}

.pro-con-item-text-container {
  flex: 1;
  min-width: 0;
}

.pro-con-item-text-content {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.5rem;
}

.pro-con-item-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #1f2937; /* gray-800 */
  line-height: 1.625;
  flex: 1;
  word-break: break-word;
}

/* Стили для бейджей влияния */
.pro-con-item-impact-badge {
  font-size: 0.75rem;
  line-height: 1rem;
  border: 1px solid;
  font-weight: 600;
  white-space: nowrap;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
}

.pro-con-item-impact-badge-low-pro {
  background-color: #f3f4f6; /* gray-100 */
  color: #374151; /* gray-700 */
  border-color: #e5e7eb; /* gray-200 */
}

.pro-con-item-impact-badge-low-con {
  background-color: #f3f4f6; /* gray-100 */
  color: #374151; /* gray-700 */
  border-color: #e5e7eb; /* gray-200 */
}

.pro-con-item-impact-badge-medium-pro {
  background-color: #fef9c3; /* yellow-100 */
  color: #a16207; /* yellow-700 */
  border-color: #fef08a; /* yellow-200 */
}

.pro-con-item-impact-badge-medium-con {
  background-color: #fef9c3; /* yellow-100 */
  color: #a16207; /* yellow-700 */
  border-color: #fef08a; /* yellow-200 */
}

.pro-con-item-impact-badge-high-pro {
  background-color: #dcfce7; /* green-100 */
  color: #166534; /* green-700 */
  border-color: #bbf7d0; /* green-200 */
}

.pro-con-item-impact-badge-high-con {
  background-color: #fee2e2; /* red-100 */
  color: #b91c1c; /* red-700 */
  border-color: #fecaca; /* red-200 */
}

/* Стили для слайдера */
.pro-con-item-slider-container {
  transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 0;
  overflow: hidden;
}

.pro-con-item-slider-container-open {
  max-height: 60px;
}

.pro-con-item-slider-content {
  margin-top: 0.5rem;
}

.pro-con-item-slider-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  width: 100%;
}

.pro-con-item-slider-label {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280; /* gray-500 */
  gap: 0.25rem;
}

.pro-con-item-slider-icon {
  width: 0.75rem;
  height: 0.75rem;
}

.pro-con-item-slider {
  flex: 1;
}

/* Кнопка удаления */
.pro-con-item-delete-button {
  height: 1.5rem;
  width: 1.5rem;
  color: #9ca3af; /* gray-400 */
  border-radius: 50%;
  flex-shrink: 0;
  padding: 0;
}

.pro-con-item-delete-button:hover {
  color: #ef4444; /* red-500 */
  background-color: #fef2f2; /* red-50 */
}

.pro-con-item-delete-icon {
  width: 0.75rem;
  height: 0.75rem;
}
</style>