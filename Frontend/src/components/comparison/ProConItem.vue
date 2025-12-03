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
        <ProConItemCard
          v-for="(item, index) in items"
          :key="index"
          :item="item"
          :is-pro="isPro"
          :readonly="readonly"
          :is-hovered="hoveredItem === index"
          @remove="() => removeItem(index)"
          @update:impact="val => updateImpact(index, val)"
          @mouseenter="() => setHoveredItem(index)"
          @mouseleave="() => setHoveredItem(null)"
          class="pro-con-item-card-wrapper"
        />
      </transition-group>
    </div>

    <transition name="pro-con-item-fade">
      <ProConItemForm
        v-if="isAdding"
        :is-pro="isPro"
        @add="addItem"
        @cancel="cancelAdd"
      />
    </transition>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import Button from '@/components/ui/Button.vue'
import ProConItemCard from './ProConItemCard.vue'
import ProConItemForm from './ProConItemForm.vue'
import { Plus } from 'lucide-vue-next'
import { useProConItem } from '@/composables/useProConItem.js'

export default {
  name: 'ProConItem',
  components: {
    Button,
    ProConItemCard,
    ProConItemForm,
    Plus
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
    const {
      isAdding,
      isPro,
      hoveredItem,
      getImpactClass,
      addItem,
      removeItem,
      updateImpact,
      cancelAdd,
      setHoveredItem
    } = useProConItem(props)

    return {
      isAdding,
      isPro,
      hoveredItem,
      getImpactClass,
      addItem,
      removeItem,
      updateImpact,
      cancelAdd,
      setHoveredItem
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

.pro-con-item-card-wrapper {
  margin-bottom: 0.25rem;
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