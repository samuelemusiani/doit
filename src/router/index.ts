import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { LOGIN_URL } from '@/consts'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue')
    }
  ]
})

router.beforeEach(async (to) => {
  const publicPages = ['/login']
  const authRequired = !publicPages.includes(to.path)

  // Probably should use the Pinia authStore
  if (authRequired) {
    try {
      const response = await fetch(LOGIN_URL, {
        credentials: 'include' // Used to send cookies; DOTO Check if this should be in production
      })
      if (!response.ok) {
        if (response.status == 401) {
          return '/login'
        } else {
          throw new Error(`Response status: ${response.status}`)
        }
      }
    } catch (error) {
      console.error(error)
    }
  }
})

export default router
