<template>
  <div class="card-header">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'CardHeader'
}
</script>

<style scoped>
.card-header {
  display: flex;
  flex-direction: column;
  gap: 0.375rem; /* 1.5 * 0.25rem = 0.375rem */
}
</style>
