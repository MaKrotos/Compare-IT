/**
 * Извлекает данные предварительного просмотра с веб-сайта через бэкенд
 */
export async function fetchPreviewData(url) {
  console.log('🔍 Начинаем получение данных через бэкенд для URL:', url);
  
  try {
    // Используем прокси через Vite для обращения к бэкенду
    const response = await fetch(`/backend/preview?url=${encodeURIComponent(url)}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    console.log('✅ Успешно получили данные через бэкенд');
    console.log('Данные:', data);
    return data;
  } catch (error) {
    console.error('❌ Ошибка при получении данных через бэкенд:', error);
    // Возвращаем минимальные данные вместо ошибки
    return {
      title: 'Ошибка загрузки',
      description: 'Не удалось загрузить данные',
      image: '',
      url: url
    };
  }
}

/**
 * Извлекает имя товара из URL или заголовка страницы
 */
export async function extractProductName(url) {
  console.log('🔍 Извлечение названия товара для:', url);
  
  try {
    const previewData = await fetchPreviewData(url);
    return previewData.title || 'Без названия';
  } catch (error) {
    console.error('❌ Ошибка при извлечении имени товара:', error);
    return 'Без названия';
  }
}

/**
 * Кэширование превью данных
 */
const previewCache = new Map();
const CACHE_DURATION = 5 * 60 * 1000; // 5 минут

/**
 * Получает данные превью с кэшированием
 */
export async function getCachedPreviewData(url) {
  // Проверяем кэш
  if (previewCache.has(url)) {
    const cached = previewCache.get(url);
    // Проверяем, не истекло ли время кэша
    if (Date.now() - cached.timestamp < CACHE_DURATION) {
      console.log('✅ Используем кэшированные данные для URL:', url);
      return cached.data;
    } else {
      // Удаляем устаревший кэш
      previewCache.delete(url);
    }
  }
  
  // Получаем новые данные
  const data = await fetchPreviewData(url);
  
  // Сохраняем в кэш
  previewCache.set(url, {
    data: data,
    timestamp: Date.now()
  });
  
  return data;
}