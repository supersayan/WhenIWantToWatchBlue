import { createRouter, createWebHistory } from 'vue-router'
import TvView from '@/views/TvView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/tv/0'
    },
    {
      path: '/tv/:id(\\d)',
      component: TvView
    }
  ]
})

export default router
