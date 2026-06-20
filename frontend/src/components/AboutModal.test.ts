import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AboutModal from './AboutModal.vue'

describe('AboutModal', () => {
  it('renders content when open', () => {
    const wrapper = mount(AboutModal, { props: { open: true } })
    expect(wrapper.text()).toContain('About')
    expect(wrapper.find('.fixed').exists()).toBe(true)
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(AboutModal, { props: { open: false } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(AboutModal, { props: { open: true } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
