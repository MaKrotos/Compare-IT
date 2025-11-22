/**
 * Утилиты для работы с коллекциями сравнений через API
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
 * Создает новую коллекцию сравнений
 */
export async function createComparisonCollection(collection) {
  try {
    const response = await fetch(`${API_BASE_URL}/comparisons`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify(collection)
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Ошибка при создании коллекции сравнений:', error);
    throw error;
  }
}

/**
 * Получает все коллекции сравнений пользователя
 */
export async function getComparisonCollections() {
  try {
    const response = await fetch(`${API_BASE_URL}/comparisons`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAuthToken()}`
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data || [];
  } catch (error) {
    console.error('Ошибка при получении коллекций сравнений:', error);
    return [];
  }
}

/**
 * Получает коллекцию сравнений по ID
 */
export async function getComparisonCollectionById(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/comparisons/get`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({ id: collectionId })
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Ошибка при получении коллекции сравнений:', error);
    throw error;
  }
}

/**
 * Обновляет коллекцию сравнений
 */
export async function updateComparisonCollection(collection) {
  try {
    const response = await fetch(`${API_BASE_URL}/comparisons/update`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify(collection)
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Ошибка при обновлении коллекции сравнений:', error);
    throw error;
  }
}

/**
 * Удаляет коллекцию сравнений
 */
export async function deleteComparisonCollection(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/comparisons/delete`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify({ id: collectionId })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return true;
  } catch (error) {
    console.error('Ошибка при удалении коллекции сравнений:', error);
    throw error;
  }
}