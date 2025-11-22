<template>
  <div class="card">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'Card'
}
</script>

<style scoped>
/* Стили для компонента Card */
.card {
  border-radius: 1rem;
  border-width: 1px;
  border-color: #e5e7eb; /* gray-200 */
  background-color: #ffffff; /* white */
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transition-property: all;
  transition-duration: 200ms;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}

.card:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}
</style>