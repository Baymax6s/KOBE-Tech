import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/features/home/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

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
      meta: { guestOnly: true },
    },
    {
      path: '/articles',
      component: () => import('@/features/articles/ArticlesView.vue'),
    },
    {
      path: '/articles/new',
      component: () => import('@/features/articles/NewArticleView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { path: '/articles' }
  }
})

export default router
