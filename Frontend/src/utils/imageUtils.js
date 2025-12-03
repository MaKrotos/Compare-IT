/**
 * Утилиты для обработки изображений
 */

/**
 * Обрабатывает вставку изображений через Ctrl+V
 * @param {string} currentUrls - Текущие URL изображений
 * @param {ClipboardEvent} event - Событие вставки
 * @param {function} onUpdate - Функция обратного вызова для обновления URL изображений
 */
export function handleImagePaste(currentUrls, event, onUpdate) {
  const items = (event.clipboardData || event.originalEvent.clipboardData).items;
  
  for (let i = 0; i < items.length; i++) {
    const item = items[i];
    
    // Если вставлен файл изображения
    if (item.type.indexOf('image') !== -1) {
      const blob = item.getAsFile();
      const reader = new FileReader();
      
      reader.onload = (e) => {
        // Создаем data URL изображения
        const dataURL = e.target.result;
        
        // Добавляем к существующим ссылкам
        const separator = currentUrls ? '\n' : '';
        const newUrls = currentUrls + separator + dataURL;
        onUpdate(newUrls);
      };
      
      reader.readAsDataURL(blob);
    }
    
    // Если вставлена ссылка
    if (item.type === 'text/plain') {
      item.getAsString((text) => {
        // Проверяем, является ли текст URL
        if (text.match(/^https?:\/\/.+\.(jpg|jpeg|png|gif|webp)/i)) {
          // Добавляем к существующим ссылкам
          const separator = currentUrls ? '\n' : '';
          const newUrls = currentUrls + separator + text;
          onUpdate(newUrls);
        }
      });
    }
  }
}

/**
 * Обрабатывает перетаскивание изображений
 * @param {string} currentUrls - Текущие URL изображений
 * @param {DragEvent} event - Событие перетаскивания
 * @param {function} onUpdate - Функция обратного вызова для обновления URL изображений
 */
export function handleImageDrop(currentUrls, event, onUpdate) {
  event.preventDefault();
  
  const files = event.dataTransfer.files;
  if (files && files.length > 0) {
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => {
          // Создаем data URL изображения
          const dataURL = e.target.result;
          
          // Добавляем к существующим ссылкам
          const separator = currentUrls ? '\n' : '';
          const newUrls = currentUrls + separator + dataURL;
          onUpdate(newUrls);
        };
        reader.readAsDataURL(file);
      }
    }
  }
}