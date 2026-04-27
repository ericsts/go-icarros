import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    component: () => import('../views/LandingView.vue'),
    meta: { public: true, guestOnly: true }
  },
  {
    path: '/login',
    component: () => import('../views/LoginView.vue'),
    meta: { public: true, guestOnly: true }
  },
  {
    path: '/register',
    component: () => import('../views/RegisterView.vue'),
    meta: { public: true, guestOnly: true }
  },
  { path: '/auctions',     component: () => import('../views/AuctionsView.vue') },
  { path: '/auctions/:id', component: () => import('../views/AuctionDetailView.vue') },
  { path: '/cars',         component: () => import('../views/CarsView.vue') },
  {
    path: '/admin/users',
    component: () => import('../views/UsersView.vue'),
    meta: { admin: true }
  },
  {
    path: '/admin/logs',
    component: () => import('../views/LogsView.vue'),
    meta: { admin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(to => {
  const auth = useAuthStore()
  if (to.meta.guestOnly && auth.isAuthenticated) return '/auctions'
  if (!to.meta.public   && !auth.isAuthenticated) return '/login'
  if (to.meta.admin     && !auth.isAdmin)         return '/auctions'
})

export default router
