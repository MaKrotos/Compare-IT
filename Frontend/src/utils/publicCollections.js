/**
 * Утилиты для работы с публичными коллекциями через API
 */

// Базовый URL для API запросов
const API_BASE_URL = '/backend';

/**
 * Получает токен авторизации из localStorage
 */
function getAuthToken() {
  // В реальном приложении токен будет сохраняться после аутентификации
  return localStorage.getItem('auth_token');
}

/**
 * Генерирует публичную ссылку для коллекции
 */
export async function generatePublicLink(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/generate-public-link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({ collection_id: collectionId })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.public_link;
  } catch (error) {
    console.error('Ошибка при генерации публичной ссылки:', error);
    throw error;
  }
}

/**
 * Удаляет публичную ссылку для коллекции
 */
export async function removePublicLink(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/remove-public-link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({ collection_id: collectionId })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Ошибка при удалении публичной ссылки:', error);
    throw error;
  }
}

/**
 * Получает публичную коллекцию по ссылке
 */
export async function getPublicCollection(publicLink) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/public`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ public_link: publicLink })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Ошибка при получении публичной коллекции:', error);
    throw error;
  }
}
