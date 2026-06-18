import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from './StatusBadge.vue'

describe('StatusBadge', () => {
  const cases: [string, string][] = [
    ['active', 'bg-s-active'],
    ['repair', 'bg-s-repair'],
    ['spare', 'bg-s-spare'],
    ['retired', 'bg-s-retired'],
    ['missing', 'bg-s-missing'],
    ['expiring', 'bg-s-expiring'],
    ['expired', 'bg-s-expired'],
  ]

  cases.forEach(([status, expectedClass]) => {
    it(`status "${status}" renders class "${expectedClass}"`, () => {
      const wrapper = mount(StatusBadge, { props: { status } })
      expect(wrapper.find('span').classes()).toContain(expectedClass)
    })
  })

  it('renders the status text', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'Active' } })
    expect(wrapper.text()).toBe('Active')
  })

  it('falls back to bg-s-retired for unknown status', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'unknown' } })
    expect(wrapper.find('span').classes()).toContain('bg-s-retired')
  })
})
