<template>
  <div class="comparison-card-title-container">
    <div v-if="isEditingTitle" class="comparison-card-title-editing">
      <Input
        v-model="title"
        class="comparison-card-title-input"
        @keyup.enter="saveTitle"
      />
      <Button size="icon" @click="saveTitle" class="comparison-card-save-title-button">
        <Save class="comparison-card-save-title-icon" />
      </Button>
    </div>
    <div v-else class="comparison-card-title-display">
      <div v-if="!readonly" class="comparison-card-title-buttons">
        <Button
          variant="ghost"
          size="icon"
          @click="isEditingTitle = true"
          class="comparison-card-title-button"
        >
          <Edit2 class="comparison-card-edit-icon" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          @click="onDelete"
          class="comparison-card-delete-button"
        >
          <Trash2 class="comparison-card-trash-icon" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          @click="startEditingImages"
          class="comparison-card-image-edit-button"
        >
          <Image class="comparison-card-image-icon" />
        </Button>
      </div>
      <h3 class="comparison-card-title-text">
        {{ title }}
      </h3>
    </div>
    <a
      :href="item.url"
      target="_blank"
      rel="noopener noreferrer"
      class="comparison-card-link"
    >
      <ExternalLink class="comparison-card-link-icon" />
      Открыть ссылку
    </a>

    <!-- Цена -->
    <div class="comparison-card-price-container">
      <div v-if="isEditingPrice" class="comparison-card-price-editing">
        <Input
          type="number"
          v-model="price"
          placeholder="Цена"
          class="comparison-card-price-input"
          @keyup.enter="savePrice"
        />
        <Button size="icon" @click="savePrice" class="comparison-card-save-price-button">
          <Save class="comparison-card-save-price-icon" />
        </Button>
      </div>
      <button
        v-else
        @click="isEditingPrice = true"
        :disabled="readonly"
        class="comparison-card-price-display"
        :class="{ 'comparison-card-price-disabled': readonly }"
      >
        <DollarSign class="comparison-card-price-icon" />
        <span v-if="item.price" class="comparison-card-price-value">{{ item.price.toLocaleString() }} {{ item.currency || '₽' }}</span>
        <span v-else class="comparison-card-price-placeholder">Указать цену</span>
      </button>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { ExternalLink, Trash2, Edit2, Save, DollarSign, Image } from 'lucide-vue-next'

export default {
  name: 'ProductHeader',
  components: {
    Button,
    Input,
    ExternalLink,
    Trash2,
    Edit2,
    Save,
    DollarSign,
    Image
  },
  props: {
    item: {
      type: Object,
      required: true
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update', 'delete', 'edit-images'],
  setup(props, { emit }) {
    const isEditingTitle = ref(false)
    const title = ref(props.item.title || 'Без названия')
    const isEditingPrice = ref(false)
    const price = ref(props.item.price || 0)
    
    // Наблюдение за изменением props.item для обновления локальных данных
    watch(() => props.item, (newItem) => {
      title.value = newItem.title || 'Без названия'
      price.value = newItem.price || 0
    }, { deep: true })
    
    const saveTitle = () => {
      emit('update', { ...props.item, title: title.value })
      isEditingTitle.value = false
    }
    
    const savePrice = () => {
      const updatedItem = { ...props.item, price: parseFloat(price.value) || 0 }
      emit('update', updatedItem)
      isEditingPrice.value = false
    }
    
    const onDelete = () => {
      emit('delete')
    }
    
    const startEditingImages = () => {
      emit('edit-images')
    }
    
    return {
      isEditingTitle,
      title,
      isEditingPrice,
      price,
      saveTitle,
      savePrice,
      onDelete,
      startEditingImages
    }
  }
}
</script>

<style scoped>
/* Заголовок */
.comparison-card-title-container {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.comparison-card-title-editing {
  display: flex;
  gap: 0.5rem;
}

.comparison-card-title-input {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  flex: 1;
}

.comparison-card-save-title-button {
  height: 2.5rem;
  width: 2.5rem;
  background-color: #2563eb;
}

.comparison-card-save-title-button:hover {
  background-color: #1d4ed8;
}

.comparison-card-save-title-icon {
  width: 1rem;
  height: 1rem;
  color: #ffffff;
}

.comparison-card-title-display {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.comparison-card-title-text {
  font-size: 1.125rem;
  line-height: 1.375;
  font-weight: 500;
  color: #111827;
  flex: 1;
  word-wrap: break-word;
  overflow-wrap: break-word;
  word-break: break-word;
}

.comparison-card-title-buttons {
  display: flex;
  gap: 0.25rem;
}

.comparison-card-title-button,
.comparison-card-delete-button,
.comparison-card-image-edit-button {
  height: 2rem;
  width: 2rem;
  color: #9ca3af;
  border-radius: 9999px;
}

.comparison-card-title-button:hover {
  color: #2563eb;
  background-color: #dbeafe;
}

.comparison-card-delete-button:hover {
  color: #dc2626;
  background-color: #fef2f2;
}

.comparison-card-image-edit-button:hover {
  color: #16a34a;
  background-color: #dcfce7;
}

.comparison-card-edit-icon,
.comparison-card-trash-icon,
.comparison-card-image-icon {
  width: 1rem;
  height: 1rem;
}

/* Ссылка на товар */
.comparison-card-link {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #2563eb;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-top: 0.5rem;
  text-decoration: none;
}

.comparison-card-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.comparison-card-link-icon {
  width: 0.75rem;
  height: 0.75rem;
}

/* Цена */
.comparison-card-price-container {
  margin-top: 0.75rem;
}

.comparison-card-price-editing {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.comparison-card-price-input {
  font-size: 1.125rem;
  line-height: 1.75rem;
  flex: 1;
}

.comparison-card-save-price-button {
  height: 2.25rem;
  width: 2.25rem;
  background-color: #2563eb;
}

.comparison-card-save-price-button:hover {
  background-color: #1d4ed8;
}

.comparison-card-save-price-icon {
  width: 1rem;
  height: 1rem;
  color: #ffffff;
}

.comparison-card-price-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.125rem;
  line-height: 1.75rem;
  font-weight: 600;
  color: #2563eb;
  border: none;
  background: none;
  cursor: pointer;
  padding: 0;
  transition: color 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-price-display:hover:not(.comparison-card-price-disabled) {
  color: #1d4ed8;
}

.comparison-card-price-disabled {
  cursor: default;
}

.comparison-card-price-icon {
  width: 1.25rem;
  height: 1.25rem;
  transition: transform 150ms cubic-bezier(0.4, 0, 0.2, 1);
}

.comparison-card-price-display:hover:not(.comparison-card-price-disabled) .comparison-card-price-icon {
  transform: scale(1.1);
}

.comparison-card-price-value,
.comparison-card-price-placeholder {
  font-weight: 600;
}
</style>