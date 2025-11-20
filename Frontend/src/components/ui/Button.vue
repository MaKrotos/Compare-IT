<template>
  <button
    :class="[
      'btn',
      variantClass,
      sizeClass
    ]"
  >
    <slot />
  </button>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'Button',
  props: {
    variant: {
      type: String,
      default: 'default'
    },
    size: {
      type: String,
      default: 'default'
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const variantClass = computed(() => {
      switch (props.variant) {
        case 'default':
          return 'btn-default'
        case 'destructive':
          return 'btn-destructive'
        case 'outline':
          return 'btn-outline'
        case 'secondary':
          return 'btn-secondary'
        case 'ghost':
          return 'btn-ghost'
        case 'link':
          return 'btn-link'
        default:
          return 'btn-default'
      }
    })

    const sizeClass = computed(() => {
      switch (props.size) {
        case 'sm':
          return 'btn-sm'
        case 'lg':
          return 'btn-lg'
        case 'icon':
          return 'btn-icon'
        default:
          return 'btn-default-size'
      }
    })

    return {
      variantClass,
      sizeClass
    }
  }
}
</script>