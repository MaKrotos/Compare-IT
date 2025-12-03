<template>
  <div class="pro-con-item-container">
    <div class="pro-con-item-header">
      <h4 :class="['pro-con-item-title', isPro ? 'pro-con-item-title-pro' : 'pro-con-item-title-con']">
        {{ isPro ? 'Плюсы' : 'Минусы' }}
      </h4>
      <Button
        v-if="!isAdding && !readonly"
        variant="ghost"
        size="sm"
        @click="isAdding = true"
        :class="['pro-con-item-add-button', isPro ? 'pro-con-item-add-button-pro' : 'pro-con-item-add-button-con']"
      >
        <Plus class="pro-con-item-add-icon" />
        Добавить
      </Button>
    </div>

    <div class="pro-con-item-list">
      <transition-group name="pro-con-item-list" tag="div">
        <div
          v-for="(item, index) in items"
          :key="index"
          :class="['pro-con-item-card', isPro ? 'pro-con-item-card-pro' : 'pro-con-item-card-con']"
          @mouseenter="hoveredItem = index"
          @mouseleave="hoveredItem = null"
        >
          <div class="pro-con-item-content">
            <div :class="['pro-con-item-bullet', isPro ? 'pro-con-item-bullet-pro' : 'pro-con-item-bullet-con']" />
            <div class="pro-con-item-text-container">
              <div class="pro-con-item-text-content">
                <p class="pro-con-item-text">{{ item.text }}</p>
                <Badge
                  variant="outline"
                  :class="['pro-con-item-impact-badge', getImpactClass(item.impact || 5)]"
                >
                  {{ item.impact || 5 }}/10
                </Badge>
              </div>
              <div class="pro-con-item-slider-container" :class="{ 'pro-con-item-slider-container-open': hoveredItem === index }">
                <div class="pro-con-item-slider-content">
                  <div class="pro-con-item-slider-header">
                    <div class="pro-con-item-slider-label">
                      <TrendingUp class="pro-con-item-slider-icon" />
                      Влияние
                    </div>
                    <Slider
                      :model-value="[item.impact || 5]"
                      @update:model-value="val => updateImpact(index, val[0])"
                      :min="1"
                      :max="10"
                      :step="1"
                      class="pro-con-item-slider"
                    />
                  </div>
                </div>
              </div>
            </div>
            <Button
              v-if="!readonly"
              variant="ghost"
              size="icon"
              @click="removeItem(index)"
              class="pro-con-item-delete-button"
            >
              <Trash2 class="pro-con-item-delete-icon" />
            </Button>
          </div>
        </div>
      </transition-group>
    </div>

    <transition name="pro-con-item-fade">
      <div v-if="isAdding" class="pro-con-item-add-container">
        <Input
          v-model="newItem"
          :placeholder="isPro ? 'Введите плюс...' : 'Введите минус...'"
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
    </transition>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import  Button  from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import  Badge from '@/components/ui/Badge.vue'
import  Slider  from '@/components/ui/Slider.vue'
import { Plus, Trash2, Check, X, TrendingUp } from 'lucide-vue-next'

