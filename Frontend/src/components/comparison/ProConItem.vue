<template>
  <div class="pro-con-item-container">
    <div class="pro-con-item-header">
      <h4 :class="['pro-con-item-title', isPro ? 'pro-con-item-title-pro' : 'pro-con-item-title-con']">
        {{ isPro ? 'Плюсы' : 'Минусы' }}
      </h4>
      <Button
        v-if="!isAdding"
        variant="ghost"
        size="sm"
        @click="isAdding = true"
        :class="['pro-con-item-add-button', isPro ? 'pro-con-item-add-button-pro' : 'pro-con-item-add-button-con']"
      >
        <Plus class="w-3 h-3 mr-1" />
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
                  :class="['pro-con-item-impact-badge', getImpactColor(item.impact || 5)]"
                >
                  {{ item.impact || 5 }}/10
                </Badge>
              </div>
              <div class="pro-con-item-slider-container" :class="{ 'pro-con-item-slider-container-open': hoveredItem === index }">
                <div class="pro-con-item-slider-content">
                  <div class="pro-con-item-slider-header">
                    <div class="pro-con-item-slider-label">
                      <TrendingUp class="w-3 h-3" />
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
              variant="ghost"
              size="icon"
              @click="removeItem(index)"
              class="pro-con-item-delete-button"
            >
              <Trash2 class="w-3 h-3" />
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
          <Check class="w-4 h-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          @click="cancelAdd"
          class="pro-con-item-cancel-button"
        >
          <X class="w-4 h-4" />
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
    }
  },
  setup(props) {
    const newItem = ref('')
    const isAdding = ref(false)
    const isPro = computed(() => props.type === 'pro')
    const hoveredItem = ref(null)

    const getImpactColor = (impact) => {
      // Для плюсов используем зеленый цвет, для минусов - красный
      if (impact <= 3) {
        return isPro.value
          ? 'bg-green-100 text-green-700 border-green-200'
          : 'bg-gray-100 text-gray-700 border-gray-200'
      }
      if (impact <= 6) {
        return isPro.value
          ? 'bg-green-100 text-green-700 border-green-200'
          : 'bg-yellow-100 text-yellow-700 border-yellow-200'
      }
      return isPro.value
        ? 'bg-green-100 text-green-700 border-green-200'
        : 'bg-red-100 text-red-700 border-red-200'
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
      getImpactColor,
      addItem,
      removeItem,
      updateImpact,
      cancelAdd
    }
  }
}
</script>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.1s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.pro-con-item-slide-container {
  @apply transition-all duration-300 ease-in-out max-h-0 overflow-hidden;
}

.pro-con-item-slide-container-open {
  @apply max-h-[60px];
}
</style>
