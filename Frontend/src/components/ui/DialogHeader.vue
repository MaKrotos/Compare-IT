<template>
  <div class="dialog-header">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'DialogHeader'
}
</script>

<style scoped>
.dialog-header {
  display: flex;
  flex-direction: column;
  gap: 0.375rem; /* space-y-1.5 = 1.5 * 0.25rem = 0.375rem */
  text-align: center;
  padding: 1.5rem; /* p-6 = 1.5rem */
  padding-bottom: 0; /* pb-0 */
}

@media (min-width: 640px) {
  .dialog-header {
    text-align: left; /* sm:text-left */
  }
}
</style>