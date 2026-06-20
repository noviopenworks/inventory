import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Sidebar palette (themed)
        'sidebar':        'var(--c-sidebar)',
        'sidebar-hover':  'var(--c-sidebar-hover)',
        'sidebar-active': 'var(--c-sidebar-active)',
        'sidebar-border': 'var(--c-sidebar-border)',
        'sidebar-text':   'var(--c-sidebar-text)',
        'sidebar-hi':     'var(--c-sidebar-hi)',
        'sidebar-accent': '#2563EB',
        // Page palette (themed)
        'page':    'var(--c-page)',
        'surface': 'var(--c-surface)',
        'border':  'var(--c-border)',
        // Status colors (semantic, constant)
        's-active':   '#2563EB',
        's-repair':   '#EA580C',
        's-spare':    '#7C3AED',
        's-retired':  '#9CA3AF',
        's-missing':  '#1E293B',
        's-expiring': '#D97706',
        's-expired':  '#DC2626',
        // Table header (themed)
        'th-bg':   'var(--c-th-bg)',
        'th-text': 'var(--c-th-text)',
        // Typography (themed)
        'text-primary':   'var(--c-text-primary)',
        'text-secondary': 'var(--c-text-secondary)',
        'text-tertiary':  'var(--c-text-tertiary)',
        // Interactive
        'accent-hover': '#1D4ED8',
        'accent':       '#2563EB',
        // Page structure (themed)
        'border-light': 'var(--c-border-light)',
        // Zebra row + row hover (themed)
        'row-alt':   'var(--c-row-alt)',
        'row-hover': 'var(--c-row-hover)',
        // Status row highlight backgrounds (themed)
        'row-expiring': 'var(--c-row-expiring)',
        'row-expired':  'var(--c-row-expired)',
      },
      fontFamily: {
        ui:   ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'system-ui', 'sans-serif'],
        mono: ['"Cascadia Code"', '"Fira Code"', 'Consolas', '"Courier New"', 'monospace'],
      },
      borderRadius: {
        badge: '2px',
      },
      width: {
        sidebar: '200px',
      },
      height: {
        topbar: '48px',
      },
    },
  },
  plugins: [],
} satisfies Config
