import { ref, computed, watch } from 'vue'

/**
 * Композабл для управления логикой элементов "Плюсы/Минусы"
 */
export function useProConItem(props) {
  const newItem = ref('')
  const isAdding = ref(false)
  const isPro = computed(() => props.type === 'pro')
  const hoveredItem = ref(null)

  // Добавим watch для отслеживания изменений newItem
  watch(newItem, (newVal) => {
    console.log('newItem changed in composable:', newVal)
  })

  /**
   * Получает класс для бейджа влияния в зависимости от значения
   * @param {number} impact - Значение влияния (1-10)
   * @returns {string} Класс для бейджа
   */
  const getImpactClass = (impact) => {
    if (impact <= 3) {
      return isPro.value ? 'pro-con-item-impact-badge-low-pro' : 'pro-con-item-impact-badge-low-con'
    }
    if (impact <= 6) {
      return isPro.value ? 'pro-con-item-impact-badge-medium-pro' : 'pro-con-item-impact-badge-medium-con'
    }
    return isPro.value ? 'pro-con-item-impact-badge-high-pro' : 'pro-con-item-impact-badge-high-con'
  }

  /**
   * Добавляет новый элемент
   */
  const addItem = () => {
   console.log('addItem called')
   console.log('newItem.value:', newItem.value)
   if (newItem.value && newItem.value.trim()) {
     const newItemObj = { text: newItem.value.trim() };
     
     // Добавляем поле влияния только для плюсов и минусов, но не для заметок
     if (props.type !== 'note') {
       newItemObj.impact = 5;
     }
     
     const newItems = [...props.items, newItemObj];
     console.log('Adding item:', newItems)
     props.onUpdate(newItems)
     newItem.value = ''
     isAdding.value = false
   }
 }

  /**
   * Удаляет элемент по индексу
   * @param {number} index - Индекс элемента для удаления
   */
  const removeItem = (index) => {
    props.onUpdate(props.items.filter((_, i) => i !== index))
  }

  /**
   * Обновляет значение влияния элемента
   * @param {number} index - Индекс элемента
   * @param {number} value - Новое значение влияния
   */
  const updateImpact = (index, value) => {
    const updated = [...props.items]
    updated[index] = { ...updated[index], impact: value }
    props.onUpdate(updated)
  }

  /**
   * Отменяет добавление нового элемента
   */
  const cancelAdd = () => {
    newItem.value = ''
    isAdding.value = false
  }

  /**
   * Устанавливает индекс элемента, над которым находится курсор
   * @param {number|null} index - Индекс элемента или null
   */
  const setHoveredItem = (index) => {
    hoveredItem.value = index
  }

  return {
    newItem,
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