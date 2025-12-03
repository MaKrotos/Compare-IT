/**
 * Утилиты для работы с подборками через API
 */

// Базовый URL для API запросов
const API_BASE_URL = '/backend';

/**
 * Получает токен авторизации из localStorage
 */
import { getToken } from './auth.js';

function getAuthToken() {
  // В реальном приложении токен будет сохраняться после аутентификации
  return getToken();
}

/**
 * Создает новую подборку
 */
export async function createCollection(collection) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections`, {
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
    console.error('Ошибка при создании подборки:', error);
    throw error;
  }
}

/**
 * Получает все подборки пользователя
 */
export async function getCollections() {
  try {
    const response = await fetch(`${API_BASE_URL}/collections`, {
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
    console.error('Ошибка при получении подборок:', error);
    return [];
  }
}

/**
 * Обновляет подборку
 */
export async function updateCollection(collection) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/update`, {
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
    console.error('Ошибка при обновлении подборки:', error);
    throw error;
  }
}

/**
 * Удаляет подборку
 */
export async function deleteCollection(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/delete`, {
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
    console.error('Ошибка при удалении подборки:', error);
    throw error;
  }
}

/**
 * Получает подборку по ID
 */
export async function getCollectionById(collectionId) {
  try {
    const response = await fetch(`${API_BASE_URL}/collections/get`, {
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
    console.error('Ошибка при получении подборки:', error);
    throw error;
  }
}

/**
 * Генерирует публичную ссылку для подборки
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
 * Удаляет публичную ссылку для подборки
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
 * Получает публичную подборку по ссылке
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
    console.error('Ошибка при получении публичной подборки:', error);
    throw error;
  }
}