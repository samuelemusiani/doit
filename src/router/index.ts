import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

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
    }
  ]
})

router.beforeEach(async (to) => {
  const publicPages = ['/login']
  const authRequired = !publicPages.includes(to.path)

  // Probably should use the Pinia authStore
  if (authRequired) {
    const loginUrl = 'http://localhost:8080/api/login' // This is orrible, PLS CHANGE ME
    try {
      const response = await fetch(loginUrl)
      if (!response.ok) {
        if (response.status == 401) {
          return '/login'
        } else {
          throw new Error(`Response status: ${response.status}`)
        }
      }
    } catch (error) {
      console.log(error)
    }
  }
})

export default router
