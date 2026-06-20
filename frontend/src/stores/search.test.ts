import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSearchStore } from './search'

describe('search store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with an empty query', () => {
    const search = useSearchStore()
    expect(search.query).toBe('')
  })

  it('setQuery updates the query', () => {
    const search = useSearchStore()
    search.setQuery('thinkpad')
    expect(search.query).toBe('thinkpad')
  })

  it('clear resets the query to empty', () => {
    const search = useSearchStore()
    search.setQuery('thinkpad')
    search.clear()
    expect(search.query).toBe('')
  })
})
