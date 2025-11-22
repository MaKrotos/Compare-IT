<template>
  <div class="slider">
    <input
      type="range"
      :min="min"
      :max="max"
      :step="step"
      :value="modelValue[0]"
      @input="handleInput"
      class="slider-input"
    >
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'Slider',
  props: {
    modelValue: {
      type: Array,
      required: true
    },
    min: {
      type: Number,
      default: 1
    },
    max: {
      type: Number,
      default: 100
    },
    step: {
      type: Number,
      default: 1
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const sliderRef = ref(null)

    const handleInput = (event) => {
      const value = parseFloat(event.target.value)
      emit('update:modelValue', [value])
    }

    const handleClick = (event) => {
      const rect = sliderRef.value.getBoundingClientRect()
      const percent = (event.clientX - rect.left) / rect.width
      const value = props.min + percent * (props.max - props.min)
      const steppedValue = Math.round(value / props.step) * props.step
      const clampedValue = Math.min(Math.max(steppedValue, props.min), props.max)
      emit('update:modelValue', [clampedValue])
    }

    return {
      sliderRef,
      handleInput,
      handleClick
    }
  }
}
</script>

<style scoped>
/* Стили для компонента Slider */
.slider {
  position: relative;
  display: flex;
  width: 100%;
  touch-action: none;
  user-select: none;
  align-items: center;
}

.slider-track {
  position: relative;
  height: 0.5rem;
  width: 100%;
  flex-grow: 1;
  overflow: hidden;
  border-radius: 9999px;
  background-color: #f3f4f6; /* secondary */
}

.slider-range {
  position: absolute;
  height: 100%;
  background-color: #3b82f6; /* primary */
}

.slider-input {
  inset: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}
</style>