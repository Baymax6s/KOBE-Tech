import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/features/home/HomeView.vue'
import { setApiErrorHandler } from '@/api/client'
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
    {
      path: '/not-found',
      name: 'not-found',
      component: () => import('@/features/errors/ErrorStatusView.vue'),
      props: {
        statusCode: 404,
        title: '何も見つかりません',
        message:
          'お探しのページまたはデータが見つかりませんでした。URLを確認するか、ホームへ戻ってください。',
      },
    },
    {
      path: '/server-error',
      name: 'server-error',
      component: () => import('@/features/errors/ErrorStatusView.vue'),
      props: {
        statusCode: 500,
        title: 'サーバーが落ちています',
        message:
          'サーバー側で問題が発生しています。時間をおいてから、もう一度お試しください。',
      },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'route-not-found',
      component: () => import('@/features/errors/ErrorStatusView.vue'),
      props: {
        statusCode: 404,
        title: '何も見つかりません',
        message:
          'お探しのページが見つかりませんでした。URLを確認するか、ホームへ戻ってください。',
      },
    },
  ],
})

setApiErrorHandler((status) => {
  const name = status === 404 ? 'not-found' : 'server-error'

  if (router.currentRoute.value.name === name) {
    return
  }

  void router.push({ name }).catch(() => undefined)
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
