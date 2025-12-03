<template>
  <h3 class="dialog-title">
    <slot />
  </h3>
</template>

<script>
export default {
  name: 'DialogTitle'
}
</script>

<style scoped>
.dialog-title {
  font-size: 1.125rem; /* text-lg */
  font-weight: 600; /* font-semibold */
  line-height: 1; /* leading-none */
  letter-spacing: -0.025em; /* tracking-tight */
}
</style>