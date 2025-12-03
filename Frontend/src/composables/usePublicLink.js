import { ref, computed } from 'vue'
import { generatePublicLink, removePublicLink } from '@/utils/collections.js'

export function usePublicLink() {
  const publicLink = ref(null)

  // Включение/выключение публичной ссылки
  const togglePublicLink = async (collectionId) => {
    try {
      if (publicLink.value) {
        // Удаляем публичную ссылку
        await removePublicLink(collectionId)
        publicLink.value = null
        alert('Публичная ссылка удалена')
      } else {
        // Создаем публичную ссылку
        const link = await generatePublicLink(collectionId)
        publicLink.value = link
        alert('Публичная ссылка создана')
      }
    } catch (error) {
      console.error('Ошибка при работе с публичной ссылкой:', error)
      alert('Не удалось изменить статус публичной ссылки: ' + error.message)
    }
  }

  // Копирование публичной ссылки в буфер обмена
  const copyPublicLink = async (fullPublicLink) => {
    try {
      await navigator.clipboard.writeText(fullPublicLink)
      alert('Ссылка скопирована в буфер обмена!')
    } catch (error) {
      console.error('Ошибка при копировании ссылки:', error)
      // Альтернативный способ копирования
      const textArea = document.createElement('textarea')
      textArea.value = fullPublicLink
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      alert('Ссылка скопирована в буфер обмена!')
    }
  }

  // Полная публичная ссылка
  const getFullPublicLink = (publicLinkValue) => {
    if (!publicLinkValue) return ''
    // В реальном приложении домен будет браться из конфигурации
    const domain = window.location.origin
    return `${domain}/public/${publicLinkValue}`
  }

  return {
    publicLink,
    togglePublicLink,
    copyPublicLink,
    getFullPublicLink
  }
}