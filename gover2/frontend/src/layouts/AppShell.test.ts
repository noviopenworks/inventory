import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { reactive } from 'vue'

const route = reactive({ path: '/computers' })
vi.mock('vue-router', () => ({
  useRoute: () => route,
}))

const getAlerts = vi.fn()
vi.mock('@/lib/api', () => ({
  getAlerts: () => getAlerts(),
}))

// Stub child components so we only exercise AppShell's own logic.
const stubs = {
  Sidebar: { template: '<div />' },
  Topbar: { template: '<div />' },
  Statusbar: { template: '<div />' },
  AlertsModal: { props: ['open', 'alerts'], template: '<div class="alerts-modal" v-if="open" />' },
  AboutModal: { props: ['open'], template: '<div class="about-modal" v-if="open" />' },
  LicenseModal: { props: ['open'], template: '<div class="license-modal" v-if="open" />' },
}

import AppShell from './AppShell.vue'
import { useSearchStore } from '@/stores/search'
import { useUiStore } from '@/stores/ui'

function mountShell() {
  return mount(AppShell, { global: { stubs } })
}

describe('AppShell', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    route.path = '/computers'
    getAlerts.mockReset().mockResolvedValue([])
    // Avoid loadFromConfig hitting the real api: stub the store method.
    const ui = useUiStore()
    vi.spyOn(ui, 'loadFromConfig').mockResolvedValue(undefined)
  })

  it('auto-opens the alerts modal on mount when getAlerts is non-empty', async () => {
    getAlerts.mockResolvedValue([
      { category: 'antivirus', id: 1, name: 'Norton', expiryDate: '2026-01-01', daysRemaining: -10, severity: 'expired' },
    ])
    const wrapper = mountShell()
    await flushPromises()
    expect(wrapper.find('.alerts-modal').exists()).toBe(true)
  })

  it('does not auto-open the alerts modal when getAlerts is empty', async () => {
    getAlerts.mockResolvedValue([])
    const wrapper = mountShell()
    await flushPromises()
    expect(wrapper.find('.alerts-modal').exists()).toBe(false)
  })

  it('clears the search query when the route changes', async () => {
    const wrapper = mountShell()
    await flushPromises()
    const search = useSearchStore()
    search.setQuery('thinkpad')
    route.path = '/users'
    await wrapper.vm.$nextTick()
    expect(search.query).toBe('')
  })
})
