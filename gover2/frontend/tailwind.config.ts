import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        // Sidebar palette
        'sidebar':        '#0E1520',
        'sidebar-hover':  '#131E2E',
        'sidebar-active': '#182337',
        'sidebar-border': '#141D2B',
        'sidebar-text':   '#6B82A0',
        'sidebar-hi':     '#DDE8F5',
        'sidebar-accent': '#2563EB',
        // Page palette
        'page':    '#EEF1F7',
        'surface': '#FFFFFF',
        'border':  '#DDE1EC',
        // Status colors
        's-active':   '#2563EB',
        's-repair':   '#EA580C',
        's-spare':    '#7C3AED',
        's-retired':  '#9CA3AF',
        's-missing':  '#1E293B',
        's-expiring': '#D97706',
        's-expired':  '#DC2626',
        // Table header
        'th-bg':   '#0E1520',
        'th-text': '#8FA5BF',
        // Typography
        'text-primary':   '#0F172A',
        'text-secondary': '#64748B',
        'text-tertiary':  '#94A3B8',
        // Interactive
        'accent-hover': '#1D4ED8',
        // Page structure
        'border-light': '#EAECF4',
        // Status row highlight backgrounds
        'row-expiring': '#FFFBEB',
        'row-expired':  '#FFF5F5',
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
