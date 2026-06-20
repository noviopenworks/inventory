import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AlertsModal from './AlertsModal.vue'
import type { Alert } from '@/lib/api'

const alerts: Alert[] = [
  { category: 'antivirus', id: 1, name: 'Norton', expiryDate: '2026-01-01', daysRemaining: -10, severity: 'expired' },
  { category: 'other_software', id: 2, name: 'Office', expiryDate: '2026-07-01', daysRemaining: 13, severity: 'expiring' },
]

describe('AlertsModal', () => {
  it('renders a row per alert when open', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    expect(wrapper.text()).toContain('Norton')
    expect(wrapper.text()).toContain('Office')
    expect(wrapper.findAll('tbody tr')).toHaveLength(2)
  })

  it('shows empty message when alerts is empty', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts: [] } })
    expect(wrapper.text()).toContain('No alerts')
  })

  it('renders nothing when closed', () => {
    const wrapper = mount(AlertsModal, { props: { open: false, alerts } })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('emits close when Close is clicked', async () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('shows expired/expiring summary counts and the warning window', () => {
    const wrapper = mount(AlertsModal, {
      props: { open: true, alerts, warningDays: 30 },
    })
    // 1 expired (Norton, severity: 'expired') + 1 expiring (Office, severity: 'expiring')
    expect(wrapper.text()).toMatch(/1\s+expired/)
    expect(wrapper.text()).toMatch(/1\s+expiring within 30 days/)
  })

  it('prettifies the severity badge labels', () => {
    const wrapper = mount(AlertsModal, { props: { open: true, alerts } })
    const cells = wrapper.findAll('tbody td:nth-child(5)')
    expect(cells).toHaveLength(2)
    expect(cells[0].text()).toBe('Expired')
    expect(cells[1].text()).toBe('Expiring soon')
  })
})
