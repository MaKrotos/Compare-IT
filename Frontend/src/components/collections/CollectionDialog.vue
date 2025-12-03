<template>
  <Dialog v-model:open="isOpen" @open-change="onOpenChange">
    <DialogContent class="collection-dialog">
      <DialogHeader>
        <DialogTitle class="collection-dialog-title">{{ isEditing ? 'Редактировать подборку' : 'Создать подборку' }}</DialogTitle>
      </DialogHeader>
      <div class="collection-dialog-content">
        <div class="collection-form">
          <div>
            <label class="collection-form-label">
              Название подборки
            </label>
            <Input
              v-model="localForm.name"
              placeholder="Введите название подборки"
              class="collection-form-input"
              @focus="$event.target.select()"
              ref="collectionNameInput"
            />
          </div>
        </div>
      </div>
      <div class="collection-dialog-footer">
        <Button variant="outline" @click="closeDialog" class="collection-dialog-cancel">Отмена</Button>
        <Button @click="saveCollection" class="collection-dialog-save">
          {{ isEditing ? 'Сохранить' : 'Создать' }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script>
import { ref, watch } from 'vue'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

export default {
  name: 'CollectionDialog',
  components: {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    Button,
    Input
  },
  props: {
    open: {
      type: Boolean,
      default: false
    },
    editingCollection: {
      type: Object,
      default: null
    }
  },
  emits: ['update:open', 'save'],
  setup(props, { emit }) {
    const isOpen = ref(false)
    const localForm = ref({
      name: ''
    })
    const collectionNameInput = ref(null)
    const isEditing = ref(false)

    // Следим за изменением props.open
    watch(() => props.open, (newVal) => {
      isOpen.value = newVal
    }, { immediate: true })

    // Следим за изменением editingCollection
    watch(() => props.editingCollection, (newVal) => {
      if (newVal) {
        isEditing.value = true
        localForm.value = { name: newVal.name }
      } else {
        isEditing.value = false
        localForm.value = { name: '' }
      }
    }, { immediate: true })

    // Обработчик открытия/закрытия диалога
    const onOpenChange = (open) => {
      isOpen.value = open
      emit('update:open', open)
      
      if (open) {
        // Фокус на поле ввода при открытии диалога
        setTimeout(() => {
          if (collectionNameInput.value) {
            collectionNameInput.value.focus()
          }
        }, 100)
      }
    }

    // Закрытие диалога
    const closeDialog = () => {
      isOpen.value = false
      emit('update:open', false)
    }

    // Сохранение подборки
    const saveCollection = () => {
      emit('save', {
        ...props.editingCollection,
        name: localForm.value.name
      })
    }

    return {
      isOpen,
      localForm,
      collectionNameInput,
      isEditing,
      onOpenChange,
      closeDialog,
      saveCollection
    }
  }
}
</script>

<style scoped>
/* Модальное окно */
.collection-dialog {
  max-width: 28rem;
}

.collection-dialog-title {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
  color: #111827;
}

.collection-dialog-content {
  padding: 1rem 0;
}

.collection-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.collection-form-label {
  display: block;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.25rem;
}

.collection-form-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.collection-form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.collection-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.collection-dialog-cancel {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  background-color: white;
  color: #374151;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  cursor: pointer;
}

.collection-dialog-cancel:hover {
  background-color: #f9fafb;
}

.collection-dialog-save {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  background-color: #2563eb;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
}

.collection-dialog-save:hover {
  background-color: #1d4ed8;
}
</style>