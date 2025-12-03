import { ref } from 'vue'

export function useRating() {
  const ratingWeights = ref({
    priceRatingWeight: 20,
    prosConsRatingWeight: 80
  })

  // Обработчик сохранения весов рейтингов
  const handleRatingWeightsSave = (weights, setRatingWeights) => {
    setRatingWeights(weights)
    // Пересчитываем рейтинги с новыми весами
    // recalculateAllRatings() // Эта функция будет в основном компоненте
    // Сохраняем коллекцию с новыми настройками рейтингов
    // saveCurrentCollection() // Эта функция будет в основном компоненте
  }

  // Пересчет всех рейтингов
  const recalculateAllRatings = (items, weights) => {
    const updatedItems = items.map(item =>
      calculateRating(item, items, weights)
    )
    return updatedItems
  }

  const createItem = (itemData, allItems, weights) => {
    const newItem = {
      ...itemData,
      id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
      created_date: new Date().toISOString(),
    }

    // Добавляем новый элемент
    const updatedItems = [newItem, ...allItems]

    // Пересчитываем рейтинги для всех элементов
    const itemsWithRatings = updatedItems.map(item =>
      calculateRating(item, updatedItems, weights)
    )

    return { newItem, itemsWithRatings }
  }

  const calculateRating = (itemData, allItems, weights) => {
    // Конфигурационные параметры для настройки веса рейтинга
    const RATING_CONFIG = {
      // Максимальный рейтинг
      MAX_RATING: 100,
      // Вес рейтинга по плюсам/минусам (в процентах от максимального рейтинга)
      PROS_CONS_WEIGHT: weights.prosConsRatingWeight,
      // Вес рейтинга по цене (в процентах от максимального рейтинга)
      PRICE_WEIGHT: weights.priceRatingWeight
    }

    const pros = itemData.pros || []
    const cons = itemData.cons || []
    const itemPrice = itemData.price || 0

    // Сумма положительных баллов
    const prosScore = pros.reduce((sum, pro) => sum + (pro.impact || 5), 0)

    // Сумма отрицательных баллов
    const consScore = cons.reduce((sum, con) => sum + (con.impact || 5), 0)

    // Базовый рейтинг (0-100)
    const maxPossibleScore = Math.max((pros.length + cons.length) * 10, 1)
    const baseRating = ((prosScore - consScore) / maxPossibleScore) * 50 + 50

    // Корректировка на цену (чем выше цена, тем меньше рейтинг)
    // Находим максимальную и минимальную цену среди всех товаров для нормализации
    let priceImpact = 0
    if (allItems.length > 0 && itemPrice > 0) {
      const maxPrice = Math.max(...allItems.map(item => item.price || 0))
      const minPrice = Math.min(...allItems.map(item => item.price || 0))
      if (maxPrice > minPrice) {
        // Чем выше цена по сравнению с минимальной, тем меньше рейтинг
        // Чем ниже цена по сравнению с максимальной, тем выше рейтинг
        const priceRange = maxPrice - minPrice
        const priceRatio = (maxPrice - itemPrice) / priceRange
        // Применяем вес цены к максимальному бонусу/штрафу
        priceImpact = priceRatio * RATING_CONFIG.PRICE_WEIGHT
      }
    }

    // Учет относительности плюсов/минусов
    let prosConsImpact = 0
    if (allItems.length > 0) {
      // Средние значения по всем карточкам
      const totalProsScore = allItems.reduce((sum, item) => sum + (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0), 0)
      const totalConsScore = allItems.reduce((sum, item) => sum + (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0), 0)
      const avgProsScore = totalProsScore / allItems.length
      const avgConsScore = totalConsScore / allItems.length

      // Сравнение с лучшей карточкой
      const maxProsScore = Math.max(...allItems.map(item => (item.pros || []).reduce((s, pro) => s + (pro.impact || 5), 0)))
      const maxConsScore = Math.max(...allItems.map(item => (item.cons || []).reduce((s, con) => s + (con.impact || 5), 0)))

      // Нормализованные значения (0-1)
      const normalizedPros = maxProsScore > 0 ? prosScore / maxProsScore : 0
      const normalizedCons = maxConsScore > 0 ? consScore / maxConsScore : 0

      // Влияние плюсов (чем больше плюсов по сравнению со средним и чем больше по сравнению с лучшей карточкой, тем выше рейтинг)
      const prosVsAvg = avgProsScore > 0 ? prosScore / avgProsScore : 1
      const prosBonus = Math.min(1.5, Math.max(0.5, prosVsAvg)) // Ограничиваем влияние

      // Влияние минусов (чем меньше минусов по сравнению со средним и чем меньше по сравнению с худшей карточкой, тем выше рейтинг)
      const consVsAvg = avgConsScore > 0 ? consScore / avgConsScore : 0
      const consBonus = Math.min(1.5, Math.max(0.5, 1 - consVsAvg)) // Ограничиваем влияние

      // Комбинированное влияние с учетом веса плюсов/минусов
      prosConsImpact = ((normalizedPros * prosBonus) - (normalizedCons * consBonus)) * RATING_CONFIG.PROS_CONS_WEIGHT * 0.15
    }

    const finalRating = Math.max(0, Math.min(RATING_CONFIG.MAX_RATING, baseRating + priceImpact + prosConsImpact))

    return { ...itemData, rating: Math.round(finalRating * 10) / 10 }
  }

  const updateItem = (id, data, allItems, weights) => {
    // Обновляем элемент
    const updatedItems = allItems.map(item =>
      item.id === id ? { ...item, ...data } : item
    )

    // Пересчитываем рейтинги для всех элементов
    const itemsWithRatings = updatedItems.map(item =>
      calculateRating(item, updatedItems, weights)
    )

    return itemsWithRatings
  }

  const deleteItem = (id, allItems) => {
    return allItems.filter(item => item.id !== id)
  }

  return {
    ratingWeights,
    handleRatingWeightsSave,
    recalculateAllRatings,
    createItem,
    calculateRating,
    updateItem,
    deleteItem
  }
}