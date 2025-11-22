<template>
  <TelegramAuth>
    <router-view />
  </TelegramAuth>
</template>

<script>
import TelegramAuth from '@/components/TelegramAuth.vue'

export default {
  name: 'App',
  components: {
    TelegramAuth
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

* {
  font-family: 'Inter', sans-serif;
}
</style>