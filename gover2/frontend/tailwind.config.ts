import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Sidebar palette (constant across themes)
        'sidebar':        '#0E1520',
        'sidebar-hover':  '#131E2E',
        'sidebar-active': '#182337',
        'sidebar-border': '#141D2B',
        'sidebar-text':   '#6B82A0',
        'sidebar-hi':     '#DDE8F5',
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
        // Zebra row (themed)
        'row-alt': 'var(--c-row-alt)',
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
