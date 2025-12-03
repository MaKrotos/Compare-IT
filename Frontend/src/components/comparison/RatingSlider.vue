<template>
  <div class="rating-control-group">
    <div class="rating-control-label">
      <component :is="icon" class="rating-control-icon" :class="iconClass" />
      <span class="rating-control-label-text">{{ label }}</span>
    </div>
    <div class="rating-control-slider-container">
      <Slider
        :model-value="[modelValue]"
        @update:model-value="updateValue"
        :min="0"
        :max="100"
        :step="1"
        class="rating-control-slider"
      />
      <span class="rating-control-value">{{ modelValue }} баллов</span>
    </div>
  </div>
</template>

<script>
import Slider from '@/components/ui/Slider.vue'

export default {
  name: 'RatingSlider',
  components: {
    Slider
  },
  props: {
    modelValue: {
      type: Number,
      required: true
    },
    icon: {
      type: Object,
      required: true
    },
    label: {
      type: String,
      required: true
    },
    iconClass: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const updateValue = (value) => {
      emit('update:modelValue', value)
    }

    return {
      updateValue
    }
  }
}
</script>

<style scoped>
.rating-control-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.rating-control-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.rating-control-icon {
  width: 1rem;
  height: 1rem;
}

.rating-control-price-icon {
  color: #f59e0b; /* amber-500 */
}

.rating-control-pros-icon {
  color: #8b5cf6; /* violet-500 */
}

.rating-control-label-text {
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
</style>