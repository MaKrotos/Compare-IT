import { createRouter, createWebHistory } from 'vue-router'
import ComparisonView from '../views/ComparisonView.vue'
import CollectionsView from '../views/CollectionsView.vue'


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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router