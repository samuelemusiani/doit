import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { LOGIN_URL } from '@/consts'
import { getCurrentUser, isLoggedIn } from '@/lib/api'

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
      component: () => import('../views/LoginView.vue')
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue')
    },
    {
      path: '/admin',
      children: [
        {
          path: '',
          name: 'admin',
          component: () => import('../views/AdminView.vue')
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../views/AdminUsersView.vue')
        },
        {
          path: 'users/:id',
          name: 'user_details',
          component: () => import('../views/AdminUserDetailsView.vue')
        }
      ]
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
    }
  }
})

export default router
