<template>
  <Card class="add-link-form-container">
    <CardHeader class="add-link-form-header">
      <CardTitle class="add-link-form-title">
        <div class="add-link-form-icon-container">
          <Link2 class="add-link-form-icon" />
        </div>
        Добавить товары для сравнения
      </CardTitle>
      <p class="add-link-form-subtitle">
        Вставьте ссылки на товары (до 10 штук)
      </p>
    </CardHeader>
    <CardContent class="add-link-form-content">
      <div class="add-link-form-fields">
        <transition-group name="list" tag="div" class="add-link-form-fields">
          <ProductLinkField v-for="(link, index) in links" :key="index" :link="link" :index="index"
            :is-loading="isLoading" :show-remove-button="links.length > 1" :is-drag-over="isDragOver"
            @update:link="(newLink) => updateLink(index, newLink)" @remove="() => removeLink(index)"
            @paste="(event) => handleImagePaste(index, event)" @dragenter="handleDragEnter" @dragleave="handleDragLeave"
            @drop="(event) => handleDrop(event, index)" class="add-link-form-field-group" />
        </transition-group>

        <div class="add-link-form-button-container">
          <Button v-if="links.length < 10" type="button" variant="outline" @click="addLinkField"
            class="add-link-form-add-button" :disabled="isLoading">
            <Plus class="add-link-form-button-icon" />
            Добавить ссылку
          </Button>
          <Button type="button" @click="() => { console.log('Кнопка загрузки нажата'); handleAddItems(onAddLinks); }"
            :disabled="isLoading || !isFormValid" class="add-link-form-submit-button">
            <Loader2 v-if="isLoading" class="add-link-form-loading-icon" />
            <span v-if="isLoading" class="add-link-form-loading-text">Загрузка...</span>
            <span v-else class="add-link-form-submit-text">Загрузить товары</span>
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<script>
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import CardContent from '@/components/ui/CardContent.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardTitle from '@/components/ui/CardTitle.vue'
import ProductLinkField from './ProductLinkField.vue'
import { Plus, Link2, Loader2 } from 'lucide-vue-next'
import { useAddLinkForm } from '@/composables/useAddLinkForm.js'

export default {
  name: 'AddLinkForm',
  components: {
    Button,
    Card,
    CardContent,
    CardHeader,
    CardTitle,
    ProductLinkField,
    Plus,
    Link2,
    Loader2
  },
  props: {
    onAddLinks: {
      type: Function,
      required: true
    }
  },
  setup(props) {
    const {
      links,
      isDragOver,
      isLoading,
      updateLink,
      addLinkField,
      removeLink,
      handleAddItems,
      handleImagePaste,
      handleDragEnter,
      handleDragLeave,
      handleDrop,
      isFormValid
    } = useAddLinkForm()

    return {
      links,
      isDragOver,
      isLoading,
      updateLink,
      addLinkField,
      removeLink,
      handleAddItems,
      handleImagePaste,
      handleDragEnter,
      handleDragLeave,
      handleDrop,
      isFormValid,
      onAddLinks: props.onAddLinks
    }
  }
}
</script>

<style scoped>
/* Стили для компонента AddLinkForm */
.add-link-form-container {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  border-radius: 1rem;
  overflow: hidden;
}

.add-link-form-header {
  border-bottom: 1px solid #f3f4f6;
  padding-bottom: 1.25rem;
  background-image: linear-gradient(to right, rgba(219, 234, 254, 0.8), rgba(224, 231, 255, 0.8));
}

.add-link-form-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  letter-spacing: -0.025em;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #1f2937;
}

.add-link-form-icon-container {
  padding: 0.5rem;
  background-color: #dbeafe;
  border-radius: 50%;
}

.add-link-form-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563eb;
}

.add-link-form-subtitle {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #4b5563;
  margin-top: 0.375rem;
}

.add-link-form-content {
  padding-top: 1.5rem;
}

.add-link-form-fields {
  margin-bottom: 1rem;
}

.add-link-form-fields> :not(:last-child) {
  margin-bottom: 1rem;
}

.add-link-form-button-container {
  display: flex;
  gap: 0.75rem;
  padding-top: 0.5rem;
}

.add-link-form-add-button {
  flex: 1;
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  height: 3rem;
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.add-link-form-add-button:hover {
  border-color: #3b82f6;
  background-color: #dbeafe;
}

.add-link-form-submit-button {
  flex: 1;
  height: 3rem;
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.add-link-form-submit-button:hover {
  background-color: #1d4ed8;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
}

.add-link-form-button-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
}

.add-link-form-loading-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
  animation: spin 1s linear infinite;
}

.add-link-form-loading-text {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
}

.add-link-form-submit-text {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>