import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/computers' },
  {
    path: '/computers',
    component: () => import('@/features/assets/ComputersView.vue'),
  },
  {
    path: '/smartphones',
    component: () => import('@/features/assets/SmartphonesView.vue'),
  },
  {
    path: '/tablets',
    component: () => import('@/features/assets/TabletsView.vue'),
  },
  {
    path: '/all',
    component: () => import('@/features/assets/AllAssetsView.vue'),
  },
  {
    path: '/windows-keys',
    component: () => import('@/features/licenses/WindowsKeysView.vue'),
  },
  {
    path: '/antivirus',
    component: () => import('@/features/licenses/AntivirusView.vue'),
  },
  {
    path: '/other-software',
    component: () => import('@/features/licenses/OtherSoftwareView.vue'),
  },
  {
    path: '/users',
    component: () => import('@/features/users/UsersView.vue'),
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
