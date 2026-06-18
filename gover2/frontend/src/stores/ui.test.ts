import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

const getConfig = vi.fn()
const setConfig = vi.fn().mockResolvedValue(undefined)

vi.mock('@/lib/api', () => ({
  getConfig: () => getConfig(),
  setConfig: (cfg: unknown) => setConfig(cfg),
}))

import { useUiStore } from './ui'

describe('ui store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    document.documentElement.classList.remove('dark')
    getConfig.mockReset().mockResolvedValue({
      dbPath: '/x', density: 'compact', darkMode: true, expiryWarningDays: 30,
    })
    setConfig.mockClear()
  })

  it('loadFromConfig applies darkMode + density and sets the dark class', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    expect(ui.darkMode).toBe(true)
    expect(ui.density).toBe('compact')
    expect(document.documentElement.classList.contains('dark')).toBe(true)
  })

  it('toggleDarkMode flips the ref, toggles the class, and persists', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    await ui.toggleDarkMode()
    expect(ui.darkMode).toBe(false)
    expect(document.documentElement.classList.contains('dark')).toBe(false)
    expect(setConfig).toHaveBeenCalledWith(
      expect.objectContaining({ darkMode: false, density: 'compact' }),
    )
  })

  it('setDensity updates and persists', async () => {
    const ui = useUiStore()
    await ui.loadFromConfig()
    await ui.setDensity('comfortable')
    expect(ui.density).toBe('comfortable')
    expect(setConfig).toHaveBeenLastCalledWith(
      expect.objectContaining({ density: 'comfortable' }),
    )
  })
})
