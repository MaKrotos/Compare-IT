<template>
  <h3 class="card-title">
    <slot />
  </h3>
</template>

<script>
export default {
  name: 'CardTitle'
}
</script>

<style scoped>
.card-title {
  font-size: 1.5rem; /* text-2xl */
  font-weight: 600; /* font-semibold */
  line-height: 1; /* leading-none */
  letter-spacing: -0.025em; /* tracking-tight */
}
</style>