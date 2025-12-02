<template>
  <Card class="rating-controls-container">
    <CardHeader class="rating-controls-header">
      <CardTitle class="rating-controls-title">
        <div class="rating-controls-icon-container">
          <Settings class="w-5 h-5 text-blue-600" />
        </div>
        Настройка рейтингов
      </CardTitle>
      <p class="rating-controls-subtitle">
        Распределите баллы между рейтингами цены и плюсов/минусов
      </p>
    </CardHeader>
    <CardContent class="rating-controls-content">
      <div class="rating-controls-sliders">
        <!-- Слайдер для рейтинга цены -->
        <div class="rating-control-group">
          <div class="rating-control-label">
            <DollarSign class="w-4 h-4 text-amber-500" />
            <span>Рейтинг цены</span>
          </div>
          <div class="rating-control-slider-container">
            <Slider
              :model-value="[priceRatingWeight]"
              @update:model-value="updatePriceRatingWeight"
              :min="0"
              :max="100"
              :step="1"
              class="rating-control-slider"
            />
            <span class="rating-control-value">{{ priceRatingWeight }} баллов</span>
          </div>
        </div>

        <!-- Слайдер для рейтинга плюсов/минусов -->
        <div class="rating-control-group">
          <div class="rating-control-label">
            <ThumbsUp class="w-4 h-4 text-violet-500" />
            <span>Рейтинг плюсов/минусов</span>
          </div>
          <div class="rating-control-slider-container">
            <Slider
              :model-value="[prosConsRatingWeight]"
              @update:model-value="updateProsConsRatingWeight"
              :min="0"
              :max="100"
              :step="1"
              class="rating-control-slider"
            />
            <span class="rating-control-value">{{ prosConsRatingWeight }} баллов</span>
          </div>
        </div>

        <!-- Кнопка сохранения -->
        <Button
          @click="saveRatingWeights"
          class="rating-controls-save-button"
        >
          <Save class="w-4 h-4 mr-2" />
          Сохранить настройки
        </Button>
      </div>
    </CardContent>
  </Card>
</template>

<script>
import { ref, watch } from 'vue'
import Card from '@/components/ui/Card.vue'
import CardContent from '@/components/ui/CardContent.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardTitle from '@/components/ui/CardTitle.vue'
import Button from '@/components/ui/Button.vue'
import Slider from '@/components/ui/Slider.vue'
import { Settings, DollarSign, ThumbsUp, Save } from 'lucide-vue-next'

export default {
  name: 'RatingControls',
  components: {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
    Button,
    Slider,
    Settings,
    DollarSign,
    ThumbsUp,
    Save
  },
  props: {
    initialPriceWeight: {
      type: Number,
      default: 20
    },
    initialProsConsWeight: {
      type: Number,
      default: 80
    },
    onSave: {
      type: Function,
      required: true
    }
  },
  setup(props) {
    const priceRatingWeight = ref(props.initialPriceWeight || 20)
    const prosConsRatingWeight = ref(props.initialProsConsWeight || 80)

    // Следим за изменением входных параметров
    watch(() => props.initialPriceWeight, (newValue) => {
      priceRatingWeight.value = newValue || 20;
    });

    watch(() => props.initialProsConsWeight, (newValue) => {
      prosConsRatingWeight.value = newValue || 80;
    });

    // Обновление веса рейтинга цены
    const updatePriceRatingWeight = (value) => {
      priceRatingWeight.value = value[0]
    }

    // Обновление веса рейтинга плюсов/минусов
    const updateProsConsRatingWeight = (value) => {
      prosConsRatingWeight.value = value[0]
    }

    // Сохранение настроек рейтингов
    const saveRatingWeights = () => {
      props.onSave({
        priceRatingWeight: priceRatingWeight.value,
        prosConsRatingWeight: prosConsRatingWeight.value
      })
    }

    return {
      priceRatingWeight,
      prosConsRatingWeight,
      updatePriceRatingWeight,
      updateProsConsRatingWeight,
      saveRatingWeights
    }
  }
}
</script>

<style scoped>
.rating-controls-container {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  background-color: #ffffff;
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  border-radius: 1rem;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.rating-controls-header {
  border-bottom-width: 1px;
  border-color: #f3f4f6;
  padding-bottom: 1.25rem;
  background-image: linear-gradient(to right, rgba(219, 234, 254, 0.8), rgba(224, 231, 255, 0.8));
}

.rating-controls-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #1f2937;
}

.rating-controls-icon-container {
  padding: 0.5rem;
  background-color: #dbeafe;
  border-radius: 9999px;
}

.rating-controls-subtitle {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
  margin-top: 0.375rem;
}

.rating-controls-content {
  padding: 1.5rem;
}

.rating-controls-sliders {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.rating-control-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.rating-control-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #374151;
}

.rating-control-slider-container {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rating-control-slider {
  flex: 1;
}

.rating-control-value {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  min-width: 3rem;
  text-align: right;
}

.rating-controls-save-button {
  width: 100%;
  height: 2.75rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
}

.rating-controls-save-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.rating-controls-save-button:hover:not(:disabled) {
  background-image: linear-gradient(to right, #2563eb, #2563eb);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
  transition-property: all;
  transition-duration: 150ms;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}
</style>