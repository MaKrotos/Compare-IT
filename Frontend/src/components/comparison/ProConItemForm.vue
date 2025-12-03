<template>
  <div class="pro-con-item-add-container">
    <Input
      v-model="newItem"
      :placeholder="type === 'pro' ? 'Введите плюс...' : type === 'con' ? 'Введите минус...' : 'Введите заметку...'"
      class="pro-con-item-input"
      @keyup.enter="addItem"
    />
    <Button
      size="icon"
      @click="addItem"
      :class="['pro-con-item-save-button', isPro ? 'pro-con-item-save-button-pro' : 'pro-con-item-save-button-con']"
    >
      <Check class="pro-con-item-save-icon" />
    </Button>
    <Button
      variant="ghost"
      size="icon"
      @click="cancelAdd"
      class="pro-con-item-cancel-button"
    >
      <X class="pro-con-item-cancel-icon" />
    </Button>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { Check, X } from 'lucide-vue-next'

export default {
  name: 'ProConItemForm',
  components: {
    Button,
    Input,
    Check,
    X
  },
  props: {
    isPro: {
      type: Boolean,
      required: true
    },
    type: {
      type: String,
      default: 'pro'
    },
    newItem: {
      type: String,
      default: ''
    }
  },
  emits: ['update:newItem', 'add', 'cancel'],
  setup(props, { emit }) {
    const newItem = ref(props.newItem || '')

    const addItem = () => {
     console.log('addItem called')
     console.log('newItem.value:', newItem.value)
     if (newItem.value && newItem.value.trim()) {
       emit('add', newItem.value.trim())
       newItem.value = ''
     }
   }

   // Добавим watch для отслеживания изменений newItem
   watch(newItem, (newVal) => {
     console.log('newItem changed:', newVal)
     emit('update:newItem', newVal)
   })

   // Добавим watch для props.newItem
   watch(() => props.newItem, (newVal) => {
     if (newVal !== undefined) {
       newItem.value = newVal
     }
   })

    const cancelAdd = () => {
      newItem.value = ''
      emit('cancel')
    }

    return {
      newItem,
      addItem,
      cancelAdd
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ProConItemForm */
.pro-con-item-add-container {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.pro-con-item-input {
  font-size: 0.875rem;
  line-height: 1.25rem;
  height: 2.25rem;
  flex: 1;
  padding: 0 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
}

.pro-con-item-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.pro-con-item-save-button {
  height: 2.25rem;
  width: 2.25rem;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pro-con-item-save-button-pro {
  background-color: #16a34a; /* green-600 */
  color: white;
}

.pro-con-item-save-button-pro:hover {
  background-color: #15803d; /* green-700 */
}

.pro-con-item-save-button-con {
  background-color: #dc2626; /* red-600 */
  color: white;
}

.pro-con-item-save-button-con:hover {
  background-color: #b91c1c; /* red-700 */
}

.pro-con-item-save-icon {
  width: 1rem;
  height: 1rem;
}

.pro-con-item-cancel-button {
  height: 2.25rem;
  width: 2.25rem;
  border-radius: 50%;
  background-color: white;
  border: 1px solid #d1d5db;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pro-con-item-cancel-button:hover {
  background-color: #f3f4f6;
}

.pro-con-item-cancel-icon {
  width: 1rem;
  height: 1rem;
  color: #6b7280;
}
</style>