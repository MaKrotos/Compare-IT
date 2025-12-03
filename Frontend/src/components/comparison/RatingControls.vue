<template>
  <Card class="rating-controls-container">
    <CardHeader class="rating-controls-header">
      <CardTitle class="rating-controls-title">
        <div class="rating-controls-icon-container">
          <Settings class="rating-controls-icon" />
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
        <RatingSlider
          :model-value="priceRatingWeight"
          @update:model-value="updatePriceRatingWeight"
          :icon="DollarSign"
          label="Рейтинг цены"
          icon-class="rating-control-price-icon"
        />

        <!-- Слайдер для рейтинга плюсов/минусов -->
        <RatingSlider
          :model-value="prosConsRatingWeight"
          @update:model-value="updateProsConsRatingWeight"
          :icon="ThumbsUp"
          label="Рейтинг плюсов/минусов"
          icon-class="rating-control-pros-icon"
        />

        <!-- Кнопка сохранения -->
        <Button
          @click="saveRatingWeights"
          class="rating-controls-save-button"
        >
          <Save class="rating-controls-save-icon" />
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
import { Settings, DollarSign, ThumbsUp, Save } from 'lucide-vue-next'
import RatingSlider from './RatingSlider.vue'

export default {
  name: 'RatingControls',
  components: {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
    Button,
    Settings,
    DollarSign,
    ThumbsUp,
    Save,
    RatingSlider
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
      priceRatingWeight.value = value
    }

    // Обновление веса рейтинга плюсов/минусов
    const updateProsConsRatingWeight = (value) => {
      prosConsRatingWeight.value = value
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
/* Стили для компонента RatingControls */
.rating-controls-container {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  border-radius: 1rem;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.rating-controls-header {
  border-bottom: 1px solid #f3f4f6;
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
  border-radius: 50%;
}

.rating-controls-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563eb;
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

.rating-controls-save-button {
  width: 100%;
  height: 2.75rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2563eb;
  color: white;
  border: none;
  cursor: pointer;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.rating-controls-save-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.rating-controls-save-button:hover:not(:disabled) {
  background-color: #1d4ed8;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}

.rating-controls-save-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
}
</style>