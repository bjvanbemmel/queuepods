import MessageView from '@/views/MessageView.vue'
import PopulationView from '@/views/PopulationView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Populations',
      component: PopulationView,
    },
    {
      path: '/messages',
      name: 'Messages',
      component: MessageView,
    },
  ],
})

export default router
