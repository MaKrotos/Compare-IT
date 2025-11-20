<template>
  <div v-if="open" class="dialog-overlay">
    <div class="dialog-backdrop" @click="onClose"></div>
    <div class="dialog-container">
      <div class="dialog-content" @click.stop>
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Dialog',
  props: {
    open: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:open'],
  setup(props, { emit }) {
    const onClose = () => {
      emit('update:open', false)
    }

    return {
      onClose
    }
  }
}
</script>

<style scoped>
.dialog-overlay {
  z-index: 100;
}
</style>