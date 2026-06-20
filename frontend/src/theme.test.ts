import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

const root = resolve(__dirname, '..')
const css = readFileSync(resolve(root, 'src/style.css'), 'utf8')
const config = readFileSync(resolve(root, 'tailwind.config.ts'), 'utf8')
const dataTable = readFileSync(resolve(root, 'src/components/DataTable.vue'), 'utf8')

function block(name: ':root' | '.dark'): string {
  const sel = name === ':root' ? ':root' : '\\.dark'
  const m = css.match(new RegExp(`${sel}\\s*\\{([^}]*)\\}`))
  return m ? m[1] : ''
}

describe('theme tokens', () => {
  it('every var(--c-*) token used in tailwind.config is defined in BOTH :root and .dark', () => {
    const used = [...config.matchAll(/var\((--c-[a-z0-9-]+)\)/g)].map((m) => m[1])
    expect(used.length).toBeGreaterThan(0)
    const light = block(':root')
    const dark = block('.dark')
    const missing: string[] = []
    for (const v of new Set(used)) {
      if (!light.includes(`${v}:`)) missing.push(`${v} (missing in :root)`)
      if (!dark.includes(`${v}:`)) missing.push(`${v} (missing in .dark)`)
    }
    expect(missing).toEqual([])
  })

  it('the sidebar palette is theme-aware (uses CSS variables, not static hex)', () => {
    for (const token of ['sidebar', 'sidebar-hover', 'sidebar-active', 'sidebar-border', 'sidebar-text', 'sidebar-hi']) {
      expect(config).toMatch(new RegExp(`'${token}':\\s*'var\\(--c-${token}\\)'`))
    }
  })

  it('the table header is light in light mode and dark in dark mode', () => {
    // th-bg must differ between the two palettes (regression: both were dark navy)
    const lightTh = block(':root').match(/--c-th-bg:\s*(#[0-9A-Fa-f]{6})/)?.[1]
    const darkTh = block('.dark').match(/--c-th-bg:\s*(#[0-9A-Fa-f]{6})/)?.[1]
    expect(lightTh).toBeDefined()
    expect(darkTh).toBeDefined()
    expect(lightTh!.toLowerCase()).not.toBe(darkTh!.toLowerCase())
    // light header should be a light color (high luminance) — first hex pair > 0x80
    expect(parseInt(lightTh!.slice(1, 3), 16)).toBeGreaterThan(0x80)
  })

  it('no opacity modifier is applied to the sidebar token (breaks with hex CSS vars)', () => {
    expect(dataTable).not.toMatch(/bg-sidebar\/\d/)
  })
})
