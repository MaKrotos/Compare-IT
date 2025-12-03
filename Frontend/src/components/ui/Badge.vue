<template>
  <div :class="['badge', variantClasses]">
    <slot />
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'Badge',
  props: {
    variant: {
      type: String,
      default: 'default'
    }
  },
  setup(props) {
    const variantClasses = computed(() => {
      switch (props.variant) {
        case 'default':
          return 'badge-default'
        case 'secondary':
          return 'badge-secondary'
        case 'destructive':
          return 'badge-destructive'
        case 'outline':
          return 'badge-outline'
        default:
          return 'badge-default'
      }
    })

    return {
      variantClasses
    }
  }
}
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  border-radius: 9999px;
  border: 1px solid;
  padding: 0.5rem 0.625rem;
  font-size: 0.75rem;
  font-weight: 600;
  transition: background-color 0.2s, border-color 0.2s;
  outline: none;
}

.badge:focus {
  outline: none;
  box-shadow: 0 0 0 2px var(--ring-color, rgba(0, 0, 0, 0.1));
}

.badge-default {
  border-color: transparent;
  background-color: var(--primary, #007bff);
  color: var(--primary-foreground, #ffffff);
}

.badge-default:hover {
  background-color: var(--primary-hover, rgba(0, 123, 255, 0.8));
}

.badge-secondary {
  border-color: transparent;
  background-color: var(--secondary, #6c757d);
  color: var(--secondary-foreground, #ffffff);
}

.badge-secondary:hover {
  background-color: var(--secondary-hover, rgba(108, 117, 125, 0.8));
}

.badge-destructive {
  border-color: transparent;
  background-color: var(--destructive, #dc3545);
  color: var(--destructive-foreground, #ffffff);
}

.badge-destructive:hover {
  background-color: var(--destructive-hover, rgba(220, 53, 69, 0.8));
}

.badge-outline {
  color: var(--foreground, #212529);
  border-color: currentColor;
  background-color: transparent;
}
</style>