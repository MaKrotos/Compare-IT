<template>
  <div class="dialog-content-inner">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'DialogContent'
}
</script>

<style scoped>
.dialog-content-inner {
  padding: 1.5rem; /* p-6 = 1.5rem */
}
</style>