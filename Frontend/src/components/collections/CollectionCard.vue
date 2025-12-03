<template>
  <div class="collection-card">
    <div class="collection-card-content">
      <div class="collection-card-header">
        <div class="collection-card-info">
          <h3 class="collection-card-title">
            {{ collection.name || 'Без названия' }}
          </h3>
          <p class="collection-card-date">
            Создана: {{ formatDate(collection.created_at) }}
          </p>
          <div v-if="collection.public_link" class="collection-card-link-indicator">
            <Link class="collection-card-link-icon" />
            <span class="collection-card-link-text">Публичная ссылка создана</span>
          </div>
        </div>
        <div class="collection-card-actions">
          <Button
            variant="ghost"
            size="icon"
            @click="editCollection(collection)"
            class="collection-card-action collection-card-edit"
          >
            <Edit3 class="collection-card-action-icon" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            @click="togglePublicLink(collection)"
            class="collection-card-action collection-card-link"
          >
            <component :is="collection.public_link ? 'Link' : 'Link2'" class="collection-card-action-icon" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            @click="togglePinCollection(collection)"
            :class="['collection-card-action', 'collection-card-pin', collection.is_pinned ? 'collection-card-pinned' : 'collection-card-unpinned']"
          >
            <component :is="collection.is_pinned ? 'Bookmark' : 'BookmarkPlus'" class="collection-card-action-icon" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            @click="deleteCollection(collection.id)"
            class="collection-card-action collection-card-delete"
          >
            <Trash2 class="collection-card-action-icon" />
          </Button>
        </div>
      </div>
      <div class="collection-card-footer">
        <Button
          @click="openCollection(collection)"
          class="collection-card-open-button"
        >
          Открыть
        </Button>
      </div>
    </div>
  </div>
</template>

<script>
import Button from '@/components/ui/Button.vue'
import { Link, Link2, Edit3, Trash2, Bookmark, BookmarkPlus } from 'lucide-vue-next'

export default {
  name: 'CollectionCard',
  components: {
    Button,
    Link,
    Link2,
    Edit3,
    Trash2,
    Bookmark,
    BookmarkPlus
  },
  props: {
    collection: {
      type: Object,
      required: true
    }
  },
  emits: ['edit', 'delete', 'toggle-public-link', 'toggle-pin', 'open'],
  setup(props, { emit }) {
    // Форматирование даты
    const formatDate = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString)
      return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      })
    }

    // Редактирование подборки
    const editCollection = (collection) => {
      emit('edit', collection)
    }

    // Удаление подборки
    const deleteCollection = (collectionId) => {
      emit('delete', collectionId)
    }

    // Включение/выключение публичной ссылки
    const togglePublicLink = (collection) => {
      emit('toggle-public-link', collection)
    }

    // Закрепление/открепление подборки
    const togglePinCollection = (collection) => {
      emit('toggle-pin', collection)
    }

    // Открытие подборки
    const openCollection = (collection) => {
      emit('open', collection)
    }

    return {
      formatDate,
      editCollection,
      deleteCollection,
      togglePublicLink,
      togglePinCollection,
      openCollection
    }
  }
}
</script>

<style scoped>
/* Карточка подборки */
.collection-card {
  background-color: white;
  border-radius: 0.75rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  transition: box-shadow 200ms ease;
}

.collection-card:hover {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1);
}

.collection-card-content {
  padding: 1.5rem;
}

.collection-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.collection-card-info {
  flex: 1;
}

.collection-card-title {
  font-size: 1.125rem;
  line-height: 1.75rem;
  font-weight: 500;
  color: #111827;
  margin-bottom: 0.25rem;
}

.collection-card-date {
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #6b7280;
}

.collection-card-link-indicator {
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #6b7280;
}

.collection-card-link-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.25rem;
  color: #6b7280;
}

.collection-card-link-text {
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.collection-card-actions {
  display: flex;
  gap: 0.25rem;
}

.collection-card-action {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
  border: none;
  cursor: pointer;
  transition: background-color 150ms ease;
}

.collection-card-edit {
  color: #6b7280;
}

.collection-card-edit:hover {
  color: #2563eb;
  background-color: #dbeafe;
}

.collection-card-link {
  color: #6b7280;
}

.collection-card-link:hover {
  color: #16a34a;
  background-color: #dcfce7;
}

.collection-card-pin {
  color: #6b7280;
}

.collection-card-pinned {
  color: #f59e0b;
}

.collection-card-unpinned:hover {
  color: #f59e0b;
  background-color: #fef3c7;
}

.collection-card-pinned:hover {
  color: #d97706;
  background-color: #fef3c7;
}

.collection-card-delete {
  color: #6b7280;
}

.collection-card-delete:hover {
  color: #dc2626;
  background-color: #fef2f2;
}

.collection-card-action-icon {
  width: 1rem;
  height: 1rem;
}

.collection-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.collection-card-open-button {
  width: 100%;
  background-color: #dbeafe;
  color: #1e40af;
  border: none;
  border-radius: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 150ms ease;
}

.collection-card-open-button:hover {
  background-color: #bfdbfe;
}
</style>