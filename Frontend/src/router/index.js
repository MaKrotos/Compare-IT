import { createRouter, createWebHistory } from 'vue-router'
import ComparisonView from '../views/ComparisonView.vue'
import CollectionsView from '../views/CollectionsView.vue'
import CollectionDetailView from '../views/CollectionDetailView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: CollectionsView
  },
  {
    path: '/compare',
    name: 'Compare',
    component: ComparisonView
  },
  {
    path: '/collections/:id',
    name: 'CollectionDetail',
    component: CollectionDetailView
  },
  {
    path: '/public-collection/:publicLink',
    name: 'PublicCollectionDetail',
    component: () => import('../views/PublicCollectionDetailView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router