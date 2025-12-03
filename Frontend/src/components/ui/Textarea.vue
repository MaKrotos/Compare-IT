<template>
  <textarea
    class="textarea"
    v-bind="$attrs"
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
  >
  </textarea>
</template>

<script>
export default {
  name: 'Textarea',
  inheritAttrs: false,
  props: {
    modelValue: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue']
}
</script>

<style scoped>
.textarea {
  display: flex;
  min-height: 80px;
  width: 100%;
  border-radius: 0.375rem;
  border: 1px solid var(--border-color, #d1d5db);
  background-color: var(--background-color, #ffffff);
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  background-color: var(--background-color, #ffffff);
}

.textarea::placeholder {
  color: var(--muted-foreground, #6b7280);
}

.textarea:focus {
  outline: none;
  box-shadow: 
    0 0 0 2px var(--ring-offset-background, #ffffff),
    0 0 0 4px var(--ring-color, rgba(59, 130, 246, 0.5));
}

.textarea:focus-visible {
  outline: none;
  box-shadow: 
    0 0 0 2px var(--ring-offset-background, #ffffff),
    0 0 0 4px var(--ring-color, rgba(59, 130, 246, 0.5));
}

.textarea:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>