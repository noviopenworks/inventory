import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import { nextTick } from 'vue'

vi.mock('@/lib/api', () => ({
  listComputers: vi.fn().mockResolvedValue([
    { id: 1, name: 'ThinkPad', model: 'X1 Carbon', status: 'Active', purchaseDate: '', warrantyExpiry: '' },
    { id: 2, name: 'Macbook', model: 'Pro 14', status: 'Active', purchaseDate: '', warrantyExpiry: '' },
  ]),
  deleteComputer: vi.fn().mockResolvedValue(undefined),
  listUsersForDropdown: vi.fn().mockResolvedValue([]),
  listDevicesForDropdown: vi.fn().mockResolvedValue([]),
}))

import ComputersView from './ComputersView.vue'
import DataTable from '@/components/DataTable.vue'
import { useSearchStore } from '@/stores/search'

describe('ComputersView search filtering', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('passes all rows to DataTable when query is empty and only matching rows when query is set', async () => {
    const wrapper = mount(ComputersView, {
      global: {
        stubs: {
          EditPanel: true,
          ConfirmDialog: true,
          DataTable: {
            props: ['columns', 'rows', 'loading'],
            template: '<div />',
          },
        },
      },
    })

    await flushPromises()

    const dataTable = wrapper.findComponent(DataTable)

    // With empty query, both rows visible
    expect((dataTable.props('rows') as unknown[]).length).toBe(2)

    // Set query to 'thinkpad' — should filter to 1 row
    const search = useSearchStore()
    search.setQuery('thinkpad')
    await nextTick()

    const rows = dataTable.props('rows') as Record<string, unknown>[]
    expect(rows.length).toBe(1)
    expect(rows[0].name).toBe('ThinkPad')
  })
})
