<template>
  <div @click="handleClick">
    <slot />
  </div>
</template>

<script>
export default {
  name: 'DialogTrigger',
  emits: ['click'],
  setup(props, { emit }) {
    const handleClick = (event) => {
      emit('click', event)
    }

    return {
      handleClick
    }
  }
}
</script>