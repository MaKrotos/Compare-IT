export function useImportExport() {
  // Функция для экспорта данных
  const handleExport = (items) => {
    try {
      const dataStr = JSON.stringify(items)
      const base64 = btoa(encodeURIComponent(dataStr))
      navigator.clipboard.writeText(base64).catch(err => {
        console.error('Не удалось скопировать в буфер обмена:', err)
      })
      alert('Данные скопированы в буфер обмена!')
    } catch (error) {
      console.error('Ошибка экспорта:', error)
    }
  }

  // Функция для импорта данных
  const handleImport = (importData, items, setItems, setShowImportField, setImportData) => {
    try {
      if (!importData.trim()) {
        alert('Нет данных для импорта')
        return
      }

      // Декодируем данные из base64
      const decoded = decodeURIComponent(atob(importData.trim()))
      const importedItems = JSON.parse(decoded)

      if (Array.isArray(importedItems)) {
        // Добавляем уникальные ID для каждого импортированного элемента
        const itemsWithIds = importedItems.map(item => ({
          ...item,
          id: item.id || Date.now().toString() + Math.random().toString(36).substr(2, 9),
          created_date: item.created_date || new Date().toISOString()
        }))

        // Добавляем импортированные элементы к существующим
        setItems([...items, ...itemsWithIds])
        setShowImportField(false)
        setImportData('')
        alert('Данные успешно импортированы!')
      }
    } catch (error) {
      console.error('Ошибка импорта:', error)
      alert('Ошибка импорта: ' + error.message)
    }
  }

  return {
    handleExport,
    handleImport
  }
}