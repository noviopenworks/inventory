import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LicenseModal from './LicenseModal.vue'

describe('LicenseModal', () => {
  it('renders content when open', () => {
    const wrapper = mount(LicenseModal, { props: { open: true } })
    expect(wrapper.text()).toContain('License')
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(LicenseModal, { props: { open: false } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(LicenseModal, { props: { open: true } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
