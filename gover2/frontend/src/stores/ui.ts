import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getConfig, setConfig } from '@/lib/api'
import type { AppConfig } from '@/lib/api'

export const useUiStore = defineStore('ui', () => {
  const density = ref<'comfortable' | 'compact'>('comfortable')
  const darkMode = ref(false)
  const sidebarCollapsed = ref(false)
  const config = ref<AppConfig>({
    dbPath: '', density: 'comfortable', darkMode: false, expiryWarningDays: 30,
  })

  function applyTheme() {
    const root = document.documentElement
    if (darkMode.value) root.classList.add('dark')
    else root.classList.remove('dark')
  }

  async function persist() {
    config.value = { ...config.value, density: density.value, darkMode: darkMode.value }
    await setConfig(config.value)
  }

  async function loadFromConfig() {
    const cfg = await getConfig()
    config.value = cfg
    density.value = cfg.density === 'compact' ? 'compact' : 'comfortable'
    darkMode.value = cfg.darkMode
    applyTheme()
  }

  async function toggleDarkMode() {
    darkMode.value = !darkMode.value
    applyTheme()
    await persist()
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  async function setDensity(d: 'comfortable' | 'compact') {
    density.value = d
    await persist()
  }

  return {
    density,
    darkMode,
    sidebarCollapsed,
    applyTheme,
    loadFromConfig,
    toggleDarkMode,
    toggleSidebar,
    setDensity,
  }
})
