<template>
  <div class="comparison-card-stats">
    <div class="comparison-card-stats-grid">
      <div class="comparison-card-stat-item">
        <div class="comparison-card-stat-value comparison-card-stat-value-pro">
          {{ (item.pros || []).reduce((sum, p) => sum + (p.impact || 5), 0) }}
        </div>
        <div class="comparison-card-stat-label">{{ (item.pros || []).length }} плюсов</div>
      </div>
      <div class="comparison-card-stat-item">
        <div class="comparison-card-stat-value comparison-card-stat-value-con">
          {{ (item.cons || []).reduce((sum, c) => sum + (c.impact || 5), 0) }}
        </div>
        <div class="comparison-card-stat-label">{{ (item.cons || []).length }} минусов</div>
      </div>
      <div class="comparison-card-stat-item">
        <div class="comparison-card-stat-value comparison-card-stat-value-rating">
          {{ item.rating?.toFixed(1) || '—' }}
        </div>
        <div class="comparison-card-stat-label">общий рейтинг</div>
      </div>
      <div class="comparison-card-stat-item">
        <div class="comparison-card-stat-value comparison-card-stat-value-price">
          {{ item.priceRating?.toFixed(1) || '—' }}
        </div>
        <div class="comparison-card-stat-label">рейтинг цены</div>
      </div>
      <div class="comparison-card-stat-item">
        <div class="comparison-card-stat-value comparison-card-stat-value-pros-cons">
          {{ item.prosConsRating?.toFixed(1) || '—' }}
        </div>
        <div class="comparison-card-stat-label">рейтинг плюсов/минусов</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProductStats',
  props: {
    item: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>
/* Статистика */
.comparison-card-stats {
  padding: 1rem;
  border-top: 1px solid #f3f4f6;
  background-color: rgba(249, 250, 251, 0.5);
}

.comparison-card-stats-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 0.75rem;
  text-align: center;
}

.comparison-card-stat-item {
  display: flex;
  flex-direction: column;
}

.comparison-card-stat-value {
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 600;
}

.comparison-card-stat-value-pro {
  color: #16a34a;
}

.comparison-card-stat-value-con {
  color: #dc2626;
}

.comparison-card-stat-value-rating {
  color: #2563eb;
}

.comparison-card-stat-value-price {
  color: #f59e0b;
}

.comparison-card-stat-value-pros-cons {
  color: #8b5cf6;
}

.comparison-card-stat-label {
  font-size: 0.75rem;
  line-height: 1rem;
  color: #6b7280;
}
</style>