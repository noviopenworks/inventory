import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

const route = { path: '/computers' }
vi.mock('vue-router', () => ({
  useRoute: () => route,
}))

const exportCSV = vi.fn().mockResolvedValue(undefined)
const newDatabaseDialog = vi.fn().mockResolvedValue(undefined)
const openDatabaseDialog = vi.fn().mockResolvedValue(undefined)
const backupDatabaseDialog = vi.fn().mockResolvedValue(undefined)
vi.mock('@/lib/api', () => ({
  exportCSV: (c: string) => exportCSV(c),
  newDatabaseDialog: () => newDatabaseDialog(),
  openDatabaseDialog: () => openDatabaseDialog(),
  backupDatabaseDialog: () => backupDatabaseDialog(),
}))

import Topbar from './Topbar.vue'
import { useSearchStore } from '@/stores/search'
import { useUiStore } from '@/stores/ui'

describe('Topbar action row', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    route.path = '/computers'
    exportCSV.mockClear()
    newDatabaseDialog.mockClear()
    openDatabaseDialog.mockClear()
    backupDatabaseDialog.mockClear()
  })

  it('shows the CSV export button on a category route', () => {
    const wrapper = mount(Topbar)
    expect(wrapper.text()).toContain('CSV')
  })

  it('hides the CSV export button on /all', () => {
    route.path = '/all'
    const wrapper = mount(Topbar)
    expect(wrapper.text()).not.toContain('CSV')
  })

  it('binds the search input to the search store', async () => {
    const wrapper = mount(Topbar)
    const search = useSearchStore()
    await wrapper.find('input[type="search"]').setValue('thinkpad')
    expect(search.query).toBe('thinkpad')
  })

  it('calls exportCSV with the active route category', async () => {
    const wrapper = mount(Topbar)
    await wrapper.find('button[title="Export CSV"]').trigger('click')
    expect(exportCSV).toHaveBeenCalledWith('computers')
  })

  it('toggles dark mode via the ui store', async () => {
    const wrapper = mount(Topbar)
    const ui = useUiStore()
    const spy = vi.spyOn(ui, 'toggleDarkMode').mockResolvedValue(undefined)
    await wrapper.find('button[title="Toggle dark mode"]').trigger('click')
    expect(spy).toHaveBeenCalled()
  })

  it('opens the New/Open Database dialogs from the overflow menu', async () => {
    const wrapper = mount(Topbar)
    await wrapper.find('button[title="More"]').trigger('click')
    const items = wrapper.findAll('[data-menu-item]')
    await items[0].trigger('click') // New Database
    expect(newDatabaseDialog).toHaveBeenCalled()
    await wrapper.find('button[title="More"]').trigger('click')
    const items2 = wrapper.findAll('[data-menu-item]')
    await items2[1].trigger('click') // Open Database
    expect(openDatabaseDialog).toHaveBeenCalled()
  })

  it('backs up the database when the menu item is clicked', async () => {
    const wrapper = mount(Topbar)
    await wrapper.find('button[title="More"]').trigger('click')
    await wrapper.get('[data-testid="backup-database"]').trigger('click')
    expect(backupDatabaseDialog).toHaveBeenCalled()
  })

  it('emits open-alerts/open-about/open-license', async () => {
    const wrapper = mount(Topbar)
    await wrapper.find('button[title="Alerts"]').trigger('click')
    expect(wrapper.emitted('open-alerts')).toBeTruthy()
  })
})
