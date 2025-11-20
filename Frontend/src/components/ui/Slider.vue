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