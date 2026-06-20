import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import ConfirmDialog from './ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  it('emits confirmed when Delete button clicked', async () => {
    const wrapper = mount(ConfirmDialog, { props: { message: 'Delete?' } })
    await wrapper.find('button.bg-red-600').trigger('click')
    expect(wrapper.emitted('confirmed')).toBeTruthy()
  })

  it('emits cancelled when Cancel button clicked', async () => {
    const wrapper = mount(ConfirmDialog, { props: { message: 'Delete?' } })
    await wrapper.find('button:not(.bg-red-600)').trigger('click')
    expect(wrapper.emitted('cancelled')).toBeTruthy()
  })
})
