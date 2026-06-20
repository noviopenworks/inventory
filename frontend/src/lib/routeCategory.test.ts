import { describe, it, expect } from 'vitest'
import { routeCategory, categoryForRoute } from './routeCategory'

describe('routeCategory map', () => {
  it('maps each known route path to its canonical category', () => {
    expect(routeCategory['/computers']).toBe('computers')
    expect(routeCategory['/smartphones']).toBe('smartphones')
    expect(routeCategory['/tablets']).toBe('tablets')
    expect(routeCategory['/windows-keys']).toBe('windowskeys')
    expect(routeCategory['/antivirus']).toBe('antivirus')
    expect(routeCategory['/other-software']).toBe('othersoftware')
    expect(routeCategory['/users']).toBe('users')
  })

  it('uses category strings matching the Go BuildCSV switch (no separators)', () => {
    expect(routeCategory['/windows-keys']).toBe('windowskeys')
    expect(routeCategory['/other-software']).toBe('othersoftware')
  })
})

describe('categoryForRoute', () => {
  it('returns the canonical category for known routes', () => {
    expect(categoryForRoute('/computers')).toBe('computers')
    expect(categoryForRoute('/windows-keys')).toBe('windowskeys')
    expect(categoryForRoute('/other-software')).toBe('othersoftware')
    expect(categoryForRoute('/users')).toBe('users')
  })

  it('returns null for /all (no single category, export hidden)', () => {
    expect(categoryForRoute('/all')).toBeNull()
  })

  it('returns null for an unknown path', () => {
    expect(categoryForRoute('/nope')).toBeNull()
    expect(categoryForRoute('')).toBeNull()
  })
})
