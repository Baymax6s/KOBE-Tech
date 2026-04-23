import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/features/home/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: HomeView,
    },
    {
      path: '/about',
      // 遅延読み込み: アクセス時に初めてJSを読み込むので初期表示が速くなる
      component: () => import('@/features/about/AboutView.vue'),
    },
     {
      path: '/login',
      // 遅延読み込み: アクセス時に初めてJSを読み込むので初期表示が速くなる
      component: () => import('@/features/auth/LoginView.vue'),
    },
  ],
})

export default router
