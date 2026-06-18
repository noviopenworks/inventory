import { describe, it, expect } from 'vitest'
import { matchesQuery } from './filter'

const columns = [{ key: 'name' }, { key: 'model' }, { key: 'status' }]
const row = { name: 'ThinkPad', model: 'X1 Carbon', status: 'active', notes: 'secret' }

describe('matchesQuery', () => {
  it('returns true for empty query', () => {
    expect(matchesQuery(row, columns, '')).toBe(true)
    expect(matchesQuery(row, columns, '   ')).toBe(true)
  })

  it('matches case-insensitively across visible columns', () => {
    expect(matchesQuery(row, columns, 'thinkpad')).toBe(true)
    expect(matchesQuery(row, columns, 'CARBON')).toBe(true)
    expect(matchesQuery(row, columns, 'active')).toBe(true)
  })

  it('returns false when no visible column matches', () => {
    expect(matchesQuery(row, columns, 'dell')).toBe(false)
  })

  it('ignores columns not in the columns list', () => {
    expect(matchesQuery(row, columns, 'secret')).toBe(false)
  })

  it('treats null/undefined cell values as empty string', () => {
    expect(matchesQuery({ name: null }, [{ key: 'name' }], 'x')).toBe(false)
  })
})
