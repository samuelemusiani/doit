import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { isLoggedIn } from '@/lib/api'

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
      component: () => import('@/views/LoginView.vue'),
      meta: {
        hide_navbar: true
      }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue')
    },
    {
      path: '/admin',
      children: [
        {
          path: '',
          name: 'admin',
          component: () => import('@/views/AdminView.vue')
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/AdminUsersView.vue')
        },
        {
          path: 'users/:id',
          name: 'user_details',
          component: () => import('@/views/AdminUserDetailsView.vue')
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

router.beforeEach(async (to) => {
  const publicPages = ['/login']
  const authRequired = !publicPages.includes(to.path)

  // Probably should use the Pinia authStore
  if (authRequired) {
    try {
      const logged = await isLoggedIn()
      if (!logged) {
        return '/login'
      }
    } catch (err) {
      console.error(err)
      return '/login'
    }
  }
})

export default router
