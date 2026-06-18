import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EditPanel from './EditPanel.vue'
import type { DropdownItem, DeviceDropdownItem } from '@/lib/types'

const users: DropdownItem[] = []
const devices: DeviceDropdownItem[] = []

const computerRow = {
  id: 1,
  name: 'TestPC',
  model: '',
  userId: null,
  status: 'active',
  purchaseDate: null,
  warrantyExpiry: null,
  notes: null,
  createdAt: '2024-01-01',
  updatedAt: '2024-01-01',
}

describe('EditPanel', () => {
  it('edit mode: name input is pre-filled from row', () => {
    const wrapper = mount(EditPanel, {
      props: {
        category: 'computers',
        mode: 'edit',
        row: computerRow as Record<string, unknown>,
        users,
        devices,
      },
    })
    const nameInput = wrapper.find('input[data-testid="field-name"]')
    expect(nameInput.exists()).toBe(true)
    expect((nameInput.element as HTMLInputElement).value).toBe('TestPC')
  })

  it('create mode: name input is empty', () => {
    const wrapper = mount(EditPanel, {
      props: {
        category: 'computers',
        mode: 'create',
        row: null,
        users,
        devices,
      },
    })
    const nameInput = wrapper.find('input[data-testid="field-name"]')
    expect(nameInput.exists()).toBe(true)
    expect((nameInput.element as HTMLInputElement).value).toBe('')
  })

  it('cancel button emits cancelled', async () => {
    const wrapper = mount(EditPanel, {
      props: {
        category: 'computers',
        mode: 'create',
        row: null,
        users,
        devices,
      },
    })
    await wrapper.find('button[data-testid="cancel-btn"]').trigger('click')
    expect(wrapper.emitted('cancelled')).toBeTruthy()
  })

  it('Escape key emits cancelled', async () => {
    const wrapper = mount(EditPanel, {
      props: {
        category: 'computers',
        mode: 'create',
        row: null,
        users,
        devices,
      },
      attachTo: document.body,
    })
    await wrapper.trigger('keydown', { key: 'Escape' })
    expect(wrapper.emitted('cancelled')).toBeTruthy()
    wrapper.unmount()
  })
})
