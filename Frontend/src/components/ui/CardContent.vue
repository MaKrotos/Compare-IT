<template>
  <div class="card-content">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'CardContent'
}
</script>

<style scoped>
.card-content {
  padding-top: 0;
}
</style>