export default {
  name: 'ProConItem',
  components: {
    Button,
    Input,
    Badge,
    Slider,
    Plus,
    Trash2,
    Check,
    X,
    TrendingUp
  },
  props: {
    type: {
      type: String,
      required: true
    },
    items: {
      type: Array,
      default: () => []
    },
    onUpdate: {
      type: Function,
      required: true
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const newItem = ref('')
    const isAdding = ref(false)
    const isPro = computed(() => props.type === 'pro')
    const hoveredItem = ref(null)

    const getImpactClass = (impact) => {
      if (impact <= 3) {
        return isPro.value ? 'pro-con-item-impact-badge-low-pro' : 'pro-con-item-impact-badge-low-con'
      }
      if (impact <= 6) {
        return isPro.value ? 'pro-con-item-impact-badge-medium-pro' : 'pro-con-item-impact-badge-medium-con'
      }
      return isPro.value ? 'pro-con-item-impact-badge-high-pro' : 'pro-con-item-impact-badge-high-con'
    }

    const addItem = () => {
      console.log('addItem called')
      console.log('newItem.value:', newItem.value)
      if (newItem.value.trim()) {
        const newItems = [...props.items, { text: newItem.value.trim(), impact: 5 }]
        console.log('Adding item:', newItems)
        props.onUpdate(newItems)
        newItem.value = ''
        isAdding.value = false
      }
    }

    const removeItem = (index) => {
      props.onUpdate(props.items.filter((_, i) => i !== index))
    }

    const updateImpact = (index, value) => {
      const updated = [...props.items]
      updated[index] = { ...updated[index], impact: value }
      props.onUpdate(updated)
    }

    const cancelAdd = () => {
      newItem.value = ''
      isAdding.value = false
    }

    return {
      newItem,
      isAdding,
      hoveredItem,
      isPro,
      getImpactClass,
      addItem,
      removeItem,
      updateImpact,
      cancelAdd
    }
  }
}
</script>

<style scoped>
/* Стили для компонента ProConItem */
.pro-con-item-container {
  margin-bottom: 0.75rem;
}

.pro-con-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.pro-con-item-title {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
}

.pro-con-item-title-pro {
  color: #166534; /* green-700 */
}

.pro-con-item-title-con {
  color: #b91c1c; /* red-700 */
}

.pro-con-item-add-button {
  height: 1.75rem;
  font-size: 0.75rem;
  line-height: 1rem;
  padding: 0 0.5rem;
  border-radius: 0.375rem;
  display: flex;
  align-items: center;
}

.pro-con-item-add-button-pro {
  background-color: #f0fdf4; /* green-50 */
  color: #16a34a; /* green-600 */
  border: 1px solid #dcfce7; /* green-100 */
}

.pro-con-item-add-button-pro:hover {
  background-color: #dcfce7; /* green-100 */
}

.pro-con-item-add-button-con {
  background-color: #fef2f2; /* red-50 */
  color: #dc2626; /* red-600 */
  border: 1px solid #fee2e2; /* red-100 */
}

.pro-con-item-add-button-con:hover {
  background-color: #fee2e2; /* red-100 */
}

.pro-con-item-add-icon {
  width: 0.75rem;
  height: 0.75rem;
  margin-right: 0.25rem;
}

.pro-con-item-list {
  margin-bottom: 0.5rem;
}

.pro-con-item-list > div {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.pro-con-item-card {
  padding: 0.5rem;
  border-radius: 0.375rem;
  border: 1px solid;
}

.pro-con-item-card-pro {
  background-color: rgba(240, 253, 244, 0.3); /* green-50 with 30% opacity */
  border-color: #dcfce7; /* green-100 */
}

.pro-con-item-card-con {
  background-color: rgba(254, 242, 242, 0.3); /* red-50 with 30% opacity */
  border-color: #fee2e2; /* red-100 */
}

.pro-con-item-content {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.pro-con-item-bullet {
  margin-top: 0.375rem;
  width: 0.375rem;
  height: 0.375rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.pro-con-item-bullet-pro {
  background-color: #22c55e; /* green-500 */
}

.pro-con-item-bullet-con {
  background-color: #ef4444; /* red-500 */
}

.pro-con-item-text-container {
  flex: 1;
  min-width: 0;
}

.pro-con-item-text-content {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 0.5rem;
}

.pro-con-item-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #1f2937; /* gray-800 */
  line-height: 1.625;
  flex: 1;
  word-break: break-word;
}

/* Стили для бейджей влияния */
.pro-con-item-impact-badge {
  font-size: 0.75rem;
  line-height: 1rem;
  border: 1px solid;
  font-weight: 600;
  white-space: nowrap;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
}

.pro-con-item-impact-badge-low-pro {
  background-color: #f3f4f6; /* gray-100 */
  color: #374151; /* gray-700 */
  border-color: #e5e7eb; /* gray-200 */
}

.pro-con-item-impact-badge-low-con {
  background-color: #f3f4f6; /* gray-100 */
  color: #374151; /* gray-700 */
  border-color: #e5e7eb; /* gray-200 */
}

.pro-con-item-impact-badge-medium-pro {
  background-color: #fef9c3; /* yellow-100 */
  color: #a16207; /* yellow-700 */
  border-color: #fef08a; /* yellow-200 */
}

.pro-con-item-impact-badge-medium-con {
  background-color: #fef9c3; /* yellow-100 */
  color: #a16207; /* yellow-700 */
  border-color: #fef08a; /* yellow-200 */
}

.pro-con-item-impact-badge-high-pro {
  background-color: #dcfce7; /* green-100 */
  color: #166534; /* green-700 */
  border-color: #bbf7d0; /* green-200 */
}

.pro-con-item-impact-badge-high-con {
  background-color: #fee2e2; /* red-100 */
  color: #b91c1c; /* red-700 */
  border-color: #fecaca; /* red-200 */
}

/* Стили для слайдера */
.pro-con-item-slider-container {
  transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 0;
  overflow: hidden;
}

.pro-con-item-slider-container-open {
  max-height: 60px;
}

.pro-con-item-slider-content {
  margin-top: 0.5rem;
}

.pro-con-item-slider-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  width: 100%;
}

.pro-con-item-slider-label {
  display: flex;
  align-items: center;
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280; /* gray-500 */
  gap: 0.25rem;
}

.pro-con-item-slider-icon {
  width: 0.75rem;
  height: 0.75rem;
}

.pro-con-item-slider {
  flex: 1;
}

/* Кнопка удаления */
.pro-con-item-delete-button {
  height: 1.5rem;
  width: 1.5rem;
  color: #9ca3af; /* gray-400 */
  border-radius: 50%;
  flex-shrink: 0;
  padding: 0;
}

.pro-con-item-delete-button:hover {
  color: #ef4444; /* red-500 */
  background-color: #fef2f2; /* red-50 */
}

.pro-con-item-delete-icon {
  width: 0.75rem;
  height: 0.75rem;
}

/* Форма добавления */
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

/* Анимации */
.pro-con-item-list-enter-active,
.pro-con-item-list-leave-active {
  transition: all 300ms cubic-bezier(0.4, 0, 0.2, 1);
}

.pro-con-item-list-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.pro-con-item-list-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.pro-con-item-list-leave-active {
  position: absolute;
}

.pro-con-item-fade-enter-active,
.pro-con-item-fade-leave-active {
  transition: opacity 100ms cubic-bezier(0.4, 0, 0.2, 1);
}

.pro-con-item-fade-enter-from,
.pro-con-item-fade-leave-to {
  opacity: 0;
}
</style>