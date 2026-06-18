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
    expect(headers).toHaveLength(4)
    expect(headers[0].text()).toBe('Name')
    expect(headers[1].text()).toBe('Status')
    expect(headers[2].text()).toBe('Model')
    expect(headers[3].text()).toBe('Delete')
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

  it('emits row-click with row data when a row is clicked', async () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const bodyRows = wrapper.findAll('tbody tr')
    await bodyRows[0].trigger('click')
    expect(wrapper.emitted('row-click')).toBeTruthy()
    expect(wrapper.emitted('row-click')![0][0]).toEqual(rows[0])
  })

  it('emits delete with row data when Delete button is clicked', async () => {
    const wrapper = mount(DataTable, { props: { columns, rows, loading: false } })
    const deleteButtons = wrapper.findAll('button')
    await deleteButtons[0].trigger('click')
    expect(wrapper.emitted('delete')).toBeTruthy()
    expect(wrapper.emitted('delete')![0][0]).toEqual(rows[0])
  })
})
