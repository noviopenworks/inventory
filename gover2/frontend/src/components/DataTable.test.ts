import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DataTable from './DataTable.vue'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'status', label: 'Status' },
  { key: 'model', label: 'Model' },
]

const rows = [
  { name: 'PC1', status: 'Active', model: 'Dell XPS' },
  { name: 'PC2', status: 'Repair', model: 'HP Elite' },
  { name: 'PC3', status: 'Spare', model: 'Lenovo T' },
]

describe('DataTable', () => {
  it('renders column headers', () => {
    const wrapper = mount(DataTable, { props: { columns, rows: [], loading: false } })
    const headers = wrapper.findAll('th')
    expect(headers).toHaveLength(3)
    expect(headers[0].text()).toBe('Name')
    expect(headers[1].text()).toBe('Status')
    expect(headers[2].text()).toBe('Model')
  })

  it('renders row data', () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const bodyRows = wrapper.findAll('tbody tr')
    expect(bodyRows).toHaveLength(3)
    expect(bodyRows[0].text()).toContain('PC1')
    expect(bodyRows[0].text()).toContain('Active')
  })

  it('shows empty state when rows is empty', () => {
    const wrapper = mount(DataTable, { props: { columns, rows: [], loading: false } })
    const bodyRows = wrapper.findAll('tbody tr')
    expect(bodyRows).toHaveLength(1)
    expect(wrapper.find('td').text()).toContain('No items')
  })

  it('shows loading state', () => {
    const wrapper = mount(DataTable, { props: { columns, rows: [], loading: true } })
    expect(wrapper.find('tbody').text()).toContain('Loading')
  })
})